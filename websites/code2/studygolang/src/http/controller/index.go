// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package controller

import (
	"logic"

	"github.com/labstack/echo"
)

type IndexController struct{}

// 注册路由
func (this *IndexController) RegisterRoute(e *echo.Echo) {
	e.Get("/", echo.HandlerFunc(this.Index))
}

// 首页
func (IndexController) Index(ctx echo.Context) error {
	// nodes := logic.GenNodes()

	num := 10
	paginator := logic.NewPaginatorWithPerPage(1, num)
	topicsList := make([]map[string]interface{}, num)

	// 置顶的topic
	topTopics := logic.DefaultTopic.FindAll(ctx, paginator, "ctime DESC", "top=1")
	if len(topTopics) < num {
		// 获取最新帖子
		paginator.SetPerPage(num - len(topTopics))
		newTopics := logic.DefaultTopic.FindAll(ctx, paginator, "ctime DESC", "top=0")

		topicsList = append(topTopics, newTopics...)
	}

	// 获得最新博文
	recentArticles := logic.DefaultArticle.FindBy(ctx, 10)
	// 获取当前用户喜欢对象信息
	var likeFlags map[int]int

	if len(recentArticles) > 0 {
		user, ok := filter.CurrentUser(req)
		if ok {
			uid := user["uid"].(int)

			likeFlags, _ = service.FindUserLikeObjects(uid, model.TYPE_ARTICLE, recentArticles[0].Id, recentArticles[len(recentArticles)-1].Id)
		}
	}

	// Golang 资源
	resources := service.FindResources("0", "10")

	return render(ctx, "index.html", map[string]interface{}{"topics": topicsList, "articles": recentArticles, "likeflags": likeFlags, "resources": resources})
}
