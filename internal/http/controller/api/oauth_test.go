package api

import "testing"

func TestIsSafeRedirect(t *testing.T) {
	tests := []struct {
		name     string
		redirect string
		expected bool
	}{
		// 合法相对路径
		{"root", "/", true},
		{"path", "/account/login", true},
		{"deep path", "/topics/123?tab=reply", true},
		{"path with fragment", "/foo#bar", true},
		// open redirect 攻击
		{"empty rejected", "", false},
		{"absolute http rejected", "http://evil.com", false},
		{"absolute https rejected", "https://evil.com/phishing", false},
		{"protocol-relative rejected (//evil.com)", "//evil.com", false},
		{"protocol-relative with path rejected", "//evil.com/login", false},
		{"javascript scheme rejected", "javascript:alert(1)", false},
		{"no leading slash rejected", "foo/bar", false},
		{"backslash trick rejected", "/\\evil.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSafeRedirect(tt.redirect); got != tt.expected {
				t.Errorf("isSafeRedirect(%q) = %v, want %v", tt.redirect, got, tt.expected)
			}
		})
	}
}
