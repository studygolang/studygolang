// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package middleware

import (
	"fmt"

	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
)

// PublishNotice 用于 echo 框架，用户发布内容邮件通知站长
func PublishNotice() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if err := next(ctx); err != nil {
				return err
			}

			curUser := ctx.Get("user").(*model.Me)
			if curUser.IsRoot {
				return nil
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

					content = fmt.Sprintf("URI:%s<br/><h1>标题：%s</h1><br/>内容：%s", requestURI, title, content)
					logic.DefaultEmail.SendMail("网站有新内容产生", content, []string{user.Email})
				}()
			}

			return nil
		}
	}
}
