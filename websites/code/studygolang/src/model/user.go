package model

import (
	"fmt"
	"math/rand"
	"time"
	"util"
)

// 用户登录信息
type UserLogin struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
	Email    string `json:"email"`
	passcode string // 加密随机串

	// 数据库访问对象
	*Dao
}

func NewUserLogin() *UserLogin {
	return &UserLogin{
		Dao: &Dao{tablename: "user_login"},
	}
}

func (this *UserLogin) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (this *UserLogin) Find() error {
	row, err := this.Dao.Find()
	if err != nil {
		return err
	}
	return row.Scan(&this.Uid, &this.Username, &this.Passwd, &this.Email, &this.passcode)
}

func (this *UserLogin) prepareInsertData() {
	this.columns = []string{"uid", "username", "passwd", "email", "passcode"}
	this.passcode = fmt.Sprintf("%x", rand.Int31())
	// 密码经过md5(passwd+passcode)加密保存
	this.Passwd = util.Md5(this.Passwd + this.passcode)
	this.colValues = []interface{}{this.Uid, this.Username, this.Passwd, this.Email, this.passcode}
}

// 用户基本信息
type User struct {
	Uid       int    `json:"uid"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	City      string `json:"city"`
	Company   string `json:"company"`
	Github    string `json:"github"`
	Weibo     string `json:"weibo"`
	Website   string `json:"website"`
	Status    string `json:"status"`
	Introduce string `json:"introduce"`
	// 不导出
	open  int
	ctime time.Time

	// 内嵌
	*Dao
}

func NewUser() *User {
	return &User{
		Dao: &Dao{tablename: "user_info"},
	}
}

func (this *User) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (this *User) Find() error {
	row, err := this.Dao.Find()
	if err != nil {
		return err
	}
	return row.Scan(&this.Uid, &this.Username, &this.Email, &this.Name, &this.open)
}

func (this *User) prepareInsertData() {
	this.columns = []string{"username", "email", "name", "avatar", "city", "company", "github", "weibo", "website", "status", "introduce"}
	this.colValues = []interface{}{this.Username, this.Email, this.Name, this.Avatar, this.City, this.Company, this.Github, this.Weibo, this.Website, this.Status, this.Introduce}
}
