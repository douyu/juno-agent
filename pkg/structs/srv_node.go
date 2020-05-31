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

import "fmt"

// RegInfo ...
type RegInfo struct {
	Name     string               `json:"name" toml:"name"`
	Scheme   string               `json:"scheme" toml:"scheme"`
	Address  string               `json:"address" toml:"address"`
	Labels   map[string]string    `json:"labels" toml:"labels"`
	Services map[string]DubboInfo `json:"services" toml:"services"`
}

// DubboInfo ...
type DubboInfo struct {
	Namespace string            `json:"namespace" toml:"namespace"`
	Name      string            `json:"name" toml:"name"`
	Labels    map[string]string `json:"labels" toml:"labels"`
	Methods   []string          `json:"methods" toml:"methods"`
}

// ServiceNode ...
type ServiceNode struct {
	AppName       string                    `json:"appName" toml:"appName"`
	Schema        string                    `json:"schema" toml:"schema"`
	IP            string                    `json:"ip" toml:"ip"`
	Port          string                    `json:"port" toml:"port"`
	Methods       []string                  `json:"methods" toml:"methods"`
	Env           string                    `json:"env" toml:"env"`
	RegInfo       RegInfo                   `json:"regInfo" toml:"regInfo"`
	Configuration *ServiceNodeConfiguration `json:"configuration" toml:"configuration"`
}

// Address ...
func (sn *ServiceNode) Address() string {
	return fmt.Sprintf("%s:%s", sn.IP, sn.Port)
}

// Key ...
func (sn *ServiceNode) Key() string {
	return fmt.Sprintf("%s://%s:%s?app=%s", sn.Schema, sn.IP, sn.Port, sn.AppName)
}

// Resource ...
type Resource struct {
}

// ServiceNodeConfiguration ...
type ServiceNodeConfiguration struct {
}

// AppConfiguration ...
type AppConfiguration struct {
	Content   string     `json:"content"`
	Resources []Resource `json:"resources"`
	Metadata  Metadata   `json:"metadata"`
}

// Metadata ...
type Metadata struct {
	Format    string `json:"format"`    // Content encoding format: toml/yaml/json
	Timestamp int64  `json:"timestamp"` //
	Version   string `json:"version"`
	FileName  string `json:"file_name"` // fileName: config-live/config-trunk
	Encoded   bool   `json:"encoded"`   // content :Whether the encryption
	AppName   string `json:"app_name" toml:"app_name" `
}
