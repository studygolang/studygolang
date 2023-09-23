// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"bytes"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/polaris1119/config"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
)

type IndexController struct{}

// 注册路由
func (self IndexController) RegisterRoute(g *echo.Group) {
	g.GET("/", self.Index)
	g.GET("/wr", self.WrapUrl)
	g.GET("/pkgdoc", self.Pkgdoc)
	g.GET("/markdown", self.Markdown)
	g.GET("/link", self.Link)
}

func (IndexController) Index(ctx echo.Context) error {
	if len(logic.WebsiteSetting.IndexNavs) == 0 {
		return render(ctx, "index.html", nil)
	}

	tab := ctx.QueryParam("tab")
	if tab == "" {
		tab = GetFromCookie(ctx, "INDEX_TAB")
	}

	if tab == "" {
		tab = logic.WebsiteSetting.IndexNavs[0].Tab
	}
	paginator := logic.NewPaginator(goutils.MustInt(ctx.QueryParam("p"), 1))

	data := logic.DefaultIndex.FindData(context.EchoContext(ctx), tab, paginator)

	SetCookie(ctx, "INDEX_TAB", data["tab"].(string))

	data["all_nodes"] = logic.GenNodes()

	if tab == model.TabAll || tab == model.TabRecommend {
		pageHtml := paginator.SetTotal(logic.DefaultFeed.GetTotalCount(context.EchoContext(ctx))).GetPageHtml(ctx.Request().URL.Path)

		data["page"] = template.HTML(pageHtml)

		data["total"] = paginator.GetTotal()
	}

	return render(ctx, "index.html", data)
}

// WrapUrl 包装链接
func (IndexController) WrapUrl(ctx echo.Context) error {
	tUrl := ctx.QueryParam("u")
	if tUrl == "" {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	// 本站
	if strings.Contains(tUrl, logic.WebsiteSetting.Domain) {
		return ctx.Redirect(http.StatusSeeOther, tUrl)
	}

	if strings.Contains(tUrl, "?") {
		tUrl += "&"
	} else {
		tUrl += "?"
	}
	tUrl += "utm_campaign=studygolang.com&utm_medium=studygolang.com&utm_source=studygolang.com"

	if CheckIsHttps(ctx) {
		return ctx.Redirect(http.StatusSeeOther, tUrl)
	}

	var (
		pUrl *url.URL
		err  error
	)

	if pUrl, err = url.Parse(tUrl); err != nil {
		return ctx.Redirect(http.StatusSeeOther, tUrl)
	}

	iframeDeny := config.ConfigFile.MustValue("crawl", "iframe_deny")
	// 检测是否禁止了 iframe 加载
	// 看是否在黑名单中
	for _, denyHost := range strings.Split(iframeDeny, ",") {
		if strings.Contains(pUrl.Host, denyHost) {
			return ctx.Redirect(http.StatusSeeOther, tUrl)
		}
	}

	// 检测会比较慢，进行异步检测，记录下来，以后分析再加黑名单
	go func() {
		resp, err := http.Head(tUrl)
		if err != nil {
			logger.Errorln("[iframe] head url:", tUrl, "error:", err)
			return
		}
		defer resp.Body.Close()
		if resp.Header.Get("X-Frame-Options") != "" {
			logger.Errorln("[iframe] deny:", tUrl)
			return
		}
	}()

	return render(ctx, "wr.html", map[string]interface{}{"url": tUrl})
}

// PkgdocHandler Go 语言文档中文版
func (IndexController) Pkgdoc(ctx echo.Context) error {
	// return render(ctx, "pkgdoc.html", map[string]interface{}{"activeDoc": "active"})
	tpl, err := template.ParseFiles(config.TemplateDir + "pkgdoc.html")
	if err != nil {
		logger.Errorln("parse file error:", err)
		return err
	}

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, nil)
	if err != nil {
		logger.Errorln("execute template error:", err)
		return err
	}

	return ctx.HTML(http.StatusOK, buf.String())
}

func (IndexController) Markdown(ctx echo.Context) error {
	return render(ctx, "markdown.html", nil)
}

// Link 用于重定向外部链接，比如广告链接
func (IndexController) Link(ctx echo.Context) error {
	tUrl := ctx.QueryParam("url")
	if strings.Contains(tUrl, "?") {
		tUrl += "&"
	} else {
		tUrl += "?"
	}
	tUrl += "utm_campaign=studygolang.com&utm_medium=studygolang.com&utm_source=studygolang.com"
	return ctx.Redirect(http.StatusSeeOther, tUrl)
}
