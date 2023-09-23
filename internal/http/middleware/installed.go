// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package middleware

import (
	"net/http"
	"strings"

	"github.com/studygolang/studygolang/db"

	echo "github.com/labstack/echo/v4"
)

// Installed 用于 echo 框架，判断是否已经安装了
func Installed(filterPrefixs []string) echo.MiddlewareFunc {
	filterPrefixs = append(filterPrefixs, "/install")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if db.MasterDB == nil {
				shouldRedirect := true

				uri := ctx.Request().RequestURI
				for _, prefix := range filterPrefixs {
					if strings.HasPrefix(uri, prefix) {
						shouldRedirect = false
						break
					}
				}

				if shouldRedirect {
					return ctx.Redirect(http.StatusSeeOther, "/install")
				}
			}
			if err := next(ctx); err != nil {
				return err
			}

			return nil
		}
	}
}
