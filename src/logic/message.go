// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"html/template"
	"model"
	"strconv"
	"strings"
	"util"

	. "db"

	"github.com/go-xorm/xorm"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
	"github.com/polaris1119/set"
	"golang.org/x/net/context"
)

type MessageLogic struct{}

var DefaultMessage = MessageLogic{}

// SendMessageTo from给to发短信息
func (MessageLogic) SendMessageTo(ctx context.Context, from, to int, content string) bool {
	objLog := GetLogger(ctx)

	message := &model.Message{
		From:    from,
		Fdel:    model.FdelNotDel,
		To:      to,
		Tdel:    model.TdelNotDel,
		Content: content,
		Hasread: model.NotRead,
	}
	if _, err := MasterDB.Insert(message); err != nil {
		objLog.Errorln("message logic SendMessageTo Error:", err)
		return false
	}

	// 通过 WebSocket 通知对方
	msg := NewMessage(WsMsgNotify, 1)
	go Book.PostMessage(to, msg)
	return true
}

// SendSystemMsgTo 给某人发系统消息
func (MessageLogic) SendSystemMsgTo(ctx context.Context, to, msgtype int, ext map[string]interface{}) bool {
	if to == 0 {
		return true
	}

	if from, ok := ext["uid"]; ok {
		// 自己的动作不发系统消息
		if to == from.(int) {
			return true
		}
	}

	message := &model.SystemMessage{
		To:      to,
		Msgtype: msgtype,
		Hasread: model.NotRead,
	}
	message.SetExt(ext)
	if _, err := MasterDB.Insert(message); err != nil {
		logger.Errorln("message logic SendSystemMsgTo Error:", err)
		return false
	}
	// 通过 WebSocket 通知对方
	msg := NewMessage(WsMsgNotify, 1)
	go Book.PostMessage(to, msg)
	return true
}

// SendSysMsgAtUids 给被@的用户发系统消息
// authors 是被评论对象的作者
func (MessageLogic) SendSysMsgAtUids(ctx context.Context, uids string, ext map[string]interface{}, author int) bool {
	if uids == "" {
		return true
	}

	message := &model.SystemMessage{
		Msgtype: model.MsgtypeAtMe,
		Hasread: model.NotRead,
	}
	message.SetExt(ext)

	msg := NewMessage(WsMsgNotify, 1)

	uidSlice := strings.Split(uids, ",")
	for _, uidStr := range uidSlice {
		uid := goutils.MustInt(strings.TrimSpace(uidStr))

		// 评论时 @ 作者了，不发通知，因为给作者已经发过一次了
		if uid == author {
			continue
		}

		if from, ok := ext["uid"]; ok {
			// 自己的动作不发系统消息
			if uid == from.(int) {
				continue
			}
		}
		message.To = uid
		if _, err := MasterDB.Insert(message); err != nil {
			logger.Errorln("message logic SendSysMsgAtUids Error:", err)
			continue
		}
		// 通过 WebSocket 通知对方
		go Book.PostMessage(uid, msg)
	}
	return true
}

// SendSysMsgAtUsernames 给被@的用户发系统消息
// ext 中可以指定 msgtype，没有指定，默认为 MsgtypeAtMe
func (MessageLogic) SendSysMsgAtUsernames(ctx context.Context, usernames string, ext map[string]interface{}, author int) bool {
	if usernames == "" {
		return true
	}
	message := &model.SystemMessage{
		Hasread: model.NotRead,
	}
	if msgtype, ok := ext["msgtype"]; ok {
		message.Msgtype = msgtype.(int)
		delete(ext, "msgtype")
	} else {
		message.Msgtype = model.MsgtypeAtMe
	}
	message.SetExt(ext)

	msg := NewMessage(WsMsgNotify, 1)

	usernameSlice := strings.Split(usernames, ",")
	for _, username := range usernameSlice {
		user := DefaultUser.FindOne(ctx, "username", strings.TrimSpace(username))
		// @ 的用户不存在
		if user == nil {
			continue
		}

		uid := user.Uid

		// 评论时 @ 作者了，不发通知，因为给作者已经发过一次了
		if uid == author {
			continue
		}

		if from, ok := ext["uid"]; ok {
			// 自己的动作不发系统消息
			if uid == from.(int) {
				continue
			}
		}
		message.To = uid
		if _, err := MasterDB.Insert(message); err != nil {
			logger.Errorln("message logic SendSysMsgAtUsernames Error:", err)
			continue
		}
		// 通过 WebSocket 通知对方
		go Book.PostMessage(uid, msg)
	}
	return true
}

