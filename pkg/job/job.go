package job

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/douyu/juno-agent/pkg/job/etcd"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.etcd.io/etcd/clientv3/concurrency"
	"go.uber.org/zap"
)

var jobExecLogger *zap.Logger

func init() {
	var err error
	jobExecLogger, err = NewLogger()
	if err != nil {
		panic(err)
	}
}

type Jobs map[string]*Job

const (
	JobTypeNormal = iota // 运行各节点都能运行任务
	JobTypeAlone         // 同一时间只允许一个节点一个任务运行
)

// 需要执行的 cron cmd 命令
// 注册到 /cronsun/cmd/<id>
type Job struct {
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	Group   string     `json:"group"`
	Command string     `json:"cmd"`
	User    string     `json:"user"`
	Rules   []*JobRule `json:"rules"`
	Pause   bool       `json:"pause"`   // 可手工控制的状态
	Timeout int64      `json:"timeout"` // 单位时间秒，任务执行时间超时设置，大于 0 时有效
	// 执行任务失败重试次数
	// 默认为 0，不重试
	Retry int `json:"retry"`
	// 执行任务失败重试时间间隔
	// 单位秒，如果不大于 0 则马上重试
	Interval int `json:"interval"`
	// 任务类型
	// 0: 普通任务，各节点均可运行
	// 1: 单机任务，同时只能单节点在线
	JobType int `json:"kind"`

	// 执行任务的结点，用于记录 job log
	runOn    string // worker id
	hostname string
	ip       string

	// 用于存储分隔后的任务
	cmd []string

	// 用于访问etcd
	*worker
}

// NewEtcdTimeoutContext return a new etcdTimeoutContext
func NewEtcdTimeoutContext(w *worker) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(w.ReqTimeout)*time.Second)
}

func (j *Job) Cmds() (cmds map[string]*Cmd) {
	cmds = make(map[string]*Cmd)
	if j.Pause {
		return
	}

	for _, r := range j.Rules {
		cmd := &Cmd{
			Job:     j,
			JobRule: r,
		}
		cmds[cmd.GetID()] = cmd
	}

	return
}

// 从 etcd 的 key 中取 id
func GetIDFromKey(key string) string {
	index := strings.LastIndex(key, "/")
	if index < 0 {
		return ""
	}

	return key[index+1:]
}

func (j *Job) Run() error {
	var (
		cmd    *exec.Cmd
		ctx    context.Context
		cancel context.CancelFunc
	)

	if j.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(j.Timeout)*time.Second)
		defer cancel()
	} else {
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
	}
	j.logger.Infof("command is : %s", j.Command)
	cmd = exec.CommandContext(ctx, "/bin/bash", "-c", j.Command)

	sysProcAttr := &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.SysProcAttr = sysProcAttr
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	if err := cmd.Start(); err != nil {
		jobExecLogger.Info(b.String())
		return err
	}

	proc := &Process{
		ID:     strconv.Itoa(cmd.Process.Pid),
		JobID:  j.ID,
		NodeID: j.runOn,
		ProcessVal: ProcessVal{
			Time: time.Now(),
		},
	}
	proc.Start(j)
	defer proc.Stop(j)

	if err := cmd.Wait(); err != nil {
		jobExecLogger.Error(fmt.Sprintf("%s\n%s", b.String(), err.Error()))
		j.logger.Error(b.String())
		return err
	}

	jobExecLogger.Info(b.String())

	return nil
}

func (j *Job) RunWithRecovery() {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			j.logger.Warnf("panic running job: %v\n%s", r, buf)
		}
	}()
	_ = j.Run()
}

func (j *Job) ValidRules() error {
	for _, r := range j.Rules {
		if err := r.Valid(); err != nil {
			return err
		}
	}
	return nil
}

type JobRule struct {
	ID    string `json:"id"`
	Timer string `json:"timer"`

	Schedule Schedule `json:"-"`
}

// 验证 timer 字段
func (rule *JobRule) Valid() error {
	// 注意 interface nil 的比较
	if rule.Schedule != nil {
		return nil
	}

	if len(rule.Timer) == 0 {
		return errors.New("invalid job rule, empty timer.")
	}

	sch, err := myParser.Parse(rule.Timer)
	if err != nil {
		return fmt.Errorf("invalid JobRule[%s], parse err: %s", rule.Timer, err.Error())
	}

	rule.Schedule = sch
	return nil
}

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	logFile := GetCurrentDirectory() + "/log"
	os.MkdirAll(logFile, os.ModePerm)

	cfg.OutputPaths = []string{
		logFile + "/job.log",
	}
	return cfg.Build()
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		panic(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

func (j *Job) Info(msg string) {
	jobExecLogger.Info(msg)
}

func (j *Job) Error(msg string) {
	jobExecLogger.Error(msg)
}

type Cmd struct {
	*Job
	*JobRule
	schEntryID EntryID
}

func (c *Cmd) GetID() string {
	return c.Job.ID + c.JobRule.ID
}

func (c *Cmd) Run() error {

	if c.Job.JobType != JobTypeNormal {
		mutex, err := etcd.NewMutex(c.Client.Client, LockKeyPrefix, concurrency.WithTTL(5))
		if err != nil {
			c.logger.Info("job get lock error : ", xlog.FieldErr(err))
			return err
		}
		err = mutex.TryLock(time.Duration(c.ReqTimeout) * time.Second)
		if err != nil {
			c.logger.Info("job lock is failed : ", xlog.FieldErr(err))
			return err
		}
		defer mutex.Unlock()
	}

	if c.Job.Retry <= 0 {
		err := c.Job.Run()
		if err != nil {
			c.logger.Info("job run failed : ", xlog.FieldErr(err))
		}
		return err
	}

	for i := 0; i <= c.Job.Retry; i++ {
		if err := c.Job.Run(); err != nil {
			c.logger.Info("job run failed : ", xlog.FieldErr(err))
			return err
		}

		if c.Job.Interval > 0 {
			time.Sleep(time.Duration(c.Job.Interval) * time.Second)
		}
	}

	return nil
}
