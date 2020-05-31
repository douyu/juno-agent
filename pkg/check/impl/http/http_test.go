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

package http

import (
	"testing"
)

const (
	GetHTTPConfigTest           = `{"http":"http://192.168.132.128:50092/openapi/testGet1","method":"get","timeout":3}`
	GetHTTPConfigWithHeaderTest = `{"http":"http://192.168.132.128:50092","method":"get","timeout":3,"header":{"aid":["123"]}}`
	PostHTTPConfigWithJSONTest  = `{"http":"http://192.168.132.128:50092/openapi/testPost1","method":"POST","timeout":3,"header":{"Content-Type":["application/json"]},"body":"{\"ctype\":2,\"cid\":1,\"cplatform\":0,\"offset\":0,\"limit\":5}"}`
	PostHTTPConfigWithFormTest  = `{"http":"http://192.168.132.128:50092/openapi/testPost2","method":"POST","timeout":3,"header":{"Content-Type":["application/x-www-form-urlencoded"]},"body":"name=testdata"}`
)

func TestGetHttpConfig(t *testing.T) {
	// get instance
	httpHealthCheck := NewHTTPHealthCheck()
	type args struct {
		extConfig string
	}
	test := struct {
		name    string
		args    args
		wantErr bool
	}{
		name:    "get http",
		args:    args{extConfig: GetHTTPConfigTest},
		wantErr: false,
	}
	t.Run(test.name, func(t *testing.T) {
		if err := httpHealthCheck.LoadExtConfig(test.args.extConfig); (err != nil) != test.wantErr {
			t.Errorf("LoadExtConfig() error = %v, wantErr %v", err, test.wantErr)
		}
	})
	t.Run(test.name, func(t *testing.T) {
		if _, err := httpHealthCheck.DoHealthCheck(); err != nil {
			t.Errorf("DoHealthCheck() error = %v ", err.Error())
		}
	})
	t.Log("name: ", test.name, "config: ", httpHealthCheck.httpConfig)

}

func TestGetHttpWithHeaderConfig(t *testing.T) {
	// get instance
	httpHealthCheck := NewHTTPHealthCheck()
	type args struct {
		extConfig string
	}
	test := struct {
		name    string
		args    args
		wantErr bool
	}{
		name:    "get http",
		args:    args{extConfig: GetHTTPConfigWithHeaderTest},
		wantErr: false,
	}
	t.Run(test.name, func(t *testing.T) {
		if err := httpHealthCheck.LoadExtConfig(test.args.extConfig); (err != nil) != test.wantErr {
			t.Errorf("LoadExtConfig() error = %v, wantErr %v", err, test.wantErr)
		}
	})
	t.Run(test.name, func(t *testing.T) {
		if _, err := httpHealthCheck.DoHealthCheck(); err != nil {
			t.Errorf("DoHealthCheck() error = %v ", err.Error())
		}
	})
	t.Log("name: ", test.name, "config: ", httpHealthCheck.httpConfig)

}

func TestPostHttpWithJsonConfig(t *testing.T) {
	// get instance
	httpHealthCheck := NewHTTPHealthCheck()
	type args struct {
		extConfig string
	}
	test := struct {
		name    string
		args    args
		wantErr bool
	}{
		name:    "get http",
		args:    args{extConfig: PostHTTPConfigWithJSONTest},
		wantErr: false,
	}
	t.Run(test.name, func(t *testing.T) {
		if err := httpHealthCheck.LoadExtConfig(test.args.extConfig); (err != nil) != test.wantErr {
			t.Errorf("LoadExtConfig() error = %v, wantErr %v", err, test.wantErr)
		}
	})
	t.Run(test.name, func(t *testing.T) {
		if _, err := httpHealthCheck.DoHealthCheck(); err != nil {
			t.Errorf("DoHealthCheck() error = %v ", err.Error())
		}
	})
	t.Log("name: ", test.name, "config: ", httpHealthCheck.httpConfig)

}

func TestPostHttpWithFormConfig(t *testing.T) {
	// get instance
	httpHealthCheck := NewHTTPHealthCheck()
	type args struct {
		extConfig string
	}
	test := struct {
		name    string
		args    args
		wantErr bool
	}{
		name:    "get http",
		args:    args{extConfig: PostHTTPConfigWithFormTest},
		wantErr: false,
	}
	t.Run(test.name, func(t *testing.T) {
		if err := httpHealthCheck.LoadExtConfig(test.args.extConfig); (err != nil) != test.wantErr {
			t.Errorf("LoadExtConfig() error = %v, wantErr %v", err, test.wantErr)
		}
	})
	t.Run(test.name, func(t *testing.T) {
		if _, err := httpHealthCheck.DoHealthCheck(); err != nil {
			t.Errorf("DoHealthCheck() error = %v ", err.Error())
		}
	})
	t.Log("name: ", test.name, "config: ", httpHealthCheck.httpConfig)

}
