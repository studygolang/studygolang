// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package util

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/studygolang/studygolang/global"

	"github.com/polaris1119/goutils"
)

// 获取头像
func Gravatar(avatar string, emailI interface{}, size uint16, isHttps bool) string {
	gravatarDomain := "http://gravatar.com"
	if isHttps {
		gravatarDomain = "https://secure.gravatar.com"
	}
	cdnDomain := global.App.CanonicalCDN(isHttps)

	if avatar != "" {
		if strings.HasPrefix(avatar, "http") {
			return fmt.Sprintf("%s&s=%d", avatar, size)
		}
		return fmt.Sprintf("%savatar/%s?imageView2/2/w/%d", cdnDomain, avatar, size)
	}

	email, ok := emailI.(string)
	if !ok {
		return fmt.Sprintf("%savatar/gopher28.png?imageView2/2/w/%d", cdnDomain, size)
	}
	return fmt.Sprintf("%s/avatar/%s?s=%d", gravatarDomain, goutils.Md5(email), size)
}

func Max(x, y int) int {
	return int(math.Max(float64(x), float64(y)))
}

// 最小值，但不会小于0
func UMin(x, y int) int {
	if x < 0 || y < 0 {
		return 0
	}
	return int(math.Min(float64(x), float64(y)))
}

// 内嵌 Wide iframe 版
func EmbedWide(content string) string {
	if !strings.Contains(content, "&lt;iframe") {
		return content
	}

	reg := regexp.MustCompile(`&lt;iframe.*src=.*(https://wide\.b3log\.org/playground.*\.go).*/iframe&gt;`)
	return reg.ReplaceAllString(content, `<iframe src="$1?embed=true" width="100%" height="600"></iframe>`)
}
