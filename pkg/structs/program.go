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

package structs

import (
	"fmt"
	parser "github.com/yangchenxing/go-nginx-conf-parser"
)

// ProgramExt the ext of the program, mainly scan the systemd or supervisor config of the app
type ProgramExt struct {
	Status       string   `json:"status" toml:"status"`
	FileName     string   `json:"file_name" toml:"file_name"` // profile name
	FilePath     string   `json:"file_path" toml:"file_path"` // profile path
	Manager      string   `json:"manager" toml:"manager"`     // systemd|supervisor
	ProgramName  string   `json:"program" toml:"program"`     // programName eg: wsd-live-srv-room-go
	StartCommand string   `json:"start_command" toml:"start_command"`
	Environments []string `json:"environments" toml:"environments"` // environment variable
	User         string   `json:"user" toml:"user"`                 // user account name
	Directory    string   `json:"directory" toml:"directory"`       // execution dir
	Restart      bool     `json:"restart" toml:"restart"`           // restart policy
	RestartSec   int      `json:"restart_sec" toml:"restart_sec"`   // start interval
	ConfData     string   `json:"conf_data" toml:"conf_data"`       // configuration data
	Config       string   `json:"config" toml:"config"`             // profile name used for specific configuration
}

// UUID the unique uuid of program
func (pe *ProgramExt) UUID() string {
	return fmt.Sprintf("%s_%s", pe.Manager, pe.Config)
}

// Program ...
type Program struct {
	PID   string   `json:"pid" toml:"pid" `
	State string   `json:"state" toml:"state" `
	Time  string   `json:"time" toml:"time" `
	Name  string   `json:"app_name" toml:"app_name" `
	Path  []string `json:"path" toml:"path" `
	Args  []string `json:"args" toml:"args" `
}

// AgentConfigStatus ...
type AgentConfigStatus struct {
	Hostname              string `json:"hostname"`
	IP                    string `json:"ip"`
	AppName               string `json:"app_name"`
	Env                   string `json:"env"`
	Directory             string `json:"directory"`
	Environment           string `json:"environment"`
	Command               string `json:"command"`
	User                  string `json:"user"`
	Autostart             string `json:"autostart"`
	Autorestart           string `json:"autorestart"`
	Startsecs             string `json:"startsecs"`
	StdoutLogfile         string `json:"stdout_logfile"`
	StdoutLogfileMaxbytes string `json:"stdout_logfile_maxbytes"`
	StdoutLogfileBackups  string `json:"stdout_logfile_backups"`
	StdoutCaptureMaxbytes string `json:"stdout_capture_maxbytes"`
	StderrLogfile         string `json:"stderr_logfile"`
	StderrLogfileMaxbytes string `json:"stderr_logfile_maxbytes"`
	StderrLogfileBackups  string `json:"stderr_logfile_backups"`
	Stopsignal            string `json:"stopsignal"`
}

// ProcessStatus ...
type ProcessStatus struct {
	User    string `json:"user"`
	PID     string `json:"pid"`
	CPU     string `json:"cpu"`
	MEM     string `json:"mem"`
	VSZ     string `json:"vsz"`
	RSS     string `json:"rss"`
	TTY     string `json:"tty"`
	Stat    string `json:"stat"`
	Start   string `json:"start"`
	Time    string `json:"time"`
	Command string `json:"command"`
}

// Unwrap ...
func (ps *ProcessStatus) Unwrap() *ProgramExt {
	return nil
}

type NginxConfExt struct {
	Name   string                     `json:"name"`
	Status string                     `json:"status"`
	Block  parser.NginxConfigureBlock `json:"block"`
}
