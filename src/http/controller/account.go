// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of self source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"html/template"
	. "http/internal/helper"
	"http/middleware"
	"logic"
	"model"
	"net/http"
	"net/url"
	"strings"
	"time"
	"util"

	. "http"

	"github.com/dchest/captcha"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/polaris1119/config"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
	guuid "github.com/twinj/uuid"
)

type AccountController struct{}

// 注册路由
func (self AccountController) RegisterRoute(g *echo.Group) {
	g.Any("/account/register", self.Register)
	g.Post("/account/send_activate_email", self.SendActivateEmail)
	g.Get("/account/activate", self.Activate)
	g.Any("/account/login", self.Login)
	g.Any("/account/edit", self.Edit, middleware.NeedLogin())
	g.Post("/account/change_avatar", self.ChangeAvatar, middleware.NeedLogin())
	g.Post("/account/changepwd", self.ChangePwd, middleware.NeedLogin())
	g.Any("/account/forgetpwd", self.ForgetPasswd)
	g.Any("/account/resetpwd", self.ResetPasswd)
	g.Get("/account/logout", self.Logout, middleware.NeedLogin())
	g.POST("/account/social/unbind", self.Unbind, middleware.NeedLogin())
}

func (self AccountController) Register(ctx echo.Context) error {
	if _, ok := ctx.Get("user").(*model.Me); ok {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	registerTpl := "register.html"
	username := ctx.FormValue("username")
	// 请求注册页面
	if username == "" || ctx.Request().Method() != "POST" {
		return render(ctx, registerTpl, map[string]interface{}{"captchaId": captcha.NewLen(4)})
	}

	data := map[string]interface{}{
		"username":  username,
		"email":     ctx.FormValue("email"),
		"captchaId": captcha.NewLen(util.CaptchaLen),
	}

	disallowUsers := config.ConfigFile.MustValueArray("account", "disallow_user", ",")
	for _, disallowUser := range disallowUsers {
		if disallowUser == username {
			data["error"] = username + " 被禁止使用，请换一个"
			return render(ctx, registerTpl, data)
		}
	}

	captchaId := ctx.FormValue("captchaid")
	// 校验验证码
	if !captcha.VerifyString(captchaId, ctx.FormValue("captchaSolution")) {
		data["error"] = "验证码错误，记得刷新验证码"
		util.SetCaptcha(captchaId)
		return render(ctx, registerTpl, data)
	}

	if ctx.FormValue("passwd") != ctx.FormValue("pass2") {
		data["error"] = "两次密码不一致"
		return render(ctx, registerTpl, data)
	}

	fields := []string{"username", "email", "passwd"}
	form := url.Values{}
	for _, field := range fields {
		form.Set(field, ctx.FormValue(field))
	}

	// 入库
	errMsg, err := logic.DefaultUser.CreateUser(ctx, form)
	if err != nil {
		// bugfix：http://studygolang.com/topics/255
		if errMsg == "" {
			errMsg = err.Error()
		}
		data["error"] = errMsg
		return render(ctx, registerTpl, data)
	}

	email := ctx.FormValue("email")

	uuid := RegActivateCode.GenUUID(email)
	var emailUrl string
	if strings.HasSuffix(email, "@gmail.com") {
		emailUrl = "http://mail.google.com"
	} else {
		pos := strings.LastIndex(email, "@")
		emailUrl = "http://mail." + email[pos+1:]
	}

	if config.ConfigFile.MustBool("account", "verify_email", true) {
		data = map[string]interface{}{
			"success": template.HTML(`
				<div style="padding:30px 30px 50px 30px;">
	 				<div style="color:#339502;font-size:22px;line-height: 2.5;">恭喜您注册成功！</div>
	 				我们已经发送一封邮件到 ` + email + `，请您根据提示信息完成邮箱验证.<br><br>
	 				<a href="` + emailUrl + `" target="_blank"><button type="button" class="btn btn-success">立即验证</button></a>&nbsp;&nbsp;<button type="button" class="btn btn-link" data-uuid="` + uuid + `" id="resend_email">未收到？再发一次</button>
				</div>`),
		}

		isHttps := CheckIsHttps(ctx)
		// 需要检验邮箱的正确性
		go logic.DefaultEmail.SendActivateMail(email, uuid, isHttps)

		return render(ctx, registerTpl, data)
	}

	// 不验证邮箱，注册完成直接登录
	// 自动登录
	SetLoginCookie(ctx, username)

	return ctx.Redirect(http.StatusSeeOther, "/balance")
}

// SendActivateEmail 发送注册激活邮件
func (self AccountController) SendActivateEmail(ctx echo.Context) error {
	isHttps := CheckIsHttps(ctx)

	uuid := ctx.FormValue("uuid")
	if uuid != "" {
		email, ok := RegActivateCode.GetEmail(uuid)
		if !ok {
			return fail(ctx, 1, "非法请求")
		}

		go logic.DefaultEmail.SendActivateMail(email, uuid, isHttps)
	} else {
		user, ok := ctx.Get("user").(*model.Me)
		if !ok {
			return fail(ctx, 1, "非法请求")
		}

		go logic.DefaultEmail.SendActivateMail(user.Email, RegActivateCode.GenUUID(user.Email), isHttps)
	}

	return success(ctx, nil)
}

// Activate 用户激活
func (AccountController) Activate(ctx echo.Context) error {
	contentTpl := "user/activate.html"

	data := map[string]interface{}{}

	param := goutils.Base64Decode(ctx.QueryParam("param"))
	values, err := url.ParseQuery(param)
	if err != nil {
		data["error"] = err.Error()
		return render(ctx, contentTpl, data)
	}

	uuid := values.Get("uuid")
	timestamp := goutils.MustInt64(values.Get("timestamp"))
	sign := values.Get("sign")
	email, ok := RegActivateCode.GetEmail(uuid)
	if !ok {
		data["error"] = "非法请求！"
		return render(ctx, contentTpl, data)
	}

	if timestamp < time.Now().Add(-4*time.Hour).Unix() {
		RegActivateCode.DelUUID(uuid)
		// TODO:可以再次发激活邮件？
		data["error"] = "链接已过期"
		return render(ctx, contentTpl, data)
	}

	user, err := logic.DefaultUser.Activate(ctx, email, uuid, timestamp, sign)
	if err != nil {
		data["error"] = err.Error()
		return render(ctx, contentTpl, data)
	}

	RegActivateCode.DelUUID(uuid)

	// 自动登录
	SetLoginCookie(ctx, user.Username)

	// return render(ctx, contentTpl, data)
	return ctx.Redirect(http.StatusSeeOther, "/balance")
}

// Login 登录
func (AccountController) Login(ctx echo.Context) error {
	if _, ok := ctx.Get("user").(*model.Me); ok {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	// 支持跳转到源页面
	uri := ctx.FormValue("redirect_uri")
	if uri == "" {
		referer := ctx.Request().Referer()
		if referer == "" {
			uri = "/"
		} else {
			uri = referer
		}
	}

	contentTpl := "login.html"
	data := make(map[string]interface{})

	username := ctx.FormValue("username")
	if username == "" || ctx.Request().Method() != "POST" {
		data["redirect_uri"] = uri
		return render(ctx, contentTpl, data)
	}

	// 处理用户登录
	passwd := ctx.FormValue("passwd")
	userLogin, err := logic.DefaultUser.Login(ctx, username, passwd)
	if err != nil {
		data["username"] = username
		data["error"] = err.Error()

		if util.IsAjax(ctx) {
			return fail(ctx, 1, err.Error())
		}

		return render(ctx, contentTpl, data)
	}

	// 登录成功，种cookie
	SetLoginCookie(ctx, userLogin.Username)

	if util.IsAjax(ctx) {
		return success(ctx, nil)
	}

	return ctx.Redirect(http.StatusSeeOther, uri)
}

// Edit 用户编辑个人信息
func (self AccountController) Edit(ctx echo.Context) error {
	me := ctx.Get("user").(*model.Me)

	if ctx.Request().Method() != "POST" {
		user := logic.DefaultUser.FindOne(ctx, "uid", me.Uid)
		bindUsers := logic.DefaultUser.FindBindUsers(ctx, me.Uid)
		return render(ctx, "user/edit.html", map[string]interface{}{
			"user":            user,
			"default_avatars": logic.DefaultAvatars,
			"has_passwd":      logic.DefaultUser.HasPasswd(ctx, me.Uid),
			"bind_users":      bindUsers,
		})
	}

	// 更新信息
	errMsg, err := logic.DefaultUser.Update(ctx, me, ctx.Request().FormParams())
	if err != nil {
		return fail(ctx, 1, errMsg)
	}

	email := ctx.FormValue("email")
	if me.Email != email {
		isHttps := CheckIsHttps(ctx)
		go logic.DefaultEmail.SendActivateMail(email, RegActivateCode.GenUUID(email), isHttps)
	}

	return success(ctx, nil)
}

// ChangeAvatar 更换头像
func (AccountController) ChangeAvatar(ctx echo.Context) error {
	objLog := getLogger(ctx)

	curUser := ctx.Get("user").(*model.Me)

	// avatar 为空时，表示使用 gravater 头像
	avatar := ctx.FormValue("avatar")
	err := logic.DefaultUser.ChangeAvatar(ctx, curUser.Uid, avatar)
	if err != nil {
		objLog.Errorln("account controller change avatar error:", err)

		return fail(ctx, 2, "更换头像失败")
	}

	return success(ctx, nil)
}

// ChangePwd 修改密码
func (AccountController) ChangePwd(ctx echo.Context) error {
	curUser := ctx.Get("user").(*model.Me)

	curPasswd := ctx.FormValue("cur_passwd")
	newPasswd := ctx.FormValue("passwd")
	errMsg, err := logic.DefaultUser.UpdatePasswd(ctx, curUser.Username, curPasswd, newPasswd)
	if err != nil {
		return fail(ctx, 1, errMsg)
	}
	return success(ctx, nil)
}

// 保存uuid和email的对应关系（TODO:重启如何处理，有效期问题）
var resetPwdMap = map[string]string{}

// ForgetPasswd 忘记密码
func (AccountController) ForgetPasswd(ctx echo.Context) error {
	if _, ok := ctx.Get("user").(*model.Me); ok {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	contentTpl := "user/forget_pwd.html"
	data := map[string]interface{}{"activeUsers": "active"}

	email := ctx.FormValue("email")
	if email == "" || ctx.Request().Method() != "POST" {
		return render(ctx, contentTpl, data)
	}

	// 校验email是否存在
	if logic.DefaultUser.UserExists(ctx, "email", email) {
		var uuid string
		for {
			uuid = guuid.NewV4().String()
			if _, ok := resetPwdMap[uuid]; !ok {
				resetPwdMap[uuid] = email
				break
			}
			logger.Infoln("forget passwd GenUUID 冲突....")
		}
		var emailUrl string
		if strings.HasSuffix(email, "@gmail.com") {
			emailUrl = "http://mail.google.com"
		} else {
			pos := strings.LastIndex(email, "@")
			emailUrl = "http://mail." + email[pos+1:]
		}

		isHttps := CheckIsHttps(ctx)
		data["success"] = template.HTML(`一封包含了重设密码链接的邮件已经发送到您的注册邮箱，按照邮件中的提示，即可重设您的密码。<a href="` + emailUrl + `" target="_blank">立即前往邮箱</a>`)
		go logic.DefaultEmail.SendResetpwdMail(email, uuid, isHttps)
	} else {
		data["error"] = "该邮箱没有在本社区注册过！"
	}

	return render(ctx, contentTpl, data)
}

// ResetPasswd 重置密码
func (AccountController) ResetPasswd(ctx echo.Context) error {
	if _, ok := ctx.Get("user").(*model.Me); ok {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	uuid := ctx.FormValue("code")
	if uuid == "" {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	contentTpl := "user/reset_pwd.html"
	data := map[string]interface{}{"activeUsers": "active"}

	method := ctx.Request().Method()

	passwd := ctx.FormValue("passwd")
	email, ok := resetPwdMap[uuid]
	if !ok {
		// 是提交重置密码
		if passwd != "" && method == "POST" {
			data["error"] = template.HTML(`非法请求！<p>将在<span id="jumpTo">3</span>秒后跳转到<a href="/" id="jump_url">首页</a></p>`)
		} else {
			data["error"] = template.HTML(`链接无效或过期，请重新操作。<a href="/account/forgetpwd">忘记密码？</a>`)
		}
		return render(ctx, contentTpl, data)
	}

	data["valid"] = true
	data["code"] = uuid
	// 提交修改密码
	if passwd != "" && method == "POST" {
		// 简单校验
		if len(passwd) < 6 || len(passwd) > 32 {
			data["error"] = "密码长度必须在6到32个字符之间"
		} else if passwd != ctx.FormValue("pass2") {
			data["error"] = "两次密码输入不一致"
		} else {
			// 更新密码
			_, err := logic.DefaultUser.ResetPasswd(ctx, email, passwd)
			if err != nil {
				data["error"] = "对不起，服务器错误，请重试！"
			} else {
				data["success"] = template.HTML(`密码重置成功，<p>将在<span id="jumpTo">3</span>秒后跳转到<a href="/account/login" id="jump_url">登录</a>页面</p>`)
			}
		}
	}
	return render(ctx, contentTpl, data)
}

// Logout 注销
func (AccountController) Logout(ctx echo.Context) error {
	// 删除cookie信息
	session := GetCookieSession(ctx)
	session.Options = &sessions.Options{Path: "/", MaxAge: -1}
	session.Save(Request(ctx), ResponseWriter(ctx))
	// 重定向得到原页面
	return ctx.Redirect(http.StatusSeeOther, ctx.Request().Referer())
}

// Unbind 第三方账号解绑
func (AccountController) Unbind(ctx echo.Context) error {
	bindId := ctx.FormValue("bind_id")
	me := ctx.Get("user").(*model.Me)
	logic.DefaultThirdUser.UnBindUser(ctx, bindId, me)

	return ctx.Redirect(http.StatusSeeOther, "/account/edit#connection")
}
