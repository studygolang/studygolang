// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"logger"
	"model"
	"strconv"
	"time"
	"util"
)

// 获得某个对象的所有评论
// owner: 被评论对象属主
// TODO:分页暂不做
func FindObjComments(objid, objtype string, owner, lastCommentUid int /*, page, pageNum int*/) (comments []map[string]interface{}, ownerUser, lastReplyUser *model.User) {
	commentList, err := model.NewComment().Where("objid=" + objid + " and objtype=" + objtype).FindAll()
	if err != nil {
		logger.Errorln("comment service FindObjComments Error:", err)
		return
	}

	commentNum := len(commentList)
	uids := make(map[int]int, commentNum+1)
	uids[owner] = owner
	// 避免某些情况下最后回复人没在回复列表中
	uids[lastCommentUid] = lastCommentUid
	for _, comment := range commentList {
		uids[comment.Uid] = comment.Uid
	}

	// 获得用户信息
	userMap := getUserInfos(uids)
	ownerUser = userMap[owner]
	if lastCommentUid != 0 {
		lastReplyUser = userMap[lastCommentUid]
	}
	comments = make([]map[string]interface{}, 0, commentNum)
	for _, comment := range commentList {
		tmpMap := make(map[string]interface{})
		util.Struct2Map(tmpMap, comment)
		tmpMap["user"] = userMap[comment.Uid]
		comments = append(comments, tmpMap)
	}
	return
}

// 获得某人在某种类型最近的评论
func FindRecentComments(uid, objtype int) []*model.Comment {
	comments, err := model.NewComment().Where("uid=" + strconv.Itoa(uid) + " AND objtype=" + strconv.Itoa(objtype)).Order("ctime DESC").Limit("0, 5").FindAll()
	if err != nil {
		logger.Errorln("comment service FindRecentComments error:", err)
		return nil
	}
	return comments
}

// 某类型的评论总数
func CommentsTotal(objtype int) (total int) {
	total, err := model.NewComment().Where("objtype=" + strconv.Itoa(objtype)).Count()
	if err != nil {
		logger.Errorln("comment service CommentsTotal error:", err)
		return
	}
	return
}

// 获取多个评论信息
func FindCommentsByIds(cids []int) []*model.Comment {
	if len(cids) == 0 {
		return nil
	}
	inCids := util.Join(cids, ",")
	comments, err := model.NewComment().Where("cid in(" + inCids + ")").FindAll()
	if err != nil {
		logger.Errorln("comment service FindCommentsByIds error:", err)
		return nil
	}
	return comments
}

// 提供给其他service调用（包内）
func getComments(cids map[int]int) map[int]*model.Comment {
	comments := FindCommentsByIds(util.MapIntKeys(cids))
	commentMap := make(map[int]*model.Comment, len(comments))
	for _, comment := range comments {
		commentMap[comment.Cid] = comment
	}
	return commentMap
}

var commenters = make(map[string]Commenter)

// 评论接口
type Commenter interface {
	// 评论回调接口，用于更新对象自身需要更新的数据
	UpdateComment(int, int, int, string)
}

// 注册评论对象，使得某种类型（帖子、博客等）可以被评论
func RegisterCommentObject(objname string, commenter Commenter) {
	if commenter == nil {
		panic("service: Register commenter is nil")
	}
	if _, dup := commenters[objname]; dup {
		panic("service: Register called twice for commenter " + objname)
	}
	commenters[objname] = commenter
}

// 发表评论（或回复）。
// objname 注册的评论对象名
func PostComment(objid, objtype, uid int, content string, objname string) error {
	comment := model.NewComment()
	comment.Objid = objid
	comment.Objtype = objtype
	comment.Uid = uid
	comment.Content = content

	// TODO:评论楼层怎么处理，避免冲突？最后的楼层信息保存在内存中？

	// 暂时只是从数据库中取出最后的评论楼层
	stringBuilder := util.NewBuffer()
	stringBuilder.Append("objid=").AppendInt(objid).Append(" AND objtype=").AppendInt(objtype)
	tmpCmt, err := model.NewComment().Where(stringBuilder.String()).Order("ctime DESC").Find()
	if err != nil {
		logger.Errorln("post comment service error:", err)
		return err
	} else {
		comment.Floor = tmpCmt.Floor + 1
	}
	// 入评论库
	cid, err := comment.Insert()
	if err != nil {
		logger.Errorln("post comment service error:", err)
		return err
	}
	// 回调，不关心处理结果（有些对象可能不需要回调）
	if commenter, ok := commenters[objname]; ok {
		logger.Debugf("评论[objid:%d] [objtype:%d] [uid:%d] 成功，通知被评论者更新", objid, objtype, uid)
		go commenter.UpdateComment(cid, objid, uid, time.Now().Format("2006-01-02 15:04:05"))
	}

	// 发评论，活跃度+5
	go IncUserWeight("uid="+strconv.Itoa(uid), 5)

	// 给被评论对象所有者发系统消息
	ext := map[string]interface{}{
		"objid":   objid,
		"objtype": objtype,
		"cid":     cid,
		"uid":     uid,
	}
	go SendSystemMsgTo(0, objtype, ext)

	// TODO: @某人 发系统消息？

	return nil
}
