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

package check

import (
	"fmt"
	"github.com/douyu/juno-agent/pkg/check/view"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/xlog"
)

// Config ...
type Config struct {
	Enable bool
}

// StdConfig ...
func StdConfig(key string) *Config {
	var config = defaultConfig()
	if err := conf.UnmarshalKey(fmt.Sprintf("plugin.%s", key), &config, conf.TagName("toml")); err != nil {
		fmt.Printf("loadSystemdConfig.err:%#v\n", err)
		panic(err)
	}
	flagConfig := flag.Bool("healCheck")
	config.Enable = flagConfig || config.Enable
	return &config
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = defaultConfig()
	if err := conf.UnmarshalKey(key, &config, conf.TagName("toml")); err != nil {
		fmt.Printf("loadSystemdConfig.err:%#v\n", err)
		panic(err)
	}
	flagConfig := flag.Bool("healCheck")
	config.Enable = flagConfig || config.Enable
	return &config
}

// defaultConfig return default config
func defaultConfig() Config {
	return Config{
		Enable: true,
	}
}

// Build new a instance
func (c *Config) Build() *HealthCheck {
	if c.Enable {
		xlog.Info("plugin", xlog.String("healCheck", "start"))
	}
	return &HealthCheck{
		enable:             c.Enable,
		resHealthCheckChan: make(chan *view.ResHealthCheck, 100),
	}
}
