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
	"util"

	. "db"

	"github.com/go-xorm/xorm"
	"github.com/polaris1119/set"
)

type FeedLogic struct{}

var DefaultFeed = FeedLogic{}

func (self FeedLogic) GetTotalCount(ctx context.Context) int64 {
	objLog := GetLogger(ctx)
	count, err := MasterDB.Where("state=0").Count(new(model.Feed))
	if err != nil {
		objLog.Errorln("FeedLogic Count error:", err)
		return 0
	}
	return count
}

func (self FeedLogic) FindRecentWithPaginator(ctx context.Context, paginator *Paginator) []*model.Feed {
	objLog := GetLogger(ctx)

	feeds := make([]*model.Feed, 0)
	err := MasterDB.Desc("updated_at").Limit(paginator.PerPage(), paginator.Offset()).Find(&feeds)
	if err != nil {
		objLog.Errorln("FeedLogic FindRecent error:", err)
		return nil
	}

	return self.fillOtherInfo(ctx, feeds, true)
}

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
	nidSet := set.New(set.NonThreadSafe)
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
			nidSet.Add(feed.Nid)
		} else if feed.Objtype == model.TypeResource {
			feed.Node = map[string]interface{}{
				"name": GetCategoryName(feed.Nid),
			}
		}

		feed.Uri = model.PathUrlMap[feed.Objtype] + strconv.Itoa(feed.Objid)
	}

	usersMap := DefaultUser.FindUserInfos(ctx, set.IntSlice(uidSet))
	nodesMap := GetNodesByNids(set.IntSlice(nidSet))
	for _, feed := range newFeeds {
		if _, ok := usersMap[feed.Uid]; ok {
			feed.User = usersMap[feed.Uid]
		}
		if _, ok := usersMap[feed.Lastreplyuid]; ok {
			feed.Lastreplyuser = usersMap[feed.Lastreplyuid]
		}

		if feed.Objtype == model.TypeTopic {
			if _, ok := nodesMap[feed.Nid]; ok {
				feed.Node = map[string]interface{}{}
				util.Struct2Map(feed.Node, nodesMap[feed.Nid])
			}
		}
	}

	return newFeeds
}

// publish 发布动态
func (FeedLogic) publish(object interface{}, objectExt interface{}) {
	go model.PublishFeed(object, objectExt)
}

// setTop 置顶或取消置顶
func (FeedLogic) setTop(session *xorm.Session, objid, objtype int, top int) error {
	_, err := session.Table(new(model.Feed)).Where("objid=? AND objtype=?", objid, objtype).
		Update(map[string]interface{}{
			"top": top,
		})

	return err
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

func (self FeedLogic) modifyTopicNode(tid, nid int) {
	go func() {
		change := map[string]interface{}{
			"nid": nid,
		}

		node := &model.TopicNode{}
		_, err := MasterDB.Id(nid).Get(node)
		if err == nil && !node.ShowIndex {
			change["state"] = model.FeedOffline
		}
		MasterDB.Table(new(model.Feed)).Where("objid=? AND objtype=?", tid, model.TypeTopic).
			Update(change)
	}()
}
