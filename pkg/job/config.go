package job

import (
	"fmt"

	"github.com/douyu/juno-agent/pkg/job/parser"
	"github.com/douyu/juno-agent/pkg/report"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/robfig/cron/v3"
)

var (
	myParser = parser.NewParser(parser.Second | parser.Minute | parser.Hour | parser.Dom | parser.Month | parser.Dow | parser.Descriptor)
)

const (
	JobsKeyPrefix = "/worker/cmd/"  // cronsun task路径
	OnceKeyPrefix = "/worker/once/" // 马上执行任务路径
	LockKeyPrefix = "/worker/lock/" // job lock 路径
	ProcKeyPrefix = "/worker/proc/" // 正在运行的Process
)

type Config struct {
	EtcdConfigKey   string // jupiter.etcdv3.xxxxxx
	ReqTimeout      int    // 请求操作ETCD的超时时间，单位秒
	RequireLockTime int64  // 抢锁等待时间，单位秒

	HostName string
	AppIP    string

	logger   *xlog.Logger
	parser   parser.Parser
	wrappers []cron.JobWrapper
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		ReqTimeout: 5,
	}
}

// StdConfig returns standard configuration information
func StdConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(fmt.Sprintf("plugin.%s", key), config, conf.TagName("toml")); err != nil {
		xlog.Error("loadWorkerConfig", xlog.Any("err", err))
		panic(err)
	}

	return config
}

// Build new a instance
func (c *Config) Build() *worker {
	c.HostName = report.ReturnHostName()
	c.AppIP = report.ReturnAppIp()

	if c.logger == nil {
		c.logger = xlog.JupiterLogger
	}
	c.logger = c.logger.With(xlog.FieldMod("worker"))

	// default
	c.parser = myParser
	// 默认前面有任务执行，则直接跳过不执行
	c.wrappers = append(c.wrappers, skipIfStillRunning(c.logger))

	return NewWorker(c)
}
