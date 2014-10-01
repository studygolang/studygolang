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
		tmpMap["content"] = template.HTML(decodeCmtContent(comment))
		tmpMap["user"] = userMap[comment.Uid]
		comments = append(comments, tmpMap)
	}
	return
}

func decodeCmtContent(comment *model.Comment) string {
	// 安全过滤
	content := template.HTMLEscapeString(comment.Content)
	// @别人
	reg := regexp.MustCompile(`@([^\s@]{4,20})`)
	content = reg.ReplaceAllString(content, `<a href="/user/$1" title="@$1">@$1</a>`)

	// 回复某一楼层
	reg = regexp.MustCompile(`#(\d+)楼`)
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
		cond = "objtype=" + strconv.Itoa(objtype)
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

	FillCommentTopics(cmtMap[model.TYPE_TOPIC])

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
// uid 评论人
func PostComment(uid, objid int, form url.Values) error {
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
	if commenter, ok := commenters[form.Get("objname")]; ok {
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

	return nil
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
