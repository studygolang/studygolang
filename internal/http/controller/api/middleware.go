package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/config"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
)

// sensitiveCheck 检查发布内容是否包含敏感词（等价于 master 的 middleware.Sensivite）
func sensitiveCheck(ctx echo.Context, me *model.Me) (ok bool) {
	ok = true

	titleSensitives := strings.Split(config.ConfigFile.MustValue("sensitive", "title"), ",")
	contentSensitives := config.ConfigFile.MustValue("sensitive", "content")

	content := ctx.FormValue("content")
	title := ctx.FormValue("title")

	if title != "" {
		for _, s := range titleSensitives {
			if hasSensitiveChar(title, s) {
				logic.DefaultUser.UpdateUserStatus(context.EchoContext(ctx), me.Uid, model.UserStatusFreeze)
				logger.Infoln("user=", me.Uid, "publish ad, title=", title, ". freeze")
				addBlackIP(ctx)
				ok = false
				return
			}
		}
	}

	if hasSensitive(title, contentSensitives) || hasSensitive(content, contentSensitives) {
		logic.DefaultUser.UpdateUserStatus(context.EchoContext(ctx), me.Uid, model.UserStatusFreeze)
		logger.Infoln("user=", me.Uid, "publish ad, title=", title, ";content=", content, ". freeze")
		addBlackIP(ctx)
		ok = false
		return
	}

	// 半夜 spam 控制
	midNightSpam := strings.Split(config.ConfigFile.MustValue("spam", "mid_night"), ",")
	num := config.ConfigFile.MustInt("spam", "num")
	if title != "" && num > 0 && len(midNightSpam) == 2 {
		curHour := time.Now().Hour()
		startHour := goutils.MustInt(midNightSpam[0])
		endHour := goutils.MustInt(midNightSpam[1])
		if startHour > endHour {
			if curHour >= startHour || curHour < endHour {
				logic.SpamRecord(context.EchoContext(ctx), me, num)
			}
		} else {
			if curHour >= startHour && curHour < endHour {
				logic.SpamRecord(context.EchoContext(ctx), me, num)
			}
		}
	}

	return
}

// balanceCheck 检查用户余额是否足够发布（等价于 master 的 middleware.BalanceCheck）
func balanceCheck(me *model.Me, isComment bool) bool {
	if me == nil {
		return false
	}
	if isComment {
		return me.Balance >= 5
	}
	return me.Balance >= 20
}

// publishNotice 发布后邮件通知站长（等价于 master 的 middleware.PublishNotice）
func publishNotice(ctx echo.Context, me *model.Me) {
	if me.IsRoot {
		return
	}

	title := ctx.FormValue("title")
	content := ctx.FormValue("content")
	if ctx.Request().Method == "POST" && (title != "" || content != "") {
		requestURI := ctx.Request().RequestURI
		go func() {
			user := logic.DefaultUser.FindOne(context.EchoContext(ctx), "is_root", 1)
			if user.Uid == 0 {
				return
			}

			emailContent := fmt.Sprintf("URI:%s<br/><h1>标题：%s</h1><br/>内容：%s", requestURI, title, content)
			logic.DefaultEmail.SendMail("网站有新内容产生", emailContent, []string{user.Email})
		}()
	}
}

// captchaCheck 检查发布时是否需要验证码（等价于 master 的 middleware.CheckCaptcha）
func captchaCheck(ctx echo.Context, me *model.Me) bool {
	if ctx.Request().Method != "POST" {
		return true
	}

	if !logic.NeedCaptcha(me) {
		return true
	}

	captchaID := ctx.FormValue("captcha_id")
	solution := ctx.FormValue("captcha_solution")
	if captchaID == "" || solution == "" {
		return false
	}

	return captcha.VerifyString(captchaID, solution)
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

// failSensitive 返回敏感词命中响应
func failSensitive(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 1,
		"msg":  "对不起，您的账号已被冻结！",
	})
}

// failBalance 返回余额不足响应
func failBalance(ctx echo.Context) error {
	return fail(ctx, "对不起，您的账号余额不足，可以领取初始资本！")
}

// failCaptcha 返回验证码错误响应
func failCaptcha(ctx echo.Context) error {
	return fail(ctx, "验证码错误，记得刷新验证码！")
}
