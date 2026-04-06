import type { Metadata } from "next"
import Link from "next/link"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Download } from "lucide-react"
import { fetchAPI } from "@/lib/api"

export const metadata: Metadata = {
  title: "Go 安装包下载 - Go语言中文网",
  description: "下载 Go 语言最新安装包，支持 Linux、macOS、Windows 等多个平台",
}

interface GoDownload {
  id: number
  version: string
  filename: string
  kind: string
  os: string
  arch: string
  size: number
  checksum: string
  category: number
  is_recommend: boolean
  seq: number
  created_at: string
}

interface DownloadsData {
  featured: GoDownload[]
  stables: Record<string, GoDownload[]>
  stable_versions: string[]
  unstables: Record<string, GoDownload[]>
  archiveds: Record<string, GoDownload[]>
  archived_versions: string[]
}

async function getDownloads(): Promise<DownloadsData | null> {
  try {
    return await fetchAPI<DownloadsData>("/downloads", { cache: "no-store" })
  } catch {
    return null
  }
}

function formatSize(bytes: number): string {
  if (bytes >= 1073741824) return (bytes / 1073741824).toFixed(1) + " GB"
  if (bytes >= 1048576) return (bytes / 1048576).toFixed(1) + " MB"
  if (bytes >= 1024) return (bytes / 1024).toFixed(1) + " KB"
  return bytes + " B"
}

function DownloadTable({ downloads }: { downloads: GoDownload[] }) {
  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-border">
            <th className="px-3 py-2 text-left font-medium">文件名</th>
            <th className="px-3 py-2 text-left font-medium">系统</th>
            <th className="px-3 py-2 text-left font-medium">架构</th>
            <th className="px-3 py-2 text-left font-medium">大小</th>
            <th className="px-3 py-2 text-left font-medium">SHA256</th>
          </tr>
        </thead>
        <tbody>
          {downloads.map((dl) => (
            <tr key={dl.id} className="border-b border-border/50 hover:bg-muted/30">
              <td className="px-3 py-2">
                <Link
                  href={`/api/v1/downloads/golang/${dl.filename}`}
                  className="text-primary hover:underline"
                >
                  {dl.filename}
                </Link>
              </td>
              <td className="px-3 py-2 text-muted-foreground">{dl.os || dl.kind}</td>
              <td className="px-3 py-2 text-muted-foreground">{dl.arch}</td>
              <td className="px-3 py-2 text-muted-foreground">
                {dl.size > 0 ? formatSize(dl.size) : "-"}
              </td>
              <td className="max-w-[200px] truncate px-3 py-2 font-mono text-xs text-muted-foreground">
                {dl.checksum || "-"}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default async function DownloadPage() {
  const data = await getDownloads()

  if (!data) {
    return (
      <PageLayout sidebar={false}>
        <PageHeader title="Go 下载" breadcrumbs={[{ label: "下载" }]} />
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-20 text-center">
          <Download className="mb-3 h-10 w-10 text-muted-foreground/40" />
          <p className="text-sm text-muted-foreground">暂时无法获取下载列表</p>
        </div>
      </PageLayout>
    )
  }

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="Go 安装包下载"
        description="下载 Go 语言最新版本安装包"
        breadcrumbs={[{ label: "Go 下载" }]}
      />

      {/* 推荐下载 */}
      {data.featured && data.featured.length > 0 && (
        <Card className="mb-6 border-primary/30 bg-primary/5">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <Download className="h-4 w-4" />
              推荐下载
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex flex-wrap gap-3">
              {data.featured.map((dl) => (
                <Link
                  key={dl.id}
                  href={`/api/v1/downloads/golang/${dl.filename}`}
                  className="inline-flex items-center gap-1.5 rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground hover:bg-primary/90"
                >
                  <Download className="h-3.5 w-3.5" />
                  {dl.filename}
                </Link>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* 稳定版本 */}
      {data.stable_versions?.map((version) => (
        <Card key={version} className="mb-4">
          <CardHeader className="pb-2">
            <CardTitle className="text-base">
              稳定版 {version}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <DownloadTable downloads={data.stables[version] || []} />
          </CardContent>
        </Card>
      ))}

      {/* 不稳定版本 */}
      {Object.keys(data.unstables || {}).length > 0 && (
        <Card className="mb-4">
          <CardHeader>
            <CardTitle className="text-base">不稳定版本</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {Object.entries(data.unstables).map(([version, downloads]) => (
              <div key={version}>
                <h3 className="mb-2 font-medium">{version}</h3>
                <DownloadTable downloads={downloads} />
              </div>
            ))}
          </CardContent>
        </Card>
      )}

      {/* 归档版本 */}
      {data.archived_versions && data.archived_versions.length > 0 && (
        <details className="mt-4">
          <summary className="cursor-pointer text-sm font-medium text-muted-foreground hover:text-primary">
            查看归档版本（{data.archived_versions.length} 个版本）
          </summary>
          <div className="mt-3 space-y-4">
            {data.archived_versions.slice(0, 10).map((version) => (
              <Card key={version}>
                <CardHeader className="pb-2">
                  <CardTitle className="text-sm">{version}</CardTitle>
                </CardHeader>
                <CardContent>
                  <DownloadTable downloads={data.archiveds[version] || []} />
                </CardContent>
              </Card>
            ))}
          </div>
        </details>
      )}
    </PageLayout>
  )
}
