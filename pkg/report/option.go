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

package report

import (
	"fmt"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/util/xtime"
	"github.com/douyu/jupiter/pkg/xlog"
	"os"
	"time"
)

// Config report config
type Config struct {
	Enable     bool          `json:"enable"`
	Debug      bool          `json:"debug"`
	Addr       string        `json:"addr"`
	Internal   time.Duration `json:"internal"`
	HostName   string        `json:"host_name"`
	RegionCode string        `json:"region_code"`
	RegionName string        `json:"region_name"`
	ZoneCode   string        `json:"zone_code"`
	ZoneName   string        `json:"zone_name"`
	Env        string        `json:"env"`
}

// StdConfig returns standard configuration information
func StdConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(fmt.Sprintf("plugin.%s", key), &config, conf.TagName("toml")); err != nil {
		xlog.Error("loadReprotConfig", xlog.Any("err", err))
		panic(err)
	}
	flagEnable := flag.Bool("agentReport")
	config.Enable = config.Enable || flagEnable
	return &config
}

// DefaultConfig return default config
func DefaultConfig() Config {
	return Config{
		Enable:   false,
		Internal: xtime.Duration("60s"),
	}
}

// Build new a instance
func (r *Config) Build() *Report {
	r.RegionCode = os.Getenv(r.RegionCode)
	r.RegionName = os.Getenv(r.RegionName)
	r.ZoneCode = os.Getenv(r.ZoneCode)
	r.ZoneName = os.Getenv(r.ZoneName)
	r.HostName = GetHostName(r.HostName)
	hostName = r.HostName
	env := os.Getenv(r.Env)
	if env == "" {
		env = "dev"
	}
	r.Env = env
	report := &Report{
		config:   r,
		Reporter: NewHTTPReport(r),
	}
	if r.Enable {
		xlog.Info("plugin", xlog.String("reportAgentStatus", "start"))
	}
	return report
}
