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

package regProxy

import (
	"fmt"
	"time"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/util/xtime"
)

// Config regConfig
type Config struct {
	EndPoints []string `json:"endpoints"`
	Timeout   time.Duration
	Secure    bool
	Enable    bool // Whether to open the open plug-in
}

// StdConfig returns standard configuration information
func StdConfig(key string) Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config, conf.TagName("toml")); err != nil {
		fmt.Printf("loadRegistryConfig.err:%#v\n", err)
		panic(err)
	}
	return config
}

// DefaultConfig return default config
func DefaultConfig() Config {
	return Config{
		EndPoints: []string{"127.0.0.1:2379"},
		Timeout:   xtime.Duration("1s"),
		Secure:    false,
		Enable:    false,
	}
}
