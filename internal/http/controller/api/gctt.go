// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"

	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/config"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
)

type GCTTController struct{}

func (self GCTTController) RegisterRoute(g *echo.Group) {
	g.GET("/gctt", self.Index)
	g.GET("/gctt/users", self.Users)
	g.GET("/gctt/issues", self.Issues)
	g.GET("/gctt/me", self.Me)
	g.POST("/gctt/apply", self.Apply)
	g.POST("/gctt/articles", self.Publish)
	g.GET("/gctt/:username", self.UserDetail)
	g.POST("/gctt/webhook", self.Webhook)
}

// Index GCTT 首页（时间线 + 核心用户 + 未翻译 Issues）
func (GCTTController) Index(ctx echo.Context) error {
	timeLines := logic.DefaultGCTT.FindTimeLines(context.EchoContext(ctx))
	coreUsers := logic.DefaultGCTT.FindCoreUsers(context.EchoContext(ctx))
	untranslatedIssues := logic.DefaultGCTT.FindUnTranslateIssues(context.EchoContext(ctx), 10)

	return success(ctx, map[string]interface{}{
		"time_lines":          timeLines,
		"core_users":          coreUsers,
		"untranslated_issues": untranslatedIssues,
	})
}

// Users GCTT 译者列表
func (GCTTController) Users(ctx echo.Context) error {
	users := logic.DefaultGCTT.FindUsers(context.EchoContext(ctx))
	return success(ctx, map[string]interface{}{
		"users": users,
	})
}

// Issues GCTT Issue 列表（分页）
func (GCTTController) Issues(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	if curPage < 1 {
		curPage = 1
	}
	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)

	state := ctx.QueryParam("state")
	queryStr := ""
	var args []interface{}
	if state == "open" {
		queryStr = "state=?"
		args = append(args, 0)
	} else if state == "closed" {
		queryStr = "state=?"
		args = append(args, 1)
	} else {
		// 默认查询未关闭的 issue
		queryStr = "state=?"
		args = append(args, 0)
	}

	issues := logic.DefaultGCTT.FindIssues(context.EchoContext(ctx), paginator, queryStr, args...)
	total := logic.DefaultGCTT.IssueCount(context.EchoContext(ctx), queryStr, args...)
	hasMore := paginator.SetTotal(total).HasMorePage()

	return success(ctx, map[string]interface{}{
		"issues":   issues,
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
	})
}

// UserDetail GCTT 译者详情
func (GCTTController) UserDetail(ctx echo.Context) error {
	username := ctx.Param("username")
	if username == "" {
		return fail(ctx, "用户名不能为空")
	}

	gcttUser := logic.DefaultGCTT.FindOne(context.EchoContext(ctx), username)
	if gcttUser == nil || gcttUser.Id == 0 {
		return fail(ctx, "用户不存在")
	}

	return success(ctx, map[string]interface{}{
		"user": gcttUser,
	})
}

// Apply 申请成为译者
func (GCTTController) Apply(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	// 检查是否已是译者
	gcttUser := logic.DefaultGCTT.FindTranslator(context.EchoContext(ctx), me)
	if gcttUser != nil && gcttUser.Id > 0 {
		return success(ctx, map[string]interface{}{
			"gctt_user": gcttUser,
		})
	}

	// 检查是否绑定 GitHub 账号
	var githubUser *model.BindUser
	bindUsers := logic.DefaultUser.FindBindUsers(context.EchoContext(ctx), me.Uid)
	for _, bindUser := range bindUsers {
		if bindUser.Type == model.BindTypeGithub {
			githubUser = bindUser
			break
		}
	}

	if githubUser == nil {
		return fail(ctx, "请先绑定 GitHub 账号")
	}

	// 查找 GitHub 用户名对应的 GCTT 记录
	existingGCTTUser := logic.DefaultGCTT.FindOne(context.EchoContext(ctx), githubUser.Username)
	if existingGCTTUser == nil || existingGCTTUser.Id == 0 {
		return fail(ctx, "未找到对应的 GCTT 记录")
	}

	// 绑定用户
	err = logic.DefaultGCTT.BindUser(context.EchoContext(ctx), existingGCTTUser, me.Uid, githubUser)
	if err != nil {
		return fail(ctx, "绑定用户失败")
	}

	return success(ctx, map[string]interface{}{
		"gctt_user": existingGCTTUser,
	})
}

// Publish 发布译文
func (GCTTController) Publish(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	title := ctx.FormValue("title")
	content := ctx.FormValue("content")

	if title == "" {
		return fail(ctx, "标题不能为空")
	}
	if content == "" {
		return fail(ctx, "内容不能为空")
	}

	// 验证译者身份
	gcttUser := logic.DefaultGCTT.FindTranslator(context.EchoContext(ctx), me)
	if gcttUser == nil || gcttUser.Id == 0 {
		return fail(ctx, "非 GCTT 译者,不允许发布")
	}

	// 获取表单参数并设置 GCTT 标记
	forms, _ := ctx.FormParams()
	forms.Set("gctt", "true")

	// 发布文章
	articleId, err := logic.DefaultArticle.Publish(context.EchoContext(ctx), me, forms)
	if err != nil {
		return fail(ctx, "发布失败: "+err.Error())
	}

	return success(ctx, map[string]interface{}{
		"id": articleId,
	})
}

// Me 当前用户 GCTT 信息
func (GCTTController) Me(ctx echo.Context) error {
	me, err := requireAuth(ctx)
	if err != nil {
		return err
	}

	gcttUser := logic.DefaultGCTT.FindTranslator(context.EchoContext(ctx), me)
	isTranslator := gcttUser != nil && gcttUser.Id > 0

	return success(ctx, map[string]interface{}{
		"gctt_user":     gcttUser,
		"is_translator": isTranslator,
	})
}

// Webhook 处理 GitHub Webhook 事件（pull_request / issue_comment / issues）
func (GCTTController) Webhook(ctx echo.Context) error {
	// 限制 body 大小为 1MB，防止内存耗尽攻击
	body, err := io.ReadAll(io.LimitReader(ctx.Request().Body, 1<<20))
	if err != nil {
		logger.Errorln("GCTTController Webhook read body error:", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "read body failed"})
	}
	// 恢复 body 供下游逻辑读取
	ctx.Request().Body = io.NopCloser(bytes.NewReader(body))

	header := ctx.Request().Header

	tokenSecret := config.ConfigFile.MustValue("gctt", "token_secret")
	if tokenSecret == "" {
		logger.Errorln("GCTTController Webhook: token_secret not configured")
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "webhook not configured"})
	}
	if !checkMAC(body, header.Get("X-Hub-Signature"), []byte(tokenSecret)) {
		logger.Errorln("GCTTController Webhook checkMAC failed")
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid signature"})
	}

	event := header.Get("X-GitHub-Event")
	logger.Infoln("GCTTController Webhook event:", event)

	switch event {
	case "pull_request":
		return logic.DefaultGithub.PullRequestEvent(context.EchoContext(ctx), body)
	case "issue_comment":
		return logic.DefaultGithub.IssueCommentEvent(context.EchoContext(ctx), body)
	case "issues":
		return logic.DefaultGithub.IssueEvent(context.EchoContext(ctx), body)
	default:
		logger.Infoln("GCTTController Webhook unhandled event:", event)
	}

	return nil
}

// checkMAC 验证 GitHub Webhook HMAC-SHA1 签名
func checkMAC(message []byte, messageMAC string, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := fmt.Sprintf("sha1=%x", mac.Sum(nil))
	return hmac.Equal([]byte(messageMAC), []byte(expectedMAC))
}
