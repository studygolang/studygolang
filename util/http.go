// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	echo "github.com/labstack/echo/v4"
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

func IsAjax(ctx echo.Context) bool {
	if ctx.Request().Header.Get(XRequestedWith) == "XMLHttpRequest" {
		return true
	}
	return false
}

func DoGet(url string, extras ...int) (body []byte, err error) {
	// 默认重试次数
	num := 3
	// 是否 sleep
	sleep := false

	switch len(extras) {
	case 0:
	case 1:
		num = extras[0]
	case 2:
		num = extras[0]
		sleep = true
	default:
	}

	for i := 0; i < num; i++ {
		body, err = doGet(url)
		if err == nil {
			break
		}
		if sleep {
			time.Sleep(time.Second * time.Duration(2*i+1))
		}
	}
	return
}

func doGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("StatusCode is not 200")
	}

	return ioutil.ReadAll(resp.Body)
}

func DoPost(url string, data url.Values, extras ...int) (body []byte, err error) {
	// 默认重试次数
	num := 3
	// 是否sleep
	sleep := false

	switch len(extras) {
	case 0:
	case 1:
		num = extras[0]
	case 2:
		num = extras[0]
		sleep = true
	default:
	}

	for i := 0; i < num; i++ {
		body, err = doPost(url, data)
		if err == nil {
			break
		}
		if sleep {
			time.Sleep(time.Second * time.Duration(2*i+1))
		}
	}
	return
}

func doPost(url string, data url.Values) ([]byte, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := httpClient.PostForm(url, data)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("StatusCode is not 200")
	}

	return ioutil.ReadAll(resp.Body)
}

func DoPostRaw(url string, bodyType string, data interface{}, extras ...int) (body []byte, err error) {
	// 默认重试次数
	num := 3
	// 是否sleep
	sleep := false

	switch len(extras) {
	case 0:
	case 1:
		num = extras[0]
	case 2:
		num = extras[0]
		sleep = true
	default:
	}

	for i := 0; i < num; i++ {
		body, err = doPostRaw(url, bodyType, data)
		if err == nil {
			break
		}
		if sleep {
			time.Sleep(time.Second * time.Duration(2*i+1))
		}
	}
	return
}

func doPostRaw(url, bodyType string, data interface{}) ([]byte, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	bodyByte, err := json.Marshal(data)

	if nil == err {
		resp, err := httpClient.Post(url, bodyType, bytes.NewReader(bodyByte))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return nil, errors.New("StatusCode is not 200")
		}

		return ioutil.ReadAll(resp.Body)
	}

	return nil, err
}
