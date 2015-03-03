// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"bytes"
	"html/template"
	"net/smtp"
	"strings"
	"time"

	. "config"
	"logger"
	"model"
	"util"
)

// 发送电子邮件功能
func SendMail(subject, content string, tos []string) error {
	message := `From: Go语言中文网 | Golang中文社区 | Go语言学习园地<` + Config["from_email"] + `>
To: ` + strings.Join(tos, ",") + `
Subject: ` + subject + `
Content-Type: text/html;charset=UTF-8

` + content

	auth := smtp.PlainAuth("", Config["smtp_username"], Config["smtp_password"], Config["smtp_host"])
	err := smtp.SendMail(Config["smtp_addr"], auth, Config["from_email"], tos, []byte(message))
	if err != nil {
		logger.Errorln("Send Mail to", strings.Join(tos, ","), "error:", err)
		return err
	}
	logger.Infoln("Send Mail to", strings.Join(tos, ","), "Successfully")
	return nil
}

// 自定义模板函数
var emailFuncMap = template.FuncMap{
	"time_format": func(s string) string {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
		if err != nil {
			return s
		}

		return t.Format("01-02")
	},
	"substring": util.Substring,
}

var emailTpl = template.Must(template.New("email.html").Funcs(emailFuncMap).ParseFiles(ROOT + "/template/email.html"))

// 订阅邮件通知
func EmailNotice() {

	beginDate := time.Now().Add(-7 * 24 * time.Hour).Format("2006-01-02")
	endDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")

	beginTime := beginDate + " 00:00:00"

	// 本周晨读（过去 7 天）
	readings, err := model.NewMorningReading().Where("ctime>? AND rtype=0", beginTime).Order("id DESC").FindAll()
	if err != nil {
		logger.Errorln("find morning reading error:", err)
	}

	// 本周精彩文章
	articles, err := model.NewArticle().Where("ctime>? AND status!=2", beginTime).Order("cmtnum DESC, likenum DESC, viewnum DESC").Limit("10").FindAll()
	if err != nil {
		logger.Errorln("find article error:", err)
	}

	// 本周热门主题
	topics, err := model.NewTopic().Where("ctime>? AND flag IN(0,1)", beginTime).Order("tid DESC").Limit("10").FindAll()
	if err != nil {
		logger.Errorln("find topic error:", err)
	}

	data := map[string]interface{}{
		"readings":  readings,
		"articles":  articles,
		"topics":    topics,
		"beginDate": beginDate,
		"endDate":   endDate,
	}

	// 给所有用户发送邮件
	userModel := model.NewUser()

	var (
		lastUid = 0
		limit   = "500"
		users   []*model.User
	)

	for {
		users, err = userModel.Where("uid>?", lastUid).Order("uid ASC").Limit(limit).FindAll()
		if err != nil {
			logger.Errorln("find user error:", err)
			continue
		}

		if len(users) == 0 {
			break
		}

		for _, user := range users {
			if user.Unsubscribe == 1 {
				logger.Infoln(user)
				continue
			}

			data["email"] = user.Email
			data["token"] = GenUnsubscribeToken(user)

			content, err := genEmailContent(data)
			if err != nil {
				logger.Errorln("from email.html gen email content error:", err)
				continue
			}

			SendMail("每周精选", content, []string{user.Email})

			if lastUid < user.Uid {
				lastUid = user.Uid
			}
		}
	}

}

// 生成 退订 邮件的 token
func GenUnsubscribeToken(user *model.User) string {
	return util.Md5(user.String() + Config["unsubscribe_token_key"])
}

func genEmailContent(data map[string]interface{}) (string, error) {
	buffer := &bytes.Buffer{}
	if err := emailTpl.Execute(buffer, data); err != nil {
		logger.Errorln("execute template error:", err)
		return "", err
	}

	return buffer.String(), nil
}
