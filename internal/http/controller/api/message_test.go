// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unit
// +build unit

package api

import (
	"testing"
)

// TestIsOriginAllowed 测试 CSRF 防护的 origin 白名单校验
// 重点验证 M4 修复：通配符子域名匹配不再被 evil<domain>.com 利用
func TestIsOriginAllowed(t *testing.T) {
	tests := []struct {
		name           string
		origin         string
		allowedOrigins []string
		expected       bool
	}{
		// 精确匹配
		{"exact match", "https://studygolang.com", []string{"https://studygolang.com"}, true},
		{"scheme mismatch rejected", "http://studygolang.com", []string{"https://studygolang.com"}, false},
		{"port mismatch rejected", "https://studygolang.com:8443", []string{"https://studygolang.com"}, false},
		// 通配符子域名
		{"subdomain match", "https://www.studygolang.com", []string{"*.studygolang.com"}, true},
		{"nested subdomain match", "https://api.v2.studygolang.com", []string{"*.studygolang.com"}, true},
		{"subdomain with port match", "https://www.studygolang.com:8443", []string{"*.studygolang.com"}, true},
		// M4 安全关键测试：evil<domain>.com 攻击
		{"evil suffix rejected (M4 fix)", "https://evilstudygolang.com", []string{"*.studygolang.com"}, false},
		{"evil prefix rejected", "https://studygolang.com.evil.com", []string{"*.studygolang.com"}, false},
		{"evil in middle rejected", "https://a.studygolang.com.evil.com", []string{"*.studygolang.com"}, false},
		// 根域名本身不匹配子域名通配符（避免攻击者注册 studygolang.com 后用 *.studygolang.com）
		{"root not matched by wildcard", "https://studygolang.com", []string{"*.studygolang.com"}, false},
		// 解析失败
		{"malformed rejected", "not a url", []string{"https://studygolang.com"}, false},
		{"empty rejected", "", []string{"https://studygolang.com"}, false},
		// referer 形式
		{"referer with path", "https://www.studygolang.com/some/path?q=1", []string{"*.studygolang.com"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOriginAllowed(tt.origin, tt.allowedOrigins); got != tt.expected {
				t.Errorf("isOriginAllowed(%q, %v) = %v, want %v",
					tt.origin, tt.allowedOrigins, got, tt.expected)
			}
		})
	}
}
