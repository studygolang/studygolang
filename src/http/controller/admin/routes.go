// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package admin

import "github.com/labstack/echo"

func RegisterRoutes(g *echo.Group) {
	new(AuthorityController).RegisterRoute(g)
	new(UserController).RegisterRoute(g)
	new(TopicController).RegisterRoute(g)
	new(NodeController).RegisterRoute(g)
	new(ArticleController).RegisterRoute(g)
	new(ProjectController).RegisterRoute(g)
	new(RuleController).RegisterRoute(g)
	new(ReadingController).RegisterRoute(g)
	new(ToolController).RegisterRoute(g)
	new(SettingController).RegisterRoute(g)
	new(MetricsController).RegisterRoute(g)
}
