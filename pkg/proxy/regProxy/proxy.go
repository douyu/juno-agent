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

package regProxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/douyu/juno-agent/pkg/proxy/regProxy/etcd"
	"github.com/douyu/juno-agent/pkg/structs"
	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/douyu/jupiter/pkg/util/xstring"
	"github.com/douyu/jupiter/pkg/xlog"
	pb "go.etcd.io/etcd/etcdserver/etcdserverpb"
	"go.etcd.io/etcd/proxy/grpcproxy"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

// RegProxy ...
type RegProxy struct {
	pb.KVServer
	pb.LeaseServer
	pb.WatchServer

	*etcdv3.Client
	nodeChan chan *structs.ServiceNode

	serviceConfigurations sync.Map
}

// NewRegProxy ...
func NewRegProxy(confClient *etcd.DataSource) *RegProxy {
	proxy := &RegProxy{
		Client:   confClient.GetClient(),
		nodeChan: make(chan *structs.ServiceNode, 100),
	}
	return proxy
}

// SayHello ...
func (proxy *RegProxy) SayHello(ctx context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	fmt.Printf("request = %+v\n", request)
	return &helloworld.HelloReply{
		Message: request.Name,
	}, nil
}

// Start ...
func (proxy *RegProxy) Start() error {
	proxy.KVServer, _ = grpcproxy.NewKvProxy(proxy.Client.Client)
	proxy.LeaseServer, _ = grpcproxy.NewLeaseProxy(proxy.Client.Client)
	proxy.WatchServer, _ = grpcproxy.NewWatchProxy(proxy.Client.Client)
	return nil
}

// Close ...
func (proxy *RegProxy) Close() {
	close(proxy.nodeChan)
}

// C ...
func (proxy *RegProxy) C() <-chan *structs.ServiceNode {
	return proxy.nodeChan
}

// LoadServiceConfiguration ...
func (proxy *RegProxy) LoadServiceConfiguration(appName string) (*structs.AppConfiguration, bool) {
	data, ok := proxy.serviceConfigurations.Load(appName)
	if !ok {
		return nil, false
	}

	conf, ok := data.(*structs.AppConfiguration)
	return conf, ok
}

// StoreServiceConfiguration ...
func (proxy *RegProxy) StoreServiceConfiguration(key string, conf *structs.AppConfiguration) {
	proxy.serviceConfigurations.Store(key, conf)
}

// Put Intercept the service registration request, register and rewrite the registration information
func (proxy *RegProxy) Put(ctx context.Context, in *pb.PutRequest) (out *pb.PutResponse, err error) {
	xlog.Info("put", xlog.Any("in", in))
	var node *structs.ServiceNode
	switch {
	// The registration information is parsed and cached
	case bytes.HasPrefix(in.GetKey(), []byte("grpc:")) ||
		bytes.HasPrefix(in.GetKey(), []byte("http:")):
		// todo(gorexlv): v1 注册信息需要额外补全，后续添加, 暂时去掉
		// node, err = extractRegInfoV1(in.GetKey(), in.GetValue())
	case bytes.HasPrefix(in.GetKey(), []byte("/wsd-reg/")):
		node, err = extractRegInfoV2(in.GetKey(), in.GetValue())
	case bytes.HasPrefix(in.GetKey(), []byte("/dubbo/")):
		// todo(gorexlv): v3 dubbo注册信息需要额外补全，后续添加, 暂时去掉
		// node, err = extractRegInfoV3(in.GetKey(), in.GetValue())
	}

	if node != nil && err == nil {
		select {
		case proxy.nodeChan <- node:
		default:
		}
	}

	if xdebug.IsDevelopmentMode() {
		defer func() {
			xdebug.PrintObject("put service err", err)
			xdebug.PrintObject("put service out", out)
		}()
	}

	return proxy.KVServer.Put(ctx, in)
}

// extractRegInfoV1 ...
func extractRegInfoV1(key, val []byte) (node *structs.ServiceNode, err error) {
	// grpc service with 'grpc:' prefix
	// http service with 'http:' prefix
	// {schema}:{app_name}:v1:{env}/{ip}:{port}
	prefix, suffix := xstring.Split(string(key), "/").Head2()
	schema, appName, _, env := xstring.Split(prefix, ":").Head4()
	ip, portStr := xstring.Split(suffix, ":").Head2()

	port, _ := strconv.Atoi(portStr)
	return &structs.ServiceNode{
		AppName: appName,
		Schema:  schema,
		IP:      ip,
		Port:    strconv.Itoa(port),
		Env:     env,
	}, nil
}

// extractRegInfoV2 ...
func extractRegInfoV2(key []byte, val []byte) (node *structs.ServiceNode, err error) {
	appName, addr := xstring.Split(strings.TrimLeft(string(key), "/wsd-reg/"), "/providers/").Head2()
	xlog.Info("extractRegInfoV2", xlog.String("appName", appName), xlog.String("addr", addr))
	uri, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	var regInfo structs.RegInfo
	if err := json.Unmarshal(val, &regInfo); err != nil {
		panic(err)
	}

	node = &structs.ServiceNode{
		AppName: appName,
		Schema:  uri.Scheme,
		Port:    uri.Port(),
		Methods: []string{},
		Env:     "",
		RegInfo: regInfo,
	}
	if strings.Contains(uri.Host, ":") {
		node.IP, node.Port = xstring.Split(uri.Host, ":").Head2()
	}

	return
}

// extractRegInfoV3 ...
func extractRegInfoV3(key, val []byte) (node *structs.ServiceNode, err error) {
	// _, srvName, _, addr := xstring.Split(string(key), "/").Head4()
	// uri, err := url.Parse(addr)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// var regInfo structs.RegInfo
	// if err := json.Unmarshal(val, &regInfo); err != nil {
	// 	panic(err)
	// }
	//
	return nil, nil
}
