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
	"encoding/json"
	"errors"
	"github.com/douyu/jupiter/pkg/xlog"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/douyu/juno-agent/pkg/structs"
	"github.com/douyu/juno-agent/util"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/ini.v1"
)

// Scanner systemd scanner
type Scanner struct {
	enable      bool
	StatusMap   sync.Map //Process configuration state
	confDir     string
	chanProgram chan *ProgramExt
	stop        chan struct{}
}

// Start start watch
func (s *Scanner) Start() error {
	if s.enable {
		return s.startWatch()
	}
	return nil
}

// Close close ...
func (s *Scanner) Close() error {
	s.stop <- struct{}{}
	return nil
}

// parseFile ...
func (s *Scanner) parseFile(file string) (*Program, []byte, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, nil, err
	}

	program, err := s.parse(content)
	return program, content, err
}

// parse ...
func (s *Scanner) parse(content []byte) (*Program, error) {
	conf, err := ini.ShadowLoad(content)
	if err != nil {
		return nil, err
	}

	var program Program

	section, err := conf.GetSection("Unit")
	if err != nil {
		return nil, err
	}

	if err := section.MapTo(&program.Unit); err != nil {
		return nil, err
	}
	section, err = conf.GetSection("Service")
	if err != nil {
		return nil, err
	}

	program.Environments = section.Key("Environment").ValueWithShadows()
	// Gets the --config parameter for the execstart directive
	program.ExecStart = section.Key("ExecStart").Value()
	kvs := strings.SplitN(program.ExecStart, " ", -1)
	for _, val := range kvs {
		if strings.Contains(val, "=") {
			params := strings.Split(strings.TrimSpace(val), "=")
			if len(params) == 2 && (params[0] == "--config" || params[0] == "-config") {
				program.Config = params[1]
			}
		}
	}

	if err := section.MapTo(&program.Service); err != nil {
		return nil, err
	}
	section, err = conf.GetSection("Install")
	if err != nil {
		return nil, err
	}

	if err := section.MapTo(&program.Install); err != nil {
		return nil, err
	}

	return &program, nil
}

// ListPrograms  show the systemd list
func (s *Scanner) ListPrograms() (tasks map[string]*ProgramExt, err error) {
	filesInfo, err := util.ReadDirFiles(s.confDir, "")
	if err != nil {
		return nil, err
	}
	tasks = make(map[string]*ProgramExt)
	for fileName, fileContent := range filesInfo {
		appName := strings.TrimRight(fileName, ".conf")

		program, err := s.parse([]byte(fileContent))
		if err != nil {
			xlog.Error("systemd.ListPrograms", xlog.String("parse err", err.Error()))
			continue
		}
		tasks[appName] = &ProgramExt{
			Status:   "list",
			FileName: fileName,
			FilePath: "",
			Program:  program,
		}
	}
	return
}

// C ...
func (s *Scanner) C() <-chan *ProgramExt {
	return s.chanProgram
}

// startWatch ...
func (s *Scanner) startWatch() error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	err = w.Add(s.confDir)
	if err != nil {
		xlog.Error("SystemdScanner add dir err", xlog.String("err", err.Error()))
		return nil
	}
	go func() {
		for {
			select {
			case ev, ok := <-w.Events:
				if !ok {
					return
				}
				if !strings.HasSuffix(ev.Name, ".conf") {
					continue
				}
				switch ev.Op {
				case fsnotify.Create:
					program, content, err := s.parseFile(ev.Name)
					if err != nil {
						xlog.Error("set status err", xlog.String("msg", err.Error()))
						continue
					}
					s.chanProgram <- &ProgramExt{
						Manager:  "systemd",
						Status:   "create",
						Program:  program,
						Content:  string(content),
						FileName: filepath.Base(ev.Name),
						FilePath: ev.Name,
					}
				case fsnotify.Write:
					program, content, err := s.parseFile(ev.Name)
					if err != nil {
						xlog.Error("set status err", xlog.String("msg", err.Error()))
						continue
					}
					s.chanProgram <- &ProgramExt{
						Status:   "update",
						Program:  program,
						Content:  string(content),
						FileName: filepath.Base(ev.Name),
						FilePath: ev.Name,
					}
				case fsnotify.Remove:
					s.chanProgram <- &ProgramExt{
						Status:   "delete",
						FileName: filepath.Base(ev.Name),
						FilePath: ev.Name,
						Program:  &Program{},
					}
				}
			case _, ok := <-w.Errors:
				if !ok {
					return
				}
			case <-s.stop:
				return
			}
		}
	}()
	return nil
}

// mapToStruct ...
func mapToStruct(data map[string]string) structs.AgentConfigStatus {
	model := structs.AgentConfigStatus{}
	buf, _ := json.Marshal(data)
	_ = json.Unmarshal(buf, &model)
	return model
}

// ProcessStatus ...
func (s *Scanner) ProcessStatus() (processList []structs.ProcessStatus, err error) {
	return nil, errors.New("not implement")
}
