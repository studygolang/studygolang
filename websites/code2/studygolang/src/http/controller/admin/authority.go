// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package admin

import (
	"logic"
	"net/http"

	"github.com/labstack/echo"
)

type AuthorityController struct{}

// 注册路由
func (self AuthorityController) RegisterRoute(g *echo.Group) {
	g.GET("/user/auth/list", self.AuthList)
	g.POST("/user/auth/query.html", self.AuthQuery)
}

// AuthList 所有权限（分页）
func (AuthorityController) AuthList(ctx echo.Context) error {
	curPage, limit := parsePage(ctx)

	total := len(logic.Authorities)
	newLimit := limit
	if total < limit {
		newLimit = total
	}

	data := map[string]interface{}{
		"datalist":   logic.Authorities[(curPage - 1):newLimit],
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	return render(ctx, "authority/list.html,authority/query.html", data)
}

func (AuthorityController) AuthQuery(ctx echo.Context) error {
	curPage, limit := parsePage(ctx)

	conds := parseConds(ctx, []string{"route", "name"})

	authorities, total := logic.DefaultAuthority.FindAuthoritiesByPage(ctx, conds, curPage, limit)

	if authorities == nil {
		return ctx.HTML(http.StatusInternalServerError, "500")
	}

	data := map[string]interface{}{
		"datalist":   authorities,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	return renderQuery(ctx, "authority/query.html", data)
}

// func NewAuthorityHandler(rw http.ResponseWriter, req *http.Request) {
// 	var data = make(map[string]interface{})

// 	if req.PostFormValue("submit") == "1" {
// 		user, _ := filter.CurrentUser(req)
// 		username := user["username"].(string)

// 		errMsg, err := service.SaveAuthority(req.PostForm, username)
// 		if err != nil {
// 			data["ok"] = 0
// 			data["error"] = errMsg
// 		} else {
// 			data["ok"] = 1
// 			data["msg"] = "添加成功"
// 		}
// 	} else {
// 		menu1, menu2 := service.GetMenus()
// 		allmenu2, _ := json.Marshal(menu2)

// 		// 设置内容模板
// 		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/authority/new.html")
// 		data["allmenu1"] = menu1
// 		data["allmenu2"] = string(allmenu2)
// 	}

// 	filter.SetData(req, data)
// }

// func ModifyAuthorityHandler(rw http.ResponseWriter, req *http.Request) {
// 	var data = make(map[string]interface{})

// 	if req.PostFormValue("submit") == "1" {
// 		user, _ := filter.CurrentUser(req)
// 		username := user["username"].(string)

// 		errMsg, err := service.SaveAuthority(req.PostForm, username)
// 		if err != nil {
// 			data["ok"] = 0
// 			data["error"] = errMsg
// 		} else {
// 			data["ok"] = 1
// 			data["msg"] = "修改成功"
// 		}
// 	} else {
// 		menu1, menu2 := service.GetMenus()
// 		allmenu2, _ := json.Marshal(menu2)

// 		authority := service.FindAuthority(req.FormValue("aid"))

// 		if authority == nil || authority.Aid == 0 {
// 			rw.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		// 设置内容模板
// 		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/authority/modify.html")
// 		data["allmenu1"] = menu1
// 		data["allmenu2"] = string(allmenu2)
// 		data["authority"] = authority
// 	}

// 	filter.SetData(req, data)
// }

// func DelAuthorityHandler(rw http.ResponseWriter, req *http.Request) {
// 	var data = make(map[string]interface{})

// 	aid := req.FormValue("aid")

// 	if _, err := strconv.Atoi(aid); err != nil {
// 		data["ok"] = 0
// 		data["error"] = "aid不是整型"

// 		filter.SetData(req, data)
// 		return
// 	}

// 	if err := service.DelAuthority(aid); err != nil {
// 		data["ok"] = 0
// 		data["error"] = "删除失败！"
// 	} else {
// 		data["ok"] = 1
// 		data["msg"] = "删除成功！"
// 	}

// 	filter.SetData(req, data)
// }
