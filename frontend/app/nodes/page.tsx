import type { Metadata } from "next"
import Link from "next/link"
import { Hash } from "lucide-react"
import { PageLayout } from "@/components/page-layout"
import { PageHeader } from "@/components/page-header"
import { Card, CardContent } from "@/components/ui/card"
import type { TopicNode } from "@/lib/types"
import { fetchAPINullable } from "@/lib/api"

export const metadata: Metadata = {
  title: "话题节点 - Go语言中文网",
  description: "Go语言中文网话题节点，浏览各类技术话题分类",
}

function getNodeColor(id: number): string {
  const colors = [
    "bg-primary/10 text-primary border-primary/20",
    "bg-chart-2/10 text-chart-2 border-chart-2/20",
    "bg-chart-3/10 text-chart-3 border-chart-3/20",
    "bg-chart-4/10 text-chart-4 border-chart-4/20",
    "bg-chart-5/10 text-chart-5 border-chart-5/20",
    "bg-accent/10 text-accent border-accent/20",
  ]
  return colors[id % colors.length]
}

export default async function NodesPage() {
  const nodes = await fetchAPINullable<TopicNode[]>("/nodes", { cache: "no-store" })

  const nodeList = nodes ?? []

  // 按 parent_id 分组，parent_id=0 为顶级
  const topLevel = nodeList.filter((n) => n.parent_id === 0)
  const childMap = new Map<number, TopicNode[]>()
  for (const node of nodeList) {
    if (node.parent_id !== 0) {
      const children = childMap.get(node.parent_id) ?? []
      childMap.set(node.parent_id, [...children, node])
    }
  }

  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="话题节点"
        description="选择你感兴趣的节点，探索相关话题"
        breadcrumbs={[{ label: "话题节点" }]}
      />

      {nodeList.length === 0 ? (
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-16 text-center">
          <Hash className="mb-3 h-10 w-10 text-muted-foreground/40" />
          <p className="text-sm text-muted-foreground">暂无节点数据</p>
        </div>
      ) : topLevel.length > 0 ? (
        <div className="space-y-8">
          {topLevel.map((parent) => {
            const children = childMap.get(parent.id) ?? []
            const displayNodes = children.length > 0 ? children : [parent]
            const isGrouped = children.length > 0

            return (
              <section key={parent.id}>
                {isGrouped && (
                  <h2 className="mb-3 flex items-center gap-2 text-base font-semibold text-foreground">
                    <Hash className="h-4 w-4 text-primary" />
                    {parent.name}
                    {parent.intro && (
                      <span className="text-xs font-normal text-muted-foreground">
                        {parent.intro}
                      </span>
                    )}
                  </h2>
                )}
                <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
                  {displayNodes.map((node) => (
                    <Link
                      key={node.id}
                      href={`/topics/node/${node.id}`}
                      className="group"
                    >
                      <Card className="h-full transition-all hover:border-primary/30 hover:shadow-sm">
                        <CardContent className="flex flex-col items-center gap-2 p-4 text-center">
                          {node.logo ? (
                            <img
                              src={node.logo}
                              alt={node.name}
                              className="h-8 w-8 rounded object-contain"
                            />
                          ) : (
                            <div
                              className={`flex h-8 w-8 items-center justify-center rounded-full border text-xs font-bold ${getNodeColor(node.id)}`}
                            >
                              {node.name.charAt(0)}
                            </div>
                          )}
                          <span className="text-sm font-medium text-foreground group-hover:text-primary">
                            {node.name}
                          </span>
                          {node.intro && (
                            <span className="line-clamp-2 text-xs text-muted-foreground">
                              {node.intro}
                            </span>
                          )}
                        </CardContent>
                      </Card>
                    </Link>
                  ))}
                </div>
              </section>
            )
          })}
        </div>
      ) : (
        // 无父级层级时，平铺所有节点
        <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
          {nodeList.map((node) => (
            <Link key={node.id} href={`/topics/node/${node.id}`} className="group">
              <Card className="h-full transition-all hover:border-primary/30 hover:shadow-sm">
                <CardContent className="flex flex-col items-center gap-2 p-4 text-center">
                  {node.logo ? (
                    <img
                      src={node.logo}
                      alt={node.name}
                      className="h-8 w-8 rounded object-contain"
                    />
                  ) : (
                    <div
                      className={`flex h-8 w-8 items-center justify-center rounded-full border text-xs font-bold ${getNodeColor(node.id)}`}
                    >
                      {node.name.charAt(0)}
                    </div>
                  )}
                  <span className="text-sm font-medium text-foreground group-hover:text-primary">
                    {node.name}
                  </span>
                  {node.intro && (
                    <span className="line-clamp-2 text-xs text-muted-foreground">
                      {node.intro}
                    </span>
                  )}
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      )}
    </PageLayout>
  )
}
