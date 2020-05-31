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
	"encoding/json"
	"fmt"

	"github.com/douyu/jupiter/pkg/store/gorm"
)

// Status ...
type Status int

const (
	// StatusInit ...
	StatusInit Status = 1
	// StatusSuccess ...
	StatusSuccess Status = 2
	// NodeTableName ...
	NodeTableName = "ma_node"
)

// InsertNode insert the node to db
func InsertNode(ctx context.Context, db *gorm.DB, node *Node) error {
	return gorm.WithContext(ctx, db).Table(NodeTableName).Create(node).Error
}

// ListNodes list the node infos
func ListNodes(ctx context.Context, db *gorm.DB) ([]Node, error) {
	var nodes []Node
	err := gorm.WithContext(ctx, db).Table(NodeTableName).
		Where("deleteTime=?", 0).
		Find(&nodes).Error
	return nodes, err
}

// Node represents service node
type Node struct {
	ID         int64  `gorm:"primary_key"`
	AppName    string `gorm:"column:aname"`
	Schema     string `gorm:"column:schema"` // http / grpc
	IP         string `gorm:"column:ip"`
	Port       int32  `gorm:"column:port"`
	CreateTime *int64 `gorm:"column:ctime"`
	UpdateTime *int64 `gorm:"column:utime"`
	DeleteTime *int64 `gorm:"column:dtime"`
}

// TableName return the db name
func (n Node) TableName() string {
	return NodeTableName
}

// String return the node string
func (n Node) String() string {
	content, _ := json.Marshal(n)
	return fmt.Sprintf("%d=>%s", n.ID, string(content))
}