// FindSysMsgsByUid 获得某人的系统消息
// 系统消息类型不同，在ext中存放的字段也不一样，如下：
//   model.MsgtypeTopicReply/MsgtypeResourceComment/MsgtypeWikiComment存放都为：
//		{"uid":xxx,"objid":xxx}
//   model.MsgtypeAtMe 为：{"uid":xxx,"cid":xxx,"objid":xxx,"objtype":xxx}
//   model.MsgtypePulishAtMe 为：{"uid":xxx,"objid":xxx,"objtype":xxx}
func (self MessageLogic) FindSysMsgsByUid(ctx context.Context, uid int, paginator *Paginator) []map[string]interface{} {
	objLog := GetLogger(ctx)

	messages := make([]*model.SystemMessage, 0)
	err := MasterDB.Where("`to`=?", uid).OrderBy("id DESC").
		Limit(paginator.PerPage(), paginator.Offset()).Find(&messages)
	if err != nil {
		objLog.Errorln("message logic FindSysMsgsByUid Error:", err)
		return nil
	}

	tidSet := set.New(set.NonThreadSafe)
	articleIdSet := set.New(set.NonThreadSafe)
	resIdSet := set.New(set.NonThreadSafe)
	wikiIdSet := set.New(set.NonThreadSafe)
	pidSet := set.New(set.NonThreadSafe)
	bookIdSet := set.New(set.NonThreadSafe)
	// 评论ID
	cidSet := set.New(set.NonThreadSafe)
	uidSet := set.New(set.NonThreadSafe)
	// subject id
	sidSet := set.New(set.NonThreadSafe)

	ids := make([]int, 0, len(messages))
	for _, message := range messages {
		ext := message.GetExt()
		if val, ok := ext["uid"]; ok {
			uidSet.Add(int(val.(float64)))
		}
		var objid int
		if val, ok := ext["objid"]; ok {
			objid = int(val.(float64))
		}
		switch message.Msgtype {
		case model.MsgtypeTopicReply:
			tidSet.Add(objid)
		case model.MsgtypeArticleComment:
			articleIdSet.Add(objid)
		case model.MsgtypeResourceComment:
			resIdSet.Add(objid)
		case model.MsgtypeWikiComment:
			wikiIdSet.Add(objid)
		case model.MsgtypeProjectComment:
			pidSet.Add(objid)
		case model.MsgtypeAtMe, model.MsgtypePublishAtMe:
			objTypeFloat := ext["objtype"].(float64)
			switch int(objTypeFloat) {
			case model.TypeTopic:
				tidSet.Add(objid)
			case model.TypeArticle:
				articleIdSet.Add(objid)
			case model.TypeResource:
				resIdSet.Add(objid)
			case model.TypeWiki:
				wikiIdSet.Add(objid)
			case model.TypeProject:
				pidSet.Add(objid)
			case model.TypeBook:
				bookIdSet.Add(objid)
			}
		case model.MsgtypeSubjectContribute:
			articleIdSet.Add(objid)
			sidSet.Add(int(ext["sid"].(float64)))
		}
		if val, ok := ext["cid"]; ok {
			cidSet.Add(int(val.(float64)))
		}
		if message.Hasread == "未读" {
			ids = append(ids, message.Id)
		}
	}
	// 标记已读
	go self.MarkHasRead(ctx, ids, true, uid)

	userMap := DefaultUser.FindUserInfos(ctx, set.IntSlice(uidSet))
	commentMap := DefaultComment.findByIds(set.IntSlice(cidSet))
	topicMap := DefaultTopic.findByTids(set.IntSlice(tidSet))
	articleMap := DefaultArticle.findByIds(set.IntSlice(articleIdSet))
	resourceMap := DefaultResource.findByIds(set.IntSlice(resIdSet))
	wikiMap := DefaultWiki.findByIds(set.IntSlice(wikiIdSet))
	projectMap := DefaultProject.findByIds(set.IntSlice(pidSet))
	bookMap := DefaultGoBook.findByIds(set.IntSlice(bookIdSet))
	subjectMap := DefaultSubject.findByIds(set.IntSlice(sidSet))

	result := make([]map[string]interface{}, len(messages))
	for i, message := range messages {
		tmpMap := make(map[string]interface{})
		// 某条信息的提示（标题）
		title := ""
		ext := message.GetExt()
		if val, ok := ext["objid"]; ok {
			objTitle := ""
			objUrl := ""
			objid := int(val.(float64))
			switch message.Msgtype {
			case model.MsgtypeTopicReply:
				objTitle = topicMap[objid].Title
				objUrl = "/topics/" + strconv.Itoa(topicMap[objid].Tid)
				title = "回复了你的主题："
			case model.MsgtypeArticleComment:
				objTitle = articleMap[objid].Title
				objUrl = "/articles/" + strconv.Itoa(articleMap[objid].Id)
				title = "回复了你的文章："
			case model.MsgtypeResourceComment:
				objTitle = resourceMap[objid].Title
				objUrl = "/resources/" + strconv.Itoa(resourceMap[objid].Id)
				title = "评论了你的资源："
			case model.MsgtypeWikiComment:
				objTitle = wikiMap[objid].Title
				objUrl = "/wiki/" + strconv.Itoa(wikiMap[objid].Id)
				title = "评论了你的Wiki页："
			case model.MsgtypeProjectComment:
				project := projectMap[objid]
				objTitle = project.Category + project.Name
				objUrl = "/p/"
				if project.Uri != "" {
					objUrl += project.Uri
				} else {
					objUrl += strconv.Itoa(project.Id)
				}
				title = "评论了你的开源项目："
			case model.MsgtypeAtMe:
				title = "评论时提到了你，在"
				switch int(ext["objtype"].(float64)) {
				case model.TypeTopic:
					topic := topicMap[objid]
					objTitle = topic.Title
					objUrl = "/topics/" + strconv.Itoa(topic.Tid) + "#commentForm"
					title += "主题："
				case model.TypeArticle:
					article := articleMap[objid]
					objTitle = article.Title
					objUrl = "/articles/" + strconv.Itoa(article.Id) + "#commentForm"
					title += "文章："
				case model.TypeResource:
					resource := resourceMap[objid]
					objTitle = resource.Title
					objUrl = "/resources/" + strconv.Itoa(resource.Id) + "#commentForm"
					title += "资源："
				case model.TypeWiki:
					wiki := wikiMap[objid]
					objTitle = wiki.Title
					objUrl = "/wiki/" + strconv.Itoa(wiki.Id) + "#commentForm"
					title += "wiki："
				case model.TypeProject:
					project := projectMap[objid]
					objTitle = project.Category + project.Name
					objUrl = "/p/"
					if project.Uri != "" {
						objUrl += project.Uri
					} else {
						objUrl += strconv.Itoa(project.Id)
					}
					objUrl += "#commentForm"
					title += "项目："
				case model.TypeBook:
					book := bookMap[objid]
					objTitle = book.Name
					objUrl = "/book/" + strconv.Itoa(book.Id) + "#commentForm"
					title += "图书："
				}

			case model.MsgtypePublishAtMe:
				title = "发布"
				switch int(ext["objtype"].(float64)) {
				case model.TypeTopic:
					topic := topicMap[objid]
					objTitle = topic.Title
					objUrl = "/topics/" + strconv.Itoa(topic.Tid)
					title += "主题"
				case model.TypeArticle:
					article := articleMap[objid]
					objTitle = article.Title
					objUrl = "/articles/" + strconv.Itoa(article.Id)
					title += "文章"
				case model.TypeResource:
					resource := resourceMap[objid]
					objTitle = resource.Title
					objUrl = "/resources/" + strconv.Itoa(resource.Id)
					title += "资源"
				case model.TypeWiki:
					wiki := wikiMap[objid]
					objTitle = wiki.Title
					objUrl = "/wiki/" + strconv.Itoa(wiki.Id)
					title += "wiki"
				case model.TypeProject:
					project := projectMap[objid]
					objTitle = project.Category + project.Name
					objUrl = "/p/"
					if project.Uri != "" {
						objUrl += project.Uri
					} else {
						objUrl += strconv.Itoa(project.Id)
					}
					title += "项目"
				}

				title += "时提到了你："

			case model.MsgtypeSubjectContribute:
				subject := subjectMap[int(ext["sid"].(float64))]
				article := articleMap[objid]
				objTitle = article.Title
				objUrl = "/articles/" + strconv.Itoa(article.Id)
				title += "收录了新文章"
				tmpMap["sprefix"] = "的专栏"
				tmpMap["surl"] = "/subject/" + strconv.Itoa(subject.Id)
				tmpMap["stitle"] = subject.Name
			}
			tmpMap["objtitle"] = objTitle
			tmpMap["objurl"] = objUrl
			tmpMap["objid"] = objid
			tmpMap["objtype"] = ext["objtype"]
		}
		tmpMap["ctime"] = message.Ctime
		tmpMap["id"] = message.Id
		tmpMap["hasread"] = message.Hasread
		if val, ok := ext["uid"]; ok {
			tmpMap["user"] = userMap[int(val.(float64))]
		}
		// content 和 cid不会同时存在
		if val, ok := ext["content"]; ok {
			tmpMap["content"] = val.(string)
		} else if val, ok := ext["cid"]; ok {
			tmpMap["content"] = template.HTML(DefaultComment.decodeCmtContent(ctx, commentMap[int(val.(float64))]))
		}
		tmpMap["title"] = title
		result[i] = tmpMap
	}
	return result
}

