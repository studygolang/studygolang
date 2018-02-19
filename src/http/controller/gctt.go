// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"html/template"
	. "http"
	"http/middleware"
	"io/ioutil"
	"logic"
	"model"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/polaris1119/config"
	"github.com/polaris1119/echoutils"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
)

type GCTTController struct{}

// 注册路由
func (self GCTTController) RegisterRoute(g *echo.Group) {
	g.Get("/gctt", self.Index)
	g.Get("/gctt-list", self.UserList)
	g.Get("/gctt-issue", self.IssueList)
	g.Get("/gctt/:username", self.User)
	g.Get("/gctt-apply", self.Apply, middleware.NeedLogin())
	g.Match([]string{"GET", "POST"}, "/gctt-new", self.Create, middleware.NeedLogin())

	g.Post("/gctt-webhook", self.Webhook)
}

func (self GCTTController) Index(ctx echo.Context) error {
	gcttTimeLines := logic.DefaultGCTT.FindTimeLines(ctx)
	gcttUsers := logic.DefaultGCTT.FindCoreUsers(ctx)
	gcttIssues := logic.DefaultGCTT.FindUnTranslateIssues(ctx, 10)

	return Render(ctx, "gctt/index.html", map[string]interface{}{
		"time_lines": gcttTimeLines,
		"users":      gcttUsers,
		"issues":     gcttIssues,
	})
}

// Apply 申请成为译者
func (GCTTController) Apply(ctx echo.Context) error {
	me := ctx.Get("user").(*model.Me)
	gcttUser := logic.DefaultGCTT.FindTranslator(ctx, me)
	if gcttUser.Id > 0 {
		return ctx.Redirect(http.StatusSeeOther, "/gctt")
	}

	// 是否绑定了 github 账号
	var githubUser *model.BindUser
	bindUsers := logic.DefaultUser.FindBindUsers(ctx, me.Uid)
	for _, bindUser := range bindUsers {
		if bindUser.Type == model.BindTypeGithub {
			githubUser = bindUser
			break
		}
	}

	// 如果已经绑定，查看是否之前已经是译者
	if githubUser != nil {
		gcttUser = logic.DefaultGCTT.FindOne(ctx, githubUser.Username)
		logic.DefaultGCTT.BindUser(ctx, gcttUser, me.Uid, githubUser)
		return ctx.Redirect(http.StatusSeeOther, "/gctt")
	}

	return render(ctx, "gctt/apply.html", map[string]interface{}{
		"activeGCTT":  "active",
		"github_user": githubUser,
	})
}

// Create 发布新译文
func (GCTTController) Create(ctx echo.Context) error {
	me := ctx.Get("user").(*model.Me)
	gcttUser := logic.DefaultGCTT.FindTranslator(ctx, me)

	title := ctx.FormValue("title")
	if title == "" || ctx.Request().Method() != "POST" {
		return render(ctx, "gctt/new.html", map[string]interface{}{
			"activeGCTT": "active",
			"gctt_user":  gcttUser,
		})
	}

	if ctx.FormValue("content") == "" {
		return fail(ctx, 1, "内容不能为空")
	}

	if gcttUser == nil {
		return fail(ctx, 2, "不允许发布！")
	}

	id, err := logic.DefaultArticle.Publish(echoutils.WrapEchoContext(ctx), me, ctx.FormParams())
	if err != nil {
		return fail(ctx, 3, "内部服务错误")
	}

	return success(ctx, map[string]interface{}{"id": id})
}

func (GCTTController) User(ctx echo.Context) error {
	username := ctx.Param("username")
	if username == "" {
		return ctx.Redirect(http.StatusSeeOther, "/gctt")
	}

	gcttUser := logic.DefaultGCTT.FindOne(ctx, username)
	if gcttUser.Id == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/gctt")
	}

	joinDays := int(gcttUser.LastAt-gcttUser.JoinedAt)/86400 + 1
	avgDays := fmt.Sprintf("%.1f", float64(gcttUser.AvgTime)/86400.0)

	articles := logic.DefaultArticle.FindTaGCTTArticles(ctx, username)

	return render(ctx, "gctt/user-info.html", map[string]interface{}{
		"gctt_user": gcttUser,
		"articles":  articles,
		"join_days": joinDays,
		"avg_days":  avgDays,
	})
}

func (GCTTController) UserList(ctx echo.Context) error {
	users := logic.DefaultGCTT.FindUsers(ctx)

	num, words := 0, 0
	for _, user := range users {
		num += user.Num
		words += user.Words
	}

	prs := logic.DefaultGCTT.FindNewestGit(ctx)

	return render(ctx, "gctt/user-list.html", map[string]interface{}{
		"users": users,
		"num":   num,
		"words": words,
		"prs":   prs,
	})
}

func (GCTTController) IssueList(ctx echo.Context) error {
	querystring, arg := "", ""

	label := ctx.QueryParam("label")

	translator := ctx.QueryParam("translator")
	if translator != "" {
		querystring = "translator=?"
		arg = translator
	} else {
		if label == model.LabelUnClaim {
			querystring = "label=?"
			arg = label
		} else if label == model.LabelClaimed {
			querystring = "label=? AND state=" + strconv.Itoa(model.IssueOpened)
			arg = label
		}
	}

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginator(curPage)

	issues := logic.DefaultGCTT.FindIssues(ctx, paginator, querystring, arg)

	total := logic.DefaultGCTT.IssueCount(ctx, querystring, arg)
	pageHTML := paginator.SetTotal(total).GetPageHtml(ctx.Request().URL().Path())

	prs := logic.DefaultGCTT.FindNewestGit(ctx)

	return render(ctx, "gctt/issue-list.html", map[string]interface{}{
		"issues":     issues,
		"prs":        prs,
		"page":       template.HTML(pageHTML),
		"translator": translator,
		"label":      label,
		"total":      total,
	})
}

func (GCTTController) Webhook(ctx echo.Context) error {
	body, err := ioutil.ReadAll(Request(ctx).Body)
	if err != nil {
		logger.Errorln("GCTTController Webhook error:", err)
		return err
	}

	header := ctx.Request().Header()

	tokenSecret := config.ConfigFile.MustValue("gctt", "token_secret")
	ok := checkMAC(body, header.Get("X-Hub-Signature"), []byte(tokenSecret))
	if !ok {
		logger.Errorln("GCTTController Webhook checkMAC error", string(body))
		return nil
	}

	event := header.Get("X-GitHub-Event")
	logger.Infoln("GCTTController Webhook event:", event)
	switch event {
	case "pull_request":
		return logic.DefaultGithub.PullRequestEvent(ctx, body)
	case "issue_comment":
		return logic.DefaultGithub.IssueCommentEvent(ctx, body)
	case "issues":
		return logic.DefaultGithub.IssueEvent(ctx, body)
	default:
		fmt.Println("not deal event:", event)
	}

	return nil
}

func checkMAC(message []byte, messageMAC string, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := fmt.Sprintf("%x", mac.Sum(nil))
	return messageMAC == "sha1="+expectedMAC
}
