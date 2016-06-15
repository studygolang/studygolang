// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"model"
	"net/url"
	"strconv"
	"time"

	. "db"

	"github.com/fatih/structs"
	"github.com/polaris1119/logger"
	"github.com/polaris1119/set"
	"golang.org/x/net/context"
)

type ResourceLogic struct{}

var DefaultResource = ResourceLogic{}

// Publish 增加（修改）资源
func (ResourceLogic) Publish(ctx context.Context, me *model.Me, form url.Values) (err error) {
	objLog := GetLogger(ctx)

	uid := me.Uid
	resource := &model.Resource{}

	if form.Get("id") != "" {
		id := form.Get("id")
		_, err = MasterDB.Id(id).Get(resource)
		if err != nil {
			logger.Errorln("ResourceLogic Publish find error:", err)
			return
		}

		if resource.Uid != uid && !me.IsAdmin {
			err = NotModifyAuthorityErr
			return
		}

		fields := []string{"title", "catid", "form", "url", "content"}
		if form.Get("form") == model.LinkForm {
			form.Set("content", "")
		} else {
			form.Set("url", "")
		}

		for _, field := range fields {
			form.Del(field)
		}

		err = schemaDecoder.Decode(resource, form)
		if err != nil {
			objLog.Errorln("ResourceLogic Publish decode error:", err)
			return
		}
		_, err = MasterDB.Id(id).Update(resource)
		if err != nil {
			objLog.Errorf("更新資源 【%s】 信息失败：%s\n", id, err)
			return
		}

		// 修改資源，活跃度+2
		go DefaultUser.IncrUserWeight("uid", uid, 2)
	} else {

		err = schemaDecoder.Decode(resource, form)
		if err != nil {
			objLog.Errorln("ResourceLogic Publish decode error:", err)
			return
		}

		resource.Uid = uid

		session := MasterDB.NewSession()
		defer session.Close()

		err = session.Begin()
		if err != nil {
			objLog.Errorln("Publish Resource begin tx error:", err)
			return
		}

		_, err = session.Insert(resource)
		if err != nil {
			session.Rollback()
			objLog.Errorln("Publish Resource insert resource error:", err)
			return
		}

		resourceEx := &model.ResourceEx{
			Id: resource.Id,
		}
		_, err = session.Insert(resourceEx)
		if err != nil {
			session.Rollback()
			objLog.Errorln("Publish Resource insert resource_ex error:", err)
			return
		}

		err = session.Commit()
		if err != nil {
			objLog.Errorln("Publish Resource commit error:", err)
		}

		// 给 被@用户 发系统消息
		ext := map[string]interface{}{
			"objid":   resource.Id,
			"objtype": model.TypeResource,
			"uid":     uid,
			"msgtype": model.MsgtypePublishAtMe,
		}
		go DefaultMessage.SendSysMsgAtUsernames(ctx, form.Get("usernames"), ext, 0)

		// 发布主题，活跃度+10
		go DefaultUser.IncrUserWeight("uid", uid, 10)
	}

	return
}

// Total 资源总数
func (ResourceLogic) Total() int64 {
	total, err := MasterDB.Count(new(model.Resource))
	if err != nil {
		logger.Errorln("CommentLogic Total error:", err)
	}
	return total
}

// FindBy 获取资源列表（分页）
func (ResourceLogic) FindBy(ctx context.Context, limit int, lastIds ...int) []*model.Resource {
	objLog := GetLogger(ctx)

	dbSession := MasterDB.OrderBy("id DESC").Limit(limit)
	if len(lastIds) > 0 {
		dbSession.Where("id>?", lastIds[0])
	}

	resourceList := make([]*model.Resource, 0)
	err := dbSession.Find(&resourceList)
	if err != nil {
		objLog.Errorln("ResourceLogic FindBy Error:", err)
		return nil
	}

	return resourceList
}

