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

package util

const DefaultConfig = `
[plugin]
    [plugin.regProxy]
        # 注册中心地址
        endpoints=["127.0.0.1:2379"]
        timeout="3s"
        secure=false
        enable = true
    [plugin.confProxy]
        # 配置中心地址
        env=["dev","live","pre"]
        prefix = "/juno-agent"
        timeout="3s"
        enable = true
        #配置中心数据源
        [pugin.confProxy.mysql]
            enable=false
            dsn=""
            secure=false
        [plugin.confProxy.etcd]
            enable=true
            endpoints=["127.0.0.1:2379"]
    [plugin.supervisor]
        enable = true
        dir = "/etc/supervisor/conf.d"
    [plugin.systemd]
        enable = true
        dir = "/etc/systemd/system"
    [plugin.nginx]
        enable = false
        dir = "/usr/local/openresty/nginx/conf"
    [plugin.report]
        enable = true
        debug = false
        addr = "http://127.0.0.1:60812/api/v1/resource/node/heartbeat"
        internal = "60s"
        hostName = "JUNO_HOST" # 环境变量的名称，或者命令行参数的名称
        regionCode = "REGION_CODE" # 环境变量的名称，或者命令行参数的名称
        regionName = "REGION_NAME"
        zoneCode = "ZONE_CODE"
        zoneName = "ZONE_NAME"
        env = "env"
    [plugin.healthCheck]
        enable = true
    [plugin.process]
        enable = true
[jupiter.logger.default]
    name = "default"
    debug = true
[jupiter.server]
  [jupiter.server.grpc]
    host = "0.0.0.0"
    port = 60813

  [jupiter.server.http]
    host = "0.0.0.0"
    port = 60814
`
