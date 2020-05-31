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

package core

import (
	"github.com/douyu/juno-agent/pkg/structs"
	"github.com/uber-go/atomic"
)

// Client represents a service node registered to agent/proxy
type Client struct {
	AppName string
	AppEnvi string
	IP      string
	Port    string
	// Configuration captured by the configuration center
	AppConfiguration *structs.AppConfiguration
	// Service nodes: An app can register multiple service nodes
	ServiceNodes map[*structs.ServiceNode]*CheckMeta
}

// CheckMeta ...
type CheckMeta struct {
	createTime    int64
	LastCheckTime *atomic.Int64
	NextCheckTime *atomic.Int64

	Disable *atomic.Bool
}

// newCheckMeta ...
func newCheckMeta() *CheckMeta {
	return &CheckMeta{
		createTime:    0,
		LastCheckTime: atomic.NewInt64(0),
		NextCheckTime: atomic.NewInt64(0),
		Disable:       atomic.NewBool(false),
	}
}
