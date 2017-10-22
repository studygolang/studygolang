// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"model"
	"net/http"
	"strings"

	. "db"

	"github.com/polaris1119/logger"
	"golang.org/x/net/context"
)

type ViewSourceLogic struct{}

var DefaultViewSource = ViewSourceLogic{}

// Record 记录浏览来源
func (ViewSourceLogic) Record(req *http.Request, objtype, objid int) {
	referer := req.Referer()
	if referer == "" || strings.Contains(referer, WebsiteSetting.Domain) {
		return
	}

	viewSource := &model.ViewSource{}
	_, err := MasterDB.Where("objid=? AND objtype=?", objid, objtype).Get(viewSource)
	if err != nil {
		logger.Errorln("ViewSourceLogic Record find error:", err)
		return
	}

	if viewSource.Id == 0 {
		viewSource.Objid = objid
		viewSource.Objtype = objtype
		_, err = MasterDB.Insert(viewSource)
		if err != nil {
			logger.Errorln("ViewSourceLogic Record insert error:", err)
			return
		}
	}

	field := "other"
	referer = strings.ToLower(referer)
	ses := []string{"google", "baidu", "bing", "sogou", "so"}
	for _, se := range ses {
		if strings.Contains(referer, se+".") {
			field = se
			break
		}
	}

	_, err = MasterDB.Id(viewSource.Id).Incr(field, 1).Update(new(model.ViewSource))
	if err != nil {
		logger.Errorln("ViewSourceLogic Record update error:", err)
		return
	}
}

// FindOne 获得浏览来源
func (ViewSourceLogic) FindOne(ctx context.Context, objid, objtype int) *model.ViewSource {
	objLog := GetLogger(ctx)

	viewSource := &model.ViewSource{}
	_, err := MasterDB.Where("objid=? AND objtype=?", objid, objtype).Get(viewSource)
	if err != nil {
		objLog.Errorln("ViewSourceLogic FindOne error:", err)
	}

	return viewSource
}
