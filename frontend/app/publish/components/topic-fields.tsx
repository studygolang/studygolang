import { useState } from "react"
import { Check, ChevronsUpDown } from "lucide-react"
import { Label } from "@/components/ui/label"
import { Button } from "@/components/ui/button"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import { cn } from "@/lib/utils"
import type { TopicNode } from "@/lib/types"

// 节点分组类型
type NodeGroup = {
  category: string
  nodes: TopicNode[]
}

interface TopicFieldsProps {
  nid: string
  setNid: (nid: string) => void
  nodeGroups: NodeGroup[]
}

export function TopicFields({ nid, setNid, nodeGroups }: TopicFieldsProps) {
  const [nodeOpen, setNodeOpen] = useState(false)

  return (
    <div className="w-56 space-y-1.5">
      <Label htmlFor="nid" className="text-sm font-medium">
        节点 <span className="text-destructive">*</span>
      </Label>
      <Popover open={nodeOpen} onOpenChange={setNodeOpen}>
        <PopoverTrigger asChild>
          <Button
            variant="outline"
            role="combobox"
            aria-expanded={nodeOpen}
            className="w-full justify-between h-10"
          >
            {nid
              ? nodeGroups
                  .flatMap((g) => g.nodes)
                  .find((node) => String(node.id) === nid)?.name
              : "选择节点"}
            <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-[300px] p-0">
          <Command>
            <CommandInput placeholder="搜索节点..." />
            <CommandList>
              <CommandEmpty>未找到节点</CommandEmpty>
              {nodeGroups.map((group) => (
                <CommandGroup key={group.category} heading={group.category}>
                  {group.nodes.map((node) => (
                    <CommandItem
                      key={node.id}
                      value={`${node.name} ${node.ename}`}
                      onSelect={() => {
                        setNid(String(node.id))
                        setNodeOpen(false)
                      }}
                    >
                      <Check
                        className={cn(
                          "mr-2 h-4 w-4",
                          nid === String(node.id) ? "opacity-100" : "opacity-0"
                        )}
                      />
                      {node.name}
                    </CommandItem>
                  ))}
                </CommandGroup>
              ))}
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
    </div>
  )
}

export type { NodeGroup }
