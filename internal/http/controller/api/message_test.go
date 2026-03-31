// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unit
// +build unit

package api

import (
	"testing"
)

// TestIsSameOrigin 测试同源检查函数
func TestIsSameOrigin(t *testing.T) {
	tests := []struct {
		name           string
		originOrReferer string
		targetHost     string
		wantSame       bool
	}{
		{
			name:           "same origin - http",
			originOrReferer: "http://studygolang.com/messages",
			targetHost:     "studygolang.com",
			wantSame:       true,
		},
		{
			name:           "same origin - https",
			originOrReferer: "https://studygolang.com/messages",
			targetHost:     "studygolang.com",
			wantSame:       true,
		},
		{
			name:           "same origin - with port",
			originOrReferer: "http://studygolang.com:8080/messages",
			targetHost:     "studygolang.com:8080",
			wantSame:       true,
		},
		{
			name:           "same origin - different port",
			originOrReferer: "http://studygolang.com:8080/messages",
			targetHost:     "studygolang.com:9090",
			wantSame:       true, // 只比较主机名，忽略端口
		},
		{
			name:           "different origin",
			originOrReferer: "http://evil.com/messages",
			targetHost:     "studygolang.com",
			wantSame:       false,
		},
		{
			name:           "invalid url",
			originOrReferer: "not a url",
			targetHost:     "studygolang.com",
			wantSame:       false,
		},
		{
			name:           "empty origin",
			originOrReferer: "",
			targetHost:     "studygolang.com",
			wantSame:       false,
		},
		{
			name:           "subdomain - different",
			originOrReferer: "http://sub.studygolang.com/messages",
			targetHost:     "studygolang.com",
			wantSame:       false,
		},
		{
			name:           "localhost - same",
			originOrReferer: "http://localhost:3000/messages",
			targetHost:     "localhost:8090",
			wantSame:       true, // 只比较主机名
		},
		{
			name:           "127.0.0.1 - same",
			originOrReferer: "http://127.0.0.1:3000/messages",
			targetHost:     "127.0.0.1:8090",
			wantSame:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isSameOrigin(tt.originOrReferer, tt.targetHost)
			if got != tt.wantSame {
				t.Errorf("isSameOrigin(%q, %q) = %v, want %v",
					tt.originOrReferer, tt.targetHost, got, tt.wantSame)
			}
		})
	}
}

// TestCSRFProtectionScenarios 测试 CSRF 保护场景
func TestCSRFProtectionScenarios(t *testing.T) {
	// 场景 1: 正常的同源请求
	t.Run("valid same-origin request", func(t *testing.T) {
		origin := "https://studygolang.com"
		host := "studygolang.com"
		if !isSameOrigin(origin, host) {
			t.Error("Same-origin request should be allowed")
		}
	})

	// 场景 2: 跨站请求伪造
	t.Run("cross-site request forgery", func(t *testing.T) {
		referer := "http://evil.com/attack-page"
		host := "studygolang.com"
		if isSameOrigin(referer, host) {
			t.Error("CSRF request should be blocked")
		}
	})

	// 场景 3: 子域名攻击
	t.Run("subdomain attack", func(t *testing.T) {
		origin := "http://attacker.studygolang.com"
		host := "studygolang.com"
		if isSameOrigin(origin, host) {
			t.Error("Different subdomain should be blocked")
		}
	})

	// 场景 4: 本地开发环境
	t.Run("localhost development", func(t *testing.T) {
		origin := "http://localhost:3000"
		host := "localhost:8090"
		if !isSameOrigin(origin, host) {
			t.Error("Localhost requests should be allowed")
		}
	})
}
