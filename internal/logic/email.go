// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strings"
	"sync"
	"time"

	. "github.com/studygolang/studygolang/db"
	"github.com/studygolang/studygolang/global"
	"github.com/studygolang/studygolang/internal/model"
	"github.com/studygolang/studygolang/util"

	"github.com/polaris1119/config"
	"github.com/polaris1119/email"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
)

type EmailLogic struct{}

var DefaultEmail = EmailLogic{}

// activateSignSalt 注册激活邮件的签名 salt。空值是严重安全漏洞：
// genActivateSign 在 salt="" 时退化成 Md5("uuid=Xemail=Ytimestamp=Z")，
// 攻击者无需任何 secret 即可伪造任意邮箱的激活签名，绕过邮箱所有权校验。
// 因此启动时强制校验：env 优先 → config 回退 → 生产环境 panic。
var activateSignSalt string

// unsubscribeTokenKey 退订邮件 token 的 secret。空值时 token 退化为
// Md5(user.String())，user 字段大多公开（username/email/avatar 等），
// 攻击者可伪造任意用户退订 token，让受害者静默收不到邮件通知。
var unsubscribeTokenKey string

func init() {
	env := config.ConfigFile.MustValue("global", "env", "dev")

	activateSignSalt = os.Getenv("ACTIVATE_SIGN_SALT")
	if activateSignSalt == "" {
		activateSignSalt = config.ConfigFile.MustValue("security", "activate_sign_salt")
	}
	if activateSignSalt == "" {
		if env == "prod" {
			panic("security.activate_sign_salt not configured! Please set ACTIVATE_SIGN_SALT env var or config/env.ini. Empty salt enables activation signature forgery.")
		}
		logger.Errorln("security.activate_sign_salt is empty; activation signatures are insecure. Fix before production.")
	}

	unsubscribeTokenKey = os.Getenv("UNSUBSCRIBE_TOKEN_KEY")
	if unsubscribeTokenKey == "" {
		unsubscribeTokenKey = config.ConfigFile.MustValue("security", "unsubscribe_token_key")
	}
	if unsubscribeTokenKey == "" {
		if env == "prod" {
			panic("security.unsubscribe_token_key not configured! Empty key enables unsubscribe token forgery (silent email suppression for any victim).")
		}
		logger.Errorln("security.unsubscribe_token_key is empty; unsubscribe tokens are insecure. Fix before production.")
	}
}

// SendMail 发送普通（通知）电子邮件
func (e EmailLogic) SendMail(subject, content string, tos []string) (err error) {
	return e.sendMail(subject, content, tos, "email")
}

// SendAuthMail 发送验证电子邮件
func (e EmailLogic) SendAuthMail(subject, content string, tos []string) error {
	return e.sendMail(subject, content, tos, "email.auth")
}

// sendMail 发送电子邮件
func (EmailLogic) sendMail(subject, content string, tos []string, section string) (err error) {
	emailConfig, _ := config.ConfigFile.GetSection(section)

	fromEmail := emailConfig["from_email"]
	smtpUsername := emailConfig["smtp_username"]
	smtpPassword := emailConfig["smtp_password"]
	smtpHost := emailConfig["smtp_host"]
	smtpPort := emailConfig["smtp_port"]

	mail := email.NewEmail()
	mail.From = WebsiteSetting.Name + ` <` + fromEmail + `>`
	mail.To = tos
	mail.Subject = subject
	mail.HTML = []byte(content)

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	smtpAddr := smtpHost + ":" + smtpPort

	if goutils.MustBool(emailConfig["tls"]) {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         smtpHost,
		}

		err = mail.SendWithTLS(smtpAddr, auth, tlsConfig)
	} else {
		err = mail.Send(smtpAddr, auth)
	}

	if err != nil {
		logger.Errorln("Send Mail to", strings.Join(tos, ","), "error:", err)
		return
	}
	logger.Infoln("Send Mail to", strings.Join(tos, ","), "Successfully")
	return
}

// 保存uuid和email的对应关系（TODO:重启如何处理）
var regActivateCodeMap = map[string]string{}

// SendActivateMail 发送激活邮件
func (self EmailLogic) SendActivateMail(email, uuid string, isHttps ...bool) {
	timestamp := time.Now().Unix()
	sign := self.genActivateSign(email, uuid, timestamp)
	param := goutils.Base64Encode(fmt.Sprintf("uuid=%s&timestamp=%d&sign=%s", uuid, timestamp, sign))

	domain := "http://" + WebsiteSetting.Domain
	if len(isHttps) > 0 && isHttps[0] {
		domain = "https://" + WebsiteSetting.Domain
	}

	activeUrl := fmt.Sprintf("%s/account/activate?param=%s", domain, param)

	global.App.SetCopyright()

	content := `
尊敬的` + WebsiteSetting.Name + `用户：<br/><br/>
感谢您选择了` + WebsiteSetting.Name + `，请点击下面的地址激活你在` + WebsiteSetting.Name + `的帐号（有效期4小时）：<br/><br/>
<a href="` + activeUrl + `">` + activeUrl + `</a><br/><br/>
<div style="text-align:right;">&copy;` + global.App.Copyright + ` ` + WebsiteSetting.Name + `</div>`
	self.SendAuthMail(WebsiteSetting.Name+"帐号激活邮件", content, []string{email})
}

