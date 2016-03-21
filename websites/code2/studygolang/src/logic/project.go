// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package logic

import (
	"model"
	"time"

	. "db"

	"github.com/polaris1119/logger"
	"golang.org/x/net/context"
)

type ProjectLogic struct{}

var DefaultProject = ProjectLogic{}

// Total 开源项目总数
func (ProjectLogic) Total() int64 {
	total, err := MasterDB.Count(new(model.OpenProject))
	if err != nil {
		logger.Errorln("ProjectLogic Total error:", err)
	}
	return total
}

// FindBy 获取开源项目列表（分页）
func (ProjectLogic) FindBy(ctx context.Context, limit int, lastIds ...int) []*model.OpenProject {
	objLog := GetLogger(ctx)

	dbSession := MasterDB.Where("status IN(?,?)", model.ProjectStatusNew, model.ProjectStatusOnline)
	if len(lastIds) > 0 {
		dbSession.And("id<?", lastIds[0])
	}

	projectList := make([]*model.OpenProject, 0)
	err := dbSession.OrderBy("id DESC").Limit(limit).Find(&projectList)
	if err != nil {
		objLog.Errorln("ProjectLogic FindBy Error:", err)
		return nil
	}

	return projectList
}

// FindByIds 获取多个项目详细信息
func (ProjectLogic) FindByIds(ids []int) []*model.OpenProject {
	if len(ids) == 0 {
		return nil
	}

	projects := make([]*model.OpenProject, 0)
	err := MasterDB.In("id", ids).Find(&projects)
	if err != nil {
		logger.Errorln("ProjectLogic FindByIds error:", err)
		return nil
	}
	return projects
}

// 项目评论
type ProjectComment struct{}

// 更新该项目的评论信息
// cid：评论id；objid：被评论对象id；uid：评论者；cmttime：评论时间
func (self ProjectComment) UpdateComment(cid, objid, uid int, cmttime time.Time) {
	// 更新评论数（TODO：暂时每次都更新表）
	_, err := MasterDB.Id(objid).Incr("cmtnum", 1).Update(new(model.OpenProject))
	if err != nil {
		logger.Errorln("更新项目评论数失败：", err)
	}
}

func (self ProjectComment) String() string {
	return "project"
}

// 实现 CommentObjecter 接口
func (self ProjectComment) SetObjinfo(ids []int, commentMap map[int][]*model.Comment) {
	projects := DefaultProject.FindByIds(ids)
	if len(projects) == 0 {
		return
	}

	for _, project := range projects {
		objinfo := make(map[string]interface{})
		objinfo["title"] = project.Category + project.Name
		objinfo["uri"] = model.PathUrlMap[model.TypeProject]
		objinfo["type_name"] = model.TypeNameMap[model.TypeProject]

		for _, comment := range commentMap[project.Id] {
			comment.Objinfo = objinfo
		}
	}
}

// 项目喜欢
type ProjectLike struct{}

// 更新该项目的喜欢数
// objid：被喜欢对象id；num: 喜欢数(负数表示取消喜欢)
// func (self ProjectLike) UpdateLike(objid, num int) {
// 	// 更新喜欢数（TODO：暂时每次都更新表）
// 	err := model.NewOpenProject().Where("id=?", objid).Increment("likenum", num)
// 	if err != nil {
// 		logger.Errorln("更新项目喜欢数失败：", err)
// 	}
// }

// func (self ProjectLike) String() string {
// 	return "project"
// }
