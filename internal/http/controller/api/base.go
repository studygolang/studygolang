// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

// Package api 提供面向 Next.js 前端的 REST API（Web 端）
package api

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/studygolang/studygolang/context"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/config"
	"github.com/polaris1119/logger"
)

const (
	perPage        = 20            // 每页默认条数
	authCookieName = "sg_token"    // HttpOnly 认证 Cookie 名
	cookieMaxAge   = 7 * 24 * 3600 // Cookie 有效期 7 天
)

// normalizeReplies 将 FindObjComments 返回的评论 map 规整为前端期望的结构。
//
// 注意 key 的大小写：cmd/studygolang/main.go 的 init() 设置了
// structs.DefaultTagName = "json"，因此 structs.Map 产出的是 json tag（小写），
// 不是 Go 字段名。按大写字段名取值会全部拿到 nil，评论的 id/uid/floor/ctime
// 就都丢了（点赞、回复、编辑全部失效）。
// uid > 0 时额外补上当前用户对每条评论的点赞状态，避免刷新后按钮状态丢失。
func normalizeReplies(ctx echo.Context, replies []map[string]interface{}, uid int) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(replies))
	cids := make([]int, 0, len(replies))
	for _, reply := range replies {
		item := map[string]interface{}{
			"id":          reply["cid"],
			"objid":       reply["objid"],
			"objtype":     reply["objtype"],
			"uid":         reply["uid"],
			"floor":       reply["floor"],
			"flag":        reply["flag"],
			"likenum":     reply["likenum"],
			"ctime":       reply["ctime"],
			"content":     reply["content"],
			"raw_content": reply["raw_content"],
		}
		if user, ok := reply["user"].(*model.User); ok && user != nil {
			item["name"] = user.Username
			item["avatar"] = user.Avatar
			item["username"] = user.Username
		}
		if cid, ok := reply["cid"].(int); ok && cid > 0 {
			cids = append(cids, cid)
		}
		result = append(result, item)
	}

	if uid > 0 && len(cids) > 0 {
		likeFlags := logic.DefaultLike.FindUserLikeFlagsByIds(
			context.EchoContext(ctx), uid, model.TypeComment, cids,
		)
		for _, item := range result {
			cid, _ := item["id"].(int)
			item["liked"] = likeFlags[cid] == model.FlagLike
		}
	}

	return result
}

// getAuthToken 读取认证 token：优先读 HttpOnly Cookie，回退到 X-Token header（兼容旧客户端）
func getAuthToken(ctx echo.Context) string {
	if cookie, err := ctx.Cookie(authCookieName); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	return ctx.Request().Header.Get("X-Token")
}

// isSecure 判断当前请求是否为 HTTPS（根据反向代理 header 或环境变量）
func isSecure(ctx echo.Context) bool {
	if ctx.Request().Header.Get("X-Forwarded-Proto") == "https" {
		return true
	}
	return config.ConfigFile.MustValue("global", "env", "prod") == "prod"
}

// setAuthCookie 设置认证 Cookie（HttpOnly, Secure, SameSite=Lax）
//
// 必须用 Lax 而不是 Strict：OAuth 回调是从 github.com/gitea 跳回来的跨站顶层
// 导航，Strict 下浏览器不携带 Cookie，optionalAuth 取不到当前用户，
// 「绑定到已有账号」会静默退化成「登录/新建账号」。
// Lax 已经能挡住跨站 POST（CSRF 的主要载体），另有 originCheck 二次校验。
func setAuthCookie(ctx echo.Context, token string) {
	cookie := new(http.Cookie)
	cookie.Name = authCookieName
	cookie.Value = token
	cookie.HttpOnly = true
	cookie.Secure = isSecure(ctx) // 根据 X-Forwarded-Proto 或环境动态决定
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Path = "/"
	cookie.MaxAge = cookieMaxAge
	ctx.SetCookie(cookie)
}

// clearAuthCookie 清除认证 Cookie
func clearAuthCookie(ctx echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = authCookieName
	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.MaxAge = -1
	// SameSite 与 setAuthCookie 保持一致
	cookie.SameSite = http.SameSiteLaxMode
	ctx.SetCookie(cookie)
}

