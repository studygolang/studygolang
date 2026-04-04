// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"github.com/dchest/captcha"
	echo "github.com/labstack/echo/v4"
)

type CaptchaController struct{}

func (self CaptchaController) RegisterRoute(g *echo.Group) {
	g.GET("/captcha/new", self.NewCaptcha)
	g.GET("/captcha/*", self.Server)
}

// NewCaptcha 生成新验证码，返回 captcha_id
func (CaptchaController) NewCaptcha(ctx echo.Context) error {
	id := captcha.New()
	return success(ctx, map[string]interface{}{
		"captcha_id": id,
	})
}

// Server 提供验证码图片（透传给 dchest/captcha 库）
func (CaptchaController) Server(ctx echo.Context) error {
	handler := captcha.Server(100, 40)
	handler.ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}
