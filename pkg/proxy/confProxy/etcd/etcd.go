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

package etcd

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/douyu/juno-agent/pkg/report"
	"github.com/douyu/juno-agent/pkg/structs"
	"github.com/douyu/juno-agent/util"
	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/util/xnet"
	"github.com/douyu/jupiter/pkg/xlog"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

var (
	// ErrEnvPass ...
	ErrEnvPass = errors.New("env pass")
)

// DataSource etcd conf datasource
type DataSource struct {
	etcdClient *etcdv3.Client
	prefix     string
	// 用于记录长轮训的应用信息
	jm list.List // *job
}

// configNode etcd node chan info
type configNode struct {
	key string
	ch  chan *structs.ConfNode
}

// NewETCDDataSource ...
func NewETCDDataSource(prefix string, etcdConfig ConfDataSourceEtcd) *DataSource {
	dataSource := &DataSource{
		etcdClient: etcdv3.RawConfig("plugin.confProxy.etcd").Build(),
		prefix:     prefix,
	}
	xgo.Go(dataSource.watch)
	return dataSource
}

// GetValues ...
func (d *DataSource) GetValues(ctx echo.Context, keys ...string) (map[string]string, error) {
	var (
		appName, appEnv, target, port = keys[0], keys[1], keys[2], keys[3]
		config                        structs.ServiceConf
		res                           = make(map[string]string)
	)
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return res, errors.New("wrong port")
	}
	if portInt <= 0 {
		return res, errors.New("wrong port")
	}

	if appName == "" || appEnv == "" {
		return res, errors.New("invalid param")
	}

	hostKey := fmt.Sprintf("/juno-agent/%s/%s/%s/static/%s/%d", report.ReturnHostName(), appName, appEnv, target, portInt)
	appKey := fmt.Sprintf("/juno-agent/%s/%s/%s/static/%s", "cluster", appName, appEnv, target)
	commonKey := fmt.Sprintf("%s/%s/%s/%s", appName, appEnv, target, port)
	data, err := d.etcdClient.GetValues(context.Background(), hostKey, appKey)
	if err != nil {
		return res, err
	}
	if _, ok := data[hostKey]; ok {
		if err := jsoniter.Unmarshal([]byte(data[hostKey]), &config); err == nil {
			res[commonKey] = config.Content
			return res, err
		}
	}

	if _, ok := data[appKey]; ok {
		if err := jsoniter.Unmarshal([]byte(data[appKey]), &config); err == nil {
			res[commonKey] = config.Content
			return res, nil
		}
	}

	xlog.Info("getAppConfigContent", xlog.String("appKey", appKey), xlog.Any("data", data))
	return res, errors.New("no etcd config is found")
}

func (d *DataSource) GetRawValues(ctx echo.Context, rawKey string) (map[string]string, error) {
	var (
		res    = make(map[string]string)
		config = structs.ServiceConf{}
	)

	data, err := d.etcdClient.GetValues(context.Background(), rawKey)
	if err != nil {
		return res, err
	}
	if _, ok := data[rawKey]; ok {
		if err := jsoniter.Unmarshal([]byte(data[rawKey]), &config); err == nil {
			res[rawKey] = config.Content
			return res, err
		}
		if err := jsoniter.Unmarshal([]byte(data[rawKey]), &config); err == nil {
			res[rawKey] = config.Content
			return res, nil
		}
	}
	xlog.Info("getAppConfigContent", xlog.String("rawKey", rawKey), xlog.Any("data", data))
	return res, errors.New("no etcd config is found")
}

// AppConfigScanner 初始化加载实例配置
func (d *DataSource) AppConfigScanner() []*structs.ConfNode {
	confuNodes := make([]*structs.ConfNode, 0)
	hostKey := strings.Join([]string{d.prefix, report.ReturnHostName()}, "/")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	resp, err := d.etcdClient.Get(ctx, hostKey, clientv3.WithPrefix())
	if err != nil {
		xlog.Error("init get hostKey error", xlog.String("plugin", "confgo"), xlog.String("msg", err.Error()), xlog.String("hostKey", hostKey))
		return confuNodes
	}
	for _, kv := range resp.Kvs {
		key, value := string(kv.Key), string(kv.Value)
		if confuNode, rr := d.update(key, value); rr != nil {
			if err == ErrEnvPass { //环境过滤
				xlog.Info("init get update env pass", xlog.String("plugin", "confgo"), xlog.String("key", key))
			} else {
				xlog.Error("init get update error", xlog.String("plugin", "confProxy"), xlog.String("msg", rr.Error()), xlog.String("key", key), xlog.String("err", rr.Error()))
			}
			continue
		} else {
			confuNodes = append(confuNodes, confuNode)
		}
		if err := d.report(key, value); err != nil {
			xlog.Error("init get report error", xlog.String("plugin", "confgo"), xlog.String("msg", err.Error()), xlog.String("hostKey", hostKey))
			continue
		}
		xlog.Debug("init update success", xlog.String("plugin", "confgo"), xlog.String("hostKey", hostKey))
	}
	return confuNodes
}