func (MessageLogic) SysMsgCount(ctx context.Context, uid int) int64 {
	total, _ := MasterDB.Where("`to`=?", uid).Count(new(model.SystemMessage))
	return total
}

func (MessageLogic) FindMsgById(ctx context.Context, id string) *model.Message {
	if id == "" {
		return nil
	}

	objLog := GetLogger(ctx)
	message := &model.Message{}
	_, err := MasterDB.Id(id).Get(message)
	if err != nil {
		objLog.Errorln("message logic FindMsgById Error:", err)
		return nil
	}

	return message
}

// 获得发给某人的短消息（收件箱）
func (self MessageLogic) FindToMsgsByUid(ctx context.Context, uid int, paginator *Paginator) []map[string]interface{} {
	objLog := GetLogger(ctx)

	messages := make([]*model.Message, 0)
	err := MasterDB.Where("`to`=? AND tdel=?", uid, model.TdelNotDel).
		Limit(paginator.PerPage(), paginator.Offset()).OrderBy("id DESC").Find(&messages)
	if err != nil {
		objLog.Errorln("message logic FindToMsgsByUid Error:", err)
		return nil
	}

	uidSet := set.New(set.NonThreadSafe)
	idSet := set.New(set.NonThreadSafe)
	for _, message := range messages {
		uidSet.Add(message.From)
		if message.Hasread == model.NotRead {
			idSet.Add(message.Id)
		}
	}

	// 标记已读
	go self.MarkHasRead(ctx, set.IntSlice(idSet), false, uid)

	userMap := DefaultUser.FindUserInfos(ctx, set.IntSlice(uidSet))
	result := make([]map[string]interface{}, len(messages))
	for i, message := range messages {
		tmpMap := make(map[string]interface{})
		util.Struct2Map(tmpMap, message)
		tmpMap["user"] = userMap[message.From]
		// 为了跟系统消息一致
		tmpMap["title"] = "发来了一条消息："
		result[i] = tmpMap
	}
	return result
}

