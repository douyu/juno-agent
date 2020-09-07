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
	"time"

	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/zap"
)

func init() {
}

type Jobs map[string]*Job

const (
	TypeNormal = 0 // 运行各节点都能运行任务
	TypeAlone  = 1 // 同一时间只允许一个节点一个任务运行
)

// 需要执行的 cron cmd 命令
// 注册到 /cronsun/cmd/<id>
type Job struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Script  string   `json:"script"`
	Timers  []*Timer `json:"timers"`
	Enable  bool     `json:"enable"`  // 可手工控制的状态
	Timeout int64    `json:"timeout"` // 单位时间秒，任务执行时间超时设置，大于 0 时有效
	Env     string   `json:"env"`
	Zone    string   `json:"zone"`

	// 执行任务失败重试次数
	// 默认为 0，不重试
	RetryCount int `json:"retry_count"`

	// 执行任务失败重试时间间隔
	// 单位秒，如果不大于 0 则马上重试
	RetryInterval int `json:"retry_interval"`

	// 任务类型
	// 0: 普通任务，各节点均可运行
	// 1: 单机任务，同时只能单节点在线
	JobType int `json:"job_type"`

	// 执行任务的结点，用于记录 job log
	runOn    string // worker id
	hostname string

	// 用于访问etcd
	*worker `json:"-"`

	mutex  *etcdv3.Mutex
	locked bool
}

// NewEtcdTimeoutContext return a new etcdTimeoutContext
func NewEtcdTimeoutContext(w *worker) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(w.ReqTimeout)*time.Second)
}

func (j *Job) Cmds() (cmds map[string]*Cmd) {
	cmds = make(map[string]*Cmd)
	if !j.Enable {
		return
	}

	for _, r := range j.Timers {
		cmd := &Cmd{
			Job:   j,
			Timer: r,
			Nodes: r.Nodes,
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

func (j *Job) Run(taskOptions ...TaskOption) error {
	var (
		cmd           *exec.Cmd
		ctx           context.Context
		cancel        context.CancelFunc
		consoleLogBuf bytes.Buffer
	)

	task := NewTask(j, taskOptions...)
	_ = task.SetStatus(CronTaskStatusProcessing, "")

	if j.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(j.Timeout)*time.Second)
		defer cancel()
	} else {
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
	}

	// check if script exists
	scriptFileState, err := os.Stat(j.Script)
	if err != nil {
		j.logger.Error("read script file failed", xlog.String("err", err.Error()))

		consoleLogBuf.WriteString("read script file failed: " + err.Error())
		_ = task.SetStatus(CronTaskStatusFailed, consoleLogBuf.String())

		return err
	} else if scriptFileState.IsDir() {
		j.logger.Error("script path is a dir", xlog.String("script", j.Script))

		consoleLogBuf.WriteString("script path is a dir, not a executable file")
		_ = task.SetStatus(CronTaskStatusFailed, consoleLogBuf.String())

		return fmt.Errorf("script is a dir, not a executable file. jobId[%s] script[%s]", j.ID, j.Script)
	}

	j.logger.Infof("command is : %s", j.Script)
	cmd = exec.CommandContext(ctx, j.Script)

	sysProcAttr := makeCmdAttr()
	cmd.SysProcAttr = sysProcAttr
	cmd.Stdout = &consoleLogBuf
	cmd.Stderr = &consoleLogBuf
	if err := cmd.Start(); err != nil {
		j.logger.Info(consoleLogBuf.String())

		consoleLogBuf.WriteString(err.Error())
		_ = task.SetStatus(CronTaskStatusFailed, consoleLogBuf.String())
		return err
	}

	proc := &Process{
		ID:     strconv.Itoa(cmd.Process.Pid),
		JobID:  j.ID,
		NodeID: j.runOn,
		TaskID: task.TaskID,
		ProcessVal: ProcessVal{
			Time: time.Now(),
		},
	}
	proc.Start(j)
	defer func() {
		go func() {
			time.Sleep(3 * time.Second)
			proc.Stop(j)
		}()
	}()

	if err := cmd.Wait(); err != nil {
		j.logger.Error(consoleLogBuf.String())
		consoleLogBuf.WriteString(err.Error())

		if ctx.Err() == context.DeadlineExceeded {
			_ = task.SetStatus(CronTaskStatusTimeout, consoleLogBuf.String())
		} else {
			_ = task.SetStatus(CronTaskStatusFailed, consoleLogBuf.String())
		}

		return err
	}

	j.logger.Info(consoleLogBuf.String())
	_ = task.SetStatus(CronTaskStatusSuccess, consoleLogBuf.String())

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
	for _, r := range j.Timers {
		if err := r.Valid(); err != nil {
			return err
		}
	}
	return nil
}

func (j *Job) Lock() error {
	var err error
	j.mutex, err = j.Client.NewMutex(LockKeyPrefix+j.ID, concurrency.WithTTL(10))
	if err != nil {
		return err
	}

	err = j.mutex.Lock(3 * time.Second)
	if err != nil {
		return err
	}

	j.locked = true

	return nil
}

func (j *Job) Unlock() {
	if j.mutex != nil {
		err := j.mutex.Unlock()
		if err != nil {
			xlog.Error("unlock failed", xlog.FieldErr(err))
		}
	}
}

type Timer struct {
	ID    string   `json:"id"`
	Cron  string   `json:"timer"`
	Nodes []string `json:"nodes"`

	Schedule Schedule `json:"-"`
}

// 验证 timer 字段
func (rule *Timer) Valid() error {
	// 注意 interface nil 的比较
	if rule.Schedule != nil {
		return nil
	}

	if len(rule.Cron) == 0 {
		return errors.New("invalid job rule, empty timer.")
	}

	sch, err := myParser.Parse(rule.Cron)
	if err != nil {
		return fmt.Errorf("invalid Timer[%s], parse err: %s", rule.Cron, err.Error())
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

type Cmd struct {
	Nodes []string

	*Job
	*Timer
	schEntryID EntryID
}

func (c *Cmd) GetID() string {
	return c.Job.ID + "-" + c.Timer.ID
}

func (c *Cmd) Run() error {
	if c.Job.RetryCount <= 0 {
		err := c.Job.Run()
		if err != nil {
			c.logger.Info("job run failed : ", xlog.FieldErr(err))
		}
		return err
	}

	for i := 0; i <= c.Job.RetryCount; i++ {
		if err := c.Job.Run(); err != nil {
			c.logger.Info("job run failed : ", xlog.FieldErr(err))
		}

		if c.Job.RetryInterval > 0 {
			time.Sleep(time.Duration(c.Job.RetryInterval) * time.Second)
		}
	}

	return nil
}
