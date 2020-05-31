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

	"github.com/douyu/juno-agent/pkg/structs"
)

// Unit systemd file struct
type Unit struct {
	Description   string
	Documentation string
	After         string
	Wants         string
}

// Service systemd service config
type Service struct {
	Environments     []string
	User             string
	Group            string
	WorkingDirectory string
	ExecStart        string
	ExecReload       string
	KillMode         string
	LimitNOFILE      int
	LimitNPROC       int
	Restart          string
	RestartSec       string
	Config           string
}

// Install install
type Install struct {
	WantedBy string
}

// Program program
type Program struct {
	Unit
	Service
	Install
}

func (program Program) String() string {
	bs, _ := json.MarshalIndent(program, "", "    ")
	return string(bs)
}

// ProgramExt program ext show the detail info
type ProgramExt struct {
	Manager  string `json:"manager" toml:"manager"`
	*Program `json:"program" toml:"program"`
	Status   string `json:"status" toml:"status"`
	Content  string `json:"content" toml:"content"`
	FileName string `json:"file_name" toml:"file_name"`
	FilePath string `json:"file_path" toml:"file_path"`
	AppName  string `json:"app_name" toml:"app_name"`
}

// Unwrap ...
func (program *ProgramExt) Unwrap() *structs.ProgramExt {
	return &structs.ProgramExt{
		FileName:     program.FileName,
		FilePath:     program.FilePath,
		Manager:      "systemd",
		Status:       program.Status,
		ProgramName:  program.Unit.Description,
		StartCommand: program.ExecStart,
		Environments: program.Environments,
		User:         program.User,
		Directory:    program.WorkingDirectory,
		Restart:      program.Restart == "always",
		RestartSec:   5, // to parse
		ConfData:     program.Content,
		Config:       program.Config,
	}
}
