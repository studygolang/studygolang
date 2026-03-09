import Link from "next/link"
import { Badge } from "@/components/ui/badge"

interface NodeGroup {
  label: string
  nodes: { name: string; href: string }[]
}

const nodeGroups: NodeGroup[] = [
  {
    label: "Go\u8bed\u8a00",
    nodes: [
      { name: "Go\u57fa\u7840", href: "/nodes/go-basics" },
      { name: "Go\u6807\u51c6\u5e93", href: "/nodes/go-stdlib" },
      { name: "Go\u6e90\u7801", href: "/nodes/go-source" },
      { name: "Go Web\u5f00\u53d1", href: "/nodes/go-web" },
      { name: "Go\u95ee\u4e0e\u7b54", href: "/nodes/go-qa" },
      { name: "Go\u52a8\u6001", href: "/nodes/go-news" },
      { name: "Go\u5f00\u53d1\u5de5\u5177", href: "/nodes/go-tools" },
      { name: "Go Web\u6846\u67b6", href: "/nodes/go-frameworks" },
      { name: "Go\u7b2c\u4e09\u65b9\u5e93", href: "/nodes/go-libs" },
      { name: "Go\u5b9e\u6218", href: "/nodes/go-practice" },
    ],
  },
  {
    label: "Geek",
    nodes: [
      { name: "\u7a0b\u5e8f\u5458", href: "/nodes/programmer" },
      { name: "\u7f16\u7a0b", href: "/nodes/coding" },
      { name: "Python", href: "/nodes/python" },
      { name: "Java", href: "/nodes/java" },
      { name: "Rust", href: "/nodes/rust" },
    ],
  },
  {
    label: "\u751f\u6d3b",
    nodes: [
      { name: "\u9177\u5de5\u4f5c", href: "/jobs" },
      { name: "\u62db\u8058", href: "/jobs" },
      { name: "\u57ce\u5e02", href: "/nodes/city" },
    ],
  },
]

export function NodeNavigation() {
  return (
    <div className="rounded-lg border border-border bg-card p-4">
      <h3 className="mb-3 text-sm font-semibold text-foreground">
        {"\u8282\u70b9\u5bfc\u822a"}
      </h3>
      <div className="space-y-3">
        {nodeGroups.map((group) => (
          <div key={group.label}>
            <span className="mb-1.5 block text-xs font-medium text-muted-foreground">
              {group.label}
            </span>
            <div className="flex flex-wrap gap-1.5">
              {group.nodes.map((node) => (
                <Link key={node.href} href={node.href}>
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
