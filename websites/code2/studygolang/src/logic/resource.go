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
)

type ResourceLogic struct{}

var DefaultResource = ResourceLogic{}

// Total 资源总数
func (ResourceLogic) Total() int64 {
	total, err := MasterDB.Count(new(model.Resource))
	if err != nil {
		logger.Errorln("CommentLogic Total error:", err)
	}
	return total
}

// FindBy 获取资源列表（分页）
func (ResourceLogic) FindBy(limit int, lastIds ...int) []*model.Resource {
	dbSession := MasterDB.OrderBy("id DESC").Limit(limit)
	if len(lastIds) > 0 {
		dbSession.Where("id>?", lastIds[0])
	}

	resourceList := make([]*model.Resource, 0)
	err := dbSession.Find(&resourceList)
	if err != nil {
		logger.Errorln("ResourceLogic FindBy Error:", err)
		return nil
	}

	return resourceList
}

// FindByIds 获取多个资源详细信息
func (ResourceLogic) FindByIds(ids []int) []*model.Resource {
	if len(ids) == 0 {
		return nil
	}
	resources := make([]*model.Resource, 0)
	err := MasterDB.In("id", ids).Find(&resources)
	if err != nil {
		logger.Errorln("ResourceLogic FindByIds error:", err)
		return nil
	}
	return resources
}

// 资源评论
type ResourceComment struct{}

// 更新该资源的评论信息
// cid：评论id；objid：被评论对象id；uid：评论者；cmttime：评论时间
func (self ResourceComment) UpdateComment(cid, objid, uid int, cmttime time.Time) {
	// 更新评论数（TODO：暂时每次都更新表）
	_, err := MasterDB.Id(objid).Incr("cmtnum", 1).Update(new(model.ResourceEx))
	if err != nil {
		logger.Errorln("更新资源评论数失败：", err)
	}
}

func (self ResourceComment) String() string {
	return "resource"
}

// 实现 CommentObjecter 接口
func (self ResourceComment) SetObjinfo(ids []int, commentMap map[int][]*model.Comment) {
	resources := DefaultResource.FindByIds(ids)
	if len(resources) == 0 {
		return
	}

	for _, resource := range resources {
		objinfo := make(map[string]interface{})
		objinfo["title"] = resource.Title
		objinfo["uri"] = model.PathUrlMap[model.TypeResource]
		objinfo["type_name"] = model.TypeNameMap[model.TypeResource]

		for _, comment := range commentMap[resource.Id] {
			comment.Objinfo = objinfo
		}
	}
}
