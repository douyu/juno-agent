package job

import (
	"fmt"
	"github.com/douyu/juno-agent/pkg/report"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/xlog"
)

type Config struct {
	Jobs          string // cmd 路径
	Once          string // 马上执行任务路径
	Lock          string // job lock 路径
	Group         string // 节点分组
	EtcdConfigKey string // jupiter.etcdv3.xxxxxx
	ReqTimeout    int    // 请求超时时间，单位秒

	HostName string
	AppIP    string
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Jobs:  "/worker/jobs/",
		Once:  "/worker/once/",
		Lock:  "/worker/lock/",
		Group: "/worker/group/",
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

	return NewWorker(c)
}