// FindByCatid 获得某个分类的资源列表，分页
func (ResourceLogic) FindByCatid(ctx context.Context, paginator *Paginator, catid int) (resources []map[string]interface{}, total int64) {
	objLog := GetLogger(ctx)

	var (
		count         = paginator.PerPage()
		resourceInfos = make([]*model.ResourceInfo, 0)
	)

	err := MasterDB.Join("INNER", "resource_ex", "resource.id=resource_ex.id").Where("catid=?", catid).
		Desc("resource.mtime").Limit(count, paginator.Offset()).Find(&resourceInfos)
	if err != nil {
		objLog.Errorln("ResourceLogic FindByCatid error:", err)
		return
	}

	total, err = MasterDB.Where("catid=?", catid).Count(new(model.Resource))
	if err != nil {
		objLog.Errorln("ResourceLogic FindByCatid count error:", err)
		return
	}

	uidSet := set.New(set.NonThreadSafe)
	for _, resourceInfo := range resourceInfos {
		uidSet.Add(resourceInfo.Uid)
	}

	usersMap := DefaultUser.FindUserInfos(ctx, set.IntSlice(uidSet))

	resources = make([]map[string]interface{}, len(resourceInfos))

	for i, resourceInfo := range resourceInfos {
		dest := make(map[string]interface{})

		structs.FillMap(resourceInfo.Resource, dest)
		structs.FillMap(resourceInfo.ResourceEx, dest)

		dest["user"] = usersMap[resourceInfo.Uid]

		// 链接的host
		if resourceInfo.Form == model.LinkForm {
			urlObj, err := url.Parse(resourceInfo.Url)
			if err == nil {
				dest["host"] = urlObj.Host
			}
		} else {
			dest["url"] = "/resources/" + strconv.Itoa(resourceInfo.Resource.Id)
		}

		resources[i] = dest
	}

	return
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

// findByIds 获取多个资源详细信息 包内使用
func (ResourceLogic) findByIds(ids []int) map[int]*model.Resource {
	if len(ids) == 0 {
		return nil
	}
	resources := make(map[int]*model.Resource)
	err := MasterDB.In("id", ids).Find(&resources)
	if err != nil {
		logger.Errorln("ResourceLogic FindByIds error:", err)
		return nil
	}
	return resources
}

// 获得资源详细信息
func (ResourceLogic) FindById(ctx context.Context, id int) (resourceMap map[string]interface{}, comments []map[string]interface{}) {
	objLog := GetLogger(ctx)

	resourceInfo := &model.ResourceInfo{}
	_, err := MasterDB.Join("INNER", "resource_ex", "resource.id=resource_ex.id").Where("resource.id=?", id).Get(resourceInfo)
	if err != nil {
		objLog.Errorln("ResourceLogic FindById error:", err)
		return
	}

	resource := &resourceInfo.Resource
	if resource.Id == 0 {
		objLog.Errorln("ResourceLogic FindById get error:", err)
		return
	}

	resourceMap = make(map[string]interface{})
	structs.FillMap(resource, resourceMap)
	structs.FillMap(resourceInfo.ResourceEx, resourceMap)

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

	// 评论信息
	comments, ownerUser, _ := DefaultComment.FindObjComments(ctx, id, model.TypeResource, resource.Uid, 0)
	resourceMap["user"] = ownerUser
	return
}

// 获取单个 Resource 信息（用于编辑）
func (ResourceLogic) FindResource(ctx context.Context, id int) *model.Resource {
	objLog := GetLogger(ctx)

	resource := &model.Resource{}
	_, err := MasterDB.Id(id).Get(resource)
	if err != nil {
		objLog.Errorf("ResourceLogic FindResource [%d] error：%s\n", id, err)
	}

	return resource
}

// 获得某个用户最近的资源
func (ResourceLogic) FindRecent(ctx context.Context, uid int) []*model.Resource {
	resourceList := make([]*model.Resource, 0)
	err := MasterDB.Where("uid=?", uid).Limit(5).OrderBy("id DESC").Find(&resourceList)
	if err != nil {
		logger.Errorln("resource logic FindRecent error:", err)
		return nil
	}

	return resourceList
}

// getOwner 通过id获得资源的所有者
func (ResourceLogic) getOwner(id int) int {
	resource := &model.Resource{}
	_, err := MasterDB.Id(id).Get(resource)
	if err != nil {
		logger.Errorln("resource logic getOwner Error:", err)
		return 0
	}
	return resource.Uid
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

// 资源喜欢
type ResourceLike struct{}

// 更新该主题的喜欢数
// objid：被喜欢对象id；num: 喜欢数(负数表示取消喜欢)
func (self ResourceLike) UpdateLike(objid, num int) {
	// 更新喜欢数（TODO：暂时每次都更新表）
	_, err := MasterDB.Where("id=?", objid).Incr("likenum", num).Update(new(model.ResourceEx))
	if err != nil {
		logger.Errorln("更新资源喜欢数失败：", err)
	}
}

func (self ResourceLike) String() string {
	return "resource"
}