func (MessageLogic) ToMsgCount(ctx context.Context, uid int) int64 {
	total, _ := MasterDB.Where("`to`=? AND tdel=?", uid, model.TdelNotDel).Count(new(model.Message))
	return total
}

// 获取某人发送的消息
func (MessageLogic) FindFromMsgsByUid(ctx context.Context, uid int, paginator *Paginator) []map[string]interface{} {
	objLog := GetLogger(ctx)

	messages := make([]*model.Message, 0)
	err := MasterDB.Where("`from`=? AND fdel=?", uid, model.FdelNotDel).
		OrderBy("id DESC").Limit(paginator.PerPage(), paginator.Offset()).Find(&messages)
	if err != nil {
		objLog.Errorln("message logic FindFromMsgsByUid Error:", err)
		return nil
	}

	uids := util.Models2Intslice(messages, "To")
	userMap := DefaultUser.FindUserInfos(ctx, uids)
	result := make([]map[string]interface{}, len(messages))
	for i, message := range messages {
		tmpMap := make(map[string]interface{})
		util.Struct2Map(tmpMap, message)
		tmpMap["user"] = userMap[message.To]
		result[i] = tmpMap
	}
	return result
}

func (MessageLogic) FromMsgCount(ctx context.Context, uid int) int64 {
	total, _ := MasterDB.Where("`from`=? AND fdel=?", uid, model.FdelNotDel).Count(new(model.Message))
	return total
}

