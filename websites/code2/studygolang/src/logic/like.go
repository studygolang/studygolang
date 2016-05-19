// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"errors"
	"fmt"

	. "db"

	"golang.org/x/net/context"

	"model"
)

type LikeLogic struct{}

var DefaultLike = LikeLogic{}

// HadLike 某个用户是否已经喜欢某个对象
func (LikeLogic) HadLike(ctx context.Context, uid, objid, objtype int) int {
	objLog := GetLogger(ctx)

	like := &model.Like{}
	_, err := MasterDB.Where("uid=? AND objid=? AND objtype=?", uid, objid, objtype).Get(like)
	if err != nil {
		objLog.Errorln("LikeLogic HadLike error:", err)
		return 0
	}

	return like.Flag
}

// FindUserLikeObjects 获取用户对一批对象是否喜欢的状态
// objids 两个值
func (LikeLogic) FindUserLikeObjects(ctx context.Context, uid, objtype int, objids ...int) (map[int]int, error) {
	objLog := GetLogger(ctx)

	if len(objids) < 2 {
		return nil, errors.New("参数错误")
	}

	littleId, greatId := objids[0], objids[1]
	if littleId > greatId {
		littleId, greatId = greatId, littleId
	}

	likes := make([]*model.Like, 0)
	err := MasterDB.Where("uid=? AND objtype=? AND objid BETWEEN ? AND ?", uid, objtype, littleId, greatId).
		Find(&likes)
	if err != nil {
		objLog.Errorln("LikeLogic FindUserLikeObjects error:", err)
		return nil, err
	}

	likeFlags := make(map[int]int, len(likes))
	for _, like := range likes {
		likeFlags[like.Objid] = like.Flag
	}

	return likeFlags, nil
}

// LikeObject 喜欢或取消喜欢
// objid 注册的喜欢对象
// uid 喜欢的人
func (LikeLogic) LikeObject(ctx context.Context, uid, objid, objtype, likeFlag int) error {
	objLog := GetLogger(ctx)

	// 点喜欢，活跃度+3
	go DefaultUser.IncrUserWeight("uid", uid, 3)

	like := &model.Like{}
	_, err := MasterDB.Where("uid=? AND objid=? AND objtype=?", uid, objid, objtype).Get(like)
	if err != nil {
		objLog.Errorln("LikeLogic LikeObject get error:", err)
		return err
	}

	// 之前喜欢过
	if like.Uid != 0 {
		// 再喜欢直接返回
		if likeFlag == model.FlagLike {
			return nil
		}

		// 取消喜欢
		if likeFlag == model.FlagCancel {
			// MasterDB.Where("uid=? AND objid=? AND objtype=?", uid, objid,objtype).Delete(like)
			_, err = MasterDB.Delete(like)
			if err != nil {
				objLog.Errorln("LikeLogic LikeObject delete error:", err)
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

	affectedRows, err := MasterDB.Insert(like)
	if err != nil {
		objLog.Errorln("LikeLogic LikeObject error:", err)
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

var likers = make(map[int]Liker)

// 喜欢接口
type Liker interface {
	fmt.Stringer
	// 喜欢 回调接口，用于更新对象自身需要更新的数据
	UpdateLike(int, int)
}

// 注册喜欢对象，使得某种类型（主题、博文等）被喜欢了可以回调
func RegisterLikeObject(objtype int, liker Liker) {
	if liker == nil {
		panic("logic: Register liker is nil")
	}
	if _, dup := likers[objtype]; dup {
		panic("logic: Register called twice for liker " + liker.String())
	}
	likers[objtype] = liker
}
