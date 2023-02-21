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
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/douyu/juno-agent/pkg/check"
	"github.com/douyu/juno-agent/pkg/job"
	"github.com/douyu/juno-agent/pkg/mbus"
	"github.com/douyu/juno-agent/pkg/mbus/rocketmq"
	"github.com/douyu/juno-agent/pkg/nginx"
	"github.com/douyu/juno-agent/pkg/pmt/supervisor"
	"github.com/douyu/juno-agent/pkg/pmt/systemd"
	"github.com/douyu/juno-agent/pkg/process"
	"github.com/douyu/juno-agent/pkg/proxy/confProxy"
	"github.com/douyu/juno-agent/pkg/proxy/regProxy"
	"github.com/douyu/juno-agent/pkg/report"
	"github.com/douyu/juno-agent/pkg/structs"
	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/core/hooks"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/xlog"
	"golang.org/x/sync/errgroup"
)

// Engine engine
type Engine struct {
	jupiter.Application
	registryClient *etcdv3.Client
	confClient     *etcdv3.Client

	clients []*Client
	tracer  mbus.Tracer

	// Depending on the configuration details, decide which plug-ins to open
	programs          sync.Map
	processMap        sync.Map
	healthCheck       *check.HealthCheck
	process           *process.Scanner
	confProxy         *confProxy.ConfProxy
	regProxy          *regProxy.RegProxy
	report            *report.Report
	supervisorScanner *supervisor.Scanner
	systemdScanner    *systemd.Scanner
	nginxScanner      *nginx.ConfScanner
	worker            *job.Worker
}

// NewEngine new the engine
func NewEngine() *Engine {
	eng := &Engine{}
	//eng.SetRegistry(
	//	compound_registry.New(
	//		etcdv3_registry.StdConfig("test").Build(),
	//	),
	//)

	if err := eng.Startup(
		eng.startReportStatus, // start report agent status
		eng.startNginxConfScanner,
		eng.loadServiceNode, // load service nodes, and init configurations
		eng.startProcessScanner,
		eng.startConfProxy,
		eng.startRegProxy,
		eng.startSupervisorScanner, // scan supervisor conf dir in agent mode
		eng.startSystemdScanner,    // scan systemd conf dir in agent mode
		eng.startEventLogger,       // start event logger,
		eng.startShellProxy,        // start shell execution proxy,
		eng.startHealthScanner,     // start health scanner,
		eng.startHealCheck,
		eng.serveGRPC,
		eng.serveHTTP,
		eng.startWorker,
	); err != nil {
		xlog.Panic("new engine", xlog.Any("err", err))
	}

	eng.RegisterHooks(hooks.Stage(jupiter.StageAfterStop), eng.cleanJobs)

	return eng
}

// loadServiceNode ... TODO
func (eng *Engine) loadServiceNode() error { // load service node from local storage
	// recover fast when run fail
	return nil
}

// startConfProxy start app conf plugin
func (eng *Engine) startConfProxy() error {
	eng.confProxy = confProxy.StdConfig("confProxy").Build()
	eng.confProxy.Start()
	xgo.Go(func() {
		for node := range eng.confProxy.C() {
			// Monitor the confNode information of confProxy and bring the node into the Engine for management
			eng.upsertConfClient(node)
			// 1.0 Prefetch the registration configuration information for the pull configuration client application
			if err := eng.loadServiceConfiguration(node.AppName); err != nil {
				xlog.Error("load service configuration")
			}
		}
	})
	return nil
}

// startRegProxy start app regist proxy plugin
func (eng *Engine) startRegProxy() error {
	eng.regProxy = regProxy.StdConfig("regProxy").Build()
	if eng.regProxy == nil {
		return nil
	}

	if err := eng.regProxy.Start(); err != nil {
		return err
	}

	xgo.Go(func() {
		for node := range eng.regProxy.C() {
			eng.upsertRegClient(node)
		}
	})
	return nil
}

// startHealthScanner check node health status
func (eng *Engine) startHealthScanner() error {
	xgo.Go(
		eng.checkServiceNodes,
	)
	return nil
}

// startSupervisorScanner check and scan supervisor config
func (eng *Engine) startSupervisorScanner() error {
	eng.supervisorScanner = supervisor.StdConfig("supervisor").Build()
	if err := eng.supervisorScanner.Start(); err != nil {
		return err
	}
	programs, err := eng.supervisorScanner.ListPrograms()
	if err != nil {
		xlog.Error("startSupervisorScanner err", xlog.String("err", err.Error()))
		return nil
	}
	for _, program := range programs {
		eng.updateProgram(program.Unwrap())
	}
	xgo.Go(func() {
		for program := range eng.supervisorScanner.C() {
			eng.updateProgram(program.Unwrap())
		}
	})
	return nil
}

// startSystemdScanner check and scan systemd cofig
func (eng *Engine) startSystemdScanner() error {
	eng.systemdScanner = systemd.StdConfig("systemd").Build()
	if err := eng.systemdScanner.Start(); err != nil {
		return err
	}
	programs, err := eng.systemdScanner.ListPrograms()
	if err != nil {
		xlog.Debug("startSystemdScanner", xlog.String("err", err.Error()))
		return nil
	}
	for _, program := range programs {
		eng.updateProgram(program.Unwrap())
	}
	xgo.Go(func() {
		for program := range eng.systemdScanner.C() {
			eng.updateProgram(program.Unwrap())
		}
	})
	return nil
}

