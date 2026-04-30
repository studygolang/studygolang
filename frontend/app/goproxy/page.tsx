import type { Metadata } from 'next'
import { PageLayout } from '@/components/page-layout'
import { PageHeader } from '@/components/page-header'
import { ExternalLink } from 'lucide-react'

export const metadata: Metadata = {
  title: 'Go模块代理 - Go语言中文网',
  description: '目前可用的 Go Proxy 推荐，极大的解决 Go 包下载的问题',
}

const proxies = [
  {
    name: 'goproxy.cn',
    url: 'https://goproxy.cn',
    description: '由七牛云赞助维护，中国大陆最稳定的 Go 模块代理',
    command: 'go env -w GOPROXY=https://goproxy.cn,direct',
  },
  {
    name: 'goproxy.io',
    url: 'https://goproxy.io',
    description: '开源的 Go 模块代理，由社区维护',
    command: 'go env -w GOPROXY=https://goproxy.io,direct',
  },
]

export default function GoProxyPage() {
  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="让下载Go语言第三方库畅通无阻"
        description="建议升级到 Go1.13 以上"
        breadcrumbs={[{ label: 'Go模块代理' }]}
      />

      <div className="space-y-6">
        <p className="text-muted-foreground">
          Go 模块代理（Go Module Proxy）可以加速 Go 第三方包的下载速度，
          解决国内网络环境下 <code className="bg-muted px-1 rounded">go get</code> 超时或失败的问题。
        </p>

        <div className="space-y-4">
          {proxies.map((proxy) => (
            <div
              key={proxy.name}
              className="rounded-lg border border-border p-6"
            >
              <div className="flex items-center gap-3 mb-3">
                <h3 className="text-xl font-semibold">{proxy.name}</h3>
                <a
                  href={proxy.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-primary hover:underline flex items-center gap-1"
                >
                  {proxy.url}
                  <ExternalLink className="h-3.5 w-3.5" />
                </a>
              </div>
              <p className="text-muted-foreground mb-4">{proxy.description}</p>
              <div className="bg-muted rounded-md p-3">
                <p className="text-sm font-medium mb-1">使用方式：</p>
                <code className="text-sm">{proxy.command}</code>
              </div>
            </div>
          ))}
        </div>

        <div className="rounded-lg bg-muted/50 p-6">
          <h3 className="text-lg font-semibold mb-3">其他配置</h3>
          <div className="space-y-3">
            <div>
              <p className="text-sm font-medium mb-1">查看当前 GOPROXY 配置：</p>
              <code className="text-sm bg-muted px-2 py-1 rounded">go env GOPROXY</code>
            </div>
            <div>
              <p className="text-sm font-medium mb-1">设置 GOSUMDB（校验和数据库）：</p>
              <code className="text-sm bg-muted px-2 py-1 rounded">go env -w GOSUMDB=sum.golang.org</code>
            </div>
            <div>
              <p className="text-sm font-medium mb-1">设置 GONOSUMDB（跳过校验的模块路径）：</p>
              <code className="text-sm bg-muted px-2 py-1 rounded">go env -w GONOSUMDB=github.com/your-company/*</code>
            </div>
          </div>
        </div>
      </div>
    </PageLayout>
  )
}
