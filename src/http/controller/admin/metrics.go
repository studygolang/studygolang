// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package admin

import (
	"expvar"
	"global"
	"logic"
	"time"

	"github.com/labstack/echo"

	. "http"
)

var (
	onlineStats    = expvar.NewMap("online_stats")
	loginUserNum   expvar.Int
	visitorUserNum expvar.Int
)

type MetricsController struct{}

// 注册路由
func (self MetricsController) RegisterRoute(g *echo.Group) {
	g.GET("/debug/vars", self.DebugExpvar)
}

func (self MetricsController) DebugExpvar(ctx echo.Context) error {
	loginUserNum.Set(int64(logic.Book.LoginLen()))
	visitorUserNum.Set(int64(logic.Book.Len()))

	onlineStats.Set("login_user_num", &loginUserNum)
	onlineStats.Set("visitor_user_num", &visitorUserNum)
	onlineStats.Set("uptime", expvar.Func(self.calculateUptime))
	onlineStats.Set("login_user_data", logic.Book.LoginUserData())

	handler := expvar.Handler()
	handler.ServeHTTP(ResponseWriter(ctx), Request(ctx))
	return nil
}

func (self MetricsController) calculateUptime() interface{} {
	return time.Since(global.App.LaunchTime).String()
}
