// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"errors"
	"fmt"

	"logger"
	"model"
)

func SaveFavorite(uid, objid, objtype int) error {
	favorite := model.NewFavorite()
	favorite.Uid = uid
	favorite.Objid = objid
	favorite.Objtype = objtype

	affectedNum, err := favorite.Insert()

	if err != nil {
		logger.Errorln("save favorite error:", err)
		return errors.New("内部服务错误")
	}

	if affectedNum == 0 {
		return errors.New("收藏失败！")
	}

	return nil
}

func CancelFavorite(uid, objid, objtype int) error {
	return model.NewFavorite().Where("uid=? AND objtype=? AND objid=?", uid, objtype, objid).Delete()
}

// 某个用户是否已经收藏某个对象
func HadFavorite(uid, objid, objtype int) int {
	favorite := model.NewFavorite()
	err := favorite.Where("uid=? AND objid=? and objtype=?", uid, objid, objtype).Find()
	if err != nil {
		logger.Errorln("favorite service HadFavorite error:", err)
		return 0
	}

	if favorite.Uid != 0 {
		return 1
	}

	return 0
}

func FindUserFavorites(uid, objtype, start, rows int) ([]*model.Favorite, int) {
	favorite := model.NewFavorite()

	limit := fmt.Sprintf("%d,%d", start, rows)
	favorites, err := favorite.Where("uid=? AND objtype=?", uid, objtype).Limit(limit).Order("objid DESC").FindAll()
	if err != nil {
		logger.Errorln("favorite service FindUserFavorites error:", err)
		return nil, 0
	}

	total, err := favorite.Count()
	if err != nil {
		logger.Errorln("favorite service FindUserFavorites count error:", err)
		return nil, 0
	}

	return favorites, total
}
