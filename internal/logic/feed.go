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
	"github.com/studygolang/studygolang/internal/dao/cache"
	"github.com/studygolang/studygolang/internal/model"
	"github.com/studygolang/studygolang/util"

	"github.com/polaris1119/set"
	"xorm.io/xorm"
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

	feeds := cache.Feed.GetList(ctx, paginator.curPage, tab)
	if len(feeds) > 0 {
		// 缓存里的列表可能是内容下线前写入的，命中后仍需按主表状态过滤
		return filterByObjectStatus(ctx, feeds)
	}

	feeds = make([]*model.Feed, 0)
	session := MasterDB.Limit(paginator.PerPage(), paginator.Offset())
	if tab == model.TabRecommend {
		session.Desc("seq")
	}
	err := session.Desc("updated_at").Find(&feeds)
	if err != nil {
		objLog.Errorln("FeedLogic FindRecent error:", err)
		return nil
	}

	feeds = self.fillOtherInfo(ctx, feeds, true)
	if len(feeds) > 0 {
		cache.Feed.SetList(ctx, paginator.curPage, tab, feeds)
	}
	return feeds
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

	feeds := cache.Feed.GetTop(ctx)
	if feeds != nil {
		return filterByObjectStatus(ctx, feeds)
	}

	feeds = make([]*model.Feed, 0)
	err := MasterDB.Where("top=1").Desc("updated_at").Find(&feeds)
	if err != nil {
		objLog.Errorln("FeedLogic FindRecent error:", err)
		return nil
	}

	feeds = self.fillOtherInfo(ctx, feeds, false)
	cache.Feed.SetTop(ctx, feeds)
	return feeds
}

// AutoUpdateSeq 自动更新动态的排序（校准）
func (self FeedLogic) AutoUpdateSeq() {
	curHour := time.Now().Hour()
	if curHour < 7 {
		return
	}
	defer cache.Feed.InvalidateAll(context.Background())

	feedDay := config.ConfigFile.MustInt("feed", "day", 3)
	cmtWeight := config.ConfigFile.MustInt("feed", "cmt_weight", 80)
	viewWeight := config.ConfigFile.MustInt("feed", "view_weight", 80)

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

			// 当天（不到24小时）发布的，不降
			elapse := int(time.Now().Sub(time.Time(feed.CreatedAt)).Hours())
			if elapse < 24 {
				continue
			}

			if feed.Uid > 0 {
				user := DefaultUser.FindOne(nil, "uid", feed.Uid)
				if DefaultUser.IsAdmin(user) {
					elapse = int(time.Now().Sub(time.Time(feed.UpdatedAt)).Hours())
				}
			}

			seq := 0
			if elapse <= feedDay*24 {
				seq = self.calcChangeSeq(feed, cmtWeight, viewWeight)
			}

			MasterDB.Table(new(model.Feed)).Where("id=?", feed.Id).Update(map[string]interface{}{
				"updated_at": time.Time(feed.UpdatedAt),
				"seq":        seq,
			})
		}
	}
}

