// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package controller

import (
	"logic"
	"net/http"

	"github.com/labstack/echo"
)

type CommentController struct{}

func (self CommentController) RegisterRoute(e *echo.Echo) {
	e.Get("/at/users", echo.HandlerFunc(self.AtUsers))
}

// AtUsers 评论或回复 @ 某人 suggest
func (CommentController) AtUsers(ctx echo.Context) error {
	term := ctx.QueryParam("term")
	users := logic.DefaultUser.GetUserMentions(term, 10)
	return ctx.JSON(http.StatusOK, users)
}
