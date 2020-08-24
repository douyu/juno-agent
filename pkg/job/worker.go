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

package job

import (
	"github.com/douyu/juno-agent/pkg/job/parser"
	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/xlog"
)

// Node 执行 cron 命令服务的结构体
type worker struct {
	*Config
	*etcdv3.Client
	//*cronsun.Node
	*Cron

	logger *xlog.Logger
	parser parser.Parser

	ImmediatelyRun bool // 是否立即执行

	//jobs           Jobs // 和结点相关的任务
	//groups         Groups
	//cmds           map[string]*cronsun.Jobs

	//link
	// 删除的 job id，用于 group 更新
	delIDs map[string]bool

	ttl  int64
	//lID  client.LeaseID // lease id
	done chan struct{}
}

func NewWorker(conf *Config) (w *worker) {

	w = &worker{
		Config:         conf,
		Client:         etcdv3.StdConfig(conf.EtcdConfigKey).Build(),
		ImmediatelyRun: false,
		delIDs:         nil,
		done:           make(chan struct{}),
	}
	// default
	w.parser = parser.NewParser(parser.Second | parser.Minute | parser.Hour | parser.Dom | parser.Month | parser.Dow | parser.Descriptor)
	w.Cron = newCron(w)

	return
}

func (w *worker) Run() error {

	w.logger.Info("worker run...")



	return nil
}