func (self FeedLogic) calcChangeSeq(feed *model.Feed, cmtWeight int, viewWeight int) int {
	seq := 0

	// 最近有评论（时间更新）的，降 1/10 个评论数
	if int(time.Now().Sub(time.Time(feed.UpdatedAt)).Hours()) < 1 {
		seq = feed.Seq - cmtWeight/10
	} else {
		// 最近有没有其他变动（赞、阅读等）
		var updatedAt time.Time
		switch feed.Objtype {
		case model.TypeTopic:
			topicEx := &model.TopicEx{}
			MasterDB.Where("tid=?", feed.Objid).Get(topicEx)
			updatedAt = topicEx.Mtime
		case model.TypeArticle:
			article := &model.Article{}
			MasterDB.ID(feed.Objid).Get(article)
			updatedAt = time.Time(article.Mtime)
		case model.TypeResource:
			resourceEx := &model.ResourceEx{}
			MasterDB.ID(feed.Objid).Get(resourceEx)
			updatedAt = resourceEx.Mtime
		case model.TypeProject:
			project := &model.OpenProject{}
			MasterDB.ID(feed.Objid).Get(project)
			updatedAt = time.Time(project.Mtime)
		case model.TypeBook:
			book := &model.Book{}
			MasterDB.ID(feed.Objid).Get(book)
			updatedAt = time.Time(book.UpdatedAt)
		}

		dynamicElapse := int(time.Now().Sub(updatedAt).Hours())

		if dynamicElapse < 1 {
			seq = feed.Seq - viewWeight*10
		} else {
			seq = feed.Seq / 2
		}
	}

	if seq < 20 {
		seq = 20
	}

	return seq
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

	// 主表状态过滤：feed.state 只在节点隐藏等少数场景维护，
	// 文章下线（status=Offline）、项目下线、话题删除（flag=FlagAuditDelete）
	// 不会同步 feed 表——读侧统一兜底，保证任何入口（含 SQL 直改）的
	// 下线内容都不会出现在首页动态里。过滤后的干净列表再进缓存。
	newFeeds = filterByObjectStatus(ctx, newFeeds)

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

// filterByObjectStatus 批量校验 feed 指向的主体是否仍有效，剔除已下线/删除的内容。
// 每次 1~3 个 IN 查询（仅涉及有状态语义的类型），成本远低于 join。
func filterByObjectStatus(ctx context.Context, feeds []*model.Feed) []*model.Feed {
	if len(feeds) == 0 {
		return feeds
	}
	objLog := GetLogger(ctx)

	topicIds := make([]int, 0, len(feeds))
	articleIds := make([]int, 0, len(feeds))
	projectIds := make([]int, 0, len(feeds))
	for _, feed := range feeds {
		switch feed.Objtype {
		case model.TypeTopic:
			topicIds = append(topicIds, feed.Objid)
		case model.TypeArticle:
			articleIds = append(articleIds, feed.Objid)
		case model.TypeProject:
			projectIds = append(projectIds, feed.Objid)
		}
	}

	valid := make(map[string]bool, len(topicIds)+len(articleIds)+len(projectIds))
	if len(topicIds) > 0 {
		topics := make([]*model.Topic, 0, len(topicIds))
		if err := MasterDB.In("tid", topicIds).Where("flag<?", model.FlagAuditDelete).Find(&topics); err != nil {
			objLog.Errorln("feed filterByObjectStatus topics error:", err)
		} else {
			for _, t := range topics {
				valid["t"+strconv.Itoa(t.Tid)] = true
			}
		}
	}
	if len(articleIds) > 0 {
		articles := make([]*model.Article, 0, len(articleIds))
		if err := MasterDB.In("id", articleIds).Where("status<?", model.ArticleStatusOffline).Find(&articles); err != nil {
			objLog.Errorln("feed filterByObjectStatus articles error:", err)
		} else {
			for _, a := range articles {
				valid["a"+strconv.Itoa(a.Id)] = true
			}
		}
	}
	if len(projectIds) > 0 {
		projects := make([]*model.OpenProject, 0, len(projectIds))
		if err := MasterDB.In("id", projectIds).
			Where("status IN(?,?)", model.ProjectStatusNew, model.ProjectStatusOnline).Find(&projects); err != nil {
			objLog.Errorln("feed filterByObjectStatus projects error:", err)
		} else {
			for _, p := range projects {
				valid["p"+strconv.Itoa(p.Id)] = true
			}
		}
	}

	result := make([]*model.Feed, 0, len(feeds))
	for _, feed := range feeds {
		var key string
		switch feed.Objtype {
		case model.TypeTopic:
			key = "t" + strconv.Itoa(feed.Objid)
		case model.TypeArticle:
			key = "a" + strconv.Itoa(feed.Objid)
		case model.TypeProject:
			key = "p" + strconv.Itoa(feed.Objid)
		default:
			// 资源/书籍/Wiki 等无下线语义，直接保留
			result = append(result, feed)
			continue
		}
		if valid[key] {
			result = append(result, feed)
		}
	}
	return result
}

// publish 发布动态
func (FeedLogic) publish(object interface{}, objectExt interface{}, me *model.Me) {
	go model.PublishFeed(object, objectExt, me)
}

func (self FeedLogic) updateSeq(objid, objtype, cmtnum, likenum, viewnum int) {
	cmtWeight := config.ConfigFile.MustInt("feed", "cmt_weight", 80)
	likeWeight := config.ConfigFile.MustInt("feed", "like_weight", 60)
	viewWeight := config.ConfigFile.MustInt("feed", "view_weight", 5)

	go func() {
		feed := &model.Feed{}
		_, err := MasterDB.Where("objid=? AND objtype=?", objid, objtype).Get(feed)
		if err != nil {
			return
		}

		if feed.State == model.FeedOffline {
			return
		}

		feedDay := config.ConfigFile.MustInt("feed", "day", 3)
		elapse := int(time.Now().Sub(time.Time(feed.CreatedAt)).Hours())

		if feed.Uid > 0 {
			user := DefaultUser.FindOne(nil, "uid", feed.Uid)
			if DefaultUser.IsAdmin(user) {
				elapse = int(time.Now().Sub(time.Time(feed.UpdatedAt)).Hours())
			}
		}

		seq := 0

		if elapse > feedDay*24 {
			if feed.Seq == 0 {
				return
			}
		} else {
			if feed.Seq == 0 {
				seq = feedDay*24 - elapse + (feed.Cmtnum+cmtnum)*cmtWeight + likenum*likeWeight + viewnum*viewWeight
			} else {
				seq = feed.Seq + cmtnum*cmtWeight + likenum*likeWeight + viewnum*viewWeight
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

		// 主动失效缓存：原 TTL 300s 期间列表页会显示旧排序/旧评论数。
		// 频繁更新场景的权衡：失效比持久陈旧更好（读取下次从 DB 拉新）。
		cache.Feed.InvalidateAll(context.Background())
	}()
}

// setTop 置顶或取消置顶
func (FeedLogic) setTop(session *xorm.Session, objid, objtype int, top int) error {
	_, err := session.Table(new(model.Feed)).Where("objid=? AND objtype=?", objid, objtype).
		Update(map[string]interface{}{
			"top": top,
		})
	if err != nil {
		return err
	}

	// 置顶变更直接影响列表头位置，必须立即失效缓存
	go cache.Feed.InvalidateAll(context.Background())
	return nil
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
		// updateSeq 内已 InvalidateAll
	}()
}

// updateLike 更新动态赞数据
func (self FeedLogic) updateLike(objid, objtype, uid, num int) {
	go func() {
		MasterDB.Where("objid=? AND objtype=?", objid, objtype).
			Incr("likenum", num).SetExpr("updated_at", "updated_at").
			Update(new(model.Feed))
	}()
	self.updateSeq(objid, objtype, 0, num, 0)
	// updateSeq 内已 InvalidateAll
}

func (self FeedLogic) modifyTopicNode(tid, nid int) {
	go func() {
		change := map[string]interface{}{
			"nid": nid,
		}

		node := &model.TopicNode{}
		_, err := MasterDB.ID(nid).Get(node)
		if err == nil && !node.ShowIndex {
			change["state"] = model.FeedOffline
		}
		MasterDB.Table(new(model.Feed)).Where("objid=? AND objtype=?", tid, model.TypeTopic).
			Update(change)
	}()
}
