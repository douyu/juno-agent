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

package process

import (
	"fmt"
	"github.com/douyu/juno-agent/pkg/structs"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/xlog"
)

// Config ...
type Config struct {
	Enable bool `json:"enable"`
}

// StdConfig returns standard configuration information
func StdConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(fmt.Sprintf("plugin.%s", key), &config, conf.TagName("toml")); err != nil {
		fmt.Printf("loadProcessConfig.err:%#v\n", err)
		panic(err)
	}
	flagConfig := flag.Bool("process")
	config.Enable = flagConfig || config.Enable
	return &config
}

// DefaultConfig return default config
func DefaultConfig() Config {
	return Config{
		Enable: false,
	}
}

// Build new a instance
func (s *Config) Build() *Scanner {
	if s.Enable {
		xlog.Info("plugin", xlog.String("process", "start"))
	}
	return &Scanner{
		enable:        s.Enable,
		chanProcesses: make(chan []structs.ProcessStatus, 8),
		stop:          make(chan struct{}),
	}
}
