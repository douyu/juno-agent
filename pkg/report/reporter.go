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

package report

import (
	"time"

	"github.com/douyu/juno-agent/pkg/model"
	"github.com/douyu/jupiter/pkg"
)

// ReporterResp ...
type ReporterResp struct {
	Err  int
	Msg  string
	Data interface{}
}

// Reporter interface
type Reporter interface {
	Report(interface{}) ReporterResp
}

// Report ...
type Report struct {
	config *Config
	Reporter
}

// ReportAgentStatus report agent status
func (r *Report) ReportAgentStatus() error {
	if !r.config.Enable {
		return nil
	}
	go func() {
		for {
			req := model.AgentReportRequest{
				Hostname:     r.config.HostName,
				IP:           appIP,
				AgentType:    1,
				VCSInfo:      pkg.AppVersion(),
				AgentVersion: "0.2.1",
				RegionCode:   r.config.RegionCode,
				RegionName:   r.config.RegionName,
				ZoneCode:     r.config.ZoneCode,
				ZoneName:     r.config.ZoneName,
				Env:          r.config.Env,
			}
			r.Reporter.Report(req)
			time.Sleep(time.Duration(r.config.Internal))
		}
	}()
	return nil
}
