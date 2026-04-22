"use client"

import { useEffect, useState, useCallback } from "react"
import { RefreshCw } from "lucide-react"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import { fetchCaptcha } from "@/lib/captcha"

interface CaptchaFieldProps {
  value: string
  onChange: (value: string) => void
  captchaId: string
  onCaptchaIdChange: (id: string) => void
}

export function CaptchaField({
  value,
  onChange,
  captchaId,
  onCaptchaIdChange,
}: CaptchaFieldProps) {
  const [imageUrl, setImageUrl] = useState("")
  const [loading, setLoading] = useState(false)

  const refreshCaptcha = useCallback(async () => {
    setLoading(true)
    try {
      const data = await fetchCaptcha()
      onCaptchaIdChange(data.captchaId)
      setImageUrl(data.imageUrl)
      onChange("")
    } catch {
      // 获取失败时保留当前状态，用户可手动重试
    } finally {
      setLoading(false)
    }
  }, [onCaptchaIdChange, onChange])

  // 挂载时自动获取验证码
  useEffect(() => {
    refreshCaptcha()
    // 仅在挂载时获取一次
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return (
    <div className="space-y-2">
      <Label htmlFor="captcha">
        验证码 <span className="text-destructive">*</span>
      </Label>
      <div className="flex items-center gap-3">
        <Input
          id="captcha"
          name="captcha"
          type="text"
          placeholder="请输入验证码"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          required
          className="w-32"
          autoComplete="off"
        />
        {/* 点击图片刷新验证码 */}
        {imageUrl && (
          <button
            type="button"
            onClick={refreshCaptcha}
            className="flex-shrink-0 cursor-pointer rounded border border-border overflow-hidden"
            title="点击刷新验证码"
          >
            <img
              src={imageUrl}
              alt="验证码"
              width={100}
              height={40}
              className="block"
            />
          </button>
        )}
        <Button
          type="button"
          variant="ghost"
          size="icon"
          onClick={refreshCaptcha}
          disabled={loading}
          title="刷新验证码"
          className="flex-shrink-0"
        >
          <RefreshCw className={`h-4 w-4 ${loading ? "animate-spin" : ""}`} />
        </Button>
      </div>
      <p className="text-xs text-muted-foreground">
        看不清？点击图片或刷新按钮更换
      </p>
    </div>
  )
}
