// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package filter

import (
	"config"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/studygolang/mux"
	"logger"
	"net/http"
	"service"
	"util"
)

// 没登陆且没有cookie，则跳转到登录页
type LoginFilter struct {
	*mux.EmptyFilter
}

func (this *LoginFilter) PreFilter(rw http.ResponseWriter, req *http.Request) bool {
	logger.Debugln("LoginFilter PreFilter...")
	if _, ok := CurrentUser(req); !ok {
		logger.Debugln("需要登录")
		// 支持跳转回原来访问的页面
		NewFlash(rw, req).AddFlash(req.RequestURI, "uri")
		util.Redirect(rw, req, "/account/login")
		return false
	}
	return true
}

// flash messages
type Flash struct {
	rw  http.ResponseWriter
	req *http.Request

	session *sessions.Session
}

func NewFlash(rw http.ResponseWriter, req *http.Request) *Flash {
	session, _ := Store.Get(req, "sg_flash")
	return &Flash{rw: rw, req: req, session: session}
}

func (this *Flash) AddFlash(value interface{}, vars ...string) {
	this.session.AddFlash(value, vars...)
	this.session.Save(this.req, this.rw)
}

func (this *Flash) Flashes(vars ...string) []interface{} {
	defer this.session.Save(this.req, this.rw)
	if flashes := this.session.Flashes(vars...); len(flashes) > 0 {
		return flashes
	}
	return nil
}

// 如果没登陆但有Cookie时，自动登录
type CookieFilter struct {
	*mux.EmptyFilter
}

func (this *CookieFilter) PreFilter(rw http.ResponseWriter, req *http.Request) bool {
	user, _ := CurrentUser(req)
	// 已登录且请求登录页面
	if user != nil && req.RequestURI == "/account/login" {
		util.Redirect(rw, req, "/")
	}
	return true
}

func (this *CookieFilter) PostFilter(rw http.ResponseWriter, req *http.Request) bool {
	// 删除设置的用户信息
	context.Delete(req, userkey)
	return true
}

// 定义key，标识存储user信息
type loginKey int

const userkey loginKey = 0

func getUser(req *http.Request) map[string]interface{} {
	if rv := context.Get(req, userkey); rv != nil {
		return rv.(map[string]interface{})
	}
	return nil
}

func setUser(req *http.Request, user map[string]interface{}) {
	context.Set(req, userkey, user)
}

var Store = sessions.NewCookieStore([]byte(config.Config["cookie_secret"]))

// 获得当前登录用户
func CurrentUser(req *http.Request) (map[string]interface{}, bool) {
	user := getUser(req)
	if len(user) != 0 {
		return user, true
	}
	session, _ := Store.Get(req, "user")
	username, ok := session.Values["username"]
	if !ok {
		return nil, false
	}
	user, err := service.FindCurrentUser(username.(string))
	if err != nil {
		return nil, false
	}
	setUser(req, user)
	return user, true
}
