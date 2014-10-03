// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"logger"
	"model"
	"net/url"
	"util"
)

// 获取动态列表（分页）
func FindDynamics(lastId, limit string) []*model.Dynamic {
	dynamic := model.NewDynamic()

	dynamicList, err := dynamic.Where("id>" + lastId).Order("seq DESC").Limit(limit).
		FindAll()
	if err != nil {
		logger.Errorln("dynamic service FindDynamics Error:", err)
		return nil
	}

	return dynamicList
}

func PublishDynamic(form url.Values) error {
	dynamic := model.NewDynamic()

	util.ConvertAssign(dynamic, form)
	_, err := dynamic.Insert()
	return err
}

// 修改动态信息
func ModifyDynamic(form url.Values) (errMsg string, err error) {

	fields := []string{
		"content", "url", "seq", "dmtype",
	}
	query, args := updateSetClause(form, fields)

	id := form.Get("id")

	err = model.NewDynamic().Set(query, args...).Where("id=" + id).Update()
	if err != nil {
		logger.Errorf("更新动态 【%s】 信息失败：%s\n", id, err)
		errMsg = "对不起，服务器内部错误，请稍后再试！"
		return
	}

	return
}
