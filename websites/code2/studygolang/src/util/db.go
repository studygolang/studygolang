// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package util

import (
	"fmt"
	"strings"
)

type Sqler interface {
	Tablename() string
	Columns() []string
	SelectCols() string // 需要查询哪些字段
	GetWhere() string
	GetOrder() string
	GetLimit() string
}

func InsertSql(sqler Sqler) string {
	columns := sqler.Columns()
	columnStr := "`" + strings.Join(columns, "`,`") + "`"
	placeHolder := strings.Repeat("?,", len(columns))
	sql := fmt.Sprintf("INSERT INTO `%s`(%s) VALUES(%s)", sqler.Tablename(), columnStr, placeHolder[:len(placeHolder)-1])
	return strings.TrimSpace(sql)
}

func UpdateSql(sqler Sqler) string {
	columnStr := strings.Join(sqler.Columns(), ",")
	if columnStr == "" {
		return ""
	}
	where := sqler.GetWhere()
	if where != "" {
		where = "WHERE " + where
	}
	sql := fmt.Sprintf("UPDATE `%s` SET %s %s", sqler.Tablename(), columnStr, where)
	return strings.TrimSpace(sql)
}

func DeleteSql(sqler Sqler) string {
	where := sqler.GetWhere()
	if where != "" {
		where = "WHERE " + where
	}
	sql := fmt.Sprintf("DELETE FROM `%s` %s", sqler.Tablename(), where)
	return strings.TrimSpace(sql)
}

func CountSql(sqler Sqler) string {
	where := sqler.GetWhere()
	if where != "" {
		where = "WHERE " + where
	}
	sql := fmt.Sprintf("SELECT COUNT(1) AS total FROM `%s` %s", sqler.Tablename(), where)
	return strings.TrimSpace(sql)
}

func SelectSql(sqler Sqler) string {
	where := sqler.GetWhere()
	if where != "" {
		where = "WHERE " + where
	}
	order := sqler.GetOrder()
	if order != "" {
		order = "ORDER BY " + order
	}
	limit := sqler.GetLimit()
	if limit != "" {
		limit = "LIMIT " + limit
	}
	sql := fmt.Sprintf("SELECT %s FROM `%s` %s %s %s", sqler.SelectCols(), sqler.Tablename(), where, order, limit)
	return strings.TrimSpace(sql)
}
