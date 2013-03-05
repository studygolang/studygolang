package util

import (
	"fmt"
	"strings"
)

type Sqler interface {
	Tablename() string
	Columns() []string
	SelectCols() string // 需要查询哪些字段
	Where() string
	Order() string
	Limit() string
}

func InsertSql(sqler Sqler) string {
	columns := sqler.Columns()
	columnStr := strings.Join(columns, ",")
	placeHolder := strings.Repeat("?,", len(columns))
	sql := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", sqler.Tablename(), columnStr, placeHolder[:len(placeHolder)-1])
	return sql
}

func SelectSql(sqler Sqler) string {
	where := sqler.Where()
	if sqler.Where() != "" {
		where = "WHERE " + where
	}
	order := sqler.Order()
	if order != "" {
		order = "ORDER BY " + order
	}
	limit := sqler.Limit()
	if limit != "" {
		limit = "LIMIT " + limit
	}
	return fmt.Sprintf("SELECT %s FROM %s %s %s %s", sqler.SelectCols(), sqler.Tablename(), where, order, limit)
}
