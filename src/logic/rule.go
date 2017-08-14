// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	. "db"
	"model"
	"net/url"

	"golang.org/x/net/context"
)

type RuleLogic struct{}

var DefaultRule = RuleLogic{}

// 获取抓取规则列表（分页）
func (RuleLogic) FindBy(ctx context.Context, conds map[string]string, curPage, limit int) ([]*model.CrawlRule, int) {
	objLog := GetLogger(ctx)

	session := MasterDB.NewSession()

	for k, v := range conds {
		session.And(k+"=?", v)
	}

	totalSession := session.Clone()

	offset := (curPage - 1) * limit
	ruleList := make([]*model.CrawlRule, 0)
	err := session.OrderBy("id DESC").Limit(limit, offset).Find(&ruleList)
	if err != nil {
		objLog.Errorln("rule find error:", err)
		return nil, 0
	}

	total, err := totalSession.Count(new(model.CrawlRule))
	if err != nil {
		objLog.Errorln("rule find count error:", err)
		return nil, 0
	}

	return ruleList, int(total)
}

func (RuleLogic) FindById(ctx context.Context, id string) *model.CrawlRule {
	objLog := GetLogger(ctx)

	rule := &model.CrawlRule{}
	_, err := MasterDB.Id(id).Get(rule)
	if err != nil {
		objLog.Errorln("find rule error:", err)
		return nil
	}

	if rule.Id == 0 {
		return nil
	}

	return rule
}

func (RuleLogic) Save(ctx context.Context, form url.Values, opUser string) (errMsg string, err error) {
	objLog := GetLogger(ctx)

	rule := &model.CrawlRule{}
	err = schemaDecoder.Decode(rule, form)
	if err != nil {
		objLog.Errorln("rule Decode error", err)
		errMsg = err.Error()
		return
	}

	rule.OpUser = opUser

	if rule.Id != 0 {
		_, err = MasterDB.Id(rule.Id).Update(rule)
	} else {
		_, err = MasterDB.Insert(rule)
	}

	if err != nil {
		errMsg = "内部服务器错误"
		objLog.Errorln("rule save:", errMsg, ":", err)
		return
	}

	return
}

func (RuleLogic) Delete(ctx context.Context, id string) error {
	_, err := MasterDB.Id(id).Delete(new(model.CrawlRule))
	return err
}
