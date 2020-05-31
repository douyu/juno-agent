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

package tcp

import (
	"encoding/json"
	"github.com/douyu/juno-agent/pkg/check/view"
	"github.com/douyu/jupiter/pkg/xlog"
	"net"
)

// TCPHealthCheck config
type TCPHealthCheck struct {
	Addr    string `json:"addr"`
	Network string `json:"network"` // tcp tcp4 tcp6
}

// NewTCPHealthCheck new instance
func NewTCPHealthCheck() *TCPHealthCheck {
	return DefaultTCPHealthCheck()
}
func DefaultTCPHealthCheck() *TCPHealthCheck {
	return &TCPHealthCheck{Network: "tcp"}
}

// LoadExtConfig parse tcp config
func (t *TCPHealthCheck) LoadExtConfig(extConfig string) (err error) {
	if err = json.Unmarshal([]byte(extConfig), &t); err != nil {
		return
	}
	return
}

// DoHealthCheck check is invoked periodically to perform the TCP check
func (t *TCPHealthCheck) DoHealthCheck() (resHealthCheck *view.ResHealthCheck, err error) {
	tcpAddr, err := net.ResolveTCPAddr(t.Network, t.Addr)
	if err != nil {
		xlog.Error("ResolveTCPAddr", xlog.Any("tcp addr err", err))
		return
	}
	conn, err := net.DialTCP(t.Network, nil, tcpAddr)
	if err != nil {
		xlog.Error("DailTcp", xlog.Any("DailTcp err", err))
		return
	}
	if err = conn.Close(); err != nil {
		xlog.Error("Conn close", xlog.Any("Connection Close err", err))
		return
	}
	resHealthCheck = view.HealthCheckResult("tcp", true, "success")
	return resHealthCheck, nil
}
