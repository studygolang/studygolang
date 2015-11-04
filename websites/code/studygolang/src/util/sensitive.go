// Copyright 2015 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package util

import (
	"strings"
)

// HasSensitive 是否有敏感词
func HasSensitive(content, sensitive string) bool {
	if content == "" {
		return false
	}

	sensitives := strings.Split(sensitive, ",")

	for _, s := range sensitives {
		if strings.Contains(content, s) {
			return true
		}
	}

	return false
}

// HasSensitiveChar 是否包含敏感字（多个词都包含）
func HasSensitiveChar(title, sensitive string) bool {
	if title == "" || sensitive == "" {
		return false
	}

	sensitives := strings.Split(sensitive, "")

	for _, s := range sensitives {
		if !strings.Contains(title, s) {
			return false
		}
	}

	return true
}
