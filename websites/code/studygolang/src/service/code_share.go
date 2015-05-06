// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package service

import (
	"net/url"
	"strconv"

	"logger"
	"model"
	"util"
)

// 分享代码片段
func PublishCode(user map[string]interface{}, form url.Values) (err error) {
	username := user["username"].(string)
	form.Set("op_user", username)

	code := model.NewCodeShare()

	if form.Get("id") != "" {
		err = code.Where("id=?", form.Get("id")).Find()
		if err != nil {
			logger.Errorln("Publish Code Find error:", err)
			return
		}

		isAdmin := false
		if _, ok := user["isadmin"]; ok {
			isAdmin = user["isadmin"].(bool)
		}
		if code.OpUser != username && !isAdmin {
			err = NotModifyAuthorityErr
			return
		}

		_, err = ModifyCode(form, username)
		if err != nil {
			logger.Errorln("Publish Code error:", err)
			return
		}
	} else {

		util.ConvertAssign(code, form)

		code.OpUser = username

		var id int64
		id, err = code.Insert()

		if err != nil {
			logger.Errorln("Publish Code error:", err)
			return
		}

		// 给 被@用户 发系统消息
		ext := map[string]interface{}{
			"objid":   id,
			"objtype": model.TYPE_CODE,
			"uid":     user["uid"],
			"msgtype": model.MsgtypePublishAtMe,
		}
		go SendSysMsgAtUsernames(form.Get("usernames"), ext)

		// 分享代码，活跃度+8
		go IncUserWeight("username="+username, 8)
	}

	return
}

// 修改代码片段
func ModifyCode(form url.Values, username string) (errMsg string, err error) {
	fields := []string{"title", "remark", "code"}
	query, args := updateSetClause(form, fields)

	id := form.Get("id")

	err = model.NewCodeShare().Set(query, args...).Where("id=?", id).Update()
	if err != nil {
		logger.Errorf("更新代码片段 【%s】 信息失败：%s\n", id, err)
		errMsg = "对不起，服务器内部错误，请稍后再试！"
		return
	}

	// 修改代码片段，活跃度+2
	go IncUserWeight("username="+username, 2)

	return
}

// 获得代码片段列表
// page 当前第几页
func FindCodes(page int) (codes []*model.CodeShare, userMap map[string]*model.User, total int) {
	var offset = 0
	if page > 1 {
		offset = (page - 1) * PAGE_NUM
	}

	codeObj := model.NewCodeShare()
	limit := strconv.Itoa(offset) + "," + strconv.Itoa(PAGE_NUM)
	codes, err := codeObj.Order("id DESC").Limit(limit).FindAll()
	if err != nil {
		logger.Errorln("code share service FindCodes error:", err)
		return
	}

	// 获得该类别总资源数
	total, err = codeObj.Count()
	if err != nil {
		logger.Errorln("code share service codeObj.Count Error:", err)
		return
	}

	count := len(codes)
	usernames := make([]string, count)
	for i, code := range codes {
		usernames[i] = code.OpUser
	}

	userMap = GetUserByUsernames(usernames)

	return
}

// 获取多个代码片段详细信息
func FindCodesByIds(ids []int) []*model.CodeShare {
	if len(ids) == 0 {
		return nil
	}
	inIds := util.Join(ids, ",")
	codes, err := model.NewCodeShare().Where("id in(" + inIds + ")").FindAll()
	if err != nil {
		logger.Errorln("code share service FindCodesByIds error:", err)
		return nil
	}
	return codes
}

// 代码片段评论
type CodeComment struct{}

// 更新该文章的评论信息
// cid：评论id；objid：被评论对象id；uid：评论者；cmttime：评论时间
func (self CodeComment) UpdateComment(cid, objid, uid int, cmttime string) {
	id := strconv.Itoa(objid)

	// 更新评论数（TODO：暂时每次都更新表）
	err := model.NewCodeShare().Where("id=?", id).Increment("cmtnum", 1)
	if err != nil {
		logger.Errorln("更新代码片段评论数失败：", err)
	}
}

func (self CodeComment) String() string {
	return "code_share"
}

// 实现 CommentObjecter 接口
func (self CodeComment) SetObjinfo(ids []int, commentMap map[int][]*model.Comment) {
	codes := FindCodesByIds(ids)
	if len(codes) == 0 {
		return
	}

	for _, code := range codes {
		objinfo := make(map[string]interface{})
		objinfo["title"] = code.Title
		objinfo["uri"] = model.PathUrlMap[model.TYPE_CODE]
		objinfo["type_name"] = model.TypeNameMap[model.TYPE_CODE]

		for _, comment := range commentMap[code.Id] {
			comment.Objinfo = objinfo
		}
	}
}

// 代码喜欢
type CodeLike struct{}

// 更新该文章的喜欢数
// objid：被喜欢对象id；num: 喜欢数(负数表示取消喜欢)
func (self CodeLike) UpdateLike(objid, num int) {
	// 更新喜欢数（TODO：暂时每次都更新表）
	err := model.NewCodeShare().Where("id=?", objid).Increment("likenum", num)
	if err != nil {
		logger.Errorln("更新代码片段喜欢数失败：", err)
	}
}

func (self CodeLike) String() string {
	return "code_share"
}
