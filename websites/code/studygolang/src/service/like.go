// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"errors"
	"fmt"
	"strconv"

	"logger"
	"model"
)

// 某个用户是否已经喜欢某个对象
func HadLike(uid, objid, objtype int) int {
	cond := fmt.Sprintf("uid=%d AND objid=%d AND objtype=%d", uid, objid, objtype)

	like := model.NewLike()
	err := like.Where(cond).Find("flag")
	if err != nil {
		logger.Errorln("like service HadLike error:", err)
		return 0
	}

	return like.Flag
}

// 获取用户对一批对象是否喜欢的状态
// objids 两个值
func FindUserLikeObjects(uid, objtype int, objids ...int) (map[int]int, error) {
	if len(objids) < 2 {
		return nil, errors.New("参数错误")
	}

	littleId, greatId := objids[0], objids[1]
	if littleId > greatId {
		littleId, greatId = greatId, littleId
	}

	query := "uid=? AND objtype=? AND objid BETWEEN ? AND ?"
	args := []interface{}{uid, objtype, littleId, greatId}
	logger.Debugln("query:", query, ";args:", args)
	likes, err := model.NewLike().Where(query, args...).FindAll()
	if err != nil {
		return nil, err
	}

	likeFlags := make(map[int]int, len(objids))
	for _, like := range likes {
		likeFlags[like.Objid] = like.Flag
	}

	return likeFlags, nil
}

var likers = make(map[int]Liker)

// 喜欢接口
type Liker interface {
	fmt.Stringer
	// 喜欢回调接口，用于更新对象自身需要更新的数据
	UpdateLike(int, int)
}

// 注册喜欢对象，使得某种类型（帖子、博客等）被喜欢了可以回调
func RegisterLikeObject(objtype int, liker Liker) {
	if liker == nil {
		panic("service: Register liker is nil")
	}
	if _, dup := likers[objtype]; dup {
		panic("service: Register called twice for liker " + liker.String())
	}
	likers[objtype] = liker
}

// 喜欢或取消喜欢
// objid 注册的喜欢对象
// uid 喜欢的人
func LikeObject(uid, objid, objtype, likeFlag int) error {
	// 点喜欢，活跃度+3
	go IncUserWeight("uid="+strconv.Itoa(uid), 3)

	cond := fmt.Sprintf("uid=%d AND objid=%d AND objtype=%d", uid, objid, objtype)
	like := model.NewLike()

	err := like.Where(cond).Find()
	if err != nil {
		return err
	}

	// 之前喜欢过
	if like.Uid != 0 {
		// 再喜欢直接返回
		if likeFlag == model.FLAG_LIKE {
			return nil
		}

		// 取消喜欢
		if likeFlag == model.FLAG_CANCEL {
			err = like.Where(cond).Delete()
			if err != nil {
				return err
			}

			// 取消喜欢成功，更新对象的喜欢数
			if liker, ok := likers[objtype]; ok {
				go liker.UpdateLike(objid, -1)
			}

			return nil
		}

		return nil
	}

	like.Uid = uid
	like.Objid = objid
	like.Objtype = objtype
	like.Flag = likeFlag

	affectedRows, err := like.Insert()
	if err != nil {
		logger.Errorln("LikeObject service error:", err)
		return err
	}

	// 喜欢成功
	if affectedRows > 0 {
		if liker, ok := likers[objtype]; ok {
			go liker.UpdateLike(objid, 1)
		}
	}

	// TODO: 给被喜欢对象所有者发系统消息
	/*
		ext := map[string]interface{}{
			"objid":   objid,
			"objtype": objtype,
			"cid":     cid,
			"uid":     uid,
		}
		go SendSystemMsgTo(0, objtype, ext)
	*/

	return nil
}