// watch 监听配置变动
func (d *DataSource) watch() {
	// etcd的key用作配置数据读取
	hostKey := strings.Join([]string{d.prefix, report.ReturnHostName()}, "/")
	// init watch
	watch, err := d.etcdClient.NewWatch(hostKey)

	if err != nil {
		panic("watch err: " + err.Error())
	}
	go func() {
		for {
			select {
			case event := <-watch.C():
				switch event.Type {
				case mvccpb.DELETE:
				case mvccpb.PUT:
					key, value := string(event.Kv.Key), string(event.Kv.Value)
					// 用于检测该key是否存在于长轮训map中
					xlog.Info("watch put", xlog.String("plugin", "confgo"), xlog.String("key", key), xlog.String("val", value))
					if confuNode, err := d.update(key, value); err != nil {
						if err == ErrEnvPass {
							xlog.Info("watch update env pass", xlog.String("plugin", "confgo"), xlog.String("key", key), xlog.String("val", value))
						} else {
							xlog.Error("watch update error", xlog.String("plugin", "confgo"), xlog.String("msg", err.Error()), xlog.String("key", key))
						}
						continue
					} else {
						commonKey := util.GetConfigKey(confuNode.AppName, confuNode.AppEnvi, confuNode.FileName, confuNode.Port)
						d.StoreAppChanInfo(commonKey, key, confuNode)
					}

					if err := d.report(key, value); err != nil {
						xlog.Error("watch report error", xlog.String("plugin", "confgo"), xlog.String("msg", err.Error()), xlog.String("key", key))
						continue
					}
					xlog.Info("watch update success", xlog.String("key", key), xlog.String("val", value))
				}
			}
		}
	}()
}

// ListenAppConfig listen the app config change
func (d *DataSource) ListenAppConfig(ctx echo.Context, key string) chan *structs.ConfNode {
	xlog.Info("confProxy", xlog.String("listenConfig", key))
	node := &configNode{
		key: key,
		ch:  make(chan *structs.ConfNode),
	}
	d.jm.PushBack(node)
	return node.ch
}

// update 更新本地文件
func (d *DataSource) update(key, value string) (*structs.ConfNode, error) {
	confNode := &structs.ConfNode{}
	// key check
	confuKeys, err := structs.ParserConfKey(key)
	if err != nil {
		return confNode, err
	}
	if err := confuKeys.CheckValid(); err != nil {
		return confNode, fmt.Errorf("key check: %s", err.Error())
	}
	xlog.Debug("file update content", xlog.String("plugin", "confgo"), xlog.Any("confuKeys", confuKeys), xlog.String("key", key), xlog.String("value", value))

	// value check
	confuValue, err := structs.ParserConfValue(value)
	if err != nil {
		return confNode, err
	}
	if err := confuValue.CheckValid(); err != nil {
		return confNode, fmt.Errorf("value check: %s", err.Error())
	}

	for _, path := range confuValue.Metadata.Paths {
		// write file
		if err := util.WriteFile(path, confuValue.Content); err != nil {
			return confNode, err
		}
	}

	confNode = &structs.ConfNode{
		AppName:  confuKeys.AppName,
		AppEnvi:  confuKeys.EnvName,
		IP:       "",
		FileName: confuKeys.FileName,
		Port:     confuKeys.Port,
		Configuration: &structs.AppConfiguration{
			Content:  confuValue.Content,
			Metadata: structs.Metadata{Format: confuValue.Metadata.Format, Timestamp: confuValue.Metadata.Timestamp, Version: confuValue.Metadata.Version},
		},
	}
	return confNode, nil
}

// report 上报配置下发状态
func (d *DataSource) report(key, value string) error {

	confuKeys, err := structs.ParserConfKey(key)
	if err != nil {
		return err
	}
	confuValue, err := structs.ParserConfValue(value)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	ip, err := xnet.GetLocalIP()
	if err != nil {
		xlog.Error("checkEffectMD5", xlog.String("xnetGetLocalIPError", err.Error()))
	}

	reportKey := strings.Join([]string{d.prefix + "/callback", confuKeys.AppName, confuKeys.FileName, confuKeys.Hostname}, "/")
	reportValue := structs.ConfReport{
		FileName:   confuKeys.FileName,
		MD5:        confuValue.Metadata.Version,
		Hostname:   confuKeys.Hostname,
		Env:        confuKeys.EnvName,
		Timestamp:  time.Now().Unix(),
		IP:         ip,
		HealthPort: confuKeys.Port,
	}
	if _, err := d.etcdClient.Put(ctx, reportKey, reportValue.JSONString()); err != nil {
		return err
	}
	return nil
}

// StoreAppChanInfo 监听到etcd的变化后，更新chan的信息
func (d *DataSource) StoreAppChanInfo(key, rawKey string, val *structs.ConfNode) {
	var n *list.Element
	for item := d.jm.Front(); nil != item; item = n {
		node := item.Value.(*configNode)
		n = item.Next()
		if node.key == key || node.key == rawKey {
			select {
			case node.ch <- val:
			default:
			}
			close(node.ch)
			d.jm.Remove(item)
		}

	}
}

// Reload 重新启动
func (d *DataSource) Reload() error {
	// 关闭监听
	if err := d.etcdClient.Watcher.Close(); err != nil {
		return err
	}
	// 建立新的watcher
	// 重新启动监听
	d.AppConfigScanner()
	d.watch()
	return nil
}

// Stop 进程退出停止监听变化
func (d *DataSource) Stop() {
	//
	if err := d.etcdClient.Watcher.Close(); err != nil {
		log.Error("confgo stop watcher error", "msg", err.Error())
		return
	}
	if err := d.etcdClient.Close(); err != nil {
		log.Error("confgo stop etcd client error", "msg", err.Error())
		return
	}
}