// originCheck 检查写操作请求来源是否合法（CSRF 防护）
// GET/HEAD/OPTIONS 请求直接放行；POST/PUT/DELETE 等写操作检查 Origin 或 Referer
func originCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		method := ctx.Request().Method
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			return next(ctx)
		}

		origin := ctx.Request().Header.Get("Origin")
		referer := ctx.Request().Header.Get("Referer")

		// Origin 和 Referer 都没有时，拒绝写操作（CSRF 防护）
		// 前端所有写请求（fetch/XHR）都会自动带上 Origin，表单提交也会带上 Referer
		if origin == "" && referer == "" {
			return fail(ctx, "缺少 Origin 或 Referer", 403)
		}

		allowed := getAllowedOrigins()

		// 优先检查 Origin
		if origin != "" {
			for _, a := range allowed {
				if origin == a {
					return next(ctx)
				}
			}
			return ctx.JSON(http.StatusForbidden, map[string]interface{}{
				"code":    1,
				"message": "非法的请求来源",
			})
		}

		// Referer 检查
		if referer != "" {
			if refURL, err := url.Parse(referer); err == nil {
				refOrigin := refURL.Scheme + "://" + refURL.Host
				for _, a := range allowed {
					if refOrigin == a {
						return next(ctx)
					}
				}
			}
			// Referer 存在但不在白名单中，拒绝（CSRF 防护）
			// 之前这里 fallthrough 到 next(ctx) 等于绕过校验
			return ctx.JSON(http.StatusForbidden, map[string]interface{}{
				"code":    1,
				"message": "非法的请求来源",
			})
		}

		return next(ctx)
	}
}

func getLogger(ctx echo.Context) *logger.Logger {
	return logic.GetLogger(context.EchoContext(ctx))
}

// success 返回成功响应
func success(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "ok",
		"data": data,
	})
}

// errResponseWritten 表示错误响应已由 fail() 写出，调用方必须立即返回。
//
// 不能直接把 ctx.JSON 的返回值当错误：Echo 写入成功时返回 nil，
// 于是 `return 0, fail(...)` 实际是 (0, nil)，调用方的 `if err != nil` 守卫全部失效，
// 未登录请求会继续执行业务逻辑（甚至 nil deref requireAuth 返回的 me）。
// middleware.HTTPError 见到 Response().Committed 为 true 会吞掉它，不会二次写入。
var errResponseWritten = errors.New("api: error response already written")

// fail 返回错误响应
func fail(ctx echo.Context, msg string, codes ...int) error {
	code := 1
	if len(codes) > 0 {
		code = codes[0]
	}

	getLogger(ctx).Errorln("api fail:", msg)

	if err := ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": code,
		"msg":  msg,
	}); err != nil {
		return err
	}
	return errResponseWritten
}

// parseAuthUID 从请求中获取并验证 token，返回 uid。
// 支持 JWT 和旧版 MD5 Token 双轨验证。
// 注意：本函数不校验 user.Status，仅适用于读操作。
// 写操作请用 parseActiveAuthUID 或 requireAuth，以拦截冻结/未激活用户。
func parseAuthUID(ctx echo.Context) (int, error) {
	token := getAuthToken(ctx)
	if token == "" {
		return 0, fail(ctx, "请先登录", NeedReLoginCode)
	}

	uid, _, valid := ValidateTokenAuto(token)
	if !valid {
		return 0, fail(ctx, "token 已过期，请重新登录", NeedReLoginCode)
	}
	if uid == 0 {
		return 0, fail(ctx, "无效的 token", NeedReLoginCode)
	}

	return uid, nil
}

