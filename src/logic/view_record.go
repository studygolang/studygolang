// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"model"

	. "db"

	"github.com/polaris1119/logger"
	"golang.org/x/net/context"
)

type ViewRecordLogic struct{}

var DefaultViewRecord = ViewRecordLogic{}

func (ViewRecordLogic) Record(objid, objtype, uid int) {

	total, err := MasterDB.Where("objid=? AND objtype=? AND uid=?", objid, objtype, uid).Count(new(model.ViewRecord))
	if err != nil {
		logger.Errorln("ViewRecord logic Record count error:", err)
		return
	}

	if total > 0 {
		return
	}

	viewRecord := &model.ViewRecord{
		Objid:   objid,
		Objtype: objtype,
		Uid:     uid,
	}

	if _, err = MasterDB.Insert(viewRecord); err != nil {
		logger.Errorln("ViewRecord logic Record insert Error:", err)
		return
	}

	return
}

func (ViewRecordLogic) FindUserNum(ctx context.Context, objid, objtype int) int64 {
	objLog := GetLogger(ctx)

	total, err := MasterDB.Where("objid=? AND objtype=?", objid, objtype).Count(new(model.ViewRecord))
	if err != nil {
		objLog.Errorln("ViewRecordLogic FindUserNum error:", err)
	}

	return total
}
