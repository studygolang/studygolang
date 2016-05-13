// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package http

import (
	"bytes"
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
}

const LayoutTpl = "common/layout.html"

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
	data["wshost"] = "127.0.0.1:8088"
	data["build"] = map[string]string{
		"version": "1.0",        // version.Version,
		"date":    "2016-01-16", // version.Date,
	}

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, data)
	if err != nil {
		objLog.Errorln("excute template error:", err)
		return err
	}

	return ctx.HTML(http.StatusOK, buf.String())
}
