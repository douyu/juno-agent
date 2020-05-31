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

package mysql

import (
	"errors"
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/jinzhu/gorm"

	// mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// SQLCommon ...
type (
	// SQLCommon alias of gorm.SQLCommon
	SQLCommon = gorm.SQLCommon
	// Callback alias of gorm.Callback
	Callback = gorm.Callback
	// CallbackProcessor alias of gorm.CallbackProcessor
	CallbackProcessor = gorm.CallbackProcessor
	// Dialect alias of gorm.Dialect
	Dialect = gorm.Dialect
	// Scope ...
	Scope = gorm.Scope
	// DB ...
	DB = gorm.DB
	// Model ...
	Model = gorm.Model
	// ModelStruct ...
	ModelStruct = gorm.ModelStruct
	// Field ...
	Field = gorm.Field
	// StructField ...
	StructField = gorm.StructField
	// RowQueryResult ...
	RowQueryResult = gorm.RowQueryResult
	// RowsQueryResult ...
	RowsQueryResult = gorm.RowsQueryResult
	// Association ...
	Association = gorm.Association
	// Errors ...
	Errors = gorm.Errors
	// Logger ...
	Logger = gorm.Logger
)

var (
	errSlowCommand = errors.New("mysql slow command")

	// IsRecordNotFoundError ...
	IsRecordNotFoundError = gorm.IsRecordNotFoundError

	// ErrRecordNotFound returns a "record not found error". Occurs only when attempting to query the database with a struct; querying with a slice won't return this error
	ErrRecordNotFound = gorm.ErrRecordNotFound
	// ErrInvalidSQL occurs when you attempt a query with invalid SQL
	ErrInvalidSQL = gorm.ErrInvalidSQL
	// ErrInvalidTransaction occurs when you are trying to `Commit` or `Rollback`
	ErrInvalidTransaction = gorm.ErrInvalidTransaction
	// ErrCantStartTransaction can't start transaction when you are trying to start one with `Begin`
	ErrCantStartTransaction = gorm.ErrCantStartTransaction
	// ErrUnaddressable unaddressable value
	ErrUnaddressable = gorm.ErrUnaddressable
)

// Open ...
func Open(dialect string, check *MysqlHealthCheck) (*DB, error) {
	inner, err := gorm.Open(dialect, check.DSN)
	if err != nil {
		return nil, err
	}
	// 设置默认连接配置
	inner.DB().SetMaxIdleConns(100)
	inner.DB().SetMaxOpenConns(50)

	if xdebug.IsDevelopmentMode() {
		inner.LogMode(true)
	}
	return inner, err
}