// optionalAuth 返回当前登录用户信息；未登录或 token 无效时返回 nil（不写错误响应）。
//
// 使用场景：GET 接口对登录状态可选（登录用户带额外字段，未登录用户看公开数据）。
// 不能用 requireAuth 替代——requireAuth 失败时会调用 fail() 写错误响应，
// 之后再调用 success() 会触发 "HTTP response already written" 双写问题，
// 导致未登录用户拿到 "请先登录" 而不是公开数据。
//
// 实现上不调用 fail()，仅返回 nil 让调用方自行处理。
func optionalAuth(ctx echo.Context) *model.Me {
	me := meFromToken(ctx)
	// 可选鉴权同样不应让冻结用户享有"登录用户"特权（如兑换状态等）
	if me == nil || me.Status != model.UserStatusAudit {
		return nil
	}
	return me
}

// oauthBindUser 返回 OAuth「绑定到当前账号」场景下的用户。
//
// 与 optionalAuth 的区别：允许未激活（UserStatusNoAudit）用户绑定。刚注册、
// 还没点激活邮件的用户正是最需要绑定第三方账号的人；若沿用 optionalAuth，
// 他们会拿到 nil 从而掉进登录分支，静默变成新建/登录另一个账号。
// 冻结/停用用户仍然拒绝。
func oauthBindUser(ctx echo.Context) *model.Me {
	me := meFromToken(ctx)
	if me == nil || me.Status > model.UserStatusAudit {
		return nil
	}
	return me
}

// meFromToken 解析 token 得到当前用户，不做状态校验、不写错误响应。
// 状态策略由调用方决定（optionalAuth / oauthBindUser）。
func meFromToken(ctx echo.Context) *model.Me {
	token := getAuthToken(ctx)
	if token == "" {
		return nil
	}
	uid, _, valid := ValidateTokenAuto(token)
	if !valid || uid == 0 {
		return nil
	}
	userInfo := logic.GetOrFetchUserInfo(uid, func() *logic.UserInfoCache {
		return fetchFullUserInfo(ctx, uid)
	})
	if userInfo == nil {
		return nil
	}
	return &model.Me{
		Uid:      userInfo.Uid,
		Username: userInfo.Username,
		Email:    userInfo.Email,
		IsRoot:   userInfo.IsRoot,
		IsAdmin:  userInfo.IsAdmin,
		IsVip:    userInfo.IsVip,
		Balance:  userInfo.Balance,
		Status:   userInfo.Status,
	}
}

// fetchFullUserInfo 从 DB 读取用户信息并构造完整的 UserInfoCache。
// parseActiveAuthUID 与 requireAuth 必须使用同一个 fetchFunc 填充缓存，
// 否则先调用的最小化版本会把 Username/Email/IsVip/Avatar/Balance 留空，
// 后续读取同一缓存的 requireAuth 拿到不完整的 Me，导致 balanceCheck 误判、
// 通知邮件 username 为空、PermissionPay 对 VIP 用户错误隐藏内容等。
func fetchFullUserInfo(ctx echo.Context, uid int) *logic.UserInfoCache {
	user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "uid", uid)
	if user == nil || user.Uid == 0 {
		return nil
	}
	// IsAdmin 必须基于 user_role 表判断（AdminMinRoleId=7），
	// 不能简化为 IsRoot，否则板块管理员/晨读管理员等角色会丢失权限。
	// FindOne 已填充 user.Roleids，UserLogic.IsAdmin 据此判断。
	return &logic.UserInfoCache{
		Uid:      user.Uid,
		Username: user.Username,
		Email:    user.Email,
		IsRoot:   user.IsRoot,
		IsAdmin:  logic.DefaultUser.IsAdmin(user),
		IsVip:    user.IsVip,
		Avatar:   user.Avatar,
		Balance:  user.Balance,
		Status:   user.Status,
	}
}

