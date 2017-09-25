// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of self source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"logic"
	"model"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/slices"
	"github.com/polaris1119/times"
)

// 侧边栏的内容通过异步请求获取
type SidebarController struct{}

func (self SidebarController) RegisterRoute(g *echo.Group) {
	g.GET("/readings/recent", self.RecentReading)
	g.GET("/topics/:nid/others", self.OtherTopics)
	g.GET("/websites/stat", self.WebsiteStat)
	g.GET("/dynamics/recent", self.RecentDynamic)
	g.GET("/topics/recent", self.RecentTopic)
	g.GET("/articles/recent", self.RecentArticle)
	g.GET("/projects/recent", self.RecentProject)
	g.GET("/resources/recent", self.RecentResource)
	g.GET("/comments/recent", self.RecentComment)
	g.GET("/nodes/hot", self.HotNodes)
	g.GET("/users/active", self.ActiveUser)
	g.GET("/users/newest", self.NewestUser)
	g.GET("/friend/links", self.FriendLinks)
	g.GET("/rank/view", self.ViewRank)
}

// RecentReading 技术晨读
func (SidebarController) RecentReading(ctx echo.Context) error {
	limit := goutils.MustInt(ctx.QueryParam("limit"), 7)
	readings := logic.DefaultReading.FindBy(ctx, limit, model.RtypeGo)
	if len(readings) == 1 {
		// 首页，三天内的晨读才显示
		if time.Time(readings[0].Ctime).Before(time.Now().Add(-3 * 24 * time.Hour)) {
			readings = nil
		}
	}
	return success(ctx, readings)
}

// OtherTopics 某节点下其他帖子
func (SidebarController) OtherTopics(ctx echo.Context) error {
	topics := logic.DefaultTopic.FindByNid(ctx, ctx.Param("nid"), ctx.QueryParam("tid"))
	topics = logic.DefaultTopic.JSEscape(topics)
	return success(ctx, topics)
}

// WebsiteStat 网站统计信息
func (SidebarController) WebsiteStat(ctx echo.Context) error {
	articleTotal := logic.DefaultArticle.Total()
	projectTotal := logic.DefaultProject.Total()
	topicTotal := logic.DefaultTopic.Total()
	cmtTotal := logic.DefaultComment.Total()
	resourceTotal := logic.DefaultResource.Total()
	bookTotal := logic.DefaultGoBook.Total()
	userTotal := logic.DefaultUser.Total()

	data := map[string]interface{}{
		"article":  articleTotal,
		"project":  projectTotal,
		"topic":    topicTotal,
		"resource": resourceTotal,
		"book":     bookTotal,
		"comment":  cmtTotal,
		"user":     userTotal,
	}

	return success(ctx, data)
}

// RecentDynamic 社区最新公告或go最新动态
func (SidebarController) RecentDynamic(ctx echo.Context) error {
	dynamics := logic.DefaultDynamic.FindBy(ctx, 0, 3)
	return success(ctx, dynamics)
}

// RecentTopic 最新帖子
func (SidebarController) RecentTopic(ctx echo.Context) error {
	limit := goutils.MustInt(ctx.QueryParam("limit"), 10)
	topicList := logic.DefaultTopic.FindRecent(limit)
	return success(ctx, topicList)
}

// RecentArticle 最新博文
func (SidebarController) RecentArticle(ctx echo.Context) error {
	limit := goutils.MustInt(ctx.QueryParam("limit"), 10)
	recentArticles := logic.DefaultArticle.FindBy(ctx, limit)
	return success(ctx, recentArticles)
}

// RecentProject 最新开源项目
func (SidebarController) RecentProject(ctx echo.Context) error {
	limit := goutils.MustInt(ctx.QueryParam("limit"), 10)
	recentProjects := logic.DefaultProject.FindBy(ctx, limit)
	return success(ctx, recentProjects)
}

// RecentResource 最新资源
func (SidebarController) RecentResource(ctx echo.Context) error {
	limit := goutils.MustInt(ctx.QueryParam("limit"), 10)
	recentResources := logic.DefaultResource.FindBy(ctx, limit)
	return success(ctx, recentResources)
}

// RecentComment 最新评论
func (SidebarController) RecentComment(ctx echo.Context) error {
	limit := goutils.MustInt(ctx.QueryParam("limit"), 10)
	recentComments := logic.DefaultComment.FindRecent(ctx, 0, -1, limit)

	uids := slices.StructsIntSlice(recentComments, "Uid")
	users := logic.DefaultUser.FindUserInfos(ctx, uids)

	result := map[string]interface{}{
		"comments": recentComments,
	}

	// json encode 不支持 map[int]...
	for uid, user := range users {
		result[strconv.Itoa(uid)] = user
	}

	return success(ctx, result)
}

// HotNodes 社区热门节点
func (SidebarController) HotNodes(ctx echo.Context) error {
	nodes := logic.DefaultTopic.FindHotNodes(ctx)
	return success(ctx, nodes)
}

// ActiveUser 活跃会员
func (SidebarController) ActiveUser(ctx echo.Context) error {
	// activeUsers := logic.DefaultUser.FindActiveUsers(ctx, 9)
	// return success(ctx, activeUsers)
	activeUsers := logic.DefaultRank.FindDAURank(ctx, 9)
	return success(ctx, activeUsers)
}

// NewestUser 新加入会员
func (SidebarController) NewestUser(ctx echo.Context) error {
	newestUsers := logic.DefaultUser.FindNewUsers(ctx, 9)
	return success(ctx, newestUsers)
}

// FriendLinks 友情链接
func (SidebarController) FriendLinks(ctx echo.Context) error {
	friendLinks := logic.DefaultFriendLink.FindAll(ctx, 5)
	return success(ctx, friendLinks)
}

// ViewRank 阅读排行榜
func (SidebarController) ViewRank(ctx echo.Context) error {
	objtype := goutils.MustInt(ctx.QueryParam("objtype"))
	rankType := ctx.QueryParam("rank_type")
	limit := goutils.MustInt(ctx.QueryParam("limit"), 10)

	var result = map[string]interface{}{
		"objtype":   objtype,
		"rank_type": rankType,
	}
	switch rankType {
	case "today":
		result["list"] = logic.DefaultRank.FindDayRank(ctx, objtype, times.Format("ymd"), limit)
	case "yesterday":
		yesterday := time.Now().Add(-1 * 24 * time.Hour)
		result["list"] = logic.DefaultRank.FindDayRank(ctx, objtype, times.Format("ymd", yesterday), limit)
	case "week":
		result["list"] = logic.DefaultRank.FindWeekRank(ctx, objtype, limit)
	case "month":
		result["list"] = logic.DefaultRank.FindMonthRank(ctx, objtype, limit)
	}

	result["path"] = model.PathUrlMap[objtype]

	return success(ctx, result)
}
