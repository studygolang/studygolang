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
	g.GET("/favorites/:objid/status", self.Status)
	g.GET("/users/:username/favorites", self.List)
}

// Status 查询当前用户是否已收藏
func (FavoriteController) Status(ctx echo.Context) error {
	token := getAuthToken(ctx)
	if token == "" {
		return success(ctx, map[string]interface{}{"has_favorite": false})
	}

	uid, _, valid := ValidateTokenAuto(token)
	if !valid || uid == 0 {
		return success(ctx, map[string]interface{}{"has_favorite": false})
	}

	objid := goutils.MustInt(ctx.Param("objid"))
	objtype := goutils.MustInt(ctx.QueryParam("objtype"))

	if objid == 0 || !isValidObjType(objtype) {
		return success(ctx, map[string]interface{}{"has_favorite": false})
	}

	hasFavorite := logic.DefaultFavorite.HadFavorite(context.EchoContext(ctx), uid, objid, objtype)
	return success(ctx, map[string]interface{}{"has_favorite": hasFavorite})
}

// Toggle 收藏/取消收藏
func (FavoriteController) Toggle(ctx echo.Context) error {
	// 写操作必须校验用户状态（与 master NeedLogin 一致）
	uid, err := parseActiveAuthUID(ctx)
	if err != nil {
		return err
	}

	objid := goutils.MustInt(ctx.Param("objid"))
	objtype := goutils.MustInt(ctx.FormValue("objtype"))
	collect := goutils.MustInt(ctx.FormValue("collect"))

	if objid == 0 || !isValidObjType(objtype) {
		return fail(ctx, "参数错误")
	}

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
// 响应字段对齐前端 favoriteAPI.listByUsername 契约：
//
//	{ favorites[], total, page, has_more }
//
// 同时保留按 objtype 分键（topics/articles/...）以便兼容老客户端。
// favorites 是统一字段，前端默认读它即可。
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

	// 统一字段：前端默认读这个 key（任何 objtype 都有）
	objList := make([]interface{}, 0)

	if total > 0 {
		objids := slices.StructsIntSlice(favorites, "Objid")

		switch objtype {
		case model.TypeTopic:
			topics := logic.DefaultTopic.FindByTids(objids)
			for i := range topics {
				objList = append(objList, topics[i])
			}
		case model.TypeArticle:
			articles := logic.DefaultArticle.FindByIds(objids)
			for i := range articles {
				objList = append(objList, articles[i])
			}
		case model.TypeResource:
			resources := logic.DefaultResource.FindByIds(objids)
			for i := range resources {
				objList = append(objList, resources[i])
			}
		case model.TypeProject:
			projects := logic.DefaultProject.FindByIds(objids)
			for i := range projects {
				objList = append(objList, projects[i])
			}
		}
	}

	data := map[string]interface{}{
		"favorites": objList, // 统一字段（前端契约）
		"objtype":   objtype,
		"total":     total,
		"page":      p,
		"has_more":  int64(p*rows) < total,
	}

	return success(ctx, data)
}
