package job

import (
	"context"
	"encoding/json"
	"strings"
	"time"
)

const (
	JobTypeNormal   = iota
	JobTypeAlone    // 任何时间段只允许单机执行
	JobTypeInterval // 一个任务执行间隔内允许执行一次
)

// 需要执行的 cron cmd 命令
// 注册到 /cronsun/cmd/groupName/<id>
type Job struct {
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	Group   string     `json:"group"`
	Command string     `json:"cmd"`
	User    string     `json:"user"`
	Rules   []*JobRule `json:"rules"`
	Pause   bool       `json:"pause"`   // 可手工控制的状态
	Timeout int64      `json:"timeout"` // 任务执行时间超时设置，大于 0 时有效
	// 设置任务在单个节点上可以同时允许多少个
	// 针对两次任务执行间隔比任务执行时间要长的任务启用
	Parallels int64 `json:"parallels"`
	// 执行任务失败重试次数
	// 默认为 0，不重试
	Retry int `json:"retry"`
	// 执行任务失败重试时间间隔
	// 单位秒，如果不大于 0 则马上重试
	Interval int `json:"interval"`
	// 任务类型
	// 0: 普通任务
	// 1: 单机任务
	// 如果为单机任务，node 加载任务的时候 Parallels 设置 1
	JobType int `json:"job_type"`

	// 用于存储分隔后的任务
	cmd []string
}

// NewEtcdTimeoutContext return a new etcdTimeoutContext
func NewEtcdTimeoutContext(w *worker) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(w.ReqTimeout)*time.Second)
}

func (w *worker) loadJobs() (jobs map[string]*Job, err error) {
	ctx, cancelFunc := NewEtcdTimeoutContext(w)
	defer cancelFunc()

	kvMaps, err := w.Client.GetPrefix(ctx, w.Config.Jobs)
	if err != nil {
		return
	}

	count := len(kvMaps)
	jobs = make(map[string]*Job, count)
	if count == 0 {
		return
	}

	for key, val := range kvMaps {
		job := &Job{}
		if e := json.Unmarshal([]byte(val), job); e != nil {
			w.logger.Warnf("job[%s] unmarshal err: %s", key, e.Error())
			continue
		}

		if err := job.ParseJob(); err != nil {
			w.logger.Warnf("job[%s] is invalid: %s", key, err.Error())
		}

		jobs[job.ID] = job
	}

	return
}

// todo 设计成interface接口
func (j *Job) ParseJob() error {
	if len(j.cmd) == 0 {
		j.splitCmd()
	}

	return nil
}

func (j *Job) splitCmd() {
	ps := strings.SplitN(j.Command, " ", 2)
	if len(ps) == 1 {
		j.cmd = ps
		return
	}

	j.cmd = make([]string, 0, 2)
	j.cmd = append(j.cmd, ps[0])
	j.cmd = append(j.cmd, parseCmdArguments(ps[1])...)
}

type JobRule struct {
	ID             string   `json:"id"`
	Timer          string   `json:"timer"`
	GroupIDs       []string `json:"gids"`
	NodeIDs        []string `json:"nids"`
	ExcludeNodeIDs []string `json:"exclude_nids"`
}
