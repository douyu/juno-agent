package job

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

// 当前执行中的任务信息
// key: /{etcd_prefix}/jobId/taskId/node/pid
// value: 开始执行时间
// key 会自动过期，防止进程意外退出后没有清除相关 key，过期时间可配置
type Process struct {
	task *Task

	// parse from key path
	ID     string `json:"id"` // pid
	JobID  string `json:"jobId"`
	NodeID string `json:"nodeId"`
	TaskID uint64 `json:"task_id"`
	// parse from value
	ProcessVal

	running int32
	hasPut  int32
}

type ProcessVal struct {
	Time   time.Time `json:"time"`   // 开始执行时间
	Killed bool      `json:"killed"` // 是否强制杀死
}

func GetProcFromKey(key string) (proc *Process, err error) {
	ss := strings.Split(key, "/")
	var sslen = len(ss)
	if sslen < 5 {
		err = fmt.Errorf("invalid proc key [%s]", key)
		return
	}

	taskId, _ := strconv.ParseUint(ss[sslen-3], 10, 64)
	proc = &Process{
		ID:     ss[sslen-1],
		NodeID: ss[sslen-2],
		TaskID: taskId,
		JobID:  ss[sslen-4],
	}
	return
}

//key: /{etcd_prefix}/jobId/taskId/node/pid
func (p *Process) Key() string {
	taskId := strconv.FormatUint(p.TaskID, 10)
	return ProcKeyPrefix + p.JobID + "/" + taskId + "/" + p.NodeID + "/" + p.ID
}

func (p *Process) Val() (string, error) {
	b, err := json.Marshal(&p.ProcessVal)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (p *Process) put(job *Job) (err error) {
	if atomic.LoadInt32(&p.running) != 1 {
		return
	}

	if !atomic.CompareAndSwapInt32(&p.hasPut, 0, 1) {
		return
	}

	val, err := p.Val()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(job.ReqTimeout)*time.Second)
	defer cancel()

	session, err := concurrency.NewSession(job.Client.Client, concurrency.WithTTL(10))
	if err != nil {
		return err
	}

	_, err = job.Client.Put(ctx, p.Key(), val, clientv3.WithLease(session.Lease()))
	return
}

func (p *Process) Start(job *Job) {
	if p == nil {
		return
	}

	if !atomic.CompareAndSwapInt32(&p.running, 0, 1) {
		return
	}

	if err := p.put(job); err != nil {
		job.logger.Warnf("proc put[%s] err: %s", p.Key(), err.Error())
	}
	return
}

func (p *Process) del(job *Job) error {
	if atomic.LoadInt32(&p.hasPut) != 1 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(job.ReqTimeout)*time.Second)
	defer cancel()

	_, err := job.Delete(ctx, p.Key())
	return err
}

func (p *Process) Stop(job *Job) {
	if p == nil {
		return
	}

	if !atomic.CompareAndSwapInt32(&p.running, 1, 0) {
		return
	}

	if err := p.del(job); err != nil {
		job.logger.Warnf("proc del[%s] err: %s", p.Key(), err.Error())
	}
}
