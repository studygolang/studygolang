package logic

import (
	"bytes"
	"io/ioutil"
	"model"

	"github.com/polaris1119/config"
	"golang.org/x/net/context"

	. "db"
)

type InstallLogic struct{}

var DefaultInstall = InstallLogic{}

func (InstallLogic) CreateTable(ctx context.Context) error {
	objLog := GetLogger(ctx)

	dbFile := config.ROOT + "/config/db.sql"
	buf, err := ioutil.ReadFile(dbFile)

	if err != nil {
		objLog.Errorln("create table, read db file error:", err)
		return err
	}

	sqlSlice := bytes.Split(buf, []byte("CREATE TABLE"))
	MasterDB.Exec("SET SQL_MODE='ALLOW_INVALID_DATES';")
	for _, oneSql := range sqlSlice {
		strSql := string(bytes.TrimSpace(oneSql))
		if strSql == "" {
			continue
		}

		strSql = "CREATE TABLE " + strSql
		_, err1 := MasterDB.Exec(strSql)
		if err1 != nil {
			objLog.Errorln("create table error:", err1)
			err = err1
		}
	}

	return err
}

// InitTable 初始化数据表
func (InstallLogic) InitTable(ctx context.Context) error {
	objLog := GetLogger(ctx)

	total, err := MasterDB.Count(new(model.Role))
	if err != nil {
		return err
	}

	if total > 0 {
		return nil
	}

	dbFile := config.ROOT + "/config/init.sql"
	buf, err := ioutil.ReadFile(dbFile)
	if err != nil {
		objLog.Errorln("init table, read init file error:", err)
		return err
	}

	sqlSlice := bytes.Split(buf, []byte("INSERT INTO"))
	for _, oneSql := range sqlSlice {
		strSql := string(bytes.TrimSpace(oneSql))
		if strSql == "" {
			continue
		}

		strSql = "INSERT INTO " + strSql
		_, err1 := MasterDB.Exec(strSql)
		if err1 != nil {
			objLog.Errorln("create table error:", err1)
			err = err1
		}
	}

	return err
}

func (InstallLogic) IsTableExist(ctx context.Context) bool {
	exists, err := MasterDB.IsTableExist(new(model.User))
	if err != nil || !exists {
		return false
	}

	return true
}

// HadRootUser 是否已经创建了超级用户
func (InstallLogic) HadRootUser(ctx context.Context) bool {
	user := &model.User{}
	_, err := MasterDB.Where("is_root=?", 1).Get(user)
	if err != nil {
		// 发生错误，认为已经创建了
		return true
	}

	return user.Uid != 0
}
