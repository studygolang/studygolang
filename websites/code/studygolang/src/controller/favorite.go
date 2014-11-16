// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"filter"
	"github.com/studygolang/mux"
	"model"
	"service"
	"util"
)

// 收藏(取消收藏)
// uri: /favorite/{objid:[0-9]+}.json
func FavoriteHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	user, _ := filter.CurrentUser(req)

	if !util.CheckInt(req.PostForm, "objtype") {
		fmt.Fprint(rw, `{"ok": 0, "error":"参数错误"}`)
		return
	}

	var err error

	objtype := util.MustInt(req.PostFormValue("objtype"))
	collect := util.MustInt(req.PostFormValue("collect"))
	if collect == 1 {
		err = service.SaveFavorite(user["uid"].(int), util.MustInt(vars["objid"]), objtype)
	} else {
		err = service.CancelFavorite(user["uid"].(int), util.MustInt(vars["objid"]), objtype)
	}

	if err != nil {
		fmt.Fprint(rw, `{"ok": 0, "error":"`+err.Error()+`""}`)
		return
	}

	fmt.Fprint(rw, `{"ok": 1, "message":"success"}`)
}

// 我的(某人的)收藏
// uri: /favorites/{username}
func SomeoneFavoritesHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]
	user := service.FindUserByUsername(username)
	if user == nil {
		util.Redirect(rw, req, "/")
		return
	}

	objtype, err := strconv.Atoi(req.FormValue("objtype"))
	if err != nil {
		objtype = model.TYPE_ARTICLE
	}

	p, err := strconv.Atoi(req.FormValue("p"))
	if err != nil {
		p = 1
	}

	data := map[string]interface{}{"objtype": objtype, "user": user}

	rows := 20
	favorites, total := service.FindUserFavorites(user.Uid, objtype, (p-1)*rows, rows)
	if total > 0 {
		objids := util.Models2Intslice(favorites, "Objid")

		switch objtype {
		case model.TYPE_TOPIC:
			data["topics"] = service.FindTopicsByIds(objids)
		case model.TYPE_ARTICLE:
			data["articles"] = service.FindArticlesByIds(objids)
		case model.TYPE_RESOURCE:
			data["resources"] = service.FindResourcesByIds(objids)
		case model.TYPE_WIKI:
			// data["wikis"] = service.FindArticlesByIds(objids)
		case model.TYPE_PROJECT:
			data["projects"] = service.FindProjectsByIds(objids)
		}

	}

	uri := fmt.Sprintf("/favorites/%s?objtype=%d&p=%d", user.Username, objtype, p)
	data["pageHtml"] = service.GenPageHtml(p, rows, total, uri)

	// 设置模板数据
	filter.SetData(req, data)
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/favorite.html")
}
