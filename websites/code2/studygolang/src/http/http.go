// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package http

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/polaris1119/config"
)

var Store = sessions.NewCookieStore([]byte(config.ConfigFile.MustValue("global", "cookie_secret")))

func SetCookie(ctx *echo.Context, username string) {
	Store.Options.HttpOnly = true

	session := GetCookieSession(ctx)
	if ctx.Form("remember_me") != "1" {
		// 浏览器关闭，cookie删除，否则保存30天(github.com/gorilla/sessions 包的默认值)
		session.Options = &sessions.Options{
			Path:     "/",
			HttpOnly: true,
		}
	}
	session.Values["username"] = username
	session.Save(ctx.Request(), ctx.Response())
}

func GetCookieSession(ctx *echo.Context) *sessions.Session {
	session, _ := Store.Get(ctx.Request(), "user")
	return session
}
