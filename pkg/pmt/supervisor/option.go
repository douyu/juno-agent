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

package supervisor

import (
	"fmt"
	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/xlog"

	"github.com/douyu/jupiter/pkg/conf"
)

// DefaultSupervisorDir ...
var DefaultSupervisorDir = "/etc/supervisor/conf.d"

// Config ...
type Config struct {
	Dir    string `json:"dir"`    // 配置中心supervisor具体配置路径
	Enable bool   `json:"enable"` // 是否开启开插件
}

// StdConfig 返回标准配置信息
func StdConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(fmt.Sprintf("plugin.%s", key), &config, conf.TagName("toml")); err != nil {
		fmt.Printf("loadSupervisorConfig.err:%#v\n", err)
		panic(err)
	}
	flagConfig := flag.Bool("supervisor")
	config.Enable = flagConfig || config.Enable
	return &config
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config, conf.TagName("toml")); err != nil {
		fmt.Printf("loadSupervisorConfig.err:%#v\n", err)
		panic(err)
	}
	flagConfig := flag.Bool("supervisor")
	config.Enable = flagConfig || config.Enable
	return &config
}

// DefaultConfig return default config
func DefaultConfig() Config {
	return Config{
		Dir:    DefaultSupervisorDir,
		Enable: false,
	}
}

// Build new a instance
func (c *Config) Build() *Scanner {
	if c.Enable {
		xlog.Info("plugin", xlog.String("supervisorScanner", "start"))
	}
	return &Scanner{
		enable:      c.Enable,
		confDir:     c.Dir,
		chanProgram: make(chan *ProgramExt, 1000),
	}
}
