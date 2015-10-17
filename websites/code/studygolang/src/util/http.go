// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package util

import (
	"errors"
	"net/http"

	"github.com/gorilla/context"
)

// Redirect 重定向到指定的 uri
func Redirect(rw http.ResponseWriter, req *http.Request, uri string) {
	// 避免跳转，context中没有清除
	context.Clear(req)

	http.Redirect(rw, req, uri, http.StatusFound)
}

// FetchRealUrl 获取链接真实的URL（获取重定向一次的结果URL）
func FetchRealUrl(uri string) (realUrl string) {

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			realUrl = req.URL.String()
			return errors.New("util fetch real url")
		},
	}

	resp, err := client.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return uri
}
