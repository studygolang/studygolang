// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package controller

import (
	"fmt"
	"http/middleware"
	"logic"
	"model"
	"net/http"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/slices"
)

type FavoriteController struct{}

// 注册路由
func (self FavoriteController) RegisterRoute(e *echo.Echo) {
	e.Post("/favorite/:objid", echo.HandlerFunc(self.Create), middleware.NeedLogin())
	e.Get("/favorites/:username", echo.HandlerFunc(self.ReadList))
}

// Create 收藏(取消收藏)
func (FavoriteController) Create(ctx echo.Context) error {
	objtype := goutils.MustInt(ctx.FormValue("objtype"))
	objid := goutils.MustInt(ctx.Param("objid"))
	collect := goutils.MustInt(ctx.FormValue("collect"))

	user := ctx.Get("user").(*model.Me)

	var err error
	if collect == 1 {
		err = logic.DefaultFavorite.Save(ctx, user.Uid, objid, objtype)
	} else {
		err = logic.DefaultFavorite.Cancel(ctx, user.Uid, objid, objtype)
	}

	if err != nil {
		return fail(ctx, 1, err.Error())
	}

	return success(ctx, nil)
}

// ReadList 我的(某人的)收藏
func (FavoriteController) ReadList(ctx echo.Context) error {
	username := ctx.Param("username")
	user := logic.DefaultUser.FindOne(ctx, "username", username)
	if user == nil || user.Uid == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	objtype := goutils.MustInt(ctx.QueryParam("objtype"), model.TypeArticle)
	p := goutils.MustInt(ctx.QueryParam("p"), 1)

	data := map[string]interface{}{"objtype": objtype, "user": user}

	rows := 20
	favorites, total := logic.DefaultFavorite.FindUserFavorites(ctx, user.Uid, objtype, (p-1)*rows, rows)
	if total > 0 {
		objids := slices.StructsIntSlice(favorites, "Objid")

		switch objtype {
		case model.TypeTopic:
			data["topics"] = logic.DefaultTopic.FindByTids(objids)
		case model.TypeArticle:
			data["articles"] = logic.DefaultArticle.FindByIds(objids)
		case model.TypeResource:
			data["resources"] = logic.DefaultResource.FindByIds(objids)
		case model.TypeWiki:
			// data["wikis"] = logic.DefaultWiki.FindWikisByIds(objids)
		case model.TypeProject:
			data["projects"] = logic.DefaultProject.FindByIds(objids)
		}
	}

	uri := fmt.Sprintf("/favorites/%s?objtype=%d&p=%d", user.Username, objtype, p)
	data["pageHtml"] = logic.GenPageHtml(p, rows, int(total), uri)

	return render(ctx, "favorite.html", data)
}
