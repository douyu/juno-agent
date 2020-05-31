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

package model

import (
	"context"

	"github.com/douyu/jupiter/pkg/store/gorm"
)

// CheckReq ...
type CheckReq struct {
	CheckDatas []struct {
		Type string `json:"type"`
		Data string `json:"data"`
	} `json:"check_datas"`
}

// CheckResult ...
type CheckResult struct {
	gorm.Model
	NodeID int64  `gorm:"column:node_id" json:"node_id"`
	Status Status `gorm:"column:status" json:"status"`
	Count  int64  `gorm:"column:count" json:"status"`
}

// TableName ...
func (cr CheckResult) TableName() string {
	return "ma_check_result"
}

// PMTShell ...
type PMTShell struct {
	Pmt     string `json:"pmt"`
	AppName string `json:"app_name"`
	Op      int    `json:"op"`
}

// InsertCheckResult ...
func InsertCheckResult(ctx context.Context, db *gorm.DB, cr CheckResult) error {
	return nil
}
