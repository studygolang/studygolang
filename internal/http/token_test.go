// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build integration
// +build integration

package http

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/polaris1119/goutils"
)

// 注意：这些单元测试不依赖外部配置，可以直接运行
// 运行方式：go test -v -run TestToken ./internal/http/

func TestGenToken(t *testing.T) {
	// 初始化 TokenSalt（测试用）
	TokenSalt = "test_salt_for_unit_testing_64_characters_long_enough_for_security"

	tests := []struct {
		name string
		uid  int
	}{
		{"regular user", 123},
		{"admin user", 1},
		{"large uid", 999999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := GenToken(tt.uid)
			if token == "" {
				t.Error("GenToken() returned empty token")
			}

			// 验证 token 格式：expireTime + md5 + "uid" + uid
			// 应该包含 "uid" 字符串
			if len(token) < 40 {
				t.Errorf("GenToken() token too short: got %d chars", len(token))
			}

			// 验证包含 "uid" 标记
			if !strings.Contains(token, "uid") {
				t.Error("GenToken() token should contain 'uid' marker")
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	TokenSalt = "test_salt_for_unit_testing_64_characters_long_enough_for_security"

	tests := []struct {
		name    string
		token   string
		wantUid int
		wantOk  bool
	}{
		{
			name:    "valid token",
			token:   GenToken(123),
			wantUid: 123,
			wantOk:  true,
		},
		{
			name:    "empty token",
			token:   "",
			wantUid: 0,
			wantOk:  false,
		},
		{
			name:    "too short token",
			token:   "short",
			wantUid: 0,
			wantOk:  false,
		},
		{
			name:    "no uid marker",
			token:   "1234567890abcdef1234567890abcdef",
			wantUid: 0,
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uid, ok := ParseToken(tt.token)
			if ok != tt.wantOk {
				t.Errorf("ParseToken() ok = %v, want %v", ok, tt.wantOk)
			}
			if uid != tt.wantUid {
				t.Errorf("ParseToken() uid = %v, want %v", uid, tt.wantUid)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	TokenSalt = "test_salt_for_unit_testing_64_characters_long_enough_for_security"

	tests := []struct {
		name       string
		setupToken func() string
		wantValid  bool
	}{
		{
			name: "valid token",
			setupToken: func() string {
				return GenToken(123)
			},
			wantValid: true,
		},
		{
			name: "empty token",
			setupToken: func() string {
				return ""
			},
			wantValid: false,
		},
		{
			name: "expired token",
			setupToken: func() string {
				// 创建一个已经过期的 token
				uid := 123
				expireTime := time.Now().Add(-24 * time.Hour).Unix() // 24小时前过期
				buffer := goutils.NewBuffer().Append(expireTime).Append(uid).Append(TokenSalt)
				md5Hash := goutils.Md5(buffer.String())
				token := fmt.Sprintf("%d%suid%d", expireTime, md5Hash, uid)
				return token
			},
			wantValid: false,
		},
		{
			name: "wrong signature",
			setupToken: func() string {
				// 创建一个签名错误的 token
				uid := 123
				expireTime := time.Now().Add(24 * time.Hour).Unix()
				wrongBuffer := goutils.NewBuffer().Append(expireTime).Append(uid).Append("wrong_salt")
				md5Hash := goutils.Md5(wrongBuffer.String())
				token := fmt.Sprintf("%d%suid%d", expireTime, md5Hash, uid)
				return token
			},
			wantValid: false,
		},
		{
			name: "too short token",
			setupToken: func() string {
				return "short_token"
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupToken()
			valid := ValidateToken(token)
			if valid != tt.wantValid {
				t.Errorf("ValidateToken() = %v, want %v", valid, tt.wantValid)
			}
		})
	}
}

func TestValidateTokenRoundTrip(t *testing.T) {
	TokenSalt = "test_salt_for_unit_testing_64_characters_long_enough_for_security"

	// 测试完整的生成-解析-验证流程
	uid := 456
	token := GenToken(uid)

	// 验证 token 有效性
	if !ValidateToken(token) {
		t.Error("GenToken() produced invalid token")
	}

	// 解析 token
	parsedUid, ok := ParseToken(token)
	if !ok {
		t.Error("ParseToken() failed on valid token")
	}
	if parsedUid != uid {
		t.Errorf("ParseToken() uid = %d, want %d", parsedUid, uid)
	}
}

func TestJWTTokenGeneration(t *testing.T) {
	TokenSalt = "test_salt_for_unit_testing_64_characters_long_enough_for_security"
	JWTSecret = []byte("test_salt_for_unit_testing_64_characters_long_enough_for_security")

	uid := 123
	username := "testuser"

	token, err := GenJWTToken(uid, username)
	if err != nil {
		t.Fatalf("GenJWTToken() failed: %v", err)
	}
	if token == "" {
		t.Error("GenJWTToken() returned empty token")
	}

	// JWT token starts with "eyJ" (base64-encoded header)
	if !strings.HasPrefix(token, "eyJ") {
		t.Errorf("GenJWTToken() token should start with 'eyJ', got prefix: %s", token[:min(10, len(token))])
	}

	claims, err := ValidateJWTToken(token)
	if err != nil {
		t.Fatalf("ValidateJWTToken() failed: %v", err)
	}
	if claims.UID != uid {
		t.Errorf("ValidateJWTToken() uid = %d, want %d", claims.UID, uid)
	}
	if claims.Username != username {
		t.Errorf("ValidateJWTToken() username = %s, want %s", claims.Username, username)
	}
}

func TestValidateTokenAuto(t *testing.T) {
	TokenSalt = "test_salt_for_unit_testing_64_characters_long_enough_for_security"
	JWTSecret = []byte("test_salt_for_unit_testing_64_characters_long_enough_for_security")

	t.Run("legacy MD5 token", func(t *testing.T) {
		oldToken := GenToken(456)
		uid, username, valid := ValidateTokenAuto(oldToken)
		if !valid {
			t.Error("ValidateTokenAuto() failed to validate old MD5 token")
		}
		if uid != 456 {
			t.Errorf("ValidateTokenAuto() uid = %d, want 456", uid)
		}
		// username is empty for legacy token (no username embedded)
		if username != "" {
			t.Errorf("ValidateTokenAuto() username should be empty for old token, got %s", username)
		}
	})

	t.Run("new JWT token", func(t *testing.T) {
		newToken, err := GenJWTToken(789, "newuser")
		if err != nil {
			t.Fatalf("GenJWTToken() failed: %v", err)
		}
		uid, username, valid := ValidateTokenAuto(newToken)
		if !valid {
			t.Error("ValidateTokenAuto() failed to validate new JWT token")
		}
		if uid != 789 {
			t.Errorf("ValidateTokenAuto() uid = %d, want 789", uid)
		}
		if username != "newuser" {
			t.Errorf("ValidateTokenAuto() username = %s, want newuser", username)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		uid, username, valid := ValidateTokenAuto("invalid-token-string")
		if valid {
			t.Error("ValidateTokenAuto() should return false for invalid token")
		}
		if uid != 0 {
			t.Errorf("ValidateTokenAuto() uid = %d, want 0", uid)
		}
		if username != "" {
			t.Errorf("ValidateTokenAuto() username should be empty, got %s", username)
		}
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
