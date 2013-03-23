// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	. "config"
	"logger"
	"net/smtp"
	"strings"
)

// 发送电子邮件功能
func SendMail(subject, content string, tos []string) error {
	message := `From: Golang中文社区 | Go语言学习园地
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
