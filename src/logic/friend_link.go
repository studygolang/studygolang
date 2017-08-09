// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	. "db"

	"model"

	"golang.org/x/net/context"
)

type FriendLinkLogic struct{}

var DefaultFriendLink = FriendLinkLogic{}

func (FriendLinkLogic) FindAll(ctx context.Context) []*model.FriendLink {
	objLog := GetLogger(ctx)

	friendLinks := make([]*model.FriendLink, 0)
	err := MasterDB.OrderBy("seq asc").Find(&friendLinks)
	if err != nil {
		objLog.Errorln("FriendLinkLogic FindAll error:", err)
		return nil
	}

	return friendLinks
}
