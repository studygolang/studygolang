// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"context"
	"strconv"
	"time"

	"github.com/polaris1119/config"
	"github.com/polaris1119/logger"

	. "github.com/studygolang/studygolang/db"
	"github.com/studygolang/studygolang/model"
	"github.com/studygolang/studygolang/util"

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

func (self FeedLogic) FindRecentWithPaginator(ctx context.Context, paginator *Paginator, tab string) []*model.Feed {
	objLog := GetLogger(ctx)

	feeds := make([]*model.Feed, 0)
	session := MasterDB.Limit(paginator.PerPage(), paginator.Offset())
	if tab == model.TabRecommend {
		session.Desc("seq")
	}
	err := session.Desc("updated_at").Find(&feeds)
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

// AutoUpdateSeq 每天自动更新一次动态的排序（校准）
func (self FeedLogic) AutoUpdateSeq() {
	feedDay := config.ConfigFile.MustInt("global", "feed_day", 7)

	var err error
	offset, limit := 0, 100
	for {
		feeds := make([]*model.Feed, 0)
		err = MasterDB.Where("seq>0").Limit(limit, offset).Find(&feeds)
		if err != nil || len(feeds) == 0 {
			return
		}

		offset += limit

		for _, feed := range feeds {
			if feed.State == model.FeedOffline {
				continue
			}

			elaspe := int(time.Now().Sub(time.Time(feed.CreatedAt)).Hours())

			if feed.Uid > 0 {
				user := DefaultUser.FindOne(nil, "uid", feed.Uid)
				if DefaultUser.IsAdmin(user) {
					elaspe = int(time.Now().Sub(time.Time(feed.UpdatedAt)).Hours())
				}
			}

			if elaspe > feedDay*24 {
				MasterDB.Table(new(model.Feed)).Where("id=?", feed.Id).Update(map[string]interface{}{
					"updated_at": time.Time(feed.UpdatedAt),
					"seq":        0,
				})
			}
		}
	}
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
func (FeedLogic) publish(object interface{}, objectExt interface{}, me *model.Me) {
	go model.PublishFeed(object, objectExt, me)
}

func (self FeedLogic) updateSeq(objid, objtype, cmtnum, likenum, viewnum int) {
	go func() {
		feed := &model.Feed{}
		_, err := MasterDB.Where("objid=? AND objtype=?", objid, objtype).Get(feed)
		if err != nil {
			return
		}

		if feed.State == model.FeedOffline {
			return
		}

		feedDay := config.ConfigFile.MustInt("global", "feed_day", 7)
		elaspe := int(time.Now().Sub(time.Time(feed.CreatedAt)).Hours())

		if feed.Uid > 0 {
			user := DefaultUser.FindOne(nil, "uid", feed.Uid)
			if DefaultUser.IsAdmin(user) {
				elaspe = int(time.Now().Sub(time.Time(feed.UpdatedAt)).Hours())
			}
		}

		seq := 0

		if elaspe > feedDay*24 {
			if feed.Seq == 0 {
				return
			}
		} else {
			if feed.Seq == 0 {
				seq = elaspe + (feed.Cmtnum+cmtnum)*100 + likenum*100 + viewnum*5
			} else {
				seq = feed.Seq - elaspe + cmtnum*100 + likenum*100 + viewnum*5
			}
		}

		_, err = MasterDB.Table(new(model.Feed)).Where("objid=? AND objtype=?", objid, objtype).Update(map[string]interface{}{
			"updated_at": time.Time(feed.UpdatedAt),
			"seq":        seq,
		})

		if err != nil {
			logger.Errorln("update feed seq error:", err)
			return
		}
	}()
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
func (self FeedLogic) updateComment(objid, objtype, uid int, cmttime time.Time) {
	go func() {
		MasterDB.Table(new(model.Feed)).Where("objid=? AND objtype=?", objid, objtype).
			Incr("cmtnum", 1).Update(map[string]interface{}{
			"lastreplyuid":  uid,
			"lastreplytime": cmttime,
		})

		self.updateSeq(objid, objtype, 1, 0, 0)
	}()
}

// updateLike 更新动态赞数据（暂时没存）
func (self FeedLogic) updateLike(objid, objtype, uid int, liketime time.Time) {
	self.updateSeq(objid, objtype, 0, 1, 0)
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
