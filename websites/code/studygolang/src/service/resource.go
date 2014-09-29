// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"logger"
	"model"
	"net/url"
	"strconv"
	"time"
	"util"
)

// 增加资源
func PublishResource(uid int, form url.Values) bool {
	resource := model.NewResource()
	err := util.ConvertAssign(resource, form)
	if err != nil {
		logger.Errorln("user ConvertAssign error", err)
		return false
	}
	resource.Ctime = time.Now().Format("2006-01-02 15:04:05")
	resource.Uid = uid
	id, err := resource.Insert()
	if err != nil {
		logger.Errorln("PublishResource error:", err)
		return false
	}

	// 发布资源，活跃度+10
	go IncUserWeight("uid="+strconv.Itoa(uid), 10)

	// 入扩展表
	resourceEx := model.NewResourceEx()
	resourceEx.Id = int(id)
	if _, err := resourceEx.Insert(); err != nil {
		logger.Errorln("PublishResource Ex error:", err)
		return false
	}

	return true
}

// 获得资源详细信息
func FindResource(id string) (resourceMap map[string]interface{}, comments []map[string]interface{}) {
	condition := "id=" + id
	resource := model.NewResource()
	err := resource.Where(condition).Find()
	if err != nil {
		logger.Errorln("resource service FindResource error:", err)
		return
	}
	resourceMap = make(map[string]interface{})
	util.Struct2Map(resourceMap, resource)
	resourceMap["catname"] = GetCategoryName(resource.Catid)
	// 链接的host
	if resource.Form == model.LinkForm {
		urlObj, err := url.Parse(resource.Url)
		if err == nil {
			resourceMap["host"] = urlObj.Host
		}
	} else {
		resourceMap["url"] = "/resources/" + strconv.Itoa(resource.Id)
	}
	resourceEx := model.NewResourceEx()
	err = resourceEx.Where(condition).Find()
	if err != nil {
		logger.Errorln("resource service FindResource Error:", err)
		return
	}
	util.Struct2Map(resourceMap, resourceEx)
	// 评论信息
	comments, ownerUser, _ := FindObjComments(id, strconv.Itoa(model.TYPE_RESOURCE), resource.Uid, 0)
	resourceMap["user"] = ownerUser
	return
}

// 通过id获得资源的所有者
func getResourceOwner(id int) int {
	resource := model.NewResource()
	err := resource.Where("id=" + strconv.Itoa(id)).Find()
	if err != nil {
		logger.Errorln("resource service getResourceOwner Error:", err)
		return 0
	}
	return resource.Uid
}

// 获得某个分类的资源列表
func FindResourcesByCatid(catid string) []map[string]interface{} {
	resourceList, err := model.NewResource().Where("catid=" + catid).Order("mtime DESC").FindAll()
	if err != nil {
		logger.Errorln("resource service FindResourcesByCatid error:", err)
		return nil
	}
	count := len(resourceList)
	ids := make([]int, count)
	uids := make(map[int]int)
	for i, resource := range resourceList {
		ids[i] = resource.Id
		uids[resource.Uid] = resource.Uid
	}

	// 获取扩展信息（计数）
	resourceExList, err := model.NewResourceEx().Where("id in(" + util.Join(ids, ",") + ")").FindAll()
	if err != nil {
		logger.Errorln("resource service FindResourcesByCatid Error:", err)
		return nil
	}
	resourceExMap := make(map[int]*model.ResourceEx, len(resourceExList))
	for _, resourceEx := range resourceExList {
		resourceExMap[resourceEx.Id] = resourceEx
	}

	userMap := getUserInfos(uids)

	resources := make([]map[string]interface{}, count)
	for i, resource := range resourceList {
		tmpMap := make(map[string]interface{})
		util.Struct2Map(tmpMap, resource)
		util.Struct2Map(tmpMap, resourceExMap[resource.Id])
		tmpMap["user"] = userMap[resource.Uid]
		// 链接的host
		if resource.Form == model.LinkForm {
			urlObj, err := url.Parse(resource.Url)
			if err == nil {
				tmpMap["host"] = urlObj.Host
			}
		} else {
			tmpMap["url"] = "/resources/" + strconv.Itoa(resource.Id)
		}
		resources[i] = tmpMap
	}
	return resources
}

// 获得最新资源
func FindRecentResources() []map[string]interface{} {
	resourceList, err := model.NewResource().Limit("0,10").Order("mtime DESC").FindAll()
	if err != nil {
		logger.Errorln("resource service FindRecentResources error:", err)
		return nil
	}
	count := len(resourceList)
	uids := make(map[int]int)
	for _, resource := range resourceList {
		uids[resource.Uid] = resource.Uid
	}
	userMap := getUserInfos(uids)
	resources := make([]map[string]interface{}, count)
	for i, resource := range resourceList {
		tmpMap := make(map[string]interface{})
		util.Struct2Map(tmpMap, resource)
		tmpMap["user"] = userMap[resource.Uid]
		resources[i] = tmpMap
	}
	return resources
}

// 获取多个资源详细信息
func FindResourcesByIds(ids []int) []*model.Resource {
	if len(ids) == 0 {
		return nil
	}
	inIds := util.Join(ids, ",")
	resources, err := model.NewResource().Where("id in(" + inIds + ")").FindAll()
	if err != nil {
		logger.Errorln("resource service FindResourcesByIds error:", err)
		return nil
	}
	return resources
}

// 提供给其他service调用（包内）
func getResources(ids map[int]int) map[int]*model.Resource {
	resources := FindResourcesByIds(util.MapIntKeys(ids))
	resourceMap := make(map[int]*model.Resource, len(resources))
	for _, resource := range resources {
		resourceMap[resource.Id] = resource
	}
	return resourceMap
}

// 资源总数
func ResourcesTotal() (total int) {
	total, err := model.NewResource().Count()
	if err != nil {
		logger.Errorln("resource service ResourcesTotal error:", err)
	}
	return
}

// 资源评论
type ResourceComment struct{}

// 更新该资源的评论信息
// cid：评论id；objid：被评论对象id；uid：评论者；cmttime：评论时间
func (self ResourceComment) UpdateComment(cid, objid, uid int, cmttime string) {
	id := strconv.Itoa(objid)
	// 更新评论数（TODO：暂时每次都更新表）
	err := model.NewResourceEx().Where("id="+id).Increment("cmtnum", 1)
	if err != nil {
		logger.Errorln("更新资源评论数失败：", err)
	}
}
