// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import "github.com/labstack/echo"

func RegisterRoutes(g *echo.Group) {
	new(IndexController).RegisterRoute(g)
	new(AccountController).RegisterRoute(g)
	new(TopicController).RegisterRoute(g)
	new(ArticleController).RegisterRoute(g)
	new(ProjectController).RegisterRoute(g)
	new(ResourceController).RegisterRoute(g)
	new(ReadingController).RegisterRoute(g)
	new(WikiController).RegisterRoute(g)
	new(UserController).RegisterRoute(g)
	new(LikeController).RegisterRoute(g)
	new(FavoriteController).RegisterRoute(g)
	new(MessageController).RegisterRoute(g)
	new(SidebarController).RegisterRoute(g)
	new(CommentController).RegisterRoute(g)
	new(SearchController).RegisterRoute(g)
	new(WideController).RegisterRoute(g)
	new(ImageController).RegisterRoute(g)
	new(CaptchaController).RegisterRoute(g)
	new(BookController).RegisterRoute(g)
	new(MissionController).RegisterRoute(g)
	new(UserRichController).RegisterRoute(g)
	new(TopController).RegisterRoute(g)
	new(GiftController).RegisterRoute(g)
	new(OAuthController).RegisterRoute(g)
	new(WebsocketController).RegisterRoute(g)
	new(DownloadController).RegisterRoute(g)
	new(LinkController).RegisterRoute(g)
	new(SubjectController).RegisterRoute(g)
	new(GCTTController).RegisterRoute(g)

	new(FeedController).RegisterRoute(g)

	new(WechatController).RegisterRoute(g)

	new(InstallController).RegisterRoute(g)
}
