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

package etcd

import (
	"container/list"
	"errors"

	"github.com/douyu/juno-agent/pkg/structs"
	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/util/xgo"
)

var (
	// ErrEnvPass ...
	ErrEnvPass = errors.New("env pass")
)

// DataSource etcd conf datasource
type DataSource struct {
	etcdClient *etcdv3.Client
	prefix     string
	// 用于记录长轮训的应用信息
	jm list.List // *job
}

// configNode etcd node chan info
type configNode struct {
	key string
	ch  chan *structs.ConfNode
}

// NewETCDDataSource ...
func NewETCDDataSource(prometheusTargetGenConfig PluginRegProxyPrometheus) *DataSource {
	dataSource := &DataSource{
		etcdClient: etcdv3.StdConfig("register").Build(),
	}
	if prometheusTargetGenConfig.Enable {
		dataSource.PrometheusConfigScanner(prometheusTargetGenConfig.Path)
		xgo.Go(func() {
			dataSource.watchPrometheus(prometheusTargetGenConfig.Path)
		})
	}
	return dataSource
}

// GetClient ..
func (d *DataSource) GetClient() *etcdv3.Client {
	return d.etcdClient
}
