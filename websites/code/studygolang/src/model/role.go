package model

import (
	"logger"
	"util"
)

// 角色分界点：roleid小于该值，则没有管理权限
const AdminMinRoleId int = 6

// 角色信息
type Role struct {
	Roleid int    `json:"roleid"`
	Name   string `json:"name"`
	ctime  string

	// 数据库访问对象
	*Dao
}

func NewRole() *Role {
	return &Role{
		Dao: &Dao{tablename: "role"},
	}
}

func (this *Role) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (this *Role) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *Role) FindAll(selectCol ...string) ([]*Role, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	roleList := make([]*Role, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		role := NewRole()
		colFieldMap := role.colFieldMap()
		scanInterface := make([]interface{}, 0, colNum)
		for _, column := range selectCol {
			scanInterface = append(scanInterface, colFieldMap[column])
		}
		err = rows.Scan(scanInterface...)
		if err != nil {
			logger.Errorln("FindAll Scan Error:", err)
			continue
		}
		roleList = append(roleList, role)
	}
	return roleList, nil
}

func (this *Role) prepareInsertData() {
	this.columns = []string{"name"}
	this.colValues = []interface{}{this.Name}
}

func (this *Role) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"roleid": &this.Roleid,
		"name":   &this.Name,
		"ctime":  &this.ctime,
	}
}

// 角色权限信息
type RoleAuthority struct {
	Roleid int    `json:"roleid"`
	Aid    int    `json:"aid"`
	Name   string `json:"name"`
	// 不导出
	ctime string

	// 内嵌
	*Dao
}

/*
func NewRoleAuthority() *RoleAuthority {
	return &RoleAuthority{
		Dao: &Dao{tablename: "role_authority"},
	}
}

func (this *RoleAuthority) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (this *RoleAuthority) Find() error {
	row, err := this.Dao.Find()
	if err != nil {
		return err
	}
	return row.Scan(&this.Uid, &this.Username, &this.Email, &this.Name, &this.open)
}

func (this *RoleAuthority) FindAll() ([]*RoleAuthority, error) {
	rows, err := this.Dao.FindAll()
	if err != nil {
		return nil, err
	}
	// TODO:
	userList := make([]*User, 0, 10)
	for rows.Next() {
		user := NewUser()
		rows.Scan(&user.Uid, &user.Email, &user.open, &user.Username, &user.Name, &user.Avatar, &user.City, &user.Company, &user.Github, &user.Weibo, &user.Website, &user.Status, &user.Introduce, &user.ctime)
		userList = append(userList, user)
	}
	return userList, nil
}

func (this *RoleAuthority) prepareInsertData() {
	this.columns = []string{"username", "email", "name", "avatar", "city", "company", "github", "weibo", "website", "status", "introduce"}
	this.colValues = []interface{}{this.Username, this.Email, this.Name, this.Avatar, this.City, this.Company, this.Github, this.Weibo, this.Website, this.Status, this.Introduce}
}
*/
