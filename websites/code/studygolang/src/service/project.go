// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"logger"
	"model"
	"util"
)

func PublishProject(username string, form url.Values) error {
	if ProjectUriExists(form.Get("uri")) {
		return errors.New("uri存在")
	}

	project := model.NewOpenProject()
	util.ConvertAssign(project, form)

	project.Username = username
	project.Ctime = util.TimeNow()
	project.Uri = strings.ToLower(project.Uri)

	github := "github.com"
	pos := strings.Index(project.Src, github)
	if pos != -1 {
		project.Repo = project.Src[pos+len(github)+1:]
	}

	_, err := project.Insert()

	if err != nil {
		logger.Errorln("Publish Project error:", err)
	}

	return err
}

func ProjectUriExists(uri string) bool {
	project := model.NewOpenProject()
	err := project.Where("uri=?", uri).Find("id")
	if err != nil {
		return false
	}

	if project.Id > 0 {
		return true
	}

	return false
}

// 获取开源项目列表（分页）
func FindProjects(lastId, limit string) []*model.OpenProject {
	project := model.NewOpenProject()

	cond := "status IN(0,1)"
	if lastId != "0" {
		cond += " AND id<" + lastId
	}

	projectList, err := project.Where(cond).Order("id DESC").Limit(limit).
		FindAll()
	if err != nil {
		logger.Errorln("project service FindProjects Error:", err)
		return nil
	}

	return projectList
}

// 获取单个项目
func FindProject(uniq string) *model.OpenProject {
	field := "id"
	_, err := strconv.Atoi(uniq)
	if err != nil {
		field = "uri"
	}

	project := model.NewOpenProject()
	err = project.Where(field+"=?", uniq).Find()

	if err != nil {
		logger.Errorln("project service FindProject error:", err)
		return nil
	}

	if project.Id == 0 {
		return nil
	}

	return project
}

// 获取多个项目详细信息
func FindProjectsByIds(ids []int) []*model.OpenProject {
	if len(ids) == 0 {
		return nil
	}
	inIds := util.Join(ids, ",")
	projects, err := model.NewOpenProject().Where("id in(" + inIds + ")").FindAll()
	if err != nil {
		logger.Errorln("project service FindProjectsByIds error:", err)
		return nil
	}
	return projects
}

// 项目评论
type ProjectComment struct{}

// 更新该项目的评论信息
// cid：评论id；objid：被评论对象id；uid：评论者；cmttime：评论时间
func (self ProjectComment) UpdateComment(cid, objid, uid int, cmttime string) {
	id := strconv.Itoa(objid)

	// 更新评论数（TODO：暂时每次都更新表）
	err := model.NewOpenProject().Where("id="+id).Increment("cmtnum", 1)
	if err != nil {
		logger.Errorln("更新项目评论数失败：", err)
	}
}

func (self ProjectComment) String() string {
	return "project"
}

// 实现 CommentObjecter 接口
func (self ProjectComment) SetObjinfo(ids []int, commentMap map[int][]*model.Comment) {
	projects := FindProjectsByIds(ids)
	if len(projects) == 0 {
		return
	}

	for _, project := range projects {
		objinfo := make(map[string]interface{})
		objinfo["title"] = project.Category + project.Name

		for _, comment := range commentMap[project.Id] {
			comment.Objinfo = objinfo
		}
	}
}

// 项目喜欢
type ProjectLike struct{}

// 更新该项目的喜欢数
// objid：被喜欢对象id；num: 喜欢数(负数表示取消喜欢)
func (self ProjectLike) UpdateLike(objid, num int) {
	id := strconv.Itoa(objid)

	// 更新喜欢数（TODO：暂时每次都更新表）
	err := model.NewOpenProject().Where("id="+id).Increment("likenum", num)
	if err != nil {
		logger.Errorln("更新项目喜欢数失败：", err)
	}
}

func (self ProjectLike) String() string {
	return "project"
}
