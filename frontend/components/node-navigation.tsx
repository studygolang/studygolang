import Link from "next/link"
import { Badge } from "@/components/ui/badge"
import type { NodeNavigationGroup, TopicNode } from "@/lib/types"

interface NodeGroup {
  label: string
  nodes: { id: string | number; name: string; href: string }[]
}

interface NodeNavigationProps {
  nodes?: TopicNode[]
  groups?: NodeNavigationGroup[]
}

export function NodeNavigation({ nodes, groups }: NodeNavigationProps) {
  let nodeGroups: NodeGroup[] = []

  if (groups && groups.length > 0) {
    nodeGroups = groups.flatMap((group) =>
      Object.entries(group).map(([label, children]) => ({
        label,
        nodes: children.map((node) => ({
          id: node.nid ?? node.ename,
          name: node.name,
          href: `/go/${encodeURIComponent(node.ename)}?fr=index_node`,
        })),
      }))
    ).filter((group) => group.nodes.length > 0)
  } else if (nodes && nodes.length > 0) {
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
              href: n.ename
                ? `/go/${encodeURIComponent(n.ename)}?fr=topic_node`
                : `/topics/node/${n.nid}`,
            })),
        }))
        .filter((group) => group.nodes.length > 0)
    } else {
      nodeGroups = [
        {
          label: "节点",
          nodes: nodes
            .filter((n) => n.nid != null)
            .map((n) => ({
              id: n.nid,
              name: n.name,
              href: n.ename
                ? `/go/${encodeURIComponent(n.ename)}?fr=topic_node`
                : `/topics/node/${n.nid}`,
            })),
        },
      ]
    }
  }

  if (nodeGroups.length === 0) return null

  return (
    <nav aria-label="节点导航" className="rounded-lg border border-border bg-card p-4">
      <div className="mb-3 flex items-center justify-between gap-3">
        <h3 className="text-sm font-semibold text-foreground">节点导航</h3>
        <Link href="/nodes" className="text-xs text-muted-foreground transition-colors hover:text-primary">
          浏览全部节点
        </Link>
      </div>
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
    </nav>
  )
}
