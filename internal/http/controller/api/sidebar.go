// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"strconv"
	"time"

	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/slices"
)

type SidebarController struct{}

// RegisterRoute 注册侧边栏相关路由
func (self SidebarController) RegisterRoute(g *echo.Group) {
	g.GET("/sidebar/readings/recent", self.RecentReading)
	g.GET("/sidebar/topics/recent", self.RecentTopic)
	g.GET("/sidebar/articles/recent", self.RecentArticle)
	g.GET("/sidebar/projects/recent", self.RecentProject)
	g.GET("/sidebar/resources/recent", self.RecentResource)
	g.GET("/sidebar/comments/recent", self.RecentComment)
	g.GET("/sidebar/nodes/hot", self.HotNodes)
	g.GET("/sidebar/users/active", self.ActiveUser)
	g.GET("/sidebar/users/newest", self.NewestUser)
	g.GET("/sidebar/friend/links", self.FriendLinks)
	g.GET("/sidebar/dynamics/recent", self.RecentDynamics)
	g.GET("/sidebar/rank/view", self.ViewRank)
}

// RecentReading 最近晨读，limit 参数默认 7
// 返回格式: { readings: [...] }
func (SidebarController) RecentReading(ctx echo.Context) error {
	limit := boundLimit(goutils.MustInt(ctx.QueryParam("limit"), 7), 7)
	readings := logic.DefaultReading.FindBy(context.EchoContext(ctx), limit, model.RtypeGo)
	if len(readings) == 1 {
		if time.Time(readings[0].Ctime).Before(time.Now().Add(-3 * 24 * time.Hour)) {
			readings = nil
		}
	}
	return success(ctx, map[string]interface{}{
		"readings": readings,
	})
}

// RecentTopic 最近话题，limit 参数默认 10
// 返回格式: { topics: [...] }
func (SidebarController) RecentTopic(ctx echo.Context) error {
	limit := boundLimit(goutils.MustInt(ctx.QueryParam("limit"), 10), 10)
	topicList := logic.DefaultTopic.FindRecent(limit)
	return success(ctx, map[string]interface{}{
		"topics": topicList,
	})
}

// RecentArticle 最近文章，limit 参数默认 10
// 返回格式: { articles: [...] }
func (SidebarController) RecentArticle(ctx echo.Context) error {
	limit := boundLimit(goutils.MustInt(ctx.QueryParam("limit"), 10), 10)
	recentArticles := logic.DefaultArticle.FindBy(context.EchoContext(ctx), limit)
	return success(ctx, map[string]interface{}{
		"articles": recentArticles,
	})
}

// RecentProject 最近项目，limit 参数默认 10
// 返回格式: { projects: [...] }
func (SidebarController) RecentProject(ctx echo.Context) error {
	limit := boundLimit(goutils.MustInt(ctx.QueryParam("limit"), 10), 10)
	recentProjects := logic.DefaultProject.FindBy(context.EchoContext(ctx), limit)
	return success(ctx, map[string]interface{}{
		"projects": recentProjects,
	})
}

// RecentResource 最近资源，limit 参数默认 10
// 返回格式: { resources: [...] }
func (SidebarController) RecentResource(ctx echo.Context) error {
	limit := boundLimit(goutils.MustInt(ctx.QueryParam("limit"), 10), 10)
	recentResources := logic.DefaultResource.FindBy(context.EchoContext(ctx), limit)
	return success(ctx, map[string]interface{}{
		"resources": recentResources,
	})
}

// RecentComment 最近评论，limit 参数默认 10
// 返回格式: { comments: [...], "<uid>": user }
func (SidebarController) RecentComment(ctx echo.Context) error {
	limit := boundLimit(goutils.MustInt(ctx.QueryParam("limit"), 10), 10)
	recentComments := logic.DefaultComment.FindRecent(context.EchoContext(ctx), 0, -1, limit)

	uids := slices.StructsIntSlice(recentComments, "Uid")
	users := logic.DefaultUser.FindUserInfos(context.EchoContext(ctx), uids)

	// Email 脱敏：未公开邮箱的用户清空 Email，与 profile.html 语义对齐。
	for uid, u := range users {
		users[uid] = sanitizeUserForPublic(u)
	}

	result := map[string]interface{}{
		"comments": recentComments,
	}
	// json encode 不支持 map[int]...，转换为字符串 key
	for uid, user := range users {
		result[strconv.Itoa(uid)] = user
	}

	return success(ctx, result)
}

