// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package admin

import (
	"net/http"

	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
)

type RuleController struct{}

// 注册路由
func (self RuleController) RegisterRoute(g *echo.Group) {
	g.GET("/crawl/rule/list", self.RuleList)
	g.POST("/crawl/rule/query.html", self.Query)
	g.Match([]string{"GET", "POST"}, "/crawl/rule/new", self.New)
	g.Match([]string{"GET", "POST"}, "/crawl/rule/modify", self.Modify)
	g.POST("/crawl/rule/del", self.Del)
}

// RuleList 所有规则（分页）
func (RuleController) RuleList(ctx echo.Context) error {
	curPage, limit := parsePage(ctx)

	rules, total := logic.DefaultRule.FindBy(context.EchoContext(ctx), nil, curPage, limit)

	if rules == nil {
		return ctx.HTML(http.StatusInternalServerError, "500")
	}

	data := map[string]interface{}{
		"datalist":   rules,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}

	return render(ctx, "rule/list.html,rule/query.html", data)
}

// Query
func (RuleController) Query(ctx echo.Context) error {
	curPage, limit := parsePage(ctx)
	conds := parseConds(ctx, []string{"domain"})

	rules, total := logic.DefaultRule.FindBy(context.EchoContext(ctx), conds, curPage, limit)

	if rules == nil {
		return ctx.HTML(http.StatusInternalServerError, "500")
	}

	data := map[string]interface{}{
		"datalist":   rules,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
		"page":       curPage,
		"limit":      limit,
	}
	return render(ctx, "rule/query.html", data)
}

// New 新建规则
func (RuleController) New(ctx echo.Context) error {
	var data = make(map[string]interface{})

	if ctx.FormValue("submit") == "1" {
		user := ctx.Get("user").(*model.Me)
		forms, _ := ctx.FormParams()
		errMsg, err := logic.DefaultRule.Save(context.EchoContext(ctx), forms, user.Username)
		if err != nil {
			return fail(ctx, 1, errMsg)
		}
		return success(ctx, nil)
	}

	return render(ctx, "rule/new.html", data)
}

// Modify 编辑规则
func (self RuleController) Modify(ctx echo.Context) error {
	var data = make(map[string]interface{})

	if ctx.FormValue("submit") == "1" {
		user := ctx.Get("user").(*model.Me)
		forms, _ := ctx.FormParams()
		errMsg, err := logic.DefaultRule.Save(context.EchoContext(ctx), forms, user.Username)
		if err != nil {
			return fail(ctx, 1, errMsg)
		}
		return success(ctx, nil)
	}

	rule := logic.DefaultRule.FindById(context.EchoContext(ctx), ctx.QueryParam("id"))
	if rule == nil {
		return ctx.Redirect(http.StatusSeeOther, ctx.Echo().URI(echo.HandlerFunc(self.RuleList)))
	}

	data["rule"] = rule

	return render(ctx, "rule/modify.html", data)
}

func (RuleController) Del(ctx echo.Context) error {
	err := logic.DefaultRule.Delete(context.EchoContext(ctx), ctx.FormValue("id"))
	if err != nil {
		return fail(ctx, 1, "删除失败")
	}
	return success(ctx, nil)
}
