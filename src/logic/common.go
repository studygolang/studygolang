// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/gorilla/schema"
	"github.com/polaris1119/logger"
	"golang.org/x/net/context"
)

var schemaDecoder = schema.NewDecoder()

func init() {
	schemaDecoder.SetAliasTag("json")
	schemaDecoder.IgnoreUnknownKeys(true)
}

var NotModifyAuthorityErr = errors.New("没有修改权限")

func GetLogger(ctx context.Context) *logger.Logger {
	if ctx == nil {
		return logger.New(os.Stdout)
	}

	_logger, ok := ctx.Value("logger").(*logger.Logger)
	if ok {
		return _logger
	}

	return logger.New(os.Stdout)
}

// parseAtUser 解析 @某人
func parseAtUser(ctx context.Context, content string) string {
	reg := regexp.MustCompile(`@([^\s@]{4,20})`)
	return reg.ReplaceAllStringFunc(content, func(matched string) string {
		username := matched[1:]

		// 校验 username 是否存在
		user := DefaultUser.FindOne(ctx, "username", username)
		if user.Username != username {
			return matched
		}
		return fmt.Sprintf(`<a href="/user/%s" title="%s">%s</a>`, username, matched, matched)
	})
}
