// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"fmt"
	"html/template"
	"logger"
	"model"
	"net/url"
	"regexp"
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

	uids := util.Models2Intslice(commentList, "Uid")

	// 避免某些情况下最后回复人没在回复列表中
	uids = append(uids, owner, lastCommentUid)

	// 获得用户信息
	userMap := GetUserInfos(uids)
	ownerUser = userMap[owner]
	if lastCommentUid != 0 {
		lastReplyUser = userMap[lastCommentUid]
	}
	comments = make([]map[string]interface{}, 0, len(commentList))
	for _, comment := range commentList {
		tmpMap := make(map[string]interface{})
		util.Struct2Map(tmpMap, comment)
		tmpMap["content"] = template.HTML(decodeCmtContent(comment))
		tmpMap["user"] = userMap[comment.Uid]
		comments = append(comments, tmpMap)
	}
	return
}

// 获得某个对象的所有评论（新版）
// TODO:分页暂不做
func FindObjectComments(objid, objtype string) (commentList []*model.Comment, err error) {
	commentList, err = model.NewComment().Where("objid=" + objid + " and objtype=" + objtype).FindAll()
	if err != nil {
		logger.Errorln("comment service FindObjectComments Error:", err)
	}

	for _, comment := range commentList {
		decodeCmtContent(comment)
	}

	return
}

func decodeCmtContent(comment *model.Comment) string {
	// 安全过滤
	content := template.HTMLEscapeString(comment.Content)
	// @别人
	content = parseAtUser(content)

	// 回复某一楼层
	reg := regexp.MustCompile(`#(\d+)楼`)
	url := fmt.Sprintf("%s%d#comment", model.PathUrlMap[comment.Objtype], comment.Objid)
	content = reg.ReplaceAllString(content, `<a href="`+url+`$1" title="$1">#$1<span>楼</span></a>`)

	comment.Content = content

	return content
}

// 获得最近的评论
// 如果 uid!=0，表示获取某人的评论；
// 如果 objtype!=-1，表示获取某类型的评论；
func FindRecentComments(uid, objtype int, limit string) []*model.Comment {
	cond := ""
	if uid != 0 {
		cond = "uid=" + strconv.Itoa(uid)
	}
	if objtype != -1 {
		if cond != "" {
			cond += " AND "
		}
		cond += "objtype=" + strconv.Itoa(objtype)
	}

	comments, err := model.NewComment().Where(cond).Order("cid DESC").Limit(limit).FindAll()
	if err != nil {
		logger.Errorln("comment service FindRecentComments error:", err)
		return nil
	}

	cmtMap := make(map[int][]*model.Comment, len(model.PathUrlMap))
	for _, comment := range comments {
		decodeCmtContent(comment)

		if _, ok := cmtMap[comment.Objtype]; !ok {
			cmtMap[comment.Objtype] = make([]*model.Comment, 0, 10)
		}

		cmtMap[comment.Objtype] = append(cmtMap[comment.Objtype], comment)
	}

	cmtObjs := []CommentObjecter{
		model.TYPE_TOPIC:    TopicComment{},
		model.TYPE_ARTICLE:  ArticleComment{},
		model.TYPE_RESOURCE: ResourceComment{},
		model.TYPE_WIKI:     nil,
		model.TYPE_PROJECT:  ProjectComment{},
	}
	for cmtType, cmts := range cmtMap {
		FillCommentObjs(cmts, cmtObjs[cmtType])
	}

	return comments
}

// 评论总数(objtype != -1 时，取某一类型的评论总数)
func CommentsTotal(objtype int) (total int) {
	var cond string
	if objtype != -1 {
		cond = "objtype=" + strconv.Itoa(objtype)
	}

	total, err := model.NewComment().Where(cond).Count()
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

// 填充评论对应的主体信息
func FillCommentObjs(comments []*model.Comment, cmtObj CommentObjecter) {
	if len(comments) == 0 {
		return
	}
	count := len(comments)
	commentMap := make(map[int][]*model.Comment, count)
	idMap := make(map[int]int, count)
	for _, comment := range comments {
		if _, ok := commentMap[comment.Objid]; !ok {
			commentMap[comment.Objid] = make([]*model.Comment, 0, count)
		}
		commentMap[comment.Objid] = append(commentMap[comment.Objid], comment)
		idMap[comment.Objid] = 1
	}
	ids := util.MapIntKeys(idMap)
	cmtObj.SetObjinfo(ids, commentMap)
}

// 填充 Comment 对象的 Objinfo 成员接口
// 评论属主应该实现该接口（以便填充 Objinfo 成员）
type CommentObjecter interface {
	// ids 是属主的主键 slice （comment 中的 objid）
	// commentMap 中的 key 是属主 id
	SetObjinfo(ids []int, commentMap map[int][]*model.Comment)
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

var commenters = make(map[int]Commenter)

// 评论接口
type Commenter interface {
	fmt.Stringer
	// 评论回调接口，用于更新对象自身需要更新的数据
	UpdateComment(int, int, int, string)
}

// 注册评论对象，使得某种类型（帖子、博客等）被评论了可以回调
func RegisterCommentObject(objtype int, commenter Commenter) {
	if commenter == nil {
		panic("service: Register commenter is nil")
	}
	if _, dup := commenters[objtype]; dup {
		panic("service: Register called twice for commenter " + commenter.String())
	}
	commenters[objtype] = commenter
}

// 发表评论（或回复）。
// objid 注册的评论对象
// uid 评论人
func PostComment(uid, objid int, form url.Values) (*model.Comment, error) {
	comment := model.NewComment()
	comment.Objid = objid
	objtype := util.MustInt(form.Get("objtype"))
	comment.Objtype = objtype
	comment.Uid = uid
	comment.Content = form.Get("content")

	// TODO:评论楼层怎么处理，避免冲突？最后的楼层信息保存在内存中？

	// 暂时只是从数据库中取出最后的评论楼层
	stringBuilder := util.NewBuffer()
	stringBuilder.Append("objid=").AppendInt(objid).Append(" AND objtype=").AppendInt(objtype)
	tmpCmt, err := model.NewComment().Where(stringBuilder.String()).Order("ctime DESC").Find()
	if err != nil {
		logger.Errorln("post comment service error:", err)
		return nil, err
	} else {
		comment.Floor = tmpCmt.Floor + 1
	}
	// 入评论库
	cid, err := comment.Insert()
	if err != nil {
		logger.Errorln("post comment service error:", err)
		return nil, err
	}
	comment.Cid = cid
	comment.Ctime = util.TimeNow()
	decodeCmtContent(comment)

	// 回调，不关心处理结果（有些对象可能不需要回调）
	if commenter, ok := commenters[objtype]; ok {
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

	// @某人 发系统消息
	go SendSysMsgAtUids(form.Get("uid"), ext)
	go SendSysMsgAtUsernames(form.Get("usernames"), ext)

	return comment, nil
}

func ModifyComment(cid, content string) (errMsg string, err error) {
	err = model.NewComment().Set("content=?", content).Where("cid=" + cid).Update()
	if err != nil {
		logger.Errorf("更新评论内容 【%s】 失败：%s\n", cid, err)
		errMsg = "对不起，服务器内部错误，请稍后再试！"
		return
	}

	return
}
