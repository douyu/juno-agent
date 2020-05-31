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
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"time"
)

// ErrCode code
type ErrCode int

const (
	// ReportErr err code
	ReportErr ErrCode = 1
)

// HTTPReport httpReport struct
type HTTPReport struct {
	Config *Config
	client *resty.Client
}

// NewHTTPReport new httpReport
func NewHTTPReport(config *Config) *HTTPReport {
	return &HTTPReport{
		Config: config,
		client: resty.New().SetDebug(config.Debug).SetTimeout(60*time.Second).SetHeader("Content-Type", "application/json;charset=utf-8"),
	}
}

// Report ...
func (r *HTTPReport) Report(info interface{}) ReporterResp {
	resp, err := r.client.R().SetBody(info).Post(r.Config.Addr)
	if err != nil {
		return ReporterResp{
			Err: int(ReportErr),
			Msg: err.Error(),
		}
	}
	var res ReporterResp
	if err := json.Unmarshal([]byte(resp.Body()), &res); err != nil {
		return ReporterResp{
			Err: int(ReportErr),
			Msg: err.Error(),
		}
	}
	return res
}
