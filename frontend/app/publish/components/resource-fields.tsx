import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"

interface ResourceFieldsProps {
  resourceForm: string
  setResourceForm: (form: string) => void
  resourceUrl: string
  setResourceUrl: (url: string) => void
}

export function ResourceFields({
  resourceForm,
  setResourceForm,
  resourceUrl,
  setResourceUrl,
}: ResourceFieldsProps) {
  return (
    <>
      {/* 资源类型选择 */}
      <div className="space-y-1.5">
        <Label className="text-sm font-medium">资源类型 <span className="text-destructive">*</span></Label>
        <div className="flex gap-4">
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="radio"
              name="resourceForm"
              value="link"
              checked={resourceForm === "link"}
              onChange={() => setResourceForm("link")}
              className="h-4 w-4"
            />
            <span className="text-sm">只是链接</span>
          </label>
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="radio"
              name="resourceForm"
              value="content"
              checked={resourceForm === "content"}
              onChange={() => setResourceForm("content")}
              className="h-4 w-4"
            />
            <span className="text-sm">包括内容</span>
          </label>
        </div>
      </div>

      {/* 链接地址（只是链接类型时显示） */}
      {resourceForm === "link" && (
        <div className="space-y-1.5">
          <Label htmlFor="resourceUrl" className="text-sm font-medium">
            资源链接 <span className="text-destructive">*</span>
          </Label>
          <Input
            id="resourceUrl"
            placeholder="https://example.com/resource"
            value={resourceUrl}
            onChange={(e) => setResourceUrl(e.target.value)}
            className="h-10"
          />
        </div>
      )}
    </>
  )
}
