// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"context"
	"model"
	"strconv"
	"time"

	. "db"

	"github.com/polaris1119/set"
)

type FeedLogic struct{}

var DefaultFeed = FeedLogic{}

func (self FeedLogic) FindRecent(ctx context.Context, num int) []*model.Feed {
	objLog := GetLogger(ctx)

	feeds := make([]*model.Feed, 0)
	err := MasterDB.Desc("updated_at").Limit(num).Find(&feeds)
	if err != nil {
		objLog.Errorln("FeedLogic FindRecent error:", err)
		return nil
	}

	return self.fillOtherInfo(ctx, feeds, true)
}

func (self FeedLogic) FindTop(ctx context.Context) []*model.Feed {
	objLog := GetLogger(ctx)

	feeds := make([]*model.Feed, 0)
	err := MasterDB.Where("top=1").Desc("updated_at").Find(&feeds)
	if err != nil {
		objLog.Errorln("FeedLogic FindRecent error:", err)
		return nil
	}

	return self.fillOtherInfo(ctx, feeds, false)
}

func (FeedLogic) fillOtherInfo(ctx context.Context, feeds []*model.Feed, filterTop bool) []*model.Feed {
	newFeeds := make([]*model.Feed, 0, len(feeds))

	uidSet := set.New(set.NonThreadSafe)
	for _, feed := range feeds {
		if feed.State == model.FeedOffline {
			continue
		}

		if filterTop && feed.Top == 1 {
			continue
		}

		newFeeds = append(newFeeds, feed)

		if feed.Uid > 0 {
			uidSet.Add(feed.Uid)
		}
		if feed.Lastreplyuid > 0 {
			uidSet.Add(feed.Lastreplyuid)
		}
		if feed.Objtype == model.TypeTopic {
			feed.Node = GetNode(feed.Nid)
		} else if feed.Objtype == model.TypeResource {
			feed.Node = map[string]interface{}{
				"name": GetCategoryName(feed.Nid),
			}
		}

		feed.Uri = model.PathUrlMap[feed.Objtype] + strconv.Itoa(feed.Objid)
	}

	usersMap := DefaultUser.FindUserInfos(ctx, set.IntSlice(uidSet))
	for _, feed := range newFeeds {
		if _, ok := usersMap[feed.Uid]; ok {
			feed.User = usersMap[feed.Uid]
		}
		if _, ok := usersMap[feed.Lastreplyuid]; ok {
			feed.Lastreplyuser = usersMap[feed.Lastreplyuid]
		}
	}

	return newFeeds
}

// publish 发布动态
func (FeedLogic) publish(object interface{}, objectExt interface{}) {
	go model.PublishFeed(object, objectExt)
}

// updateComment 更新动态评论数据
func (FeedLogic) updateComment(objid, objtype, uid int, cmttime time.Time) {
	go func() {
		MasterDB.Table(new(model.Feed)).Where("objid=? AND objtype=?", objid, objtype).
			Incr("cmtnum", 1).Update(map[string]interface{}{
			"lastreplyuid":  uid,
			"lastreplytime": cmttime,
		})
	}()
}
