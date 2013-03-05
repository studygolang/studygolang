package model

import (
	"database/sql"
	_ "github.com/Go-SQL-Driver/MySQL"
	"strings"
	"util"
)

type Dao struct {
	*sql.DB
	// 构造sql语句相关
	tablename string
	where     string
	whereVal  []interface{} // where条件对应中字段对应的值
	limit     string
	order     string
	// 插入需要
	columns   []string      // 需要插入数据的字段
	colValues []interface{} // 需要插入字段对应的值
	// 查询需要
	selectCols string // 想要查询那些字段，接在SELECT之后的，默认为"*"
}

func (this *Dao) Open() (err error) {
	this.DB, err = sql.Open("mysql", "root:@/studygolang?charset=utf8")
	return
}

// Insert 插入数据
func (this *Dao) Insert() (sql.Result, error) {
	strSql := util.InsertSql(this)
	err := this.Open()
	if err != nil {
		return nil, err
	}
	defer this.Close()
	stmt, err := this.Prepare(strSql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(this.ColValues()...)
}

// Find 查找单条数据
func (this *Dao) Find() (*sql.Row, error) {
	strSql := util.SelectSql(this)
	err := this.Open()
	if err != nil {
		return nil, err
	}
	defer this.Close()
	stmt, err := this.Prepare(strSql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.QueryRow(this.whereVal...), nil
}

func (this *Dao) Columns() []string {
	return this.columns
}

func (this *Dao) ColValues() []interface{} {
	return this.colValues
}

func (this *Dao) SetWhere(args ...string) {
	fields := make([]string, len(args))
	this.whereVal = make([]interface{}, len(args))
	for i, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			// TODO:怎么处理
		}
		fields[i] = parts[0] + "=?"
		this.whereVal[i] = parts[1]
	}
	this.where = strings.Join(fields, " AND ")
}

func (this *Dao) SelectCols() string {
	if this.selectCols == "" {
		return "*"
	}
	return this.selectCols
}

func (this *Dao) Where() string {
	return this.where
}

func (this *Dao) SetOrder(order string) {
	this.order = order
}

func (this *Dao) Order() string {
	return this.order
}

func (this *Dao) SetLimit(limit string) {
	this.limit = limit
}

func (this *Dao) Limit() string {
	return this.limit
}

func (this *Dao) Tablename() string {
	return this.tablename
}

type ORMer interface {
	Insert()
	Update()
	Find()
	FindAll()
	Delete()
}
