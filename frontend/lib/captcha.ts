// 验证码相关 API 封装

export interface CaptchaData {
  captchaId: string
  imageUrl: string
}

/**
 * 获取新验证码
 * 调用后端 /api/v1/captcha/new 接口，返回 captcha_id
 * 验证码图片通过 /api/v1/captcha/{captchaId}.png 获取
 */
export async function fetchCaptcha(): Promise<CaptchaData> {
  const res = await fetch("/api/v1/captcha/new", { credentials: "include" })
  const json = await res.json()
  if (json.code === 0) {
    const captchaId = json.data.captcha_id as string
    return {
      captchaId,
      imageUrl: `/api/v1/captcha/${captchaId}.png`,
    }
  }
  throw new Error("获取验证码失败")
}
