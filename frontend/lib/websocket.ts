/**
 * WebSocket 连接管理器
 *
 * 安全要点:
 * - 认证: 通过 HttpOnly Cookie 携带 session token，握手时由后端 requireAuth 验证
 * - 协议: 生产环境强制 wss://（通过 Nginx SSL 终止 + proxy_pass 到后端）
 * - 重连: 指数退避策略（1s → 2s → 4s → ... → 60s cap），最多 10 次
 */
type MessageHandler = (data: unknown) => void

interface WSOptions {
  onMessage: MessageHandler
  onOpen?: () => void
  onClose?: () => void
  onError?: (err: Event) => void
}

interface WSController {
  disconnect: () => void
}

export function createWebSocket(opts: WSOptions): WSController {
  // WebSocket 端点在 /api/v1/ws（通过 Next.js rewrite 或 Nginx 代理到 Go 后端）
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:"
  const wsUrl = `${protocol}//${window.location.host}/api/v1/ws`

  let ws: WebSocket | null = null
  let retryCount = 0
  const maxRetries = 10
  let retryTimer: ReturnType<typeof setTimeout>

  function connect() {
    ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      retryCount = 0
      opts.onOpen?.()
    }

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        opts.onMessage(data)
      } catch {
        // 非 JSON 消息忽略
      }
    }

    ws.onclose = () => {
      opts.onClose?.()
      reconnect()
    }

    ws.onerror = (err) => {
      opts.onError?.(err)
    }
  }

  function reconnect() {
    if (retryCount >= maxRetries) return
    // 指数退避: 1s, 2s, 4s, 8s, 16s... 最大 60s
    const delay = Math.min(1000 * Math.pow(2, retryCount), 60000)
    retryCount++
    retryTimer = setTimeout(connect, delay)
  }

  function disconnect() {
    clearTimeout(retryTimer)
    ws?.close()
    ws = null
  }

  connect()

  return { disconnect }
}
