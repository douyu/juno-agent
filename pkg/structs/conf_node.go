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

package structs

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// ConfNode confuNode 主要代表部署在应用机器上的应用的具体配置信息
type ConfNode struct {
	AppName       string            `json:"appName" toml:"appName"`             // 应用名称
	AppEnvi       string            `json:"appEnvi" toml:"appEnvi"`             // 应用部署环境
	IP            string            `json:"ip" toml:"ip"`                       // 应用部署机器ip
	Port          string            `json:"port" toml:"port"`                   // 应用部署机器port
	FileName      string            `json:"file_name" toml:"file_name"`         // 应用部署配置文件名称
	Configuration *AppConfiguration `json:"configuration" toml:"configuration"` // 应用部署配置具体信息
}

// ContentNode ...
type ContentNode struct {
	Content string `json:"content"` // 应用部署配置内容
	Version int64  `json:"version"` // 应用部署配置版本
}

// ConfKey 存储配置的key字段 for instance /wsd-sider/wsd.go-141-66.stress.unp/wsd-ocr-task-intercepter-go/stress/static/config-stress.toml
type ConfKey struct {
	Prefix   string `json:"prefix"`
	Hostname string `json:"hostname"`
	AppName  string `json:"app_name"`
	EnvName  string `json:"env_name"`
	Rest     string `json:"rest"`
	FileName string `json:"file_name"`
	Port     string `json:"port"`
}

// CheckValid check valid
func (c *ConfKey) CheckValid() error {
	if c.Prefix == "" {
		return errors.New("prefix empty")
	}
	if c.AppName == "" {
		return errors.New("appname empty")
	}
	if c.Hostname == "" {
		return errors.New("hostname empty")
	}
	if c.FileName == "" {
		return errors.New("filename empty")
	}
	if c.Port == "" {
		return errors.New("port empty")
	}
	return nil
}

// ParserConfKey parse conf key
func ParserConfKey(key string) (keyData ConfKey, err error) {
	arr := strings.Split(key, "/")
	if len(arr) != 8 {
		err = fmt.Errorf("key invaild")
		return
	}
	keyData = ConfKey{
		Prefix:   arr[1],
		Hostname: arr[2],
		AppName:  arr[3],
		EnvName:  arr[4],
		Rest:     arr[5],
		FileName: arr[6],
		Port:     arr[7],
	}
	return
}

// ConfValue {"content":"","metadata":{"timestamp":1560354378,"version":"v1","format":"toml"}}
type ConfValue struct {
	Content  string   `json:"content"`
	Metadata MetaData `json:"metadata"`
}

// MetaData ...
type MetaData struct {
	Timestamp int64  `json:"timestamp"`
	Version   string `json:"version"`
	Format    string `json:"format"`
	Path      string `json:"path"`
}

// CheckValid ...
func (c *ConfValue) CheckValid() error {
	if c.Metadata.Timestamp == 0 {
		return errors.New("timestamp empty")
	}
	if c.Metadata.Format == "" {
		return errors.New("format empty")
	}
	if c.Metadata.Version == "" {
		return errors.New("version empty")
	}
	if c.Content == "" {
		return errors.New("content empty")
	}
	return nil
}

// ParserConfValue ...
func ParserConfValue(value string) (valueData ConfValue, err error) {
	if err = json.Unmarshal([]byte(value), &valueData); err != nil {
		return
	}
	return
}

// ConfReport {"file_name":"config-live.toml","md5":"0f07572ba1212a75d8b5a0167c5507c2","hostname":"wsd-go.a1-41-116.live.unp","env":"live","timestamp":1560477493}
type ConfReport struct {
	FileName   string `json:"file_name"`
	MD5        string `json:"md5"`
	Hostname   string `json:"hostname"`
	Env        string `json:"env"`
	Timestamp  int64  `json:"timestamp"`
	IP         string `json:"ip"`
	HealthPort string `json:"health_port"`
}

// JSONString json
func (c *ConfReport) JSONString() string {
	buf, _ := json.Marshal(c)
	return string(buf)
}

// ServiceConf ...
type ServiceConf struct {
	Content   string     `json:"content"`
	Resources []Resource `json:"resources"`
	Metadata  struct {
		Format    string `json:"format"`    // content编码格式: toml/yaml/json
		Timestamp int64  `json:"timestamp"` //
		Version   string `json:"version"`
		FileName  string `json:"file_name"` // 文件名称 config-live/config-trunk
		Encoded   bool   `json:"encoded"`   // content是否加密
		AppName   string `json:"app_name" toml:"app_name" `
	} `json:"metadata"`
}
