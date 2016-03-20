// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package util

import (
	"errors"
	"net/http"
)

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

const XRequestedWith = "X-Requested-With"

func IsAjax(req *http.Request) bool {
	return req.Header.Get(XRequestedWith) == "XMLHttpRequest"
}
