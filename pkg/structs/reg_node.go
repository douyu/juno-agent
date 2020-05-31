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
	"fmt"
	"time"

	"github.com/uber-go/atomic"
)

// Result ...
type Result int

// SuccessResult ...
var SuccessResult Result = 1

// FailResult ...
var FailResult Result = 2

// Status ...
type Status int

var (
	// Unregistered ...
	Unregistered Status = 1
	// Registered ...
	Registered Status = 2
)

// checkResult check reg result
type checkResult struct {
	result Result
	count  int64
	status Status
}

// RegNode 检查点
type RegNode struct {
	ServiceNode
	createTime      int64
	LastCheckTime   *atomic.Int64
	NextCheckTime   *atomic.Int64
	Disable         *atomic.Bool
	lastCheckResult *checkResult
}

// Address reg addr
func (n RegNode) Address() string {
	return fmt.Sprintf("%s:%s", n.IP, n.Port)
}

// NewRegNode ...
func NewRegNode(node ServiceNode) *RegNode {
	return &RegNode{
		ServiceNode:   node,
		createTime:    time.Now().Unix(),
		Disable:       atomic.NewBool(false),
		LastCheckTime: atomic.NewInt64(0),
		NextCheckTime: atomic.NewInt64(0),
	}
}
