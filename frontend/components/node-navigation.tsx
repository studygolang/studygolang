import Link from "next/link"
import { Badge } from "@/components/ui/badge"
import type { TopicNode } from "@/lib/types"

interface NodeGroup {
  label: string
  nodes: { id: number; name: string; href: string }[]
}

interface NodeNavigationProps {
  nodes?: TopicNode[]
}

export function NodeNavigation({ nodes }: NodeNavigationProps) {
  if (!nodes || nodes.length === 0) return null

  let nodeGroups: NodeGroup[]

  // 找出父节点，按层级分组
  const parentNodes = nodes.filter((n) => n.parent === 0)
  if (parentNodes.length > 0) {
    nodeGroups = parentNodes
      .map((parent) => ({
        label: parent.name,
        nodes: nodes
          .filter((n) => n.parent === parent.nid && n.nid != null)
          .map((n) => ({
            id: n.nid,
            name: n.name,
            // 节点话题列表路由为 /topics/node/[nid]，使用数字 id
            href: `/topics/node/${n.nid}`,
          })),
      }))
      .filter((g) => g.nodes.length > 0)
  } else {
    // 无层级关系，平铺显示
    nodeGroups = [
      {
        label: "节点",
        nodes: nodes
          .filter((n) => n.nid != null)
          .map((n) => ({
            id: n.nid,
            name: n.name,
            // 节点话题列表路由为 /topics/node/[nid]，使用数字 id
            href: `/topics/node/${n.nid}`,
          })),
      },
    ]
  }

  if (nodeGroups.length === 0) return null

  return (
    <div className="rounded-lg border border-border bg-card p-4">
      <h3 className="mb-3 text-sm font-semibold text-foreground">{"节点导航"}</h3>
      <div className="space-y-3">
        {nodeGroups.map((group) => (
          <div key={group.label}>
            <span className="mb-1.5 block text-xs font-medium text-muted-foreground">{group.label}</span>
            <div className="flex flex-wrap gap-1.5">
              {group.nodes.map((node) => (
                <Link key={String(node.id)} href={node.href}>
                  <Badge
                    variant="secondary"
                    className="cursor-pointer text-xs transition-colors hover:bg-primary/10 hover:text-primary"
                  >
                    {node.name}
                  </Badge>
                </Link>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
