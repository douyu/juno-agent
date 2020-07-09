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

package confProxy

import (
	"fmt"
	"time"

	"github.com/douyu/juno-agent/pkg/proxy/confProxy/etcd"
	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/xlog"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/util/xtime"
)

// DefaultConfDir ...
var DefaultConfDir = "/home/www/.config/juno-agent"

// Config confProxy config
type Config struct {
	Dir     string        `json:"dir"` // 配置中心具体配置路径
	Prefix  string        `json:"prefix"`
	Env     []string      `json:"env"`
	Timeout time.Duration // etcd连接超时时间
	Secure  bool
	Enable  bool                // 是否开启开插件
	Mysql   ConfDataSourceMysql `json:"mysql"`
	Etcd    ConfDataSourceEtcd  `json:"etcd"`
}

// ConfDataSourceMysql mysql dataSource
type ConfDataSourceMysql struct {
	Enable bool   // 是否开启用该数据源
	Dsn    string `json:"dsn"`
}

// ConfDataSourceEtcd ETCD dataSource config
type ConfDataSourceEtcd struct {
	Enable                        bool // 是否开启用该数据源
	Secure                        bool
	IsWatchPrometheusTargetConfig bool
	EndPoints                     []string `json:"endpoints"` // 注册中心etcd节点信息
}

// StdConfig 返回标准配置信息
func StdConfig(key string) *Config {
	var config = DefaultConfig()

	if err := conf.UnmarshalKey(fmt.Sprintf("plugin.%s", key), &config, conf.TagName("toml")); err != nil {
		fmt.Printf("loadRegistryConfig.err:%#v\n", err)
		xlog.Error("confProxy", xlog.String("parse config err", err.Error()))
		panic(err)
	}
	flagConfig := flag.Bool("confProxy")
	config.Enable = flagConfig || config.Enable
	return &config
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()

	if err := conf.UnmarshalKey(key, &config, conf.TagName("toml")); err != nil {
		fmt.Printf("loadRegistryConfig.err:%#v\n", err)
		xlog.Error("confProxy", xlog.String("parse config err", err.Error()))
		panic(err)
	}
	flagConfig := flag.Bool("confProxy")
	config.Enable = flagConfig || config.Enable
	return &config
}

// DefaultConfig default config info
func DefaultConfig() Config {
	return Config{
		Dir:     DefaultConfDir,
		Timeout: xtime.Duration("1s"),
		Enable:  false,
		Mysql: ConfDataSourceMysql{
			Enable: false,
			Dsn:    "127.0.0.1:6379",
		},
		Etcd: ConfDataSourceEtcd{
			Enable:                        false,
			IsWatchPrometheusTargetConfig: false,
			Secure:                        false,
			EndPoints:                     []string{"127.0.0.1:2379"},
		},
	}
}

// Build  new the instance
func (c *Config) Build() *ConfProxy {
	if c.Enable {
		switch c.Etcd.Enable {
		case true:
			xlog.Info("plugin", xlog.String("appConf.etcd", "start"))
			return NewConfProxy(c.Enable, etcd.NewETCDDataSource(c.Prefix, c.Etcd.IsWatchPrometheusTargetConfig))
		default:
			xlog.Info("plugin", xlog.String("appConf.mysql", "start"))
		}
		// todo mysql implement
	}
	return nil
}
