// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package job

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/douyu/juno-agent/pkg/job/etcd"
	"github.com/douyu/juno-agent/util"
	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/sony/sonyflake"
)

// Node 执行 cron 命令服务的结构体
type Worker struct {
	*Config
	*etcdv3.Client
	*Cron

	ID             string
	ImmediatelyRun bool // 是否立即执行

	jobs        Jobs // 和结点相关的任务
	cmds        map[string]*Cmd
	runningJobs map[string]context.CancelFunc

	done      chan struct{}
	taskIdGen *sonyflake.Sonyflake
}

func NewWorker(conf *Config) (w *Worker) {
	w = &Worker{
		Config:         conf,
		ID:             conf.HostName,
		Client:         etcdv3.StdConfig("default").Build(),
		ImmediatelyRun: false,
		cmds:           make(map[string]*Cmd),
		runningJobs:    make(map[string]context.CancelFunc),
		done:           make(chan struct{}),
		taskIdGen:      sonyflake.NewSonyflake(sonyflake.Settings{}), // default setting
	}

	w.Cron = newCron(w)

	w.logger.Info("agent info :", xlog.String("name", conf.AppIP+":"+conf.HostName))

	return
}

func (w *Worker) Run() error {
	w.logger.Info("Worker run...")

	w.Cron.Run()

	lockWCh := w.Client.Watch(context.Background(), LockKeyPrefix, clientv3.WithPrefix())
	onceWch := w.Client.Watch(context.Background(), OnceKeyPrefix+w.HostName, clientv3.WithPrefix())
	procWch := w.Client.Watch(context.Background(), ProcKeyPrefix, clientv3.WithPrefix())
	jobWch, err := etcd.WatchPrefix(w.Client, context.Background(), JobsKeyPrefix)
	if err != nil {
		panic(err)
	}

	// load prev jobs
	w.loadJobs(jobWch.IncipientKeyValues())

	for {
		select {
		case ev := <-lockWCh:
			w.handleLockEv(ev)

		case ev := <-onceWch:
			w.handleOnceEv(ev)

		case ev := <-procWch:
			w.handleProcEv(ev)

		case ev := <-jobWch.C():
			w.handleJobEv(ev)

		}
	}

	return nil
}

func (w *Worker) loadJobs(keyValue []*mvccpb.KeyValue) {
	w.jobs = make(map[string]*Job)
	if len(keyValue) == 0 {
		return
	}

	for _, val := range keyValue {
		job, err := w.GetJobContentFromKv(val.Key, val.Value)
		if err != nil {
			w.logger.Warnf("job[%s] is invalid: %s", val.Key, err.Error())
			continue
		}

		job.runOn = w.ID
		if _, ok := w.jobs[job.ID]; !ok {
			w.addJob(job)
		}
	}

	return
}

func (w *Worker) delJob(id string) {
	job, ok := w.jobs[id]
	// 之前此任务没有在当前结点执行
	if !ok {
		return
	}

	xlog.Error("Worker.delJob:delete a job", xlog.String("jobId", id))

	delete(w.jobs, id)
	job.Unlock()

	cmds := job.Cmds()
	if len(cmds) == 0 {
		return
	}

	for _, cmd := range cmds {
		w.delCmd(cmd)
	}
	return
}

func (w *Worker) modJob(job *Job) {
	oJob, ok := w.jobs[job.ID]
	if !ok {
		w.addJob(job)
		return
	}

	job.Worker = w
	job.locked = oJob.locked

	if util.InStringArray(job.Nodes, w.HostName) < 0 {
		w.delJob(job.ID)
		return
	}

	if job.JobType != oJob.JobType { // if job-type modified
		if job.JobType == TypeNormal {
			if job.mutex != nil && job.locked {
				job.mutex.Unlock()
			}
		} else if job.JobType == TypeAlone {
			w.delJob(job.ID)
			w.addJob(job)
			return
		}
	}

	prevCmds := oJob.Cmds()
	*oJob = *job
	cmds := oJob.Cmds()

	// 筛选出需要删除的任务
	for id, cmd := range cmds {
		w.modCmd(cmd)
		delete(prevCmds, id)
	}

	for _, cmd := range prevCmds {
		w.delCmd(cmd)
	}
}

func (w *Worker) addJob(job *Job) {
	job.Worker = w

	if util.InStringArray(job.Nodes, w.HostName) < 0 {
		// ignore
		xlog.Info("Worker.addJob: Nodes do not contain current node, skip it.", xlog.String("jobId", job.ID))
		return
	}

	if job.JobType == TypeAlone {
		err := job.Lock()
		if err != nil {
			xlog.Info("failed to lock job. ignore it", xlog.String("jobId", job.ID))
			return
		}
	}

	xlog.Info("Worker.addJob: add a job", xlog.String("jobId", job.ID), xlog.Any("job", job))

	// 添加任务到当前节点
	w.jobs[job.ID] = job

	cmds := job.Cmds()
	if len(cmds) == 0 {
		return
	}

	for _, cmd := range cmds {
		w.addCmd(cmd)
	}
	return
}

func (w *Worker) delCmd(cmd *Cmd) {
	c, ok := w.cmds[cmd.GetID()]
	if ok {
		delete(w.cmds, cmd.GetID())
		w.Cron.Remove(c.schEntryID)
	}
	w.logger.Infof("job[%s] rule[%s] timer[%s] has deleted", cmd.Job.ID, cmd.Timer.ID, cmd.Timer.Cron)
}