// HotNodes 热门节点
// 返回格式: { nodes: [...] }
func (SidebarController) HotNodes(ctx echo.Context) error {
	nodes := logic.DefaultTopic.FindHotNodes(context.EchoContext(ctx))
	return success(ctx, map[string]interface{}{
		"nodes": nodes,
	})
}

// ActiveUser 活跃用户
// 返回格式: { users: [...] }
func (SidebarController) ActiveUser(ctx echo.Context) error {
	activeUsers := sanitizeUsersForPublic(logic.DefaultRank.FindDAURank(context.EchoContext(ctx), 9))
	return success(ctx, map[string]interface{}{
		"users": activeUsers,
	})
}

// NewestUser 最新用户
// 返回格式: { users: [...] }
func (SidebarController) NewestUser(ctx echo.Context) error {
	newestUsers := sanitizeUsersForPublic(logic.DefaultUser.FindNewUsers(context.EchoContext(ctx), 9))
	return success(ctx, map[string]interface{}{
		"users": newestUsers,
	})
}

// FriendLinks 友情链接
// 返回格式: { links: [...] }
func (SidebarController) FriendLinks(ctx echo.Context) error {
	friendLinks := logic.DefaultFriendLink.FindAll(context.EchoContext(ctx), 10)
	return success(ctx, map[string]interface{}{
		"links": friendLinks,
	})
}

// RecentDynamics 最近综合动态（话题、文章、评论、资源等)
// 返回格式: { dynamics: [...] }
func (SidebarController) RecentDynamics(ctx echo.Context) error {
	limit := boundLimit(goutils.MustInt(ctx.QueryParam("limit"), 10), 10)

	// 获取最近的话题
	recentTopics := logic.DefaultTopic.FindRecent(limit)
	// 获取最近的文章
	recentArticles := logic.DefaultArticle.FindBy(context.EchoContext(ctx), limit)
	// 获取最近的评论
	recentComments := logic.DefaultComment.FindRecent(context.EchoContext(ctx), 0, -1, limit)
	// 获取最近的资源
	recentResources := logic.DefaultResource.FindRecent(context.EchoContext(ctx), 0)

	// 合并成综合动态列表
	dynamics := make([]map[string]interface{}, 0, len(recentTopics)+len(recentArticles)+len(recentComments)+len(recentResources))

	for _, t := range recentTopics {
		dynamics = append(dynamics, map[string]interface{}{
			"type":    "topic",
			"objid":   t.Tid,
			"objtype": model.TypeTopic,
			"title":   t.Title,
			"uid":     t.Uid,
			"ctime":   t.Ctime,
		})
	}
	for _, a := range recentArticles {
		dynamics = append(dynamics, map[string]interface{}{
			"type":    "article",
			"objid":   a.Id,
			"objtype": model.TypeArticle,
			"title":   a.Title,
			"author":  a.AuthorTxt,
			"ctime":   a.Ctime,
		})
	}
	for _, c := range recentComments {
		dynamics = append(dynamics, map[string]interface{}{
			"type":    "comment",
			"objid":   c.Cid,
			"objtype": model.TypeComment,
			"uid":     c.Uid,
			"content": c.Content,
			"ctime":   c.Ctime,
		})
	}
	for _, res := range recentResources {
		dynamics = append(dynamics, map[string]interface{}{
			"type":    "resource",
			"objid":   res.Id,
			"objtype": model.TypeResource,
			"name":    res.Title,
			"uid":     res.Uid,
			"ctime":   res.Ctime,
		})
	}

	return success(ctx, map[string]interface{}{
		"dynamics": dynamics,
	})
}

// ViewRank 浏览量排行榜
// 返回格式: { topics: [...], articles: [...], resources: [...] }
func (SidebarController) ViewRank(ctx echo.Context) error {
	limit := boundLimit(goutils.MustInt(ctx.QueryParam("limit"), 10), 10)

	// 获取话题浏览量排行（本周）
	weekTopics := logic.DefaultRank.FindWeekRank(context.EchoContext(ctx), model.TypeTopic, limit, true)

	// 获取文章浏览量排行（本周）
	weekArticles := logic.DefaultRank.FindWeekRank(context.EchoContext(ctx), model.TypeArticle, limit, true)

	// 获取资源浏览量排行（本周）
	weekResources := logic.DefaultRank.FindWeekRank(context.EchoContext(ctx), model.TypeResource, limit, true)

	return success(ctx, map[string]interface{}{
		"topics":    weekTopics,
		"articles":  weekArticles,
		"resources": weekResources,
	})
}
