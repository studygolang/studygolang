// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package logic

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	gio "io"
	"io/ioutil"
	"mime"
	"model"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	. "db"

	"github.com/polaris1119/config"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
	"github.com/polaris1119/times"
	"github.com/qiniu/api.v6/conf"
	"github.com/qiniu/api.v6/io"
	"github.com/qiniu/api.v6/rs"
)

const (
	MaxImageSize = 5 << 20 // 5M
)

// http://developer.qiniu.com/code/v6/sdk/go-sdk-6.html
type UploaderLogic struct {
	bucketName string

	uptoken   string
	tokenTime time.Time
	locker    sync.RWMutex
}

var DefaultUploader = &UploaderLogic{}

func (this *UploaderLogic) InitQiniu() {
	conf.ACCESS_KEY = config.ConfigFile.MustValue("qiniu", "access_key")
	conf.SECRET_KEY = config.ConfigFile.MustValue("qiniu", "secret_key")
	conf.UP_HOST = config.ConfigFile.MustValue("qiniu", "up_host", conf.UP_HOST)
	this.bucketName = config.ConfigFile.MustValue("qiniu", "bucket_name")
}

// 生成上传凭证
func (this *UploaderLogic) genUpToken() {
	// 避免服务器时间不同步，45分钟就更新 token
	if this.uptoken != "" && this.tokenTime.Add(45*time.Minute).After(time.Now()) {
		return
	}

	putPolicy := rs.PutPolicy{
		Scope: this.bucketName,
		// CallbackUrl:  callbackUrl,
		// CallbackBody: callbackBody,
		// ReturnUrl:    returnUrl,
		// ReturnBody:   returnBody,
		// AsyncOps:     asyncOps,
		// EndUser:      endUser,
		// 指定上传凭证有效期（默认1小时）
		// Expires:      expires,
	}

	this.locker.Lock()
	this.uptoken = putPolicy.Token(nil)
	this.locker.Unlock()
	this.tokenTime = time.Now()
}

func (this *UploaderLogic) uploadLocalFile(localFile, key string) (err error) {
	this.genUpToken()

	var ret io.PutRet
	var extra = &io.PutExtra{
		// Params:   params,
		// MimeType: mieType,
		// Crc32:    crc32,
		// CheckCrc: CheckCrc,
	}

	// ret       变量用于存取返回的信息，详情见 io.PutRet
	// uptoken   为业务服务器生成的上传口令
	// key       为文件存储的标识(文件名)
	// localFile 为本地文件名
	// extra     为上传文件的额外信息，详情见 io.PutExtra，可选
	err = io.PutFile(nil, &ret, this.uptoken, key, localFile, extra)

	if err != nil {
		//上传产生错误
		logger.Errorln("io.PutFile failed:", err)
		return
	}

	//上传成功，处理返回值
	logger.Debugln(ret.Hash, ret.Key)

	return
}

func (this *UploaderLogic) uploadMemoryFile(r gio.Reader, key string, size int) (err error) {
	this.genUpToken()

	var ret io.PutRet
	var extra = &io.PutExtra{
		// Params:   params,
		// MimeType: mieType,
		// Crc32:    crc32,
		// CheckCrc: CheckCrc,
	}

	// ret       变量用于存取返回的信息，详情见 io.PutRet
	// uptoken   为业务服务器端生成的上传口令
	// key       为文件存储的标识
	// r         为io.Reader类型，用于从其读取数据
	// extra     为上传文件的额外信息,可为空， 详情见 io.PutExtra, 可选
	err = io.Put2(nil, &ret, this.uptoken, key, r, int64(size), extra)

	// 上传产生错误
	if err != nil {
		logger.Errorln("io.Put failed:", err)

		errInfo := make(map[string]interface{})
		err = json.Unmarshal([]byte(err.Error()), &errInfo)
		if err != nil {
			logger.Errorln("io.Put Unmarshal failed:", err)
			return
		}

		code, ok := errInfo["code"]
		if ok && code == 614 {
			err = nil
		}

		return
	}

	// 上传成功，处理返回值
	logger.Debugln(ret.Hash, ret.Key)

	return
}

func (this *UploaderLogic) UploadImage(ctx context.Context, reader gio.Reader, imgDir string, buf []byte, ext string) (string, error) {
	objLogger := GetLogger(ctx)

	md5 := goutils.Md5Buf(buf)
	objImage, err := this.findImage(md5)
	if err != nil {
		objLogger.Errorln("find image:", md5, "error:", err)
		return "", err
	}

	if objImage.Pid > 0 {
		return objImage.Path, nil
	}

	path := imgDir + "/" + md5 + ext
	if err = this.uploadMemoryFile(reader, path, len(buf)); err != nil {
		return "", err
	}

	go this.saveImage(buf, path)

	return path, nil
}

// TransferUrl 将外站图片URL转为本站，如果失败，返回原图
func (this *UploaderLogic) TransferUrl(ctx context.Context, origUrl string, prefixs ...string) (string, error) {
	if origUrl == "" || strings.Contains(origUrl, WebsiteSetting.Domain) {
		return origUrl, errors.New("origin image is empty or is " + WebsiteSetting.Domain)
	}

	resp, err := http.Get(origUrl)
	if err != nil {
		return origUrl, errors.New("获取图片失败")
	}
	defer resp.Body.Close()

	buf, _ := ioutil.ReadAll(resp.Body)

	md5 := goutils.Md5Buf(buf)
	objImage, err := this.findImage(md5)
	if err != nil {
		logger.Errorln("find image:", md5, "error:", err)
		return origUrl, err
	}

	if objImage.Pid > 0 {
		return objImage.Path, nil
	}

	ext := filepath.Ext(origUrl)
	if ext == "" {
		contentType := http.DetectContentType(buf)
		exts, _ := mime.ExtensionsByType(contentType)
		if len(exts) > 0 {
			ext = exts[0]
		}
	}

	prefix := times.Format("ymd")
	if len(prefixs) > 0 {
		prefix = prefixs[0]
	}
	path := prefix + "/" + md5 + ext
	reader := bytes.NewReader(buf)

	if len(buf) > MaxImageSize {
		return origUrl, errors.New("文件太大")
	}

	err = this.uploadMemoryFile(reader, path, len(buf))
	if err != nil {
		return origUrl, err
	}

	go this.saveImage(buf, path)

	return path, nil
}

func (this *UploaderLogic) findImage(md5 string) (*model.Image, error) {
	objImage := &model.Image{}
	_, err := MasterDB.Where("md5=?", md5).Get(objImage)
	if err != nil {
		return nil, err
	}

	return objImage, nil
}

func (this *UploaderLogic) saveImage(buf []byte, path string) {
	objImage := &model.Image{
		Path: path,
		Md5:  goutils.Md5Buf(buf),
		Size: len(buf),
	}

	reader := bytes.NewReader(buf)
	img, _, err := image.Decode(reader)
	if err != nil {
		logger.Errorln("image decode err:", err)
	} else {
		objImage.Width = img.Bounds().Dx()
		objImage.Height = img.Bounds().Dy()
	}

	_, err = MasterDB.Insert(objImage)
	if err != nil {
		logger.Errorln("image insert err:", err)
	}
}
