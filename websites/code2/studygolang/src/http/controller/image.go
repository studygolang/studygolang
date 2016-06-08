// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of self source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"io/ioutil"
	"logic"
	"os"
	"path/filepath"

	. "http"

	"github.com/labstack/echo"
	"github.com/polaris1119/times"
)

// 图片处理
type ImageController struct{}

func (self ImageController) RegisterRoute(g *echo.Group) {
	g.POST("/image/upload", self.Upload)
	g.Match([]string{"GET", "POST"}, "/image/transfer", self.Transfer)
}

// Upload 上传图片
func (ImageController) Upload(ctx echo.Context) error {
	objLogger := getLogger(ctx)

	file, fileHeader, err := Request(ctx).FormFile("img")
	if err != nil {
		objLogger.Errorln("upload error:", err)
		return fail(ctx, 1, "非法文件上传！")
	}
	defer file.Close()

	// 如果是临时文件，存在硬盘中，则是 *os.File（大于32M），直接报错
	if _, ok := file.(*os.File); ok {
		objLogger.Errorln("upload error:file too large!")
		return fail(ctx, 2, "文件太大！")
	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return fail(ctx, 3, "文件读取失败！")
	}
	if len(buf) > logic.MaxImageSize {
		return fail(ctx, 4, "文件太大！")
	}

	imgDir := times.Format("ymd")
	if ctx.FormValue("avatar") != "" {
		imgDir = "avatar"
	}

	path, err := logic.DefaultUploader.UploadImage(ctx, file, imgDir, buf, filepath.Ext(fileHeader.Filename))
	if err != nil {
		return fail(ctx, 5, "文件上传失败！")
	}

	return success(ctx, map[string]interface{}{"uri": path})
}

// Transfer 转换图片：通过 url 从远程下载图片然后转存到七牛
func (ImageController) Transfer(ctx echo.Context) error {
	origUrl := ctx.FormValue("url")
	if origUrl == "" {
		return fail(ctx, 1, "url不能为空！")
	}

	path, err := logic.DefaultUploader.TransferUrl(ctx, origUrl)
	if err != nil {
		return fail(ctx, 2, "文件上传失败！")
	}

	return success(ctx, map[string]interface{}{"uri": logic.ImageDomain + path})
}
