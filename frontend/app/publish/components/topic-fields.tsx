"use client"

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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { useAuth } from "@/lib/auth-context"
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
  permission: string
  setPermission: (permission: string) => void
}

export function TopicFields({
  nid,
  setNid,
  nodeGroups,
  permission,
  setPermission,
}: TopicFieldsProps) {
  const [nodeOpen, setNodeOpen] = useState(false)
  const { user } = useAuth()

  return (
    <>
    <div className="w-full space-y-1.5 md:w-56 md:shrink-0">
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
                  .find((node) => String(node.nid) === nid)?.name
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
                      key={node.nid}
                      value={`${node.name} ${node.ename}`}
                      onSelect={() => {
                        setNid(String(node.nid))
                        setNodeOpen(false)
                      }}
                    >
                      <Check
                        className={cn(
                          "mr-2 h-4 w-4",
                          nid === String(node.nid) ? "opacity-100" : "opacity-0"
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

    {/* 可见性：与旧版 template/topics/new.html 一致，付费可见仅 root/VIP 可选 */}
    <div className="w-full space-y-1.5 md:w-56 md:shrink-0">
      <Label htmlFor="permission" className="text-sm font-medium">
        可见性
      </Label>
      <Select value={permission} onValueChange={setPermission}>
        <SelectTrigger id="permission" className="h-10">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="0">公开</SelectItem>
          <SelectItem value="1">登录可见</SelectItem>
          {(user?.is_root || user?.is_vip) && (
            <SelectItem value="3">付费可见</SelectItem>
          )}
        </SelectContent>
      </Select>
    </div>
    </>
  )
}

export type { NodeGroup }
