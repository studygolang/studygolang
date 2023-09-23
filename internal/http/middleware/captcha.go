// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package middleware

import (
	"net/http"

	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"
	"github.com/studygolang/studygolang/util"

	"github.com/dchest/captcha"
	echo "github.com/labstack/echo/v4"
)

// CheckCaptcha 用于 echo 框架校验发布验证码
func CheckCaptcha() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			curUser := ctx.Get("user").(*model.Me)

			if ctx.Request().Method == "POST" {
				if logic.NeedCaptcha(curUser) {
					captchaId := ctx.FormValue("captchaid")
					if !captcha.VerifyString(captchaId, ctx.FormValue("captchaSolution")) {
						util.SetCaptcha(captchaId)
						return ctx.String(http.StatusOK, `{"ok":0,"error":"验证码错误，记得刷新验证码！"}`)
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