// parseActiveAuthUID 在 parseAuthUID 基础上额外校验用户状态。
// 用于写操作（点赞/收藏/发消息/改资料等），与 master 的 NeedLogin 中间件一致：
// 仅 UserStatusAudit（已激活）允许；冻结/未激活/拒绝/停用 均拒绝。
//
// 为什么需要这个函数：requireAuth 会返回完整 Me（带 Balance/IsAdmin 等），
// 但有些写 handler 只需要 uid（如 like/favorite toggle）。直接用 parseAuthUID
// 又会绕过状态校验，导致冻结用户仍可写。本函数补齐这一缺口。
func parseActiveAuthUID(ctx echo.Context) (int, error) {
	uid, err := parseAuthUID(ctx)
	if err != nil {
		return 0, err
	}

	userInfo := logic.GetOrFetchUserInfo(uid, func() *logic.UserInfoCache {
		return fetchFullUserInfo(ctx, uid)
	})

	if userInfo == nil {
		return 0, fail(ctx, "用户不存在")
	}
	if userInfo.Status != model.UserStatusAudit {
		return 0, fail(ctx, "账号已被冻结或未激活，请重新登录", NeedReLoginCode)
	}

	return uid, nil
}
func requireAuth(ctx echo.Context) (*model.Me, error) {
	token := getAuthToken(ctx)
	if token == "" {
		return nil, fail(ctx, "请先登录", NeedReLoginCode)
	}

	// 使用统一验证方法，自动识别 JWT 和旧 MD5 Token
	uid, _, valid := ValidateTokenAuto(token)
	if !valid {
		return nil, fail(ctx, "token 已过期，请重新登录", NeedReLoginCode)
	}

	// 优先从缓存获取用户信息（与 parseActiveAuthUID 共享 fetchFullUserInfo，
	// 避免任一路径用最小字段集污染缓存导致另一路径拿到不完整 Me）
	userInfo := logic.GetOrFetchUserInfo(uid, func() *logic.UserInfoCache {
		return fetchFullUserInfo(ctx, uid)
	})

	if userInfo == nil {
		return nil, fail(ctx, "用户不存在")
	}

	// 状态校验：与 master 的 NeedLogin 中间件保持一致。
	// 仅 UserStatusAudit（已激活）允许写操作；UserStatusOutage（冻结）等需重新登录。
	// master 在中间件层统一拦截；refactor 把 requireAuth 作为唯一认证入口，故在此检查。
	if userInfo.Status != model.UserStatusAudit {
		return nil, fail(ctx, "账号已被冻结或未激活，请重新登录", NeedReLoginCode)
	}

	return &model.Me{
		Uid:      userInfo.Uid,
		Username: userInfo.Username,
		IsRoot:   userInfo.IsRoot,
		IsAdmin:  userInfo.IsAdmin,
		IsVip:    userInfo.IsVip,
		Balance:  userInfo.Balance,
		Status:   userInfo.Status,
	}, nil
}

// isValidObjType 校验对象类型是否在合法枚举内。
// TypeTopic=0 是个陷阱（不能简单用 objtype==0 判断无效）。
func isValidObjType(objtype int) bool {
	switch objtype {
	case model.TypeTopic, model.TypeArticle, model.TypeResource,
		model.TypeWiki, model.TypeProject, model.TypeBook, model.TypeInterview:
		return true
	}
	return false
}

// maxSidebarLimit 侧边栏/列表类接口允许的 limit 上限。
// 防止恶意构造 limit=N（如 100000）触发 DB 全表扫描 + 大对象序列化 DoS。
const maxSidebarLimit = 50

// boundLimit 将客户端传入的 limit 规整到 [1, maxSidebarLimit]。
// def 为 limit 缺省/非法时的默认值（会被 clamp 到上限内）。
func boundLimit(limit, def int) int {
	if limit <= 0 {
		return def
	}
	if limit > maxSidebarLimit {
		return maxSidebarLimit
	}
	return limit
}

// maxPageNumber 分页接口允许的最大页码。
// 防止恶意构造 p=N（如 100000）触发搜索引擎/DB 深度分页 DoS：
// - Solr/Lucene 的 start+offset 模式在深分页时性能急剧下降
// - MySQL OFFSET 在大偏移量下需要扫描跳过大量行
// 100 页 * 50 行/页 = 5000 条，对社区搜索已足够；超出视为越界，回退到第 1 页。
const maxPageNumber = 100