func (EmailLogic) genActivateSign(email, uuid string, ts int64) string {
	// 使用启动时校验过的全局 salt；若部署时漏配，init() 已 panic 或告警，
	// 这里再多一道防御：salt 为空时返回空串，让所有签名校验失败（fail-closed）
	// 而不是退化成无 secret 的可伪造签名。
	if activateSignSalt == "" {
		return ""
	}
	origStr := fmt.Sprintf("uuid=%semail=%stimestamp=%d%s", uuid, email, ts, activateSignSalt)
	return goutils.Md5(origStr)
}

// SendResetpwdMail 发重置密码邮件
func (self EmailLogic) SendResetpwdMail(email, uuid string, isHttps ...bool) {
	global.App.SetCopyright()

	domain := "http://" + WebsiteSetting.Domain
	if len(isHttps) > 0 && isHttps[0] {
		domain = "https://" + WebsiteSetting.Domain
	}
	content := `您好，` + email + `,<br/><br/>
&nbsp;&nbsp;&nbsp;&nbsp;我们的系统收到一个请求，说您希望通过电子邮件重新设置您在 <a href="` + domain + `">` + WebsiteSetting.Name + `</a> 的密码。您可以点击下面的链接重设密码：<br/><br/>

&nbsp;&nbsp;&nbsp;&nbsp;` + domain + `/account/resetpwd?code=` + uuid + ` <br/><br/>

如果这个请求不是由您发起的，那没问题，您不用担心，您可以安全地忽略这封邮件。<br/><br/>

如果您有任何疑问，可以回复这封邮件向我们提问。谢谢！<br/><br/>

<div style="text-align:right;">&copy;` + global.App.Copyright + ` ` + WebsiteSetting.Name + `</div>`
	self.SendAuthMail("【"+WebsiteSetting.Name+"】重设密码 ", content, []string{email})
}

// 自定义模板函数
var emailFuncMap = template.FuncMap{
	"time_format": func(t model.OftenTime) string {
		return time.Time(t).Format("01-02")
	},
	"substring": util.Substring,
}

var (
	emailTplOnce sync.Once
	emailTpl     *template.Template
)

func getEmailTpl() *template.Template {
	emailTplOnce.Do(func() {
		emailTpl = template.Must(template.New("email.html").Funcs(emailFuncMap).ParseFiles(config.TemplateDir + "email.html"))
	})
	return emailTpl
}

// 订阅邮件通知
func (self EmailLogic) EmailNotice() {

	beginDate := time.Now().Add(-7 * 24 * time.Hour).Format("2006-01-02")
	endDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")

	beginTime := beginDate + " 00:00:00"

	// 本周晨读（过去 7 天）
	readings, err := DefaultReading.FindLastList(beginTime)
	if err != nil {
		logger.Errorln("find morning reading error:", err)
	}

	// 本周精彩文章
	articles, err := DefaultArticle.FindLastList(beginTime, 10)
	if err != nil {
		logger.Errorln("find article error:", err)
	}

	// 本周热门主题
	topics, err := DefaultTopic.FindLastList(beginTime, 10)
	if err != nil {
		logger.Errorln("find topic error:", err)
	}

	global.App.SetCopyright()

	data := map[string]interface{}{
		"readings":  readings,
		"articles":  articles,
		"topics":    topics,
		"beginDate": beginDate,
		"endDate":   endDate,
		"setting":   WebsiteSetting,
		"app":       global.App,
	}

	// 给所有用户发送邮件
	var (
		lastUid = 0
		limit   = 500
		users   = make([]*model.User, 0)
	)

	day := time.Now().Day()
	monthDayNum := util.MonthDayNum(time.Now())

	for {
		err = MasterDB.Where("uid>?", lastUid).Asc("uid").Limit(limit).Find(&users)
		if err != nil {
			logger.Errorln("find user error:", err)
			continue
		}

		if len(users) == 0 {
			break
		}

		for _, user := range users {
			if lastUid < user.Uid {
				lastUid = user.Uid
			}

			if user.Uid%monthDayNum != day {
				continue
			}

			if user.Unsubscribe == 1 {
				logger.Infoln("user unsubscribe", user)
				continue
			}

			if user.Status != model.UserStatusAudit {
				logger.Infoln("user is not normal:", user)
				continue
			}

			if user.IsThird == 1 && strings.HasSuffix(user.Email, "github.com") {
				logger.Infoln("the email is not exists:", user)
				continue
			}

			data["email"] = user.Email
			data["token"] = self.GenUnsubscribeToken(user)

			content, err := self.genEmailContent(data)
			if err != nil {
				logger.Errorln("from email.html gen email content error:", err)
				continue
			}

			self.SendMail("每周精选", content, []string{user.Email})

			// 控制发信速度
			time.Sleep(60 * time.Second)
		}

		users = make([]*model.User, 0)
	}

}

// 生成 退订 邮件的 token
func (EmailLogic) GenUnsubscribeToken(user *model.User) string {
	// fail-closed：key 为空时返回空串，让所有 token 校验失败而非可伪造
	if unsubscribeTokenKey == "" {
		return ""
	}
	return goutils.Md5(user.String() + unsubscribeTokenKey)
}

func (EmailLogic) genEmailContent(data map[string]interface{}) (string, error) {
	buffer := &bytes.Buffer{}
	if err := getEmailTpl().Execute(buffer, data); err != nil {
		logger.Errorln("email logic execute template error:", err)
		return "", err
	}

	return buffer.String(), nil
}
