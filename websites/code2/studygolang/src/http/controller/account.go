// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package controller

import (
	"html/template"
	"logic"
	"model"
	"net/http"
	"net/url"
	"strings"
	"time"

	. "http"

	"github.com/dchest/captcha"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
	guuid "github.com/twinj/uuid"
)

type AccountController struct{}

// 注册路由
func (this *AccountController) RegisterRoute(e *echo.Echo) {
	e.Any("/account/register", this.Register)
	e.Post("/account/send_activate_email", this.SendActivateEmail)
	e.Get("/account/activate", this.Activate)
	e.Any("/account/login", this.Login)
	e.Any("/account/edit", this.Edit)
	e.Post("/account/change_avatar", this.ChangeAvatar)
	e.Post("/account/changepwd", this.ChangePwd)
	e.Any("/account/forgetpwd", this.ForgetPasswd)
	e.Any("/account/resetpwd", this.ResetPasswd)
	e.Get("/account/logout", this.Logout)
}

// 保存uuid和email的对应关系（TODO:重启如何处理，有效期问题）
var regActivateCodeMap = map[string]string{}

func (AccountController) Register(ctx *echo.Context) error {
	if _, ok := ctx.Get("user").(*model.User); ok {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	registerTpl := "register.html"
	username := ctx.Form("username")
	// 请求注册页面
	if username == "" || ctx.Request().Method != "POST" {
		return render(ctx, registerTpl, map[string]interface{}{"captchaId": captcha.NewLen(4)})
	}

	// 校验验证码
	if !captcha.VerifyString(ctx.Form("captchaid"), ctx.Form("captchaSolution")) {
		return render(ctx, registerTpl, map[string]interface{}{"error": "验证码错误", "captchaId": captcha.NewLen(4)})
	}

	// 入库
	errMsg, err := logic.DefaultUser.CreateUser(ctx, ctx.Request().Form)
	if err != nil {
		// bugfix：http://studygolang.com/topics/255
		if errMsg == "" {
			errMsg = err.Error()
		}
		return render(ctx, registerTpl, map[string]interface{}{"error": errMsg})
	}

	var (
		uuid  string
		email = ctx.Form("email")
	)
	for {
		uuid = guuid.NewV4().String()
		if _, ok := regActivateCodeMap[uuid]; !ok {
			regActivateCodeMap[uuid] = email
			break
		}
		logger.Errorln("GenUUID 冲突....")
	}
	var emailUrl string
	if strings.HasSuffix(email, "@gmail.com") {
		emailUrl = "http://mail.google.com"
	} else {
		pos := strings.LastIndex(email, "@")
		emailUrl = "http://mail." + email[pos+1:]
	}
	data := map[string]interface{}{
		"success": template.HTML(`
			<div style="padding:30px 30px 50px 30px;">
 				<div style="color:#339502;font-size:22px;line-height: 2.5;">恭喜您注册成功！</div>
 				我们已经发送一封邮件到 ` + email + `，请您根据提示信息完成邮箱验证.<br><br>
 				<a href="` + emailUrl + `" target="_blank"><button type="button" class="btn btn-success">立即验证</button></a>&nbsp;&nbsp;<button type="button" class="btn btn-link" data-uuid="` + uuid + `" id="resend_email">未收到？再发一次</button>
			</div>`),
	}
	// 需要检验邮箱的正确性
	go logic.DefaultEmail.SendActivateMail(email, uuid)

	return render(ctx, registerTpl, data)
}

// SendActivateEmail 发送注册激活邮件
func (AccountController) SendActivateEmail(ctx *echo.Context) error {
	uuid := ctx.Form("uuid")
	email, ok := regActivateCodeMap[uuid]
	if !ok {
		return fail(ctx, 1, "非法请求")
	}

	go logic.DefaultEmail.SendActivateMail(email, uuid)

	return success(ctx, nil)
}

// Activate 用户激活
func (AccountController) Activate(ctx *echo.Context) error {
	contentTpl := "user/activate.html"

	data := map[string]interface{}{}

	param := goutils.Base64Decode(ctx.Query("param"))
	values, err := url.ParseQuery(param)
	if err != nil {
		data["error"] = err.Error()
		return render(ctx, contentTpl, data)
	}

	uuid := values.Get("uuid")
	timestamp := goutils.MustInt64(values.Get("timestamp"))
	sign := values.Get("sign")
	email, ok := regActivateCodeMap[uuid]
	if !ok {
		data["error"] = "非法请求！"
		return render(ctx, contentTpl, data)
	}

	if timestamp < time.Now().Add(-4*time.Hour).Unix() {
		delete(regActivateCodeMap, uuid)
		// TODO:可以再次发激活邮件？
		data["error"] = "链接已过期"
		return render(ctx, contentTpl, data)
	}

	user, err := logic.DefaultUser.Activate(ctx, email, uuid, timestamp, sign)
	if err != nil {
		data["error"] = err.Error()
		return render(ctx, contentTpl, data)
	}

	delete(regActivateCodeMap, uuid)

	// 自动登录
	SetCookie(ctx, user.Username)

	return render(ctx, contentTpl, data)
}

// Login 登录
func (AccountController) Login(ctx *echo.Context) error {
	if _, ok := ctx.Get("user").(*model.User); ok {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	contentTpl := "login.html"

	data := make(map[string]interface{})

	username := ctx.Form("username")
	if username == "" || ctx.Request().Method != "POST" {
		return render(ctx, contentTpl, data)
	}

	// 处理用户登录
	passwd := ctx.Form("passwd")
	userLogin, err := logic.DefaultUser.Login(ctx, username, passwd)
	if err != nil {
		data["username"] = username
		data["error"] = err.Error()
		return render(ctx, contentTpl, data)
	}

	// 登录成功，种cookie
	SetCookie(ctx, userLogin.Username)

	// 支持跳转到源页面
	uri := ctx.Query("redirect_uri")
	if uri == "" {
		uri = "/"
	}

	return ctx.Redirect(http.StatusSeeOther, uri)
}

// Edit 用户编辑个人信息
func (AccountController) Edit(ctx *echo.Context) error {
	curUser := ctx.Get("user").(*model.User)

	if ctx.Request().Method != "POST" {
		return render(ctx, "user/edit.html", map[string]interface{}{
			"user":            curUser,
			"default_avatars": logic.DefaultAvatars,
		})
	}

	// 更新信息
	errMsg, err := logic.DefaultUser.Update(ctx, curUser.Uid, ctx.Request().PostForm)
	if err != nil {
		return fail(ctx, 1, errMsg)
	}

	return success(ctx, nil)
}

// ChangeAvatar 更换头像
func (AccountController) ChangeAvatar(ctx *echo.Context) error {
	objLog := getLogger(ctx)

	curUser := ctx.Get("user").(*model.User)

	avatar := ctx.Form("avatar")
	if avatar == "" {
		return fail(ctx, 1, "非法请求")
	}
	err := logic.DefaultUser.ChangeAvatar(ctx, curUser.Uid, avatar)
	if err != nil {
		objLog.Errorln("account controller change avatar error:", err)

		return fail(ctx, 2, "更换头像失败")
	}

	return success(ctx, nil)
}

// ChangePwd 修改密码
func (AccountController) ChangePwd(ctx *echo.Context) error {
	curUser := ctx.Get("user").(*model.User)

	curPasswd := ctx.Form("cur_passwd")
	newPasswd := ctx.Form("passwd")
	errMsg, err := logic.DefaultUser.UpdatePasswd(ctx, curUser.Username, curPasswd, newPasswd)
	if err != nil {
		return fail(ctx, 1, errMsg)
	}
	return success(ctx, nil)
}

// 保存uuid和email的对应关系（TODO:重启如何处理，有效期问题）
var resetPwdMap = map[string]string{}

// ForgetPasswd 忘记密码
func (AccountController) ForgetPasswd(ctx *echo.Context) error {
	if _, ok := ctx.Get("user").(*model.User); ok {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	contentTpl := "user/forget_pwd.html"
	data := map[string]interface{}{"activeUsers": "active"}

	email := ctx.Form("email")
	if email == "" || ctx.Request().Method != "POST" {
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
		data["success"] = template.HTML(`一封包含了重设密码链接的邮件已经发送到您的注册邮箱，按照邮件中的提示，即可重设您的密码。<a href="` + emailUrl + `" target="_blank">立即前往邮箱</a>`)
		go logic.DefaultEmail.SendResetpwdMail(email, uuid)
	} else {
		data["error"] = "该邮箱没有在本社区注册过！"
	}

	return render(ctx, contentTpl, data)
}

// ResetPasswd 重置密码
func (AccountController) ResetPasswd(ctx *echo.Context) error {
	if _, ok := ctx.Get("user").(*model.User); ok {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	uuid := ctx.Form("code")
	if uuid == "" {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	contentTpl := "user/reset_pwd.html"
	data := map[string]interface{}{"activeUsers": "active"}

	method := ctx.Request().Method

	passwd := ctx.Form("passwd")
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
		} else if passwd != ctx.Form("pass2") {
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
func (AccountController) Logout(ctx *echo.Context) error {
	// 删除cookie信息
	session := GetCookieSession(ctx)
	session.Options = &sessions.Options{Path: "/", MaxAge: -1}
	session.Save(ctx.Request(), ctx.Response())
	// 重定向得到登录页（TODO:重定向到什么页面比较好？）
	return ctx.Redirect(http.StatusSeeOther, "/account/login")
}
