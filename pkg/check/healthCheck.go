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

package check

import (
	"errors"
	"github.com/douyu/juno-agent/pkg/check/impl/http"
	"github.com/douyu/juno-agent/pkg/check/impl/tcp"
	"github.com/douyu/juno-agent/pkg/check/view"
	"sync"

	"github.com/douyu/juno-agent/pkg/check/impl/mysql"
	"github.com/douyu/juno-agent/pkg/check/impl/redis"
	"github.com/douyu/juno-agent/pkg/model"
)

var container sync.Map

// MysqlHealthCheck ...
type HealthCheck struct {
	enable             bool
	wg                 sync.WaitGroup
	resHealthCheckChan chan *view.ResHealthCheck
}

// HealthCheck detection of service dependency
func (h *HealthCheck) HealthCheck(req model.CheckReq) (checkResults []*view.ResHealthCheck) {
	checkResults = make([]*view.ResHealthCheck, 0)
	for _, info := range req.CheckDatas {
		componentType := info.Type
		extConfig := info.Data
		h.wg.Add(1)
		go func(componentType string, extConfig string) {
			checkResult, err := h.DoHealthCheck(componentType, extConfig)
			if err != nil {
				h.resHealthCheckChan <- view.HealthCheckResult(componentType, false, err.Error())
			} else {
				h.resHealthCheckChan <- checkResult
			}
		}(componentType, extConfig)
	}
	h.wg.Wait()
	for {
		select {
		case res := <-h.resHealthCheckChan:
			checkResults = append(checkResults, res)
			if len(checkResults) == len(req.CheckDatas) {
				return
			}
		default:
		}
	}
	return
}
func (h *HealthCheck) DoHealthCheck(componentType string, extConfig string) (checkResult *view.ResHealthCheck, err error) {
	defer h.wg.Done()
	healthCheck, ok := container.Load(componentType)
	if !ok || healthCheck == nil {
		err = errors.New("can not load healthCheck instance")
		return
	}
	// get instance
	healthCheckImpl, err := Get(invoke(healthCheck))
	if err != nil {
		return
	}
	err = healthCheckImpl.LoadExtConfig(extConfig)
	if err != nil {
		return
	}
	checkResult, err = healthCheckImpl.DoHealthCheck()
	if err != nil {
		return
	}
	return
}

func init() {
	container.Store("mysql", mysql.NewMysqlHealthCheck)
	container.Store("redis", redis.NewRedisHealthCheck)
	container.Store("tcp", tcp.NewTCPHealthCheck)
	container.Store("http", http.NewHTTPHealthCheck)
}
