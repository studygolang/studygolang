// Copyright 2026 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/studygolang/studygolang/internal/model"

	echo "github.com/labstack/echo/v4"
)

// newFormContext 构造携带表单的 echo 测试上下文
func newFormContext(method, form string) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(method, "/", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return e.NewContext(req, httptest.NewRecorder())
}

func TestHasSensitive(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		sensitive string
		want      bool
	}{
		{"命中单个敏感词", "低价代开发票", "发票", true},
		{"命中逗号分隔列表", "加微信xxx", "发票,加微信", true},
		{"正常内容", "Go 是一门简洁的语言", "发票", false},
		{"空内容", "", "发票", false},
		{"空敏感词配置", "任意内容", "", false},
		{"空词项不误报（多余逗号）", "正常内容", ",,", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasSensitive(tt.content, tt.sensitive); got != tt.want {
				t.Errorf("hasSensitive(%q, %q) = %v, want %v", tt.content, tt.sensitive, got, tt.want)
			}
		})
	}
}

func TestHasSensitiveChar(t *testing.T) {
	// 敏感词按字符拆分，标题需同时包含所有字符才算命中（防乱序变体）
	tests := []struct {
		name      string
		title     string
		sensitive string
		want      bool
	}{
		{"所有字符都包含", "【北京发票】开票", "发票", true},
		{"字符顺序打乱也命中", "票据发put", "发票", true},
		{"缺少一个字符", "发发发", "发票", false},
		{"空标题", "", "发票", false},
		{"空敏感词", "标题", "", false},
		{"单字符敏感词", "票", "票", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasSensitiveChar(tt.title, tt.sensitive); got != tt.want {
				t.Errorf("hasSensitiveChar(%q, %q) = %v, want %v", tt.title, tt.sensitive, got, tt.want)
			}
		})
	}
}

func TestBalanceCheck(t *testing.T) {
	tests := []struct {
		name      string
		balance   int
		isComment bool
		want      bool
	}{
		{"发布余额不足", 10, false, false},
		{"发布余额刚好达标", 20, false, true},
		{"评论门槛较低达标", 5, true, true},
		{"评论余额不足", 4, true, false},
		{"零余额评论不过", 0, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			me := &model.Me{Balance: tt.balance}
			if got := balanceCheck(me, tt.isComment); got != tt.want {
				t.Errorf("balanceCheck(balance=%d, isComment=%v) = %v, want %v", tt.balance, tt.isComment, got, tt.want)
			}
		})
	}

	t.Run("nil 用户拒绝", func(t *testing.T) {
		if balanceCheck(nil, false) {
			t.Error("balanceCheck(nil) = true, want false")
		}
	})
}

func TestCaptchaCheckMethodGuard(t *testing.T) {
	// 非 POST 请求直接跳过检查（不触发 NeedCaptcha 的 Redis 依赖）
	ctx := newFormContext("PUT", "title=x")
	me := &model.Me{}
	if !captchaCheck(ctx, me) {
		t.Error("captchaCheck PUT 应跳过检查，got false")
	}
}
