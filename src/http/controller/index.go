// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"bytes"
	"html/template"
	"logic"
	"math/rand"
	"model"
	"net/http"
	"net/url"
	"strings"

	. "http"

	"github.com/labstack/echo"
	"github.com/polaris1119/config"
	"github.com/polaris1119/logger"
)

type IndexController struct{}

// 注册路由
func (self IndexController) RegisterRoute(g *echo.Group) {
	g.GET("/", self.NewIndex)
	g.GET("/wr", self.WrapUrl)
	g.GET("/pkgdoc", self.Pkgdoc)
}

func (IndexController) NewIndex(ctx echo.Context) error {
	tab := ctx.QueryParam("tab")
	if tab == "" {
		tab = GetFromCookie(ctx, "INDEX_TAB")
	}

	if tab == "" {
		tab = logic.WebsiteSetting.IndexNavs[0].Tab
	}
	SetCookie(ctx, "INDEX_TAB", tab)

	data := logic.DefaultIndex.FindData(ctx, tab)
	data["all_nodes"] = logic.GenNodes()

	return render(ctx, "new_index.html", data)
}

// Index 首页
func (IndexController) Index(ctx echo.Context) error {
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
		curUser, ok := ctx.Get("user").(*model.Me)
		if ok {
			likeFlags, _ = logic.DefaultLike.FindUserLikeObjects(ctx, curUser.Uid, model.TypeArticle, recentArticles[0].Id, recentArticles[len(recentArticles)-1].Id)
		}
	}

	// 资源
	resources := logic.DefaultResource.FindBy(ctx, 10)

	books := logic.DefaultGoBook.FindBy(ctx, 24)
	if len(books) > 8 {
		bookNum := 8
		bookStart := rand.Intn(len(books) - bookNum)
		books = books[bookStart : bookStart+bookNum]
	}

	// 学习资料
	materials := logic.DefaultLearningMaterial.FindAll(ctx)

	return render(ctx, "index.html",
		map[string]interface{}{
			"topics":    topicsList,
			"articles":  recentArticles,
			"likeflags": likeFlags,
			"resources": resources,
			"books":     books,
			"materials": materials,
		})
}

// WrapUrl 包装链接
func (IndexController) WrapUrl(ctx echo.Context) error {
	tUrl := ctx.QueryParam("u")
	if tUrl == "" {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	if pUrl, err := url.Parse(tUrl); err != nil {
		return ctx.Redirect(http.StatusSeeOther, tUrl)
	} else {
		if !pUrl.IsAbs() {
			return ctx.Redirect(http.StatusSeeOther, tUrl)
		}

		// 本站
		if strings.Contains(pUrl.Host, logic.WebsiteSetting.Domain) {
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
	}

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
