// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: momaek momaek17@gmail.com

package middleware

import (
	"net/http"
	"net/url"
	"util"

	"github.com/labstack/echo"
)

// ErrorRet 如果是 ajax 请求，返回前端错误信息的通用结构体
type ErrorRet struct {
	OK    int    `json:"ok"`
	Error string `json:"error"`
}

// CsrfRefererFilter 通过 referer 过滤csrf请求
func CsrfRefererFilter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			req := ctx.Request()
			can := false

			defer func() {
				if can {
					err = next(ctx)
				} else {
					if util.IsAjax(ctx) {
						ctx.JSON(499, &ErrorRet{0, "CSRF Detected"})
					} else {
						ctx.String(499, "CSRF Detected")
					}
				}
			}()

			switch req.Method() {
			case http.MethodGet, http.MethodHead:
				can = true
				return
			}

			referer := req.Referer()
			if len(referer) == 0 {
				return
			}

			u, err := url.Parse(referer)
			if err != nil {
				return
			}

			if u.Host != req.Host() {
				return
			}

			can = true
			return
		}
	}
}
