package db

import (
	"database/sql"
	"fmt"

	. "github.com/polaris1119/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var MasterDB *xorm.Engine

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
	if err = open(); err != nil {
		panic(err)
	}
}

func open() error {
	var err error

	MasterDB, err = xorm.NewEngine("mysql", dns)
	if err != nil {
		return err
	}

	maxIdle := ConfigFile.MustInt("mysql", "max_idle", 2)
	maxConn := ConfigFile.MustInt("mysql", "max_conn", 10)

	MasterDB.SetMaxIdleConns(maxIdle)
	MasterDB.SetMaxOpenConns(maxConn)

	showSQL := ConfigFile.MustBool("xorm", "show_sql", false)
	logLevel := ConfigFile.MustInt("xorm", "log_level", 0)

	MasterDB.ShowSQL(showSQL)
	MasterDB.Logger().SetLevel(core.LogLevel(logLevel))

	// 启用缓存
	// cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	// MasterDB.SetDefaultCacher(cacher)

	return nil
}

func StdMasterDB() *sql.DB {
	return MasterDB.DB().DB
}
