// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"encoding/json"
	"logic"
	"net/http"
	"strings"

	"github.com/polaris1119/goutils"

	. "http"

	"github.com/labstack/echo"
	"github.com/polaris1119/logger"
	"github.com/polaris1119/nosql"
)

func getLogger(ctx echo.Context) *logger.Logger {
	return logic.GetLogger(ctx)
}

// render html 输出
func render(ctx echo.Context, contentTpl string, data map[string]interface{}) error {
	return Render(ctx, contentTpl, data)
}

func success(ctx echo.Context, data interface{}) error {
	result := map[string]interface{}{
		"ok":   1,
		"msg":  "操作成功",
		"data": data,
	}

	b, err := json.Marshal(result)
	if err != nil {
		return err
	}

	oldETag := ctx.Request().Header().Get("If-None-Match")
	if strings.HasPrefix(oldETag, "W/") {
		oldETag = oldETag[2:]
	}
	newETag := goutils.Md5Buf(b)
	if oldETag == newETag {
		return ctx.NoContent(http.StatusNotModified)
	}

	go func(b []byte) {
		if cacheKey := ctx.Get(nosql.CacheKey); cacheKey != nil {
			nosql.DefaultLRUCache.CompressAndAdd(cacheKey, b, nosql.NewCacheData())
		}
	}(b)

	if ctx.Response().Committed() {
		getLogger(ctx).Flush()
		return nil
	}

	ctx.Response().Header().Add("ETag", newETag)

	return ctx.JSONBlob(http.StatusOK, b)
}

func fail(ctx echo.Context, code int, msg string) error {
	if ctx.Response().Committed() {
		getLogger(ctx).Flush()
		return nil
	}

	result := map[string]interface{}{
		"ok":    0,
		"error": msg,
	}

	getLogger(ctx).Errorln("operate fail:", result)

	return ctx.JSON(http.StatusOK, result)
}