func (w *Worker) modCmd(cmd *Cmd) {
	c, ok := w.cmds[cmd.GetID()]
	if !ok {
		w.addCmd(cmd)
		return
	}

	entryID := c.schEntryID
	sch := c.Timer.Cron
	*c = *cmd
	c.schEntryID = entryID

	// 节点执行时间改变，更新 cron
	// 否则不用更新 cron
	if c.Timer.Cron != sch {
		w.Cron.Remove(entryID)
		c.schEntryID = w.Cron.Schedule(c.Timer.Schedule, c)
	}

	w.logger.Infof("job[%s]rule[%s] timer[%s] has updated", c.Job.ID, c.Timer.ID, c.Timer.Cron)
}

func (w *Worker) addCmd(cmd *Cmd) {
	cmd.schEntryID = w.Cron.Schedule(cmd.Timer.Schedule, cmd)
	w.cmds[cmd.GetID()] = cmd

	w.logger.Infof("job[%s] rule[%s] timer[%s] has added",
		cmd.Job.ID, cmd.Timer.ID, cmd.Timer.Cron)
	return
}

func (w *Worker) GetJobContentFromKv(key []byte, value []byte) (*Job, error) {
	var err error
	job := &Job{}

	if err := json.Unmarshal(value, job); err != nil {
		w.logger.Warnf("job[%s] unmarshal err: %s", key, err.Error())
		return nil, err
	}
	if err := job.ValidRules(); err != nil {
		w.logger.Warnf("valid rules [%s] err: %s", key, err.Error())
		return nil, err
	}

	job.mutex, err = w.Client.NewMutex(LockKeyPrefix+job.ID, concurrency.WithTTL(30))
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (w *Worker) GetOnceJobFromKv(key []byte, value []byte) (*OnceJob, error) {
	job := &OnceJob{}

	if err := json.Unmarshal(value, job); err != nil {
		w.logger.Warnf("job[%s] unmarshal err: %s", key, err.Error())
		return nil, err
	}
	if err := job.ValidRules(); err != nil {
		w.logger.Warnf("valid rules [%s] err: %s", key, err.Error())
		return nil, err
	}

	return job, nil
}

func (w *Worker) KillExecutingProc(process *Process) {
	pid, _ := strconv.Atoi(process.ID)
	if err := killProcess(pid); err != nil {
		w.logger.Warnf("process:[%d] force kill failed, error:[%s]", pid, err)
		return
	}
}

func (w *Worker) handleLockEv(ev clientv3.WatchResponse) {
	for _, ev := range ev.Events {
		switch {
		case ev.Type == clientv3.EventTypeDelete:
			// watch deleted job and try to lock that job
			jobId := getJobIDFromLockKey(string(ev.Kv.Key))

			resp, err := w.Client.Get(context.Background(), LockKeyPrefix+jobId+"/", clientv3.WithPrefix())
			if err != nil {
				return
			}
			if resp.Count > 0 {
				return
			}

			resp, err = w.Client.Get(context.Background(), JobsKeyPrefix+jobId)
			if err != nil {
				return
			}
			if len(resp.Kvs) == 0 {
				return
			}

			jobKv := resp.Kvs[0]
			job, err := w.GetJobContentFromKv(jobKv.Key, jobKv.Value)
			if err != nil {
				return
			}

			job.runOn = w.ID
			if _, ok := w.jobs[jobId]; !ok { // job not exists
				w.addJob(job)
			}
		}
	}
}

func (w *Worker) handleOnceEv(ev clientv3.WatchResponse) {
	for _, event := range ev.Events {
		switch {
		case event.IsCreate(), event.IsModify():
			w.logger.Info("once task...")

			job, err := w.GetOnceJobFromKv(event.Kv.Key, event.Kv.Value)
			if err != nil {
				xlog.Error("get job from kv failed", xlog.String("err", err.Error()))
				continue
			}

			job.Worker = w
			go job.RunWithRecovery(WithTaskID(job.TaskID))
		}
	}
}

func (w *Worker) handleProcEv(ev clientv3.WatchResponse) {
	for _, event := range ev.Events {
		switch {
		case event.IsModify():
			w.logger.Info("exec process task...")

			key := string(event.Kv.Key)
			process, err := GetProcFromKey(key)
			if err != nil {
				w.logger.Warnf("err: %s, kv: %s", err.Error(), event.Kv.String())
				continue
			}

			if process.NodeID != w.ID {
				continue
			}

			val := string(event.Kv.Value)
			pv := &ProcessVal{}
			err = json.Unmarshal([]byte(val), pv)
			if err != nil {
				continue
			}
			process.ProcessVal = *pv
			if process.Killed {
				w.KillExecutingProc(process)
			}
		}
	}
}

func (w *Worker) handleJobEv(event *clientv3.Event) {
	switch {
	case event.IsCreate():
		w.logger.Info("is create..")
		job, err := w.GetJobContentFromKv(event.Kv.Key, event.Kv.Value)
		if err != nil {
			return
		}

		job.runOn = w.ID
		w.addJob(job)
	case event.IsModify():
		w.logger.Info("is IsModify..")
		job, err := w.GetJobContentFromKv(event.Kv.Key, event.Kv.Value)
		if err != nil {
			return
		}

		job.runOn = w.ID
		w.modJob(job)
	case event.Type == clientv3.EventTypeDelete:
		w.logger.Info("is EventTypeDelete..")
		w.delJob(GetIDFromKey(string(event.Kv.Key)))
	default:
		w.logger.Warnf("unknown event type[%v] from job[%s]", event.Type, string(event.Kv.Key))
	}
}

func (w *Worker) CleanJobs() {
	w.logger.Info("Worker: start clean jobs")

	for _, job := range w.jobs {
		w.delJob(job.ID)
	}
}

func getJobIDFromLockKey(key string) (jobId string) {
	key = strings.TrimLeft(key, LockKeyPrefix)
	return strings.Split(key, "/")[0]
}
