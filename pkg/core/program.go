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
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/douyu/jupiter/pkg/xlog"
)

//updateProgram update process information to local cache
func (eng *Engine) updateProgram(program *structs.ProgramExt) {
	switch program.Status {
	case "create":
		eng.programs.Store(program.UUID(), program)
	case "delete":
		eng.programs.Delete(program.UUID())
	case "list":
		eng.programs.Store(program.UUID(), program)
	case "update":
		eng.programs.Store(program.UUID(), program)
	}
	xdebug.PrintObject("programs", program)
}

func (eng *Engine) updateProcesses(processes ...structs.ProcessStatus) {
	for _, info := range processes {
		xlog.Info("process", xlog.Any("info", info))
		eng.processMap.Store(info.Command, info)
	}
}

// updateNginxProgram  update nginx information to local cache
func (eng *Engine) updateNginxProgram(conf *structs.NginxConfExt) {
	switch conf.Status {
	case "create":
		eng.programs.Store(conf.Name, conf)
	case "delete":
		eng.programs.Delete(conf.Name)
	case "list":
		eng.programs.Store(conf.Name, conf)
	case "update":
		eng.programs.Store(conf.Name, conf)
	}
}
