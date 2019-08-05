// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package app

import "github.com/labstack/echo"

func RegisterRoutes(g *echo.Group) {
	new(IndexController).RegisterRoute(g)
	new(ArticleController).RegisterRoute(g)
	new(TopicController).RegisterRoute(g)
	new(ResourceController).RegisterRoute(g)
	new(ProjectController).RegisterRoute(g)
	new(UserController).RegisterRoute(g)
	new(WechatController).RegisterRoute(g)
	new(CommentController).RegisterRoute(g)
}
