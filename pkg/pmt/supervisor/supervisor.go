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

package supervisor

import (
	"errors"
	"fmt"
	"github.com/douyu/juno-agent/pkg/structs"
	"github.com/douyu/juno-agent/util"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Scanner supervisor sacnner
type Scanner struct {
	enable      bool
	confDir     string // conf文件路径
	stop        chan struct{}
	chanProgram chan *ProgramExt
}

// Show ...
func (s *Scanner) Show() map[string]interface{} {
	return nil
}

// GetStatus ...
func (s *Scanner) GetStatus(appName string) ([]structs.AgentConfigStatus, error) {
	return nil, nil
}

// ProcessStatus ...
func (s *Scanner) ProcessStatus() (processList []structs.ProcessStatus, err error) {
	return nil, nil
}

// Start ...
func (s *Scanner) Start() error {
	if s.enable {
		return s.startWatch()
	}
	return nil
}

// parseFile parse supervisor file
func (s *Scanner) parseFile(file string) (*Program, []byte, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, nil, err
	}
	program, err := s.parse(content)
	return program, content, err
}

// parse parse...
func (s *Scanner) parse(content []byte) (*Program, error) {
	conf, err := ini.LoadSources(ini.LoadOptions{
		AllowBooleanKeys: true,
	}, content)
	if err != nil {
		return nil, err
	}
	for _, section := range conf.Sections() {
		if section.Name() != ini.DEFAULT_SECTION {
			program := new(Program)
			err = section.MapTo(program)
			if program.Command != "" {
				kvs := strings.SplitN(program.Command, " ", -1)
				for _, val := range kvs {
					if strings.Contains(val, "=") {
						params := strings.Split(strings.TrimSpace(val), "=")
						if len(params) == 2 && (params[0] == "--config" || params[0] == "-config") {
							program.Config = params[1]
							return program, nil
						}
					}
				}
			}
			return program, nil
		}
	}

	return nil, errors.New("invalid supervisor conf")
}

// ListPrograms list the supervisor list
func (s *Scanner) ListPrograms() (tasks map[string]*ProgramExt, err error) {
	tasks = make(map[string]*ProgramExt)
	if s.enable {
		filesInfo, err := util.ReadDirFiles(s.confDir, "")
		if err != nil {
			return nil, err
		}

		for fileName, fileContent := range filesInfo {
			appName := strings.TrimRight(fileName, ".conf")

			program, err := s.parse([]byte(fileContent))
			if err != nil {
				xlog.Error("supervisor.ListPrograms", xlog.String("parse err", err.Error()))
				continue
			}
			tasks[appName] = &ProgramExt{
				Manager:  "supervisor",
				Status:   "list",
				Program:  program,
				Content:  fileContent,
				FileName: fileName,
				FilePath: fmt.Sprintf("%s/%s", s.confDir, fileName),
				AppName:  appName,
			}
		}
	}
	return tasks, nil
}

// C consume the chan program
func (s *Scanner) C() <-chan *ProgramExt {
	return s.chanProgram
}

// startWatch Monitor folder changes to update process configuration
func (s *Scanner) startWatch() error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		xlog.Error("fsnotify new watcher err", xlog.String("msg", err.Error()))
		return err
	}
	if err := w.Add(s.confDir); err != nil {
		xlog.Error("fsnotify add dir err", xlog.String("msg", err.Error()))
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
				xlog.Debug("supervisorConfigChange", xlog.String("events", ev.String()))
				switch ev.Op {
				case fsnotify.Create:
					program, content, err := s.parseFile(ev.Name)
					if err != nil {
						xlog.Error("set status err", xlog.String("msg", err.Error()))
						continue
					}
					s.chanProgram <- &ProgramExt{
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
			case err, ok := <-w.Errors:
				xlog.Error("watch err", xlog.String("msg", err.Error()))
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
