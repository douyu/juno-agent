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
	"reflect"

	"github.com/douyu/juno-agent/pkg/check/view"
)

// DataSource interface
type DataSource interface {
	DoHealthCheck() (resHealthCheck *view.ResHealthCheck, err error)
	LoadExtConfig(extConfig string) (err error)
}

// invoke
func invoke(f interface{}, params ...interface{}) []reflect.Value {
	fv := reflect.ValueOf(f)
	realParams := make([]reflect.Value, len(params))
	for i, item := range params {
		realParams[i] = reflect.ValueOf(item)
	}
	rs := fv.Call(realParams)
	return rs
}

//Get  return HealthCheck instance
func Get(reflect []reflect.Value) (DataSource, error) {
	if len(reflect) > 0 {
		return reflect[0].Interface().(DataSource), nil
	}
	return nil, errors.New("get implement err")
}
