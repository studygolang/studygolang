// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package app

import (
	"logic"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"

	. "http"
	"model"
)

type ProjectController struct{}

// 注册路由
func (self ProjectController) RegisterRoute(g *echo.Group) {
	g.GET("/projects", self.ReadList)
	g.GET("/project/detail", self.Detail)
}

// ReadList 开源项目列表页
func (ProjectController) ReadList(ctx echo.Context) error {
	limit := 20

	lastId := goutils.MustInt(ctx.QueryParam("base"))
	projects := logic.DefaultProject.FindBy(ctx, limit+5, lastId)
	if projects == nil {
		return fail(ctx, "获取失败")
	}

	projectList := make([]map[string]interface{}, 0, len(projects))
	for _, project := range projects {
		if lastId > 0 {
			// TODO: 推荐？
			// if project.Top == 1 {
			// 	continue
			// }
		}
		projectList = append(projectList, map[string]interface{}{
			"id":       project.Id,
			"name":     project.Name,
			"category": project.Category,
			"logo":     project.Logo,
			"tags":     project.Tags,
			"viewnum":  project.Viewnum,
			"cmtnum":   project.Cmtnum,
			"likenum":  project.Likenum,
			"author":   project.Author,
			"ctime":    project.Ctime,
		})
	}

	hasMore := false
	if len(projectList) > limit {
		hasMore = true
		projectList = projectList[:limit]
	}

	data := map[string]interface{}{
		"has_more": hasMore,
		"projects": projectList,
	}

	return success(ctx, data)
}

// Detail 项目详情
func (ProjectController) Detail(ctx echo.Context) error {
	id := goutils.MustInt(ctx.QueryParam("id"))
	project := logic.DefaultProject.FindOne(ctx, id)
	if project == nil || project.Id == 0 {
		return fail(ctx, "获取失败或已下线")
	}

	logic.Views.Incr(Request(ctx), model.TypeProject, project.Id)

	// 为了阅读数即时看到
	project.Viewnum++

	return success(ctx, map[string]interface{}{"project": project})
}
