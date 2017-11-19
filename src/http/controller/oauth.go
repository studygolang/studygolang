// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"logic"
	"model"
	"net/http"

	. "http"

	"github.com/labstack/echo"
)

type OAuthController struct{}

// 注册路由
func (self OAuthController) RegisterRoute(g *echo.Group) {
	g.Get("/oauth/github/callback", self.GithubCallback)
	g.Get("/oauth/github/login", self.GithubLogin)
}

func (OAuthController) GithubLogin(ctx echo.Context) error {
	uri := ctx.QueryParam("uri")
	url := logic.DefaultThirdUser.GithubAuthCodeUrl(ctx, uri)
	return ctx.Redirect(http.StatusSeeOther, url)
}

func (OAuthController) GithubCallback(ctx echo.Context) error {
	code := ctx.FormValue("code")

	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		// 已登录用户，绑定 github
		logic.DefaultThirdUser.BindGithub(ctx, code, me)

		redirectURL := ctx.QueryParam("redirect_url")
		if redirectURL == "" {
			redirectURL = "/account/edit#connection"
		}
		return ctx.Redirect(http.StatusSeeOther, redirectURL)
	}

	user, err := logic.DefaultThirdUser.LoginFromGithub(ctx, code)
	if err != nil || user.Uid == 0 {
		var errMsg = ""
		if err != nil {
			errMsg = err.Error()
		} else {
			errMsg = "服务内部错误"
		}

		return render(ctx, "login.html", map[string]interface{}{"error": errMsg})
	}

	// 登录成功，种cookie
	SetLoginCookie(ctx, user.Username)

	if user.Balance == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/balance")
	}

	return ctx.Redirect(http.StatusSeeOther, "/")
}
