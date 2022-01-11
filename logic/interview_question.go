// Copyright 2022 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"bytes"
	"context"
	"net/url"
	"strconv"
	"time"

	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
	"github.com/polaris1119/nosql"
	. "github.com/studygolang/studygolang/db"
	"github.com/studygolang/studygolang/model"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const questionIDKey = "question:id"

type InterviewLogic struct{}

var DefaultInterview = InterviewLogic{}

func (InterviewLogic) Publish(ctx context.Context, form url.Values) (*model.InterviewQuestion, error) {
	objLog := GetLogger(ctx)

	var err error

	id := form.Get("id")
	isModify := id != ""

	interview := &model.InterviewQuestion{}

	if isModify {
		_, err = MasterDB.ID(id).Get(interview)
		if err != nil {
			objLog.Errorln("Publish interview question error:", err)
			return nil, err
		}

		err = schemaDecoder.Decode(interview, form)
		if err != nil {
			objLog.Errorln("Publish interview question schema decode error:", err)
			return nil, err
		}
	} else {
		err = schemaDecoder.Decode(interview, form)
		if err != nil {
			objLog.Errorln("Publish interview question schema decode error:", err)
			return nil, err
		}
	}

	// 生成 sn
	interview.Sn = snowFlake.NextID()

	if isModify {
		_, err = MasterDB.Update(interview)
	} else {
		_, err = MasterDB.Insert(interview)
	}

	if err != nil {
		objLog.Errorln("Publish interview error:", err)
		return nil, err
	}

	return interview, nil
}

func (iq InterviewLogic) TodayQuestion(ctx context.Context) *model.InterviewQuestion {
	objLog := GetLogger(ctx)

	redis := nosql.NewRedisFromPool()
	defer redis.Close()

	id := goutils.MustInt(redis.GET(questionIDKey), 1)

	question := &model.InterviewQuestion{}
	_, err := MasterDB.ID(id).Get(question)
	if err != nil {
		objLog.Errorln("InterviewLogic TodayQuestion error:", err)
		return nil
	}

	err = iq.parseMarkdown(ctx, question)
	if err != nil {
		return nil
	}
	return question
}

func (iq InterviewLogic) FindOne(ctx context.Context, sn int64) (*model.InterviewQuestion, error) {
	question := &model.InterviewQuestion{}
	_, err := MasterDB.Where("sn=?", sn).Get(question)
	if err != nil {
		logger.Errorln("interview logic FindOne Error:", err)
		return nil, err
	}

	err = iq.parseMarkdown(ctx, question)
	return question, err
}

func (InterviewLogic) UpdateTodayQuestionID() {
	question := &model.InterviewQuestion{}
	_, err := MasterDB.Desc("id").Get(question)
	if err != nil {
		return
	}

	redis := nosql.NewRedisFromPool()
	defer redis.Close()

	id := goutils.MustInt(redis.GET(questionIDKey), 0)
	if id == 0 {
		id = 1
	} else {
		id = (id + 1) % question.Id
	}
	redis.SET(questionIDKey, id, 0)
}

// findByIds 获取多个问题详细信息 包内使用
func (InterviewLogic) findByIds(ids []int) map[int]*model.InterviewQuestion {
	if len(ids) == 0 {
		return nil
	}

	questions := make(map[int]*model.InterviewQuestion)
	err := MasterDB.In("id", ids).Find(&questions)
	if err != nil {
		logger.Errorln("InterviewLogic findByIds error:", err)
		return nil
	}
	return questions
}

func (InterviewLogic) parseMarkdown(ctx context.Context, question *model.InterviewQuestion) error {
	objLog := GetLogger(ctx)

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(question.Question), &buf); err != nil {
		objLog.Errorln("InterviewLogic TodayQuestion markdown convert error:", err)
		return err
	}
	question.Question = buf.String()

	buf.Reset()
	if err := md.Convert([]byte(question.Answer), &buf); err != nil {
		objLog.Errorln("InterviewLogic TodayQuestion markdown convert error:", err)
		return err
	}
	question.Answer = buf.String()

	return nil
}

// 面试题回复（评论）
type InterviewComment struct{}

// UpdateComment 更新该面试题的回复信息
// cid：评论id；objid：被评论对象id；uid：评论者；cmttime：评论时间
func (self InterviewComment) UpdateComment(cid, objid, uid int, cmttime time.Time) {
	// 更新回复数（TODO：暂时每次都更新表）
	_, err := MasterDB.ID(objid).Incr("cmtnum", 1).Update(new(model.InterviewQuestion))
	if err != nil {
		logger.Errorln("更新主题回复数失败：", err)
		return
	}
}

func (self InterviewComment) String() string {
	return "interview"
}

// 实现 CommentObjecter 接口
func (self InterviewComment) SetObjinfo(ids []int, commentMap map[int][]*model.Comment) {
	questions := DefaultInterview.findByIds(ids)
	if len(questions) == 0 {
		return
	}

	for _, question := range questions {
		strID := strconv.Itoa(question.Id)
		objinfo := make(map[string]interface{})
		objinfo["title"] = "Go每日一题（" + strID + "）"
		objinfo["uri"] = "/interview/question/" + question.ShowSn
		objinfo["type_name"] = model.TypeNameMap[model.TypeInterview]

		for _, comment := range commentMap[question.Id] {
			comment.Objinfo = objinfo
		}
	}
}

// 面试题喜欢
type InterviewLike struct{}

// 更新该面试题的喜欢数（赞数）
// objid：被喜欢对象id；num: 喜欢数(负数表示取消喜欢)
func (self InterviewLike) UpdateLike(objid, num int) {
	// 更新喜欢数（TODO：暂时每次都更新表）
	_, err := MasterDB.Where("id=?", objid).Incr("likenum", num).Update(new(model.InterviewQuestion))
	if err != nil {
		logger.Errorln("更新面试题喜欢数失败：", err)
	}
}

func (self InterviewLike) String() string {
	return "interview"
}
