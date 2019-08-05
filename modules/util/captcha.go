package util

import "github.com/dchest/captcha"

const CaptchaLen = 4

var captchaStorage captcha.Store

func init() {
	// 存储在内存中，不考虑程序退出会清理掉的问题
	captchaStorage = captcha.NewMemoryStore(captcha.CollectNum, captcha.Expiration)

	captcha.SetCustomStore(captchaStorage)
}

// SetCaptcha 设置验证码
func SetCaptcha(id string) {
	digits := captcha.RandomDigits(CaptchaLen)
	captchaStorage.Set(id, digits)
}
