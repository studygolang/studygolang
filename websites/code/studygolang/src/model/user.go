package model

import (
	"time"
)

const (
	tablename = "user_login"
)

type UserLogin struct {
	Uid      int `json:"uid"`
	Username string	`json:"username"`
	Passwd   string	`json:"passwd"`
	Email    string	`json:"email"`
	passcode string // 加密随机串

	// 数据库访问对象
	*Dao
}

func NewUserLogin() *UserLogin {
	return &UserLogin{
		Dao: &Dao{tablename: tablename},
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
	this.passcode = "123ff"
	this.colValues = []interface{}{this.Uid, this.Username, this.Passwd, this.Email, this.passcode}
}

type User struct {
	Uid       int
	Username  string
	Passwd    string
	Email     string
	Name      string
	Avatar    string
	City      string
	Company   string
	Github    string
	Weibo     string
	Website   string
	Status    string
	Introduce string
	// 不导出
	open  int
	ctime time.Time
}

func NewUser() *User {
	return nil
}