// boundPage 将客户端传入的页码规整到 [1, maxPageNumber]。
// 缺省/非法（非正整数）时返回 1；超过上限时回退到 1（视为越界，不展示数据）。
func boundPage(pageStr string) int {
	p := 1
	if pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil && v > 0 {
			p = v
		}
	}
	if p > maxPageNumber {
		// 越界页码回退到 1，避免向搜索引擎/DB 传递超大 offset
		return 1
	}
	return p
}

// sanitizeUserForPublic 对外公开接口返回 User 前的字段脱敏。
//
// master 模板 profile.html 中 Email 显示受 `{{if .user.Open}}` 控制：
// 仅当用户勾选"公开 Email"(Open==1) 时才在主页展示邮箱。
// 重构为 JSON API 后，model.User 全字段带 JSON tag 自动序列化，
// 若不显式过滤，不公开邮箱的用户主页 /api/v1/user/:username 与
// 用户列表 /api/v1/users 都会把邮箱泄露给任意访客（含未登录用户）。
//
// 规则：Open==0 时清空 Email，与 master 模板语义对齐。
// 就地修改（不复制）以避免大量分配；调用方已不再使用原 Email。
func sanitizeUserForPublic(user *model.User) *model.User {
	if user == nil {
		return nil
	}
	if user.Open == 0 {
		user.Email = ""
	}
	return user
}

// sanitizeUsersForPublic 批量脱敏，便于 List 类接口使用。
func sanitizeUsersForPublic(users []*model.User) []*model.User {
	for _, u := range users {
		sanitizeUserForPublic(u)
	}
	return users
}

// sanitizeFeedForPublic 对 Feed 的内嵌 User / Lastreplyuser 做脱敏。
//
// Feed struct 无 JSON tag，User / Lastreplyuser 会带 Email 等敏感字段
// 自动序列化。master 模板通过 `{{if .user.Open}}` 控制邮箱展示，
// 这里在 API 边界统一脱敏，未勾选"公开 Email"的用户清空 Email。
func sanitizeFeedForPublic(feed *model.Feed) *model.Feed {
	if feed == nil {
		return nil
	}
	sanitizeUserForPublic(feed.User)
	sanitizeUserForPublic(feed.Lastreplyuser)
	return feed
}

// sanitizeFeedsForPublic 批量脱敏 Feed 切片。
func sanitizeFeedsForPublic(feeds []*model.Feed) []*model.Feed {
	for _, f := range feeds {
		sanitizeFeedForPublic(f)
	}
	return feeds
}

// sanitizeSubjectForPublic 对 Subject 的内嵌 User 做脱敏。
// Subject.User 字段带 `json:"user"`，会自动序列化暴露 Email。
func sanitizeSubjectForPublic(subject *model.Subject) *model.Subject {
	if subject == nil {
		return nil
	}
	sanitizeUserForPublic(subject.User)
	return subject
}

// sanitizeSubjectsForPublic 批量脱敏 Subject 切片。
func sanitizeSubjectsForPublic(subjects []*model.Subject) []*model.Subject {
	for _, s := range subjects {
		sanitizeSubjectForPublic(s)
	}
	return subjects
}

// sanitizeProjectForPublic 对 OpenProject 的内嵌 User / LastReplyUser 做脱敏。
// 两个字段都带 JSON tag，会自动序列化暴露 Email。
func sanitizeProjectForPublic(project *model.OpenProject) *model.OpenProject {
	if project == nil {
		return nil
	}
	sanitizeUserForPublic(project.User)
	sanitizeUserForPublic(project.LastReplyUser)
	return project
}

// sanitizeProjectsForPublic 批量脱敏 OpenProject 切片。
func sanitizeProjectsForPublic(projects []*model.OpenProject) []*model.OpenProject {
	for _, p := range projects {
		sanitizeProjectForPublic(p)
	}
	return projects
}

// sanitizeSubjectFollowersForPublic 批量脱敏 SubjectFollower。
// SubjectFollower.User 无 JSON tag，Go 默认按字段名序列化，
// 会暴露 Email。
func sanitizeSubjectFollowersForPublic(followers []*model.SubjectFollower) []*model.SubjectFollower {
	for _, f := range followers {
		if f != nil {
			sanitizeUserForPublic(f.User)
		}
	}
	return followers
}
