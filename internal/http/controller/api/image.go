// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/global"
	. "github.com/studygolang/studygolang/internal/http"
	"github.com/studygolang/studygolang/internal/logic"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/times"
)

type ImageController struct{}

func (self ImageController) RegisterRoute(g *echo.Group) {
	g.POST("/image/upload", self.Upload)
	g.POST("/image/paste_upload", self.PasteUpload)
	g.POST("/image/quick_upload", self.QuickUpload)
	g.POST("/image/transfer", self.Transfer)
}

// allowedMIMETypes 图片上传允许的 MIME 类型白名单
var allowedMIMETypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// handleImageUpload 提取公共图片上传逻辑（DRY）
// imgDir 为空时自动使用日期目录
func handleImageUpload(ctx echo.Context, fieldName string, imgDir string) (string, error) {
	file, fileHeader, err := Request(ctx).FormFile(fieldName)
	if err != nil {
		return "", fail(ctx, "非法文件上传")
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		return "", fail(ctx, "文件读取失败")
	}
	if len(buf) > logic.MaxImageSize {
		return "", fail(ctx, "文件太大")
	}

	// 检测上传文件的实际 MIME type，防止伪造扩展名攻击
	mimeType := http.DetectContentType(buf)
	if !allowedMIMETypes[mimeType] {
		return "", fail(ctx, "不支持的图片格式，仅支持 JPEG、PNG、GIF、WebP")
	}

	if imgDir == "" {
		imgDir = times.Format("ymd")
	}
	// Seek 失败时 reader 处于 EOF，UploadImage 上传会得到空文件，
	// 但 MD5（基于 buf）正确，导致图片记录与存储内容不一致。
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fail(ctx, "文件重置失败")
	}
	return logic.DefaultUploader.UploadImage(context.EchoContext(ctx), file, imgDir, buf, filepath.Ext(fileHeader.Filename))
}

// Upload 通用图片上传（multipart form: img，可选 avatar 参数指定头像目录）
func (ImageController) Upload(ctx echo.Context) error {
	if _, err := requireAuth(ctx); err != nil {
		return err
	}

	// 头像使用专用目录，否则使用日期目录
	imgDir := ""
	if ctx.FormValue("avatar") != "" {
		imgDir = "avatar"
	}

	path, err := handleImageUpload(ctx, "img", imgDir)
	if err != nil {
		return err
	}

	cdnDomain := global.App.CanonicalCDN(CheckIsHttps(ctx))

	return success(ctx, map[string]interface{}{
		"url": cdnDomain + path,
		"uri": path,
	})
}

// PasteUpload 粘贴上传图片（multipart form: imageFile）
func (ImageController) PasteUpload(ctx echo.Context) error {
	if _, err := requireAuth(ctx); err != nil {
		return err
	}

	path, err := handleImageUpload(ctx, "imageFile", "")
	if err != nil {
		return err
	}

	cdnDomain := global.App.CanonicalCDN(CheckIsHttps(ctx))

	return success(ctx, map[string]interface{}{
		"success": 1,
		"message": cdnDomain + path,
	})
}

// QuickUpload 编辑器快速上传（multipart form: upload）
func (ImageController) QuickUpload(ctx echo.Context) error {
	if _, err := requireAuth(ctx); err != nil {
		return err
	}

	path, err := handleImageUpload(ctx, "upload", "")
	if err != nil {
		return err
	}

	cdnDomain := global.App.CanonicalCDN(CheckIsHttps(ctx))

	return success(ctx, map[string]interface{}{
		"uploaded": 1,
		"url":      cdnDomain + path,
	})
}

// Transfer 通过 URL 转存图片（form: url）
func (ImageController) Transfer(ctx echo.Context) error {
	if _, err := requireAuth(ctx); err != nil {
		return err
	}

	origUrl := ctx.FormValue("url")
	if origUrl == "" {
		return fail(ctx, "url 不能为空")
	}

	path, err := logic.DefaultUploader.TransferUrl(context.EchoContext(ctx), origUrl)
	if err != nil {
		return fail(ctx, "文件上传失败")
	}

	cdnDomain := global.App.CanonicalCDN(CheckIsHttps(ctx))

	return success(ctx, map[string]interface{}{
		"url": cdnDomain + path,
	})
}