// startProcessScanner check go process
func (eng *Engine) startProcessScanner() error {
	eng.process = process.StdConfig("process").Build()
	if err := eng.process.Start(); err != nil {
		return err
	}

	processes, err := eng.process.Scan()
	if err != nil {
		return err
	}

	eng.updateProcesses(processes...)
	xgo.Go(func() {
		for processes := range eng.process.C() {
			eng.updateProcesses(processes...)
		}
	})

	return nil
}

// startShellProxy ...
func (eng *Engine) startShellProxy() error {
	return nil
}

// startEventLogger ...
func (eng *Engine) startEventLogger() error {
	tracer := rocketmq.New()
	if err := tracer.Start(); err != nil {
		return nil
	}
	// eng.tracer = tracer
	return nil
}

// startNginxConfScanner ...
func (eng *Engine) startNginxConfScanner() error {
	eng.nginxScanner = nginx.StdConfig("nginx").Build()
	confs, err := eng.nginxScanner.Scan()
	if err != nil {
		xlog.Error("startNginxScanner", xlog.String("err", err.Error()))
		return nil
	}
	for _, conf := range confs {
		eng.updateNginxProgram(conf)
	}
	if err := eng.nginxScanner.Start(); err != nil {
		return err
	}
	xgo.Go(func() {
		for update := range eng.nginxScanner.C() {
			eng.updateNginxProgram(update)
		}
	})
	return nil
}

// startReportStatus monitor the health status of the machine deployed by the agent, and regularly report the health information to the caller
func (eng *Engine) startReportStatus() error {
	eng.report = report.StdConfig("report").Build()
	err := eng.report.ReportAgentStatus()
	return err
}

// startHealCheck health status (including tcp, mysql, redis, http)
func (eng *Engine) startHealCheck() error {
	eng.healthCheck = check.StdConfig("healthCheck").Build()
	return nil
}

func (eng *Engine) startWorker() error {
	worker := job.StdConfig("worker").Build()
	if worker == nil {
		return nil
	}
	eng.worker = worker
	return worker.Run()
}

func (eng *Engine) loadServiceConfiguration(name string) interface{} {
	return nil
}

func (eng *Engine) getConfigStatus(configPath string) map[string]bool {
	var confStatusMap = map[string]bool{
		"supervisor": false,
		"systemd":    false,
	}
	if _, ok := eng.programs.Load(fmt.Sprintf("%s_%s", "supervisor", configPath)); ok {
		confStatusMap["supervisor"] = true
	}
	if _, ok := eng.programs.Load(fmt.Sprintf("%s_%s", "systemd", configPath)); ok {
		confStatusMap["systemd"] = true
	}
	return confStatusMap
}

func (eng *Engine) upsertRegClient(node *structs.ServiceNode) {
	for _, c := range eng.clients {
		if c.AppName != node.AppName || c.IP != node.IP {
			continue
		}
		if c.ServiceNodes == nil {
			c.ServiceNodes = make(map[*structs.ServiceNode]*CheckMeta)
		}
		for sn := range c.ServiceNodes {
			if sn.IP == node.IP && sn.Port == node.Port && sn.Schema == node.Schema {
				// have registered
				// todo(gorexlv): fulfill the register info?
				return
			}
		}

		c.ServiceNodes[node] = newCheckMeta()
		return
	}

	var client = &Client{
		AppName: node.AppName,
		IP:      node.IP,
		Port:    "",
		ServiceNodes: map[*structs.ServiceNode]*CheckMeta{
			node: newCheckMeta(),
		},
	}
	eng.clients = append(eng.clients, client)
}

// If the confNode already exists in the client managed by Engine, update the AppConfiguration in time;
// otherwise, add the confNode to the client managed by Engine
func (eng *Engine) upsertConfClient(node *structs.ConfNode) {
	for _, c := range eng.clients {
		if c.AppName == node.AppName &&
			c.AppEnvi == node.AppEnvi &&
			c.IP == node.IP {
			c.AppConfiguration = node.Configuration
			return
		}
	}
	var client = &Client{
		AppName:          node.AppName,
		AppEnvi:          node.AppEnvi,
		IP:               node.IP,
		Port:             "",
		AppConfiguration: node.Configuration,
		ServiceNodes:     nil,
	}
	eng.clients = append(eng.clients, client)
}

func (eng *Engine) checkServiceNodes() {
	for {
		time.Sleep(time.Second)
		var eg errgroup.Group
		// ping all registered nodes, then flush all node's status into storage
		for _, client := range eng.clients {
			for node, meta := range client.ServiceNodes {
				node := node
				meta := meta
				if meta.NextCheckTime.Load() > now().Unix() {
					continue
				}
				eg.Go(func() (err error) {
					fmt.Printf("node = %+v\n", node)
					if err := ping(context.TODO(), node); err != nil {
						fmt.Printf("err = %+v\n", err)
						meta.NextCheckTime.Add(3) // retry after 3s
						// todo(gorexlv): 超过一定次数, deregister service node
					} else {
						fmt.Printf("err = %+v\n", err)
						meta.NextCheckTime.Add(60) // retry after 60s
					}
					return
				})
			}
		}
		if err := eg.Wait(); err != nil {
			xlog.Error("group wait", xlog.Any("err", err))
		}
	}
}

func (eng *Engine) cleanJobs() {
	if eng.worker == nil {
		return
	}
	eng.worker.CleanJobs()
	return
}
