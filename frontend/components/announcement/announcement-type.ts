export type AnnouncementVariant = "default" | "secondary" | "destructive" | "outline"

export interface AnnouncementTypeInfo {
  label: string
  variant: AnnouncementVariant
}

export const ANNOUNCEMENT_TYPE_CONFIG: Record<number, AnnouncementTypeInfo> = {
  1: { label: "公告", variant: "default" },
  2: { label: "活动", variant: "secondary" },
  3: { label: "警告", variant: "destructive" },
}

export const DEFAULT_ANNOUNCEMENT_TYPE: AnnouncementTypeInfo = {
  label: "公告",
  variant: "default",
}

export function getAnnouncementType(type: number): AnnouncementTypeInfo {
  return ANNOUNCEMENT_TYPE_CONFIG[type] ?? DEFAULT_ANNOUNCEMENT_TYPE
}

// 列表页用的过滤器：id=0 表示"全部"（无 type 参数）
export const ANNOUNCEMENT_TYPE_FILTERS: ReadonlyArray<{
  id: number
  label: string
}> = [
  { id: 0, label: "全部" },
  ...Object.entries(ANNOUNCEMENT_TYPE_CONFIG).map(([id, info]) => ({
    id: Number(id),
    label: info.label,
  })),
]

export const ANNOUNCEMENT_ALL_TYPE_ID = 0

