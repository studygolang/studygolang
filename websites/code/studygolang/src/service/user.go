package service

import (
	"model"
	"net/url"
	"util"
)

func CreateUser(form url.Values) (errMsg string, err error) {
	user := model.NewUser()
	err = util.ConvertAssign(user, form)
	if err != nil {
		errMsg = err.Error()
		return
	}
	uid, err := user.Insert()
	if err != nil {
		// TODO:记日志
		errMsg = "内部服务器错误"
		return
	}

	userLogin := model.NewUserLogin()
	err = util.ConvertAssign(userLogin, form)
	if err != nil {
		errMsg = err.Error()
		return
	}
	userLogin.Uid = int(uid)
	_, err = userLogin.Insert()
	if err != nil {
		errMsg = "内部服务器错误"
		return
	}
	return
}
