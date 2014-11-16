// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"net/url"
	"strconv"

	"logger"
	"model"
	"util"
)

// 增加（修改）资源
func PublishResource(user map[string]interface{}, form url.Values) (err error) {
	uid := user["uid"].(int)

	resource := model.NewResource()

	if form.Get("id") != "" {
		err = resource.Where("id=?", form.Get("id")).Find()
		if err != nil {
			logger.Errorln("Publish Resource find error:", err)
			return
		}

		isAdmin := false
		if _, ok := user["isadmin"]; ok {
			isAdmin = user["isadmin"].(bool)
		}
		if resource.Uid != uid && !isAdmin {
			err = NotModifyAuthorityErr
			return
		}

		fields := []string{"title", "catid", "form", "url", "content"}
		if form.Get("form") == model.LinkForm {
			form.Set("content", "")
		} else {
			form.Set("url", "")
		}

		id := form.Get("id")
		query, args := updateSetClause(form, fields)
		err = resource.Set(query, args...).Where("id=?", id).Update()
		if err != nil {
			logger.Errorf("更新資源 【%s】 信息失败：%s\n", id, err)
			return
		}

		// 修改資源，活跃度+2
		go IncUserWeight("uid="+strconv.Itoa(uid), 2)
	} else {

		util.ConvertAssign(resource, form)

		resource.Uid = uid
		resource.Ctime = util.TimeNow()

		var id int64
		id, err = resource.Insert()

		if err != nil {
			logger.Errorln("Publish Resource error:", err)
			return
		}

		// 存扩展信息
		resourceEx := model.NewResourceEx()
		resourceEx.Id = int(id)
		if _, err = resourceEx.Insert(); err != nil {
			logger.Errorln("PublishResource Ex error:", err)
			return
		}

		// 给 被@用户 发系统消息
		/*
			ext := map[string]interface{}{
				"objid":   id,
				"objtype": model.TYPE_RESOURCE,
				"uid":     user["uid"],
				"msgtype": model.MsgtypePublishAtMe,
			}
			go SendSysMsgAtUsernames(form.Get("usernames"), ext)
		*/

		// 发布主题，活跃度+10
		go IncUserWeight("uid="+strconv.Itoa(uid), 10)
	}

	return
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

// 获取单个 Resource 信息（用于编辑）
func FindResourceById(id string) *model.Resource {
	resource := model.NewResource()
	err := resource.Where("id=?", id).Find()
	if err != nil {
		logger.Errorf("FindResourceById [%s] error：%s\n", id, err)
	}

	return resource
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
// page 当前第几页
func FindResourcesByCatid(catid string, page int) (resources []map[string]interface{}, total int) {
	var offset = 0
	if page > 1 {
		offset = (page - 1) * PAGE_NUM
	}

	resourceObj := model.NewResource()
	limit := strconv.Itoa(offset) + "," + strconv.Itoa(PAGE_NUM)
	resourceList, err := resourceObj.Where("catid=?", catid).Order("mtime DESC").Limit(limit).FindAll()
	if err != nil {
		logger.Errorln("resource service FindResourcesByCatid error:", err)
		return
	}

	// 获得该类别总资源数
	total, err = resourceObj.Count()
	if err != nil {
		logger.Errorln("resource service resourceObj.Count Error:", err)
		return
	}

	count := len(resourceList)
	ids := make([]int, count)
	uids := make([]int, count)
	for i, resource := range resourceList {
		ids[i] = resource.Id
		uids[i] = resource.Uid
	}

	// 获取扩展信息（计数）
	resourceExList, err := model.NewResourceEx().Where("id in(" + util.Join(ids, ",") + ")").FindAll()
	if err != nil {
		logger.Errorln("resource service FindResourcesByCatid Error:", err)
		return
	}
	resourceExMap := make(map[int]*model.ResourceEx, len(resourceExList))
	for _, resourceEx := range resourceExList {
		resourceExMap[resourceEx.Id] = resourceEx
	}

	userMap := GetUserInfos(uids)

	resources = make([]map[string]interface{}, count)
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
	return
}

// 获得最新资源
func FindRecentResources() []map[string]interface{} {
	resourceList, err := model.NewResource().Limit("0,10").Order("mtime DESC").FindAll()
	if err != nil {
		logger.Errorln("resource service FindRecentResources error:", err)
		return nil
	}

	uids := util.Models2Intslice(resourceList, "Uid")
	userMap := GetUserInfos(uids)

	count := len(resourceList)
	resources := make([]map[string]interface{}, count)
	for i, resource := range resourceList {
		tmpMap := make(map[string]interface{})
		util.Struct2Map(tmpMap, resource)
		tmpMap["user"] = userMap[resource.Uid]
		resources[i] = tmpMap
	}
	return resources
}

// 获得某个用户最近的资源
func FindUserRecentResources(uid int) []*model.Resource {
	resourceList, err := model.NewResource().Where("uid=?", uid).Limit("0,5").Order("mtime DESC").FindAll()
	if err != nil {
		logger.Errorln("resource service FindUserRecentResources error:", err)
		return nil
	}

	return resourceList
}

// 获取资源列表（分页）
func FindResources(lastId, limit string) []*model.Resource {
	resource := model.NewResource()

	resourceList, err := resource.Where("id>" + lastId).Order("id DESC").Limit(limit).
		FindAll()
	if err != nil {
		logger.Errorln("resource service FindResources Error:", err)
		return nil
	}

	return resourceList
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

func (self ResourceComment) String() string {
	return "resource"
}

// 实现 CommentObjecter 接口
func (self ResourceComment) SetObjinfo(ids []int, commentMap map[int][]*model.Comment) {
	resources := FindResourcesByIds(ids)
	if len(resources) == 0 {
		return
	}

	for _, resource := range resources {
		objinfo := make(map[string]interface{})
		objinfo["title"] = resource.Title
		objinfo["uri"] = model.PathUrlMap[model.TYPE_RESOURCE]
		objinfo["type_name"] = model.TypeNameMap[model.TYPE_RESOURCE]

		for _, comment := range commentMap[resource.Id] {
			comment.Objinfo = objinfo
		}
	}
}

// 资源喜欢
type ResourceLike struct{}

// 更新该主题的喜欢数
// objid：被喜欢对象id；num: 喜欢数(负数表示取消喜欢)
func (self ResourceLike) UpdateLike(objid, num int) {
	// 更新喜欢数（TODO：暂时每次都更新表）
	err := model.NewResourceEx().Where("id=?", objid).Increment("likenum", num)
	if err != nil {
		logger.Errorln("更新资源喜欢数失败：", err)
	}
}

func (self ResourceLike) String() string {
	return "resource"
}
