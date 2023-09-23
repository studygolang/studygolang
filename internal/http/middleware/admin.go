// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package middleware

import (
	"net/http"

	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
)

// AdminAuth 用于 echo 框架的判断用户是否有管理后台权限
func AdminAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			user := ctx.Get("user").(*model.Me)
			if !user.IsAdmin {
				return ctx.HTML(http.StatusForbidden, `403 Forbidden`)
			}

			if err := next(ctx); err != nil {
				return err
			}

			return nil
		}
	}
}
