// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package model

import (
	"strings"

	"github.com/polaris1119/config"
	"github.com/polaris1119/keyword"
)

func init() {
	go keyword.Extractor.Init(keyword.DefaultProps, true, config.ROOT+"/data/programming.txt,"+config.ROOT+"/data/dictionary.txt")
}

// AutoTag 自动生成 tag
func AutoTag(title, content string, num int) string {
	return strings.Join(keyword.ExtractWithTitle(title, content, num), ",")
}
