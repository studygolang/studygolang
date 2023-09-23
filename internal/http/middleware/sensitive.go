// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/config"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
)

var (
	titleSensitives   []string
	contentSensitives string

	midNightSpam []string
	num          int
)

func init() {
	titleSensitives = strings.Split(config.ConfigFile.MustValue("sensitive", "title"), ",")
	contentSensitives = config.ConfigFile.MustValue("sensitive", "content")
	midNightSpam = strings.Split(config.ConfigFile.MustValue("spam", "mid_night"), ",")
	num = config.ConfigFile.MustInt("spam", "num")
}

// Sensivite 用于 echo 框架的过滤发布敏感词（广告）
func Sensivite() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			content := ctx.FormValue("content")
			title := ctx.FormValue("title")

			user := ctx.Get("user").(*model.Me)

			if title != "" {
				for _, s := range titleSensitives {
					if hasSensitiveChar(title, s) {
						// 把账号冻结
						logic.DefaultUser.UpdateUserStatus(context.EchoContext(ctx), user.Uid, model.UserStatusFreeze)
						logger.Infoln("user=", user.Uid, "publish ad, title=", title, ". freeze")
						// IP 加入黑名单
						addBlackIP(ctx)
						return ctx.String(http.StatusOK, `{"ok":0,"error":"对不起，您的账号已被冻结！"}`)
					}
				}
			}

			if hasSensitive(title, contentSensitives) || hasSensitive(content, contentSensitives) {
				// 把账号冻结
				logic.DefaultUser.UpdateUserStatus(context.EchoContext(ctx), user.Uid, model.UserStatusFreeze)
				logger.Infoln("user=", user.Uid, "publish ad, title=", title, ";content=", content, ". freeze")
				// IP 加入黑名单
				addBlackIP(ctx)
				return ctx.String(http.StatusOK, `{"ok":0,"error":"对不起，您的账号已被冻结！"}`)
			}

			// 半夜 spam 控制；评论不算
			if title != "" && num > 0 && len(midNightSpam) == 2 {
				curHour := time.Now().Hour()
				startHour := goutils.MustInt(midNightSpam[0])
				endHour := goutils.MustInt(midNightSpam[1])
				// 比如 23 ~ 8（不包括 8 点）
				if startHour > endHour {
					if curHour >= startHour || curHour < endHour {
						logic.SpamRecord(context.EchoContext(ctx), user, num)
					}
				} else {
					// 比如 0 ~ 8（不包括 8 点）
					if curHour >= startHour && curHour < endHour {
						logic.SpamRecord(context.EchoContext(ctx), user, num)
					}
				}
			}

			if err := next(ctx); err != nil {
				return err
			}

			return nil
		}
	}
}

// hasSensitive 是否有敏感词
func hasSensitive(content, sensitive string) bool {
	if content == "" || sensitive == "" {
		return false
	}

	sensitives := strings.Split(sensitive, ",")

	for _, s := range sensitives {
		if strings.Contains(content, s) {
			return true
		}
	}

	return false
}

// hasSensitiveChar 是否包含敏感字（多个词都包含）
func hasSensitiveChar(title, sensitive string) bool {
	if title == "" || sensitive == "" {
		return false
	}

	sensitives := strings.Split(sensitive, "")

	for _, s := range sensitives {
		if !strings.Contains(title, s) {
			return false
		}
	}

	return true
}

func addBlackIP(ctx echo.Context) {
	ip := goutils.RemoteIp(ctx.Request())

	logic.DefaultRisk.AddBlackIP(ip)
}
