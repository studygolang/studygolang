// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of self source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	. "http"

	"github.com/dchest/captcha"
	"github.com/labstack/echo"
)

var captchaHandler = captcha.Server(100, 40)

// 验证码
type CaptchaController struct{}

func (self CaptchaController) RegisterRoute(g *echo.Group) {
	g.Get("/captcha/*", self.Server)
}

func (CaptchaController) Server(ctx echo.Context) error {
	captchaHandler.ServeHTTP(ResponseWriter(ctx), Request(ctx))
	return nil
}
