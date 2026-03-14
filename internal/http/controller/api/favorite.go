// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/slices"
)

type FavoriteController struct{}

func (self *FavoriteController) RegisterRoute(g *echo.Group) {
	g.POST("/favorites/:objid", self.Toggle)
	g.GET("/users/:username/favorites", self.List)
}

// Toggle 收藏/取消收藏
func (FavoriteController) Toggle(ctx echo.Context) error {
	token := getAuthToken(ctx)
	if token == "" {
		return fail(ctx, "未登录", NeedReLoginCode)
	}

	if !ValidateToken(token) {
		return fail(ctx, "token 已过期，请重新登录", NeedReLoginCode)
	}

	uid, ok := ParseToken(token)
	if !ok || uid == 0 {
		return fail(ctx, "无效的 token", NeedReLoginCode)
	}

	objid := goutils.MustInt(ctx.Param("objid"))
	objtype := goutils.MustInt(ctx.QueryParam("objtype"))
	collect := goutils.MustInt(ctx.QueryParam("collect"))

	if objid == 0 || objtype == 0 {
		return fail(ctx, "参数错误")
	}

	var err error
	if collect == 1 {
		err = logic.DefaultFavorite.Save(context.EchoContext(ctx), uid, objid, objtype)
	} else {
		err = logic.DefaultFavorite.Cancel(context.EchoContext(ctx), uid, objid, objtype)
	}

	if err != nil {
		return fail(ctx, err.Error())
	}

	return success(ctx, map[string]interface{}{
		"collected": collect == 1,
	})
}

// List 用户收藏列表
func (FavoriteController) List(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "username", username)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}

	objtype := goutils.MustInt(ctx.QueryParam("objtype"), model.TypeArticle)
	p := goutils.MustInt(ctx.QueryParam("p"), 1)
	if p < 1 {
		p = 1
	}

	rows := goutils.MustInt(ctx.QueryParam("rows"), perPage)
	if rows > 50 {
		rows = 50
	}

	favorites, total := logic.DefaultFavorite.FindUserFavorites(
		context.EchoContext(ctx), user.Uid, objtype, (p-1)*rows, rows,
	)

	data := map[string]interface{}{
		"user":     user,
		"objtype":  objtype,
		"total":    total,
		"page":     p,
		"has_more": int64(p*rows) < total,
	}

	if total > 0 {
		objids := slices.StructsIntSlice(favorites, "Objid")

		switch objtype {
		case model.TypeTopic:
			data["topics"] = logic.DefaultTopic.FindByTids(objids)
		case model.TypeArticle:
			data["articles"] = logic.DefaultArticle.FindByIds(objids)
		case model.TypeResource:
			data["resources"] = logic.DefaultResource.FindByIds(objids)
		case model.TypeProject:
			data["projects"] = logic.DefaultProject.FindByIds(objids)
		}
	}

	return success(ctx, data)
}
