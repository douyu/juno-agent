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

package confProxy

import (
	"errors"
	"github.com/douyu/juno-agent/util"
	"github.com/labstack/echo/v4"
	"time"

	"github.com/douyu/juno-agent/pkg/structs"
	"github.com/douyu/jupiter/pkg/xlog"
)

// ConfProxy confProxy struct
type ConfProxy struct {
	enable     bool
	dataSource DataSource
	nodeInput  chan *structs.ConfNode
}

// NewConfProxy new instance
func NewConfProxy(enable bool, confClient DataSource) *ConfProxy {
	return &ConfProxy{
		enable:     enable,
		dataSource: confClient,
		nodeInput:  make(chan *structs.ConfNode, 100),
	}
}

// Start ...
func (cp *ConfProxy) Start() {
	if cp.enable {
		for _, node := range cp.dataSource.AppConfigScanner() {
			select {
			case cp.nodeInput <- node:
			default:
				xlog.Warn("ConfProxy.AppConfigScanner", xlog.String("nodeInput chan", "err"))
			}
		}
	}
}

// Close ...
func (cp *ConfProxy) Close() {
	close(cp.nodeInput)
}

// C ...
func (cp *ConfProxy) C() <-chan *structs.ConfNode {
	return cp.nodeInput
}

// GetValues ...
func (cp *ConfProxy) GetValues(ctx echo.Context, appName, appEnv, target, port string) (config string, err error) {
	data, err := cp.dataSource.GetValues(ctx, appName, appEnv, target, port)
	commonKey := util.GetConfigKey(appName, appEnv, target, port)
	if err != nil {
		return "nil", err
	}
	if len(data[commonKey]) > 0 {
		return data[commonKey], nil
	}

	return "nil", errors.New("unknown config")
}

// GetRawValues ...
func (cp *ConfProxy) GetRawValues(ctx echo.Context, rawKey string) (config string, err error) {
	data, err := cp.dataSource.GetRawValues(ctx, rawKey)
	if err != nil {
		return "nil", err
	}
	if len(data[rawKey]) > 0 {
		return data[rawKey], nil
	}

	return "nil", errors.New("unknown config")
}

// ListenAppConfig ...
func (cp *ConfProxy) ListenAppConfig(ctx echo.Context, appName, appEnv, target, port string, watch bool, internal int) (structs.ContentNode, error) {
	commonKey := util.GetConfigKey(appName, appEnv, target, port)
	switch watch {
	case true:
		ch := cp.dataSource.ListenAppConfig(ctx, commonKey)
		select {
		case info := <-ch:
			if info.Configuration != nil {
				return structs.ContentNode{
					Content: info.Configuration.Content,
					Version: time.Now().Unix(),
				}, nil
			}
			return structs.ContentNode{}, errors.New("get app config nil")
		case <-time.After(time.Second * time.Duration(internal)):
		}
		return structs.ContentNode{}, errors.New("no change")
	default:
		content, err := cp.GetValues(ctx, appName, appEnv, target, port)
		if err != nil {
			return structs.ContentNode{}, err
		}
		return structs.ContentNode{
			Content: content,
			Version: time.Now().Unix(),
		}, nil
	}
}

// ListenRawKeyAppConfig ...
func (cp *ConfProxy) ListenRawKeyAppConfig(ctx echo.Context, rawKey string, watch bool, internal int) (structs.ContentNode, error) {
	switch watch {
	case true:
		ch := cp.dataSource.ListenAppConfig(ctx, rawKey)
		select {
		case info := <-ch:
			if info.Configuration != nil {
				return structs.ContentNode{
					Content: info.Configuration.Content,
					Version: time.Now().Unix(),
				}, nil
			}
			return structs.ContentNode{}, errors.New("get raw key, app config nil")
		case <-time.After(time.Second * time.Duration(internal)):
		}
		return structs.ContentNode{}, errors.New("no change")
	default:
		content, err := cp.GetRawValues(ctx, rawKey)
		if err != nil {
			return structs.ContentNode{}, err
		}
		return structs.ContentNode{
			Content: content,
			Version: time.Now().Unix(),
		}, nil
	}
}

// Reload ...
func (cp *ConfProxy) Reload() error {
	return cp.dataSource.Reload()
}

// extractConfNode ...
func (cp *ConfProxy) extractConfNode(appName, appEnv string, ip string) {
	select {
	case cp.nodeInput <- &structs.ConfNode{AppName: appName, AppEnvi: appEnv, IP: ip}:
	default:
		xlog.Warn("extract conf node timeout",
			xlog.String("appName", appName),
			xlog.String("appEnv", appEnv),
			xlog.Any("chan size", len(cp.nodeInput)),
		)
	}
}
