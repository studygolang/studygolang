// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"fmt"
	"html/template"
	"logic"
	"net/http"

	"vendor/github.com/polaris1119/goutils"

	"github.com/labstack/echo"

	"filter"
	"model"
	"service"
	"util"

	"github.com/studygolang/mux"
)

// 在需要评论（喜欢）且要回调的地方注册评论（喜欢）对象
func init() {
	// 注册评论（喜欢）对象
	logic.RegisterCommentObject(model.TypeResource, logic.ResourceComment{})
	// service.RegisterLikeObject(model.TYPE_RESOURCE, service.ResourceLike{})
}

type ResourceController struct{}

// 注册路由
func (this *ResourceController) RegisterRoute(e *echo.Echo) {
	e.Get("/resources", echo.HandlerFunc(this.ReadList))
	e.Get("/resources/cat/:catid", echo.HandlerFunc(this.ReadCatResources))
}

// ReadList 资源索引页
func (ResourceController) ReadList(ctx echo.Context) error {
	return ctx.Redirect(http.StatusSeeOther, "/resources/cat/1")
}

// ReadCatResources 某个分类的资源列表
func (ResourceController) ReadCatResources(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.Query("p"), 1)
	paginator := logic.NewPaginator(curPage)
	catid := goutils.MustInt(ctx.Param("catid"))

	resources, total := logic.DefaultResource.FindByCatid(ctx, paginator, catid)
	pageHtml := paginator.SetTotal(total).GetPageHtml(Request(ctx).URL.Path)

	return render(ctx, "resources/index.html", map[string]interface{}{"activeResources": "active", "resources": resources, "categories": logic.AllCategory, "page": template.HTML(pageHtml), "curCatid": catid})
}

// 某个资源详细页
// uri: /resources/{id:[0-9]+}
func ResourceDetailHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	resource, comments := service.FindResource(vars["id"])

	if len(resource) == 0 {
		util.Redirect(rw, req, "/resources")
		return
	}

	likeFlag := 0
	hadCollect := 0
	user, ok := filter.CurrentUser(req)
	if ok {
		uid := user["uid"].(int)
		id := resource["id"].(int)
		likeFlag = service.HadLike(uid, id, model.TYPE_RESOURCE)
		hadCollect = service.HadFavorite(uid, id, model.TYPE_RESOURCE)
	}

	service.Views.Incr(req, model.TYPE_RESOURCE, util.MustInt(vars["id"]))

	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/resources/detail.html,/template/common/comment.html")
	filter.SetData(req, map[string]interface{}{"activeResources": "active", "resource": resource, "comments": comments, "likeflag": likeFlag, "hadcollect": hadCollect})
}

// 发布新资源
// uri: /resources/new{json:(|.json)}
func NewResourceHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	title := req.PostFormValue("title")
	// 请求新建资源页面
	if title == "" || req.Method != "POST" || vars["json"] == "" {
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/resources/new.html")
		filter.SetData(req, map[string]interface{}{"activeResources": "active", "categories": service.AllCategory})
		return
	}

	errMsg := ""
	resForm := req.PostFormValue("form")
	if resForm == model.LinkForm {
		if req.PostFormValue("url") == "" {
			errMsg = "url不能为空"
		}
	} else {
		if req.PostFormValue("content") == "" {
			errMsg = "内容不能为空"
		}
	}
	if errMsg != "" {
		fmt.Fprint(rw, `{"ok": 0, "error":"`+errMsg+`"}`)
		return
	}

	user, _ := filter.CurrentUser(req)
	err := service.PublishResource(user, req.PostForm)
	if err != nil {
		fmt.Fprint(rw, `{"ok": 0, "error":"内部服务错误，请稍候再试！"}`)
		return
	}

	fmt.Fprint(rw, `{"ok": 1, "data":""}`)
}

// 修改資源
// uri: /resources/modify{json:(|.json)}
func ModifyResourceHandler(rw http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		util.Redirect(rw, req, "/resources")
		return
	}

	vars := mux.Vars(req)
	// 请求编辑資源页面
	if req.Method != "POST" || vars["json"] == "" {
		resource := service.FindResourceById(id)
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/resources/new.html")
		filter.SetData(req, map[string]interface{}{"resource": resource, "activeResources": "active", "categories": service.AllCategory})
		return
	}

	user, _ := filter.CurrentUser(req)
	err := service.PublishResource(user, req.PostForm)
	if err != nil {
		if err == service.NotModifyAuthorityErr {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		fmt.Fprint(rw, `{"ok": 0, "error":"内部服务错误！"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":""}`)
}