// MarkHasRead 标记消息已读
func (MessageLogic) MarkHasRead(ctx context.Context, ids []int, isSysMsg bool, uid int) bool {
	if len(ids) == 0 {
		return true
	}

	var session *xorm.Session
	if isSysMsg {
		session = MasterDB.Table(new(model.SystemMessage))
	} else {
		session = MasterDB.Table(new(model.Message))
	}

	if len(ids) > 1 {
		session.In("id", ids)
	} else {
		session.Id(ids[0])
	}

	_, err := session.Update(map[string]interface{}{"hasread": model.HasRead})
	if err != nil {
		logger.Errorln("message logic MarkHasRead Error:", err)
		return false
	}
	// 将显示的消息数减少
	msg := NewMessage(WsMsgNotify, -len(ids))
	go Book.PostMessage(uid, msg)
	return true
}

// DeleteMessage 删除消息
// msgtype -> system(系统消息)/inbox(outbox)(短消息)
func (MessageLogic) DeleteMessage(ctx context.Context, id, msgtype string) bool {
	var err error
	if msgtype == "system" {
		_, err = MasterDB.Id(id).Delete(&model.SystemMessage{})
	} else if msgtype == "inbox" {
		// 打标记
		_, err = MasterDB.Table(new(model.Message)).Id(id).Update(map[string]interface{}{"tdel": model.TdelHasDel})
	} else {
		_, err = MasterDB.Table(new(model.Message)).Id(id).Update(map[string]interface{}{"fdel": model.FdelHasDel})
	}
	if err != nil {
		logger.Errorln("message logic DeleteMessage Error:", err)
		return false
	}
	return true
}

// FindNotReadMsgNum 获得某个用户未读消息数（系统消息和短消息）
func (MessageLogic) FindNotReadMsgNum(ctx context.Context, uid int) int {
	sysMsgNum, err := MasterDB.Where("`to`=? AND hasread=?", uid, model.NotRead).Count(new(model.SystemMessage))
	if err != nil {
		logger.Errorln("SystemMessage logic FindNotReadMsgNum Error:", err)
	}

	msgNum, err := MasterDB.Where("`to`=? AND hasread=?", uid, model.NotRead).Count(new(model.Message))
	if err != nil {
		logger.Errorln("Message logic FindNotReadMsgNum Error:", err)
	}
	return int(sysMsgNum + msgNum)
}
