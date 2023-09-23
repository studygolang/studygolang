// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of self source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/global"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/times"
)

// 图片处理
type ImageController struct{}

func (self ImageController) RegisterRoute(g *echo.Group) {
	// todo 这三个upload差不多啊
	g.POST("/image/upload", self.Upload)
	g.POST("/image/paste_upload", self.PasteUpload)
	g.POST("/image/quick_upload", self.QuickUpload)
	g.Match([]string{"GET", "POST"}, "/image/transfer", self.Transfer)
}

// PasteUpload jquery 粘贴上传图片
func (self ImageController) PasteUpload(ctx echo.Context) error {
	objLogger := getLogger(ctx)

	file, fileHeader, err := Request(ctx).FormFile("imageFile")
	if err != nil {
		objLogger.Errorln("upload error:", err)
		return self.pasteUploadFail(ctx, err.Error())
	}
	defer file.Close()

	// 如果是临时文件，存在硬盘中，则是 *os.File（大于32M），直接报错
	if _, ok := file.(*os.File); ok {
		objLogger.Errorln("upload error:file too large!")
		return self.pasteUploadFail(ctx, "文件太大！")
	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return self.pasteUploadFail(ctx, "文件读取失败！")
	}
	if len(buf) > logic.MaxImageSize {
		return self.pasteUploadFail(ctx, "文件太大！")
	}

	imgDir := times.Format("ymd")
	file.Seek(0, io.SeekStart)
	path, err := logic.DefaultUploader.UploadImage(context.EchoContext(ctx), file, imgDir, buf, filepath.Ext(fileHeader.Filename))
	if err != nil {
		return self.pasteUploadFail(ctx, "文件上传失败！")
	}

	cdnDomain := global.App.CanonicalCDN(CheckIsHttps(ctx))

	data := map[string]interface{}{
		"success": 1,
		"message": cdnDomain + path,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ctx.JSONBlob(http.StatusOK, b)
}

// QuickUpload CKEditor 编辑器，上传图片，支持粘贴方式上传
func (self ImageController) QuickUpload(ctx echo.Context) error {
	objLogger := getLogger(ctx)

	file, fileHeader, err := Request(ctx).FormFile("upload")
	if err != nil {
		objLogger.Errorln("upload error:", err)
		return self.quickUploadFail(ctx, err.Error())
	}
	defer file.Close()

	// 如果是临时文件，存在硬盘中，则是 *os.File（大于32M），直接报错
	if _, ok := file.(*os.File); ok {
		objLogger.Errorln("upload error:file too large!")
		return self.quickUploadFail(ctx, "文件太大！")
	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return self.quickUploadFail(ctx, "文件读取失败！")
	}
	if len(buf) > logic.MaxImageSize {
		return self.quickUploadFail(ctx, "文件太大！")
	}

	fileName := goutils.Md5Buf(buf) + filepath.Ext(fileHeader.Filename)
	imgDir := times.Format("ymd")
	file.Seek(0, io.SeekStart)
	path, err := logic.DefaultUploader.UploadImage(context.EchoContext(ctx), file, imgDir, buf, filepath.Ext(fileHeader.Filename))
	if err != nil {
		return self.quickUploadFail(ctx, "文件上传失败！")
	}

	cdnDomain := global.App.CanonicalCDN(CheckIsHttps(ctx))

	data := map[string]interface{}{
		"uploaded": 1,
		"fileName": fileName,
		"url":      cdnDomain + path,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ctx.JSONBlob(http.StatusOK, b)
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

	cdnDomain := global.App.CanonicalCDN(CheckIsHttps(ctx))

	file.Seek(0, io.SeekStart)
	path, err := logic.DefaultUploader.UploadImage(context.EchoContext(ctx), file, imgDir, buf, filepath.Ext(fileHeader.Filename))
	if err != nil {
		return fail(ctx, 5, "文件上传失败！")
	}

	return success(ctx, map[string]interface{}{"url": cdnDomain + path, "uri": path})
}

// Transfer 转换图片：通过 url 从远程下载图片然后转存到七牛
func (ImageController) Transfer(ctx echo.Context) error {
	origUrl := ctx.FormValue("url")
	if origUrl == "" {
		return fail(ctx, 1, "url不能为空！")
	}

	path, err := logic.DefaultUploader.TransferUrl(context.EchoContext(ctx), origUrl)
	if err != nil {
		return fail(ctx, 2, "文件上传失败！")
	}

	cdnDomain := global.App.CanonicalCDN(CheckIsHttps(ctx))

	return success(ctx, map[string]interface{}{"url": cdnDomain + path})
}

func (ImageController) quickUploadFail(ctx echo.Context, message string) error {
	data := map[string]interface{}{
		"uploaded": 0,
		"error": map[string]string{
			"message": message,
		},
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ctx.JSONBlob(http.StatusOK, b)
}

func (ImageController) pasteUploadFail(ctx echo.Context, message string) error {
	data := map[string]interface{}{
		"success": 0,
		"message": message,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ctx.JSONBlob(http.StatusOK, b)
}
