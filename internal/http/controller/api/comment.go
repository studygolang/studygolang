// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"net/http"
	"time"

	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type CommentController struct{}

func (self CommentController) RegisterRoute(g *echo.Group) {
	g.GET("/comments", self.List)
	g.GET("/comments/:cid/detail", self.Detail)
	g.POST("/comments/:objid", self.Create)
	g.PUT("/comments/:cid", self.Modify)
	g.GET("/at/users", self.AtUsers)
}

// List 评论列表（objid、objtype 查询参数）
func (CommentController) List(ctx echo.Context) error {
	objid := goutils.MustInt(ctx.QueryParam("objid"))
	objtype := goutils.MustInt(ctx.QueryParam("objtype"))

	if objid == 0 || !isValidObjType(objtype) {
		return fail(ctx, "参数有误")
	}

	comments, _, _ := logic.DefaultComment.FindObjComments(
		context.EchoContext(ctx), objid, objtype, 0, 0,
	)

	return success(ctx, map[string]interface{}{
		"comments": normalizeReplies(comments),
	})
}

// Detail 评论详情（包含附近评论）
func (CommentController) Detail(ctx echo.Context) error {
	cid := goutils.MustInt(ctx.Param("cid"))
	objid := goutils.MustInt(ctx.QueryParam("objid"))
	objtype := goutils.MustInt(ctx.QueryParam("objtype"))

	if cid == 0 || objid == 0 {
		return fail(ctx, "参数有误")
	}

	// 获取当前评论和附近评论（2条）
	comment, nearbyComments := logic.DefaultComment.FindComment(context.EchoContext(ctx), cid, objid, objtype)

	if comment.Cid == 0 {
		return fail(ctx, "评论不存在")
	}

	// 获取所有相关用户信息
	uids := []int{comment.Uid}
	for _, c := range nearbyComments {
		uids = append(uids, c.Uid)
	}
	users := logic.DefaultUser.FindUserInfos(context.EchoContext(ctx), uids)

	return success(ctx, map[string]interface{}{
		"comment":         comment,
		"nearby_comments": nearbyComments,
		"users":           users,
	})
}

// Create 创建评论（支持 Cookie 和 X-Token header）
func (CommentController) Create(ctx echo.Context) error {
	token := getAuthToken(ctx)
	if token == "" {
		return fail(ctx, "未登录", NeedReLoginCode)
	}

	uid, _, valid := ValidateTokenAuto(token)
	if !valid || uid == 0 {
		return fail(ctx, "token 已过期，请重新登录", NeedReLoginCode)
	}

	objid := goutils.MustInt(ctx.Param("objid"))
	if objid == 0 {
		return fail(ctx, "参数有误")
	}

	forms, _ := ctx.FormParams()
	objtype := goutils.MustInt(forms.Get("objtype"))
	// 校验 objtype 必须在合法枚举内：
	//   1) 避免 Publish 内部把缺省/非法值写入 comment.objtype，造成评论挂错对象或孤儿评论
	//   2) 保护楼层计数器（key 含 objtype），防止恶意构造的 objtype 污染计数器空间
	if !isValidObjType(objtype) {
		return fail(ctx, "objtype 参数非法")
	}

	// 获取完整用户信息（余额检查等需要 Balance 字段）
	me := &model.Me{Uid: uid}
	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", uid)
	if user == nil || user.Uid == 0 {
		return fail(ctx, "用户不存在")
	}
	me.Username = user.Username
	me.Balance = user.Balance
	me.IsRoot = user.IsRoot
	// IsAdmin 必须基于 user_role 表判断（AdminMinRoleId=7），不能简化为 IsRoot，
	// 否则板块管理员/晨读管理员等角色会丢失权限（见 base.go requireAuth 同款修复）
	me.IsAdmin = logic.DefaultUser.IsAdmin(user)
	me.CreatedAt = time.Time(user.Ctime)

	// 敏感词检查
	if !sensitiveCheck(ctx, me) {
		return failSensitive(ctx)
	}
	// 余额检查（评论要求余额 >= 5）
	if !balanceCheck(me, true) {
		return failBalance(ctx)
	}

	comment, err := logic.DefaultComment.Publish(context.EchoContext(ctx), uid, objid, forms)
	if err != nil {
		return fail(ctx, "发布评论失败", 2)
	}

	// 发布后邮件通知站长
	publishNotice(ctx, me)

	return success(ctx, map[string]interface{}{
		"comment": comment,
	})
}

// AtUsers 获取可 @ 的用户列表（评论 @ 自动补全）
func (CommentController) AtUsers(ctx echo.Context) error {
	term := ctx.QueryParam("term")
	if term == "" {
		return ctx.JSON(http.StatusOK, []map[string]string{})
	}
	users := logic.DefaultUser.GetUserMentions(term, 10, false)
	if users == nil {
		users = make([]map[string]string, 0)
	}
	return ctx.JSON(http.StatusOK, users)
}

// Modify 修改评论（需要登录，且只有作者可修改）
func (CommentController) Modify(ctx echo.Context) error {
	uid, err := parseAuthUID(ctx)
	if err != nil {
		return err
	}

	cid := goutils.MustInt(ctx.Param("cid"))
	if cid == 0 {
		return fail(ctx, "参数有误")
	}

	content := ctx.FormValue("content")
	if content == "" {
		return fail(ctx, "评论内容不能为空")
	}

	comment, findErr := logic.DefaultComment.FindById(cid)
	if findErr != nil {
		return fail(ctx, "评论不存在")
	}

	// 使用 CanEdit 进行权限校验（包含时间限制检查）
	me := &model.Me{Uid: uid}
	// 补充完整用户信息（CanEdit 需要 IsAdmin/IsRoot）
	if user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", uid); user != nil && user.Uid > 0 {
		me.Username = user.Username
		me.IsRoot = user.IsRoot
		// IsAdmin 必须基于 user_role 表判断（见 base.go requireAuth 同款修复）
		me.IsAdmin = logic.DefaultUser.IsAdmin(user)
		me.CreatedAt = time.Time(user.Ctime)
	}
	if !logic.CanEdit(me, comment) {
		return fail(ctx, "没有修改权限")
	}

	// 敏感词检查
	if !sensitiveCheck(ctx, me) {
		return failSensitive(ctx)
	}

	errMsg, modifyErr := logic.DefaultComment.Modify(context.EchoContext(ctx), cid, content)
	if modifyErr != nil {
		return fail(ctx, errMsg)
	}

	return success(ctx, map[string]interface{}{"cid": cid})
}
