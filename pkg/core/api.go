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

package core

import (
	"strconv"

	pb "github.com/coreos/etcd/etcdserver/etcdserverpb"
	"github.com/douyu/juno-agent/pkg/file"
	"github.com/douyu/juno-agent/pkg/model"
	"github.com/douyu/juno-agent/pkg/pmt"
	"github.com/douyu/juno-agent/pkg/structs"
	"github.com/douyu/juno-agent/util"
	"github.com/douyu/jupiter/pkg/server/xecho"
	"github.com/douyu/jupiter/pkg/server/xgrpc"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

func (eng *Engine) serveHTTP() error {

	s := xecho.StdConfig("http").Build()

	group := s.Group("/api")
	group.GET("/agent/reload", eng.agentReload)           // restart confd monitoring
	group.GET("/agent/process/status", eng.processStatus) // real time process status
	group.POST("/agent/process/shell", eng.pmtShell)
	group.GET("/agent/file", eng.readFile) // 文件读取

	v1Group := s.Group("/api/v1")
	v1Group.GET("/agent/:target", eng.getAppConfig) // get app config
	v1Group.GET("/agent/config", eng.listenConfig)  // listenConfig
	v1Group.POST("/agent/check", eng.agentCheck)    // 加入依赖探活检测
	v1Group.POST("/conf/command_line/status", eng.confStatus)

	v1Group.GET("/agent/rawKey/getConfig", eng.getRawAppConfig)       // 根据原生key获取配置信息
	v1Group.GET("/agent/rawKey/listenConfig", eng.listenRawKeyConfig) // 根据原生key长轮训监听配置

	return eng.Serve(s)
}

func (eng *Engine) serveGRPC() error {
	config := xgrpc.StdConfig("grpc")
	server := config.Build()
	pb.RegisterKVServer(server.Server, eng.regProxy)
	pb.RegisterWatchServer(server.Server, eng.regProxy)
	pb.RegisterLeaseServer(server.Server, eng.regProxy)
	helloworld.RegisterGreeterServer(server.Server, eng.regProxy)

	return eng.Serve(server)
}

// agentReload reload agent watch
func (eng *Engine) agentReload(ctx echo.Context) error {
	if err := eng.confProxy.Reload(); err != nil {
		return reply400(ctx, "reload err "+err.Error())
	}
	return reply200(ctx, nil)
}

type confStatusBind struct {
	Config string `json:"config"` //path to profile
}

// confStatus returns whether the configuration corresponding to systemd or supervisor is connected to the configuration center
// input： config path
// output：return whether the corresponding access
func (eng *Engine) confStatus(ctx echo.Context) error {
	bind := confStatusBind{}
	if err := ctx.Bind(&bind); err != nil {
		return reply400(ctx, err.Error())
	}
	if bind.Config == "" {
		return reply400(ctx, "input parameter is empty")
	}

	confStatus := eng.getConfigStatus(bind.Config)
	return reply200(ctx, confStatus)
}

// listenConfig record the change of config
// if the config change in internal time(default 60s),return the changed config
// other return the status 400
func (eng *Engine) listenConfig(ctx echo.Context) error {
	var defaultListenInternal = 60
	appName := ctx.QueryParam("name")
	appEnv := ctx.QueryParam("env")
	target := ctx.QueryParam("target")
	port := ctx.QueryParam("port")
	enableWatch, _ := strconv.ParseBool(ctx.QueryParam("watch"))
	listenInternal, _ := strconv.Atoi(ctx.QueryParam("internal"))
	if listenInternal > 0 {
		defaultListenInternal = listenInternal
	}
	var (
		config structs.ContentNode
		err    error
	)
	if config, err = eng.confProxy.ListenAppConfig(ctx, appName, appEnv, target, port, enableWatch, defaultListenInternal); err != nil {
		return reply400(ctx, "no data change")
	}
	return reply200(ctx, config)
}

// listenRawKeyConfig record the change of config
// if the config change in internal time(default 60s),return the changed config
// other return the status 400
func (eng *Engine) listenRawKeyConfig(ctx echo.Context) error {
	var defaultListenInternal = 60
	rawKey := ctx.QueryParam("rawKey")
	if rawKey == "" {
		return reply400(ctx, "listen config, raw key is null")
	}
	enableWatch, _ := strconv.ParseBool(ctx.QueryParam("watch"))
	listenInternal, _ := strconv.Atoi(ctx.QueryParam("internal"))
	if listenInternal > 0 {
		defaultListenInternal = listenInternal
	}

	var (
		config structs.ContentNode
		err    error
	)
	if config, err = eng.confProxy.ListenRawKeyAppConfig(ctx, rawKey, enableWatch, defaultListenInternal); err != nil {
		return reply400(ctx, "no data change")
	}
	return reply200(ctx, config)
}

// getAppConfig get the app config data
func (eng *Engine) getAppConfig(ctx echo.Context) error {
	appName := ctx.QueryParam("name")
	appEnv := ctx.QueryParam("env")
	port := ctx.QueryParam("port")
	target := ctx.Param("target")
	res, err := eng.confProxy.GetValues(ctx, appName, appEnv, target, port)
	if err != nil {
		return reply400(ctx, err.Error())
	}
	return reply200(ctx, res)
}

func (eng *Engine) getRawAppConfig(ctx echo.Context) error {
	rawKey := ctx.QueryParam("rawKey")
	if rawKey == "" {
		return reply400(ctx, "get raw app config, the raw key is null")
	}
	res, err := eng.confProxy.GetRawValues(ctx, rawKey)
	if err != nil {
		return reply400(ctx, err.Error())
	}
	return reply200(ctx, res)
}

// processStatus show the process status of machine
func (eng *Engine) processStatus(ctx echo.Context) error {
	list, err := eng.process.GetProcessStatus()
	if err != nil {
		return reply400(ctx, "process list parser err:"+err.Error())
	}
	return reply200(ctx, list)
}

// agentCheck add the health check of relies
func (eng *Engine) agentCheck(ctx echo.Context) error {
	checkDatas := model.CheckReq{}
	if err := ctx.Bind(&checkDatas); err != nil {
		return ctx.JSON(400, err.Error())
	}
	res := eng.healthCheck.HealthCheck(checkDatas)
	return reply200(ctx, res)
}

// PMTShell pmt shell exec
func (eng *Engine) pmtShell(ctx echo.Context) error {
	pmtShellReq := model.PMTShell{}
	if err := ctx.Bind(&pmtShellReq); err != nil {
		return reply400(ctx, err.Error())
	}
	args, err := pmt.GenCommand(pmtShellReq.Pmt, pmtShellReq.AppName, pmtShellReq.Op)
	if err != nil {
		return reply400(ctx, err.Error())
	}
	reply, err := pmt.Exec(args)
	if err != nil {
		return reply400(ctx, err.Error())
	}
	return reply200(ctx, reply)
}

func (eng *Engine) readFile(c echo.Context) error {
	var param model.GetFileReq
	err := c.Bind(&param)
	if err != nil {
		return reply400(c, err.Error())
	}

	content, err := file.ReadFile(param.FileName)
	if err != nil {
		return reply400(c, err.Error())
	}

	content = util.EncryptAPIResp(content)

	return reply200(c, map[string]interface{}{
		"content": content,
	})
}

func reply200(ctx echo.Context, data interface{}) error {
	return ctx.JSON(200, map[string]interface{}{
		"code": 200,
		"data": data,
		"msg":  "success",
	})
}

func reply400(ctx echo.Context, msg string) error {
	return ctx.JSON(200, map[string]interface{}{
		"code": 400,
		"msg":  msg,
	})
}
