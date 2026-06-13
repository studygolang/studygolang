"use client"

import { useCallback, useRef, useState } from "react"
import { toast } from "sonner"
import { Upload, X, ImagePlus, Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Progress } from "@/components/ui/progress"
import { cn } from "@/lib/utils"
import {
  uploadImage,
  isValidImageType,
  isValidFileSize,
  formatFileSize,
  type ProgressCallback,
} from "@/lib/upload"

export interface ImageUploaderProps {
  /** 上传成功回调 */
  onUploaded: (url: string) => void
  /** 接受的文件类型，默认 "image/*" */
  accept?: string
  /** 最大文件大小（字节），默认 5MB（与后端 MaxImageSize 一致） */
  maxSize?: number
  /** 自定义类名 */
  className?: string
  /** 上传模式：upload 通用上传 / avatar 头像上传 */
  mode?: "upload" | "avatar"
  /** 是否禁用 */
  disabled?: boolean
}

type UploadStatus = "idle" | "uploading" | "success" | "error"

export function ImageUploader({
  onUploaded,
  accept = "image/*",
  maxSize = 5 * 1024 * 1024,
  className,
  mode = "upload",
  disabled = false,
}: ImageUploaderProps) {
  const inputRef = useRef<HTMLInputElement>(null)
  const [status, setStatus] = useState<UploadStatus>("idle")
  const [progress, setProgress] = useState(0)
  const [dragOver, setDragOver] = useState(false)
  const [preview, setPreview] = useState<string | null>(null)

  // 校验并上传文件
  const handleFile = useCallback(
    async (file: File) => {
      // 类型校验
      if (!isValidImageType(file)) {
        toast.error("不支持的图片格式，请上传 JPG/PNG/GIF/WebP 等图片文件")
        return
      }

      // 大小校验
      if (!isValidFileSize(file, maxSize)) {
        toast.error(`文件过大，最大允许 ${formatFileSize(maxSize)}`)
        return
      }

      // 生成预览
      const previewUrl = URL.createObjectURL(file)
      setPreview(previewUrl)
      setStatus("uploading")
      setProgress(0)

      const onProgress: ProgressCallback = (percent) => {
        setProgress(percent)
      }

      try {
        const result = await uploadImage(file, {
          avatar: mode === "avatar",
          onProgress,
        })
        setStatus("success")
        setProgress(100)
        toast.success("图片上传成功")
        onUploaded(result.url)
      } catch (err) {
        setStatus("error")
        setProgress(0)
        toast.error(err instanceof Error ? err.message : "上传失败")
      }
    },
    [maxSize, mode, onUploaded]
  )

  // 文件选择
  const handleInputChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const file = e.target.files?.[0]
      if (file) handleFile(file)
      // 清空 input 以便再次选择同一文件
      e.target.value = ""
    },
    [handleFile]
  )

  // 拖拽事件
  const handleDragOver = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault()
      e.stopPropagation()
      if (!disabled) setDragOver(true)
    },
    [disabled]
  )

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setDragOver(false)
  }, [])

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault()
      e.stopPropagation()
      setDragOver(false)
      if (disabled) return

      const file = e.dataTransfer.files?.[0]
      if (file) handleFile(file)
    },
    [disabled, handleFile]
  )

  // 清除预览
  const handleClear = useCallback(() => {
    if (preview) URL.revokeObjectURL(preview)
    setPreview(null)
    setStatus("idle")
    setProgress(0)
  }, [preview])

  // 头像模式：圆形裁剪区域
  const isAvatar = mode === "avatar"

  return (
    <div className={cn("relative", className)}>
      <input
        ref={inputRef}
        type="file"
        accept={accept}
        onChange={handleInputChange}
        className="hidden"
        disabled={disabled || status === "uploading"}
      />

      {/* 上传区域 */}
      <div
        role="button"
        tabIndex={0}
        onClick={() => inputRef.current?.click()}
        onKeyDown={(e) => {
          if (e.key === "Enter" || e.key === " ") {
            e.preventDefault()
            inputRef.current?.click()
          }
        }}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        className={cn(
          "relative flex cursor-pointer flex-col items-center justify-center border-2 border-dashed transition-colors",
          isAvatar
            ? "h-32 w-32 rounded-full"
            : "min-h-[120px] w-full rounded-lg",
          dragOver
            ? "border-primary bg-primary/5"
            : "border-muted-foreground/25 hover:border-primary/50 hover:bg-muted/50",
          (disabled || status === "uploading") && "pointer-events-none opacity-60"
        )}
      >
        {/* 预览图 */}
        {preview && (
          <div
            className={cn(
              "absolute inset-0 overflow-hidden",
              isAvatar ? "rounded-full" : "rounded-lg"
            )}
          >
            <img
              src={preview}
              alt="预览"
              className="h-full w-full object-cover"
            />
            {/* 上传中遮罩 */}
            {status === "uploading" && (
              <div className="absolute inset-0 flex flex-col items-center justify-center bg-black/40">
                <Loader2 className="mb-1 h-6 w-6 animate-spin text-white" />
                <span className="text-xs text-white">{progress}%</span>
              </div>
            )}
          </div>
        )}

        {/* 默认占位内容 */}
        {!preview && (
          <div className="flex flex-col items-center gap-1.5 p-4 text-center">
            {status === "uploading" ? (
              <>
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
                <p className="text-xs text-muted-foreground">上传中...</p>
              </>
            ) : (
              <>
                <ImagePlus className="h-8 w-8 text-muted-foreground/50" />
                <p className="text-xs text-muted-foreground">
                  点击或拖拽图片上传
                </p>
                <p className="text-xs text-muted-foreground/60">
                  最大 {formatFileSize(maxSize)}
                </p>
              </>
            )}
          </div>
        )}
      </div>

      {/* 上传进度条 */}
      {status === "uploading" && !isAvatar && (
        <div className="mt-2">
          <Progress value={progress} className="h-1.5" />
        </div>
      )}

      {/* 清除按钮 */}
      {preview && status !== "uploading" && (
        <Button
          variant="destructive"
          size="icon-sm"
          className={cn(
            "absolute -right-1 -top-1 z-10 h-5 w-5 rounded-full",
            isAvatar && "-right-0.5 -top-0.5"
          )}
          onClick={(e) => {
            e.stopPropagation()
            handleClear()
          }}
        >
          <X className="h-3 w-3" />
        </Button>
      )}
    </div>
  )
}
