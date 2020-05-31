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

package systemd

import (
	"fmt"
	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/xlog"

	"github.com/douyu/jupiter/pkg/conf"
)

// DefaultSystemdDir ...
var DefaultSystemdDir = "/etc/systemd/system"

// Config systemd config
type Config struct {
	Dir    string `json:"dir"`    // Configure the supervisor specific configuration path
	Enable bool   `json:"enable"` // Whether to open the open plug-in
}

// StdConfig returns standard configuration information
func StdConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(fmt.Sprintf("plugin.%s", key), &config, conf.TagName("toml")); err != nil {
		fmt.Printf("loadSystemdConfig.err:%#v\n", err)
		panic(err)
	}
	flagConfig := flag.Bool("systemd")
	config.Enable = flagConfig || config.Enable
	return &config
}

// DefaultConfig return default config
func DefaultConfig() Config {
	return Config{
		Dir:    DefaultSystemdDir,
		Enable: false,
	}
}

// Build new a instance
func (s *Config) Build() *Scanner {
	if s.Enable {
		xlog.Info("plugin", xlog.String("systemdScanner", "start"))
	}
	return &Scanner{
		enable:      s.Enable,
		confDir:     s.Dir,
		chanProgram: make(chan *ProgramExt, 1000),
	}
}
