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

package pmt

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

var (
	// Start shell start
	Start = 0
	// Stop shell stop
	Stop = 1
	// ReStart shell restart
	ReStart = 2

	// StartOP start
	StartOP = "start"
	// StopOP stop
	StopOP = "stop"
	// ReStartOP restart
	ReStartOP = "restart"
)

// GenCommand Script returns a command to execute a script through a shell.
func GenCommand(pmt, app string, op int) (args []string, err error) {
	args = make([]string, 0)
	args = append(args, "-c")

	switch pmt {
	case "systemd":
		args = append(args, "systemctl")
	case "supervisor":
		args = append(args, "supervisorctl")
	default:
		return nil, errors.New("cmd shell is not found")
	}

	switch op {
	case Start:
		args = append(args, StartOP)
	case Stop:
		args = append(args, StopOP)
	case ReStart:
		args = append(args, ReStartOP)
	default:
		return nil, errors.New("cmd op is not found")
	}

	// check app is valid
	if strings.Contains(app, "|") || len(strings.TrimSpace(app)) == 0 {
		return nil, errors.New("cmd app is not correct")
	}

	args = append(args, strings.TrimSpace(app))

	return args, nil
}

// Exec exec shell command
func Exec(command []string) (string, error) {
	cmd := exec.Command("bash", command...)
	err := cmd.Run()
	var out bytes.Buffer
	cmd.Stdout = &out
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
