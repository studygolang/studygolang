// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package db

import (
	"fmt"

	. "config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB gorm.DB

var dns string

func init() {
	mysqlConfig, err := ConfigFile.GetSection("mysql")
	if err != nil {
		panic(err)
	}

	dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlConfig["user"],
		mysqlConfig["password"],
		mysqlConfig["host"],
		mysqlConfig["port"],
		mysqlConfig["dbname"],
		mysqlConfig["charset"])

	// 启动时就打开数据库连接
	open()
}

func open() error {
	maxIdle := ConfigFile.MustInt("mysql", "max_idle", 2)
	maxConn := ConfigFile.MustInt("mysql", "max_conn", 10)

	var err error

	DB, err = gorm.Open("mysql", dns)
	if err != nil {
		return err
	}

	DB.DB().SetMaxIdleConns(maxIdle)
	DB.DB().SetMaxOpenConns(maxConn)

	// Disable table name's pluralization
	DB.SingularTable(true)

	DB.LogMode(ConfigFile.MustBool("mysql", "gorm_log", false))

	return nil
}
