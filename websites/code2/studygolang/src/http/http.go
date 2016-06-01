// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package http

import (
	"bytes"
	"global"
	"html/template"
	"logic"
	"model"
	"net/http"
	"strings"
	"time"
	"util"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/polaris1119/config"
)

var Store = sessions.NewCookieStore([]byte(config.ConfigFile.MustValue("global", "cookie_secret")))

func SetCookie(ctx echo.Context, username string) {
	Store.Options.HttpOnly = true

	session := GetCookieSession(ctx)
	if ctx.FormValue("remember_me") != "1" {
		// 浏览器关闭，cookie删除，否则保存30天(github.com/gorilla/sessions 包的默认值)
		session.Options = &sessions.Options{
			Path:     "/",
			HttpOnly: true,
		}
	}
	session.Values["username"] = username
	req := Request(ctx)
	resp := ResponseWriter(ctx)
	session.Save(req, resp)
}

func GetCookieSession(ctx echo.Context) *sessions.Session {
	session, _ := Store.Get(Request(ctx), "user")
	return session
}

func Request(ctx echo.Context) *http.Request {
	return ctx.Request().(*standard.Request).Request
}

func ResponseWriter(ctx echo.Context) http.ResponseWriter {
	return ctx.Response().(*standard.Response).ResponseWriter
}

// 自定义模板函数
var funcMap = template.FuncMap{
	// 获取gravatar头像
	"gravatar": util.Gravatar,
	// 转为前端显示需要的时间格式
	"formatTime": func(i interface{}) string {
		ctime, ok := i.(string)
		if !ok {
			return ""
		}
		t, _ := time.Parse("2006-01-02 15:04:05", ctime)
		return t.Format(time.RFC3339) + "+08:00"
	},
	"substring": util.Substring,
	"add": func(nums ...interface{}) int {
		total := 0
		for _, num := range nums {
			if n, ok := num.(int); ok {
				total += n
			}
		}
		return total
	},
	"explode": func(s, sep string) []string {
		return strings.Split(s, sep)
	},
	"noescape": func(s string) template.HTML {
		return template.HTML(s)
	},
	"timestamp": func() int64 {
		return time.Now().Unix()
	},
}

const (
	LayoutTpl      = "common/layout.html"
	AdminLayoutTpl = "common.html"
)

// Render html 输出
func Render(ctx echo.Context, contentTpl string, data map[string]interface{}) error {
	if data == nil {
		data = map[string]interface{}{}
	}

	objLog := logic.GetLogger(ctx)

	contentTpl = LayoutTpl + "," + contentTpl
	// 为了使用自定义的模板函数，首先New一个以第一个模板文件名为模板名。
	// 这样，在ParseFiles时，新返回的*Template便还是原来的模板实例
	htmlFiles := strings.Split(contentTpl, ",")
	for i, contentTpl := range htmlFiles {
		htmlFiles[i] = config.TemplateDir + contentTpl
	}
	tpl, err := template.New("layout.html").Funcs(funcMap).ParseFiles(htmlFiles...)
	if err != nil {
		objLog.Errorf("解析模板出错（ParseFiles）：[%q] %s\n", Request(ctx).RequestURI, err)
		return err
	}

	return executeTpl(ctx, tpl, data)
}

// RenderAdmin html 输出
func RenderAdmin(ctx echo.Context, contentTpl string, data map[string]interface{}) error {
	if data == nil {
		data = map[string]interface{}{}
	}

	objLog := logic.GetLogger(ctx)

	contentTpl = AdminLayoutTpl + "," + contentTpl
	// 为了使用自定义的模板函数，首先New一个以第一个模板文件名为模板名。
	// 这样，在ParseFiles时，新返回的*Template便还是原来的模板实例
	htmlFiles := strings.Split(contentTpl, ",")
	for i, contentTpl := range htmlFiles {
		htmlFiles[i] = config.TemplateDir + "admin/" + contentTpl
	}

	requestURI := Request(ctx).RequestURI
	tpl, err := template.New("common.html").Funcs(funcMap).ParseFiles(htmlFiles...)
	if err != nil {
		objLog.Errorf("解析模板出错（ParseFiles）：[%q] %s\n", requestURI, err)
		return err
	}

	// 当前用户信息
	curUser := ctx.Get("user").(*model.Me)

	if menu1, menu2, curMenu1 := logic.DefaultAuthority.GetUserMenu(ctx, curUser, requestURI); menu2 != nil {
		data["menu1"] = menu1
		data["menu2"] = menu2
		data["uri"] = requestURI
		data["cur_menu1"] = curMenu1
	}

	return executeTpl(ctx, tpl, data)
}

// 后台 query 查询返回结果
func RenderQuery(ctx echo.Context, contentTpl string, data map[string]interface{}) error {
	objLog := logic.GetLogger(ctx)

	contentTpl = "common_query.html," + contentTpl
	contentTpls := strings.Split(contentTpl, ",")
	for i, contentTpl := range contentTpls {
		contentTpls[i] = config.TemplateDir + "admin/" + strings.TrimSpace(contentTpl)
	}

	requestURI := Request(ctx).RequestURI
	tpl, err := template.New("common_query.html").Funcs(funcMap).ParseFiles(contentTpls...)
	if err != nil {
		objLog.Errorf("解析模板出错（ParseFiles）：[%q] %s\n", requestURI, err)
		return err
	}

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, data)
	if err != nil {
		objLog.Errorf("执行模板出错（Execute）：[%q] %s\n", requestURI, err)
		return err
	}

	return ctx.HTML(http.StatusOK, buf.String())
}

func executeTpl(ctx echo.Context, tpl *template.Template, data map[string]interface{}) error {
	objLog := logic.GetLogger(ctx)

	// 如果没有定义css和js模板，则定义之
	if jsTpl := tpl.Lookup("js"); jsTpl == nil {
		tpl.Parse(`{{define "js"}}{{end}}`)
	}
	if jsTpl := tpl.Lookup("css"); jsTpl == nil {
		tpl.Parse(`{{define "css"}}{{end}}`)
	}

	// 当前用户信息
	curUser, ok := ctx.Get("user").(*model.Me)
	if ok {
		data["me"] = curUser
	} else {
		data["me"] = map[string]interface{}{}
	}

	// websocket主机
	if global.OnlineEnv() {
		data["wshost"] = config.ConfigFile.MustValue("global", "domain")
	} else {
		data["wshost"] = global.App.Host + ":" + global.App.Port
	}
	global.App.SetUptime()
	data["app"] = global.App

	buf := new(bytes.Buffer)
	err := tpl.Execute(buf, data)
	if err != nil {
		objLog.Errorln("excute template error:", err)
		return err
	}

	return ctx.HTML(http.StatusOK, buf.String())
}
