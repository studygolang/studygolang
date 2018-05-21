// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/axgle/mahonia"
	"github.com/polaris1119/keyword"
	"github.com/polaris1119/nosql"
)

type item struct {
	Score float32
	Tag   string
}
type resData struct {
	Log_id int
	Items  []item
}

type resTokenData struct {
	Access_token   string
	Scope          string
	Session_key    string
	Refresh_token  string
	Session_secret string
	Expires_in     int
}

// AutoTag 自动生成 tag
func AutoTag(title, content string, num int) string {
	defer func() {
		recover()
	}()
	key := "baidu_access_token"
	client_id := "geogVB0En5UM936L6Llf5EWr"
	client_secret := "ec120xF6SItrEU4sjZk5s3av61eWde2X&"

	// 取百度token
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()
	token := redisClient.GET(key)
	if token == "" {
		resp, err := http.Get("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=" + client_id + "client_secret=" + client_secret)
		if err != nil {
			fmt.Println(err)
			return strings.Join(keyword.ExtractWithTitle(title, content, num), ",")
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return strings.Join(keyword.ExtractWithTitle(title, content, num), ",")
		}
		var data resTokenData
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println(err)
			return strings.Join(keyword.ExtractWithTitle(title, content, num), ",")
		}
		token = data.Access_token
		err = redisClient.SET(key, token, data.Expires_in)
		if err != nil {
			fmt.Println(err)
			return strings.Join(keyword.ExtractWithTitle(title, content, num), ",")
		}
	}

	// 转成GBK
	titleGBK := mahonia.NewEncoder("gbk").ConvertString(title)
	contentGBK := mahonia.NewEncoder("gbk").ConvertString(content)
	post := "{\"title\":\"" + string(titleGBK) + "\",\"content\":\"" + string(contentGBK) + "\"}"
	url := "https://aip.baidubce.com/rpc/2.0/nlp/v1/keyword?access_token=" + token
	jsonStr := []byte(post)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return strings.Join(keyword.ExtractWithTitle(title, content, num), ",")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// 解析返回值
	var data resData
	err = json.Unmarshal(ConvertToByte(string(body), "gbk", "utf8"), &data)
	if err != nil {
		return strings.Join(keyword.ExtractWithTitle(title, content, num), ",")
	}
	var words []string
	var length int
	if len(data.Items) > num {
		length = num
	} else {
		length = len(data.Items)
	}
	for i := 0; i < length; i++ {
		word := data.Items[i].Tag
		words = append(words, string(word))
	}
	return strings.Join(words, ",")
}

func ConvertToByte(src string, srcCode string, targetCode string) []byte {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(targetCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	return cdata
}
