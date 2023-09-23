// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package admin

import (
	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type UserController struct{}

// 注册路由
func (self UserController) RegisterRoute(g *echo.Group) {
	g.GET("/user/user/list", self.UserList)
	g.POST("/user/user/query.html", self.UserQuery)
	g.GET("/user/user/detail", self.Detail)
	g.POST("/user/user/modify", self.Modify)
	g.POST("/user/user/add_black", self.AddBlack)
}

// UserList 所有用户（分页）
func (UserController) UserList(ctx echo.Context) error {
	curPage, limit := parsePage(ctx)

	users, total := logic.DefaultUser.FindUserByPage(context.EchoContext(ctx), nil, curPage, limit)

	data := map[string]interface{}{
		"datalist":   users,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	return render(ctx, "user/list.html,user/query.html", data)
}

func (UserController) UserQuery(ctx echo.Context) error {
	curPage, limit := parsePage(ctx)
	conds := parseConds(ctx, []string{"uid", "username", "email"})

	users, total := logic.DefaultUser.FindUserByPage(context.EchoContext(ctx), conds, curPage, limit)

	data := map[string]interface{}{
		"datalist":   users,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	return renderQuery(ctx, "user/query.html", data)
}

func (UserController) Detail(ctx echo.Context) error {
	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", ctx.QueryParam("uid"))

	data := map[string]interface{}{
		"user": user,
	}

	return render(ctx, "user/detail.html", data)
}

func (UserController) Modify(ctx echo.Context) error {
	uid := ctx.FormValue("uid")

	amount := goutils.MustInt(ctx.FormValue("amount"))
	forms, _ := ctx.FormParams()
	if amount > 0 {
		logic.DefaultUserRich.Recharge(context.EchoContext(ctx), uid, forms)
	}
	logic.DefaultUser.AdminUpdateUser(context.EchoContext(ctx), uid, forms)

	return success(ctx, nil)
}

func (UserController) AddBlack(ctx echo.Context) error {
	uid := goutils.MustInt(ctx.FormValue("uid"))
	err := logic.DefaultUser.UpdateUserStatus(context.EchoContext(ctx), uid, model.UserStatusOutage)
	if err != nil {
		return fail(ctx, 1, err.Error())
	}

	// 将用户 IP 加入黑名单
	logic.DefaultRisk.AddBlackIPByUID(uid)

	truncate := goutils.MustBool(ctx.FormValue("truncate"))
	if truncate {
		err = logic.DefaultUser.DeleteUserContent(context.EchoContext(ctx), uid)
		if err != nil {
			return fail(ctx, 1, err.Error())
		}
	}

	return success(ctx, nil)
}
