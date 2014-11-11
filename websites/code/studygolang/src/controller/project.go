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

// 在需要评论（喜欢）且要回调的地方注册评论（喜欢）对象
func init() {
	// 注册评论（喜欢）对象
	service.RegisterCommentObject(model.TYPE_PROJECT, service.ProjectComment{})
	service.RegisterLikeObject(model.TYPE_PROJECT, service.ProjectLike{})
}

// 开源项目列表页
// uri: /projects
func ProjectsHandler(rw http.ResponseWriter, req *http.Request) {
	limit := 20

	lastId := req.FormValue("lastid")
	if lastId == "" {
		lastId = "0"
	}

	projects := service.FindProjects(lastId, "25")
	if projects == nil {
		// TODO:服务暂时不可用？
	}

	num := len(projects)
	if num == 0 {
		if lastId == "0" {
			util.Redirect(rw, req, "/")
		} else {
			util.Redirect(rw, req, "/projects")
		}

		return
	}

	var (
		hasPrev, hasNext bool
		prevId, nextId   int
	)

	if lastId != "0" {
		prevId, _ = strconv.Atoi(lastId)

		// 避免因为项目下线，导致判断错误（所以 > 5）
		if prevId-projects[0].Id > 5 {
			hasPrev = false
		} else {
			prevId += limit
			hasPrev = true
		}
	}

	if num > limit {
		hasNext = true
		projects = projects[:limit]
		nextId = projects[limit-1].Id
	} else {
		nextId = projects[num-1].Id
	}

	pageInfo := map[string]interface{}{
		"has_prev": hasPrev,
		"prev_id":  prevId,
		"has_next": hasNext,
		"next_id":  nextId,
	}

	// 获取当前用户喜欢对象信息
	user, ok := filter.CurrentUser(req)
	var likeFlags map[int]int
	if ok {
		uid := user["uid"].(int)
		likeFlags, _ = service.FindUserLikeObjects(uid, model.TYPE_PROJECT, projects[0].Id, nextId)
	}

	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/projects/list.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"projects": projects, "activeProjects": "active", "page": pageInfo, "likeflags": likeFlags})
}

// 新建项目
// uri: /project/new{json:(|.json)}
func NewProjectHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	name := req.PostFormValue("name")
	// 请求新建项目页面
	if name == "" || req.Method != "POST" || vars["json"] == "" {
		project := model.NewOpenProject()
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/projects/new.html")
		filter.SetData(req, map[string]interface{}{"project": project, "activeProjects": "active"})
		return
	}

	user, _ := filter.CurrentUser(req)
	err := service.PublishProject(user, req.PostForm)
	if err != nil {
		fmt.Fprint(rw, `{"ok": 0, "error":"内部服务错误！"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":""}`)
}

// 修改项目
// uri: /project/modify{json:(|.json)}
func ModifyProjectHandler(rw http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		util.Redirect(rw, req, "/projects")
		return
	}

	vars := mux.Vars(req)
	// 请求编辑项目页面
	if req.Method != "POST" || vars["json"] == "" {
		project := service.FindProject(id)
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/projects/new.html")
		filter.SetData(req, map[string]interface{}{"project": project, "activeProjects": "active"})
		return
	}

	user, _ := filter.CurrentUser(req)
	err := service.PublishProject(user, req.PostForm)
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

// 项目详情
// uri: /p/{uniq}
func ProjectDetailHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	project := service.FindProject(vars["uniq"])
	if project == nil {
		util.Redirect(rw, req, "/projects")
		return
	}

	likeFlag := 0
	hadCollect := 0
	user, ok := filter.CurrentUser(req)
	if ok {
		uid := user["uid"].(int)
		likeFlag = service.HadLike(uid, project.Id, model.TYPE_PROJECT)
		hadCollect = service.HadFavorite(uid, project.Id, model.TYPE_PROJECT)
	}

	service.Views.Incr(req, model.TYPE_PROJECT, project.Id)

	// 为了阅读数即时看到
	project.Viewnum++

	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/projects/detail.html")
	filter.SetData(req, map[string]interface{}{"activeProjects": "active", "project": project, "likeflag": likeFlag, "hadcollect": hadCollect})
}

// 检测 uri 对应的项目是否存在(验证，true表示不存在；false表示存在)
// uri: /project/uri.json
func ProjectUriHandler(rw http.ResponseWriter, req *http.Request) {
	uri := req.FormValue("uri")
	if uri == "" {
		fmt.Fprint(rw, `true`)
		return
	}

	if service.ProjectUriExists(uri) {
		fmt.Fprint(rw, `false`)
		return
	}
	fmt.Fprint(rw, `true`)
}
