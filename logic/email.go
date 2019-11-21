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
	"strings"
	"time"

	. "github.com/studygolang/studygolang/db"
	"github.com/studygolang/studygolang/global"
	"github.com/studygolang/studygolang/model"
	"github.com/studygolang/studygolang/util"

	"github.com/polaris1119/config"
	"github.com/polaris1119/email"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
)

type EmailLogic struct{}

var DefaultEmail = EmailLogic{}

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
	emailSignSalt := config.ConfigFile.MustValue("security", "activate_sign_salt")
	origStr := fmt.Sprintf("uuid=%semail=%stimestamp=%d%s", uuid, email, ts, emailSignSalt)
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

var emailTpl = template.Must(template.New("email.html").Funcs(emailFuncMap).ParseFiles(config.TemplateDir + "email.html"))

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
	return goutils.Md5(user.String() + config.ConfigFile.MustValue("security", "unsubscribe_token_key"))
}

func (EmailLogic) genEmailContent(data map[string]interface{}) (string, error) {
	buffer := &bytes.Buffer{}
	if err := emailTpl.Execute(buffer, data); err != nil {
		logger.Errorln("email logic execute template error:", err)
		return "", err
	}

	return buffer.String(), nil
}
