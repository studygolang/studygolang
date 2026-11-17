package logic

import "testing"

func TestGetWechatConfigEnvironmentNames(t *testing.T) {
	tests := []struct {
		section string
		key     string
		envKey  string
	}{
		{section: "wechat", key: "appid", envKey: "WECHAT_APPID"},
		{section: "wechat", key: "appsecret", envKey: "WECHAT_APPSECRET"},
		{section: "wechat", key: "subscribe", envKey: "WECHAT_SUBSCRIBE"},
		{section: "wechat.xcx", key: "appid", envKey: "WECHAT_XCX_APPID"},
		{section: "wechat.xcx", key: "appsecret", envKey: "WECHAT_XCX_APPSECRET"},
	}

	for _, tt := range tests {
		t.Run(tt.envKey, func(t *testing.T) {
			t.Setenv(tt.envKey, "from-environment")
			if got := getWechatConfig(tt.section, tt.key); got != "from-environment" {
				t.Errorf("getWechatConfig(%q, %q) = %q, want environment value", tt.section, tt.key, got)
			}
		})
	}
}
