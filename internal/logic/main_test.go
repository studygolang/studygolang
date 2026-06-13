package logic

import (
	"os"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"xorm.io/xorm"
	"xorm.io/xorm/core"

	. "github.com/studygolang/studygolang/db"
)

// TestMain 为 logic 包测试初始化一个 mock 数据库，避免 MasterDB 为 nil 导致 panic。
// sqlmock 对未期望的查询返回 error，logic 代码走错误处理分支（return err/false），
// 使"不 panic 就过"类测试可通过。验证具体数据的测试需另行精确 mock 或标记为集成测试。
func TestMain(m *testing.M) {
	sqlDB, _, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	engine, err := xorm.NewEngineWithDB("mysql", "", core.FromDB(sqlDB))
	if err != nil {
		panic(err)
	}

	MasterDB = engine

	os.Exit(m.Run())
}
