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
	"github.com/douyu/juno-agent/pkg/structs"
)

// Program ...
type Program struct {
	Directory             string `ini:"directory"`
	Environment           string `ini:"environment"`
	Command               string `ini:"command"`
	User                  string `ini:"user"`
	Autostart             bool   `ini:"autostart"`
	Autorestart           bool   `ini:"autorestart"`
	Startsecs             int    `ini:"startsecs"`
	StdoutLogfile         string `ini:"stdout_logfile"`
	StdoutLogfileMaxbytes string `ini:"stdout_logfile_maxbytes"`
	StdoutLogfileBackups  int    `ini:"stdout_logfile_backups"`
	StdoutCaptureMaxbytes string `ini:"stdout_capture_maxbytes"`
	StderrLogfile         string `ini:"stderr_logfile"`
	StderrLogfileMaxbytes string `ini:"stderr_logfile_maxbytes"`
	StderrLogfileBackups  int    `ini:"stderr_logfile_backups"`
	Stopsignal            string `ini:"stopsignal"`
	Config                string `ini:"-"`
}

// ProgramExt ...
type ProgramExt struct {
	Manager  string `json:"manager" toml:"manager"` // systemd|supervisor
	Status   string `json:"status" toml:"status"`
	*Program `json:"program" toml:"program"`
	Content  string `json:"content" toml:"content"`
	FileName string `json:"file_name" toml:"file_name"`
	FilePath string `json:"file_path" toml:"file_path"`
	AppName  string `json:"app_name" toml:"app_name"`
}

// Unwrap ...
func (program *ProgramExt) Unwrap() *structs.ProgramExt {
	return &structs.ProgramExt{
		FileName: program.FileName,
		FilePath: program.FilePath,
		Manager:  "supervisor",
		Status:   program.Status,
		// ProgramName:  program.Unit.Description,
		// StartCommand: program.ExecStart,
		// Environments: program.Environments,
		User: program.User,
		// Directory:    program.WorkingDirectory,
		Restart:    program.Autorestart,
		RestartSec: program.Startsecs, // to parse
		ConfData:   program.Content,
		Config:     program.Config,
	}
}
