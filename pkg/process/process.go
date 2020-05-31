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

package process

import (
	"os/exec"
	"strings"
	"time"

	"github.com/douyu/juno-agent/pkg/structs"
	"github.com/douyu/jupiter/pkg/util/xgo"
)

// Scanner ...
type Scanner struct {
	enable        bool
	chanProcesses chan []structs.ProcessStatus
	stop          chan struct{}
}

//Start after a delay of the specified time,
// the loop scan is started for status information and the data is written to the channel of the processScanner
func (ps *Scanner) Start() error {
	xgo.DelayGo(time.Second*10, ps.monitor)
	return nil
}

// C ...
func (ps *Scanner) C() <-chan []structs.ProcessStatus {
	return ps.chanProcesses
}

// Close ...
func (ps *Scanner) Close() error {
	ps.stop <- struct{}{}
	close(ps.chanProcesses)
	return nil
}

// Scan ...
func (ps *Scanner) Scan() ([]structs.ProcessStatus, error) {
	return ps.scan()
}

// Scan the process state and returns
func (ps *Scanner) scan() ([]structs.ProcessStatus, error) {
	processList := make([]structs.ProcessStatus, 0)
	if ps.enable {
		const statusCmdStr = "ps aux | grep \"go/bin\" | grep -v \"grep\" | tr -s \" \""
		cmd := exec.Command("/bin/bash", "-c", statusCmdStr)
		resp, err := cmd.CombinedOutput()
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(resp), "\n")
		for _, line := range lines {
			items := strings.SplitAfterN(line, " ", 11)

			if len(items) != 11 {
				continue
			}
			processList = append(processList, structs.ProcessStatus{
				User:    items[0],
				PID:     items[1],
				CPU:     items[2],
				MEM:     items[3],
				VSZ:     items[4],
				RSS:     items[5],
				TTY:     items[6],
				Stat:    items[7],
				Start:   items[8],
				Time:    items[9],
				Command: items[10],
			})
		}
	}
	return processList, nil
}

// GetProcessStatus ...
func (ps Scanner) GetProcessStatus() ([]structs.ProcessStatus, error) {
	return ps.scan()
}

// monitor Periodically monitor the process status and write to the channel
func (ps *Scanner) monitor() {
	var ticker = time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			processes, err := ps.scan()
			if err != nil {
				// log.Errord("scan process")
				continue
			}
			ps.chanProcesses <- processes
		case <-ps.stop:
			return
		}
	}
}
