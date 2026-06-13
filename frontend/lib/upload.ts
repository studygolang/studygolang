// 图片上传工具函数
// 后端 API 端点：
//   POST /api/v1/image/upload       — 普通上传（field: img）
//   POST /api/v1/image/paste_upload  — 粘贴上传（field: imageFile）
//   POST /api/v1/image/quick_upload  — 快捷上传（field: upload）
//   POST /api/v1/image/transfer      — URL 转存（field: url）

import type { APIResponse } from './types'

// 客户端调用，直接访问后端（与 api.ts 保持一致）
function getAPIBase(): string {
  return process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8090'
}

/** 上传结果 */
export interface UploadResult {
  url: string
  filename?: string
}

/** 上传进度回调 */
export type ProgressCallback = (percent: number) => void

/**
 * 通用文件上传（后端 field 名为 "img"）
 * 可选 avatar 参数传 "1" 上传到头像目录
 */
export async function uploadImage(
  file: File,
  options?: { avatar?: boolean; onProgress?: ProgressCallback }
): Promise<UploadResult> {
  const base = getAPIBase()
  const formData = new FormData()
  formData.append('img', file)
  if (options?.avatar) {
    formData.append('avatar', '1')
  }

  return postFormData(`${base}/api/v1/image/upload`, formData, options?.onProgress)
}

/**
 * 粘贴上传图片（后端 field 名为 "imageFile"）
 */
export async function pasteUploadImage(
  file: File,
  onProgress?: ProgressCallback
): Promise<UploadResult> {
  const base = getAPIBase()
  const formData = new FormData()
  formData.append('imageFile', file)

  return postFormData(`${base}/api/v1/image/paste_upload`, formData, onProgress)
}

/**
 * 编辑器快捷上传（后端 field 名为 "upload"）
 */
export async function quickUploadImage(
  file: File,
  onProgress?: ProgressCallback
): Promise<UploadResult> {
  const base = getAPIBase()
  const formData = new FormData()
  formData.append('upload', file)

  return postFormData(`${base}/api/v1/image/quick_upload`, formData, onProgress)
}

/**
 * 通过 URL 转存图片
 */
export async function transferImage(url: string): Promise<UploadResult> {
  const base = getAPIBase()
  const formData = new FormData()
  formData.append('url', url)

  return postFormData(`${base}/api/v1/image/transfer`, formData)
}

/**
 * 通用 FormData POST 请求
 * 使用 XMLHttpRequest 支持上传进度
 */
function postFormData(
  url: string,
  formData: FormData,
  onProgress?: ProgressCallback
): Promise<UploadResult> {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()

    xhr.withCredentials = true

    // 上传进度
    if (onProgress) {
      xhr.upload.addEventListener('progress', (e) => {
        if (e.lengthComputable) {
          const percent = Math.round((e.loaded / e.total) * 100)
          onProgress(percent)
        }
      })
    }

    xhr.addEventListener('load', () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        try {
          const json: APIResponse<{ url: string; uri?: string; message?: string }> = JSON.parse(xhr.responseText)
          if (json.code !== 0) {
            reject(new Error(json.msg || '上传失败'))
            return
          }
          // 粘贴上传返回 { success: 1, message: url }
          const data = json.data as Record<string, unknown>
          const imageUrl = (data?.url as string) || (data?.message as string) || ''
          if (!imageUrl) {
            reject(new Error('服务器未返回图片地址'))
            return
          }
          resolve({ url: imageUrl })
        } catch {
          reject(new Error('解析响应失败'))
        }
      } else {
        reject(new Error(`上传失败 (HTTP ${xhr.status})`))
      }
    })

    xhr.addEventListener('error', () => {
      reject(new Error('网络错误，请检查网络连接'))
    })

    xhr.addEventListener('abort', () => {
      reject(new Error('上传已取消'))
    })

    xhr.open('POST', url)
    xhr.send(formData)
  })
}

/** 允许的图片 MIME 类型（与后端 allowedMIMETypes 保持一致） */
const ALLOWED_IMAGE_TYPES = [
  'image/jpeg',
  'image/png',
  'image/gif',
  'image/webp',
]

/** 默认最大文件大小 5MB（与后端 MaxImageSize 保持一致） */
const DEFAULT_MAX_SIZE = 5 * 1024 * 1024

/**
 * 校验文件是否为允许的图片类型
 */
export function isValidImageType(file: File): boolean {
  return ALLOWED_IMAGE_TYPES.includes(file.type)
}

/**
 * 校验文件大小
 */
export function isValidFileSize(file: File, maxSize: number = DEFAULT_MAX_SIZE): boolean {
  return file.size <= maxSize
}

/**
 * 格式化文件大小为人类可读格式
 */
export function formatFileSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}
