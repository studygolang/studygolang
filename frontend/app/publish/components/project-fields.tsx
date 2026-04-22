import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

interface ProjectFieldsProps {
  src: string
  setSrc: (src: string) => void
  category: string
  setCategory: (category: string) => void
  lang: string
  setLang: (lang: string) => void
  home: string
  setHome: (home: string) => void
  doc: string
  setDoc: (doc: string) => void
  licence: string
  setLicence: (licence: string) => void
  os: string
  setOs: (os: string) => void
}

export function ProjectFields({
  src,
  setSrc,
  category,
  setCategory,
  lang,
  setLang,
  home,
  setHome,
  doc,
  setDoc,
  licence,
  setLicence,
  os,
  setOs,
}: ProjectFieldsProps) {
  return (
    <>
      {/* 源码地址 */}
      <div className="space-y-1.5">
        <Label htmlFor="src" className="text-sm font-medium">
          源码地址 <span className="text-destructive">*</span>
        </Label>
        <Input
          id="src"
          placeholder="例如：https://github.com/username/repo"
          value={src}
          onChange={(e) => setSrc(e.target.value)}
          className="h-10"
        />
      </div>

      {/* 项目分类 + 开发语言 */}
      <div className="flex gap-4">
        <div className="flex-1 space-y-1.5">
          <Label htmlFor="category" className="text-sm font-medium">项目分类</Label>
          <Input
            id="category"
            placeholder="例如：Web框架、CLI工具"
            value={category}
            onChange={(e) => setCategory(e.target.value)}
            className="h-10"
          />
        </div>
        <div className="flex-1 space-y-1.5">
          <Label htmlFor="lang" className="text-sm font-medium">开发语言</Label>
          <Select value={lang} onValueChange={setLang}>
            <SelectTrigger className="h-10">
              <SelectValue placeholder="选择语言" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="Go">Go</SelectItem>
              <SelectItem value="Rust">Rust</SelectItem>
              <SelectItem value="Python">Python</SelectItem>
              <SelectItem value="JavaScript">JavaScript</SelectItem>
              <SelectItem value="TypeScript">TypeScript</SelectItem>
              <SelectItem value="C">C</SelectItem>
              <SelectItem value="C++">C++</SelectItem>
              <SelectItem value="Java">Java</SelectItem>
              <SelectItem value="其他">其他</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* 项目主页 + 文档地址 */}
      <div className="flex gap-4">
        <div className="flex-1 space-y-1.5">
          <Label htmlFor="home" className="text-sm font-medium">项目主页</Label>
          <Input
            id="home"
            placeholder="https://example.com"
            value={home}
            onChange={(e) => setHome(e.target.value)}
            className="h-10"
          />
        </div>
        <div className="flex-1 space-y-1.5">
          <Label htmlFor="doc" className="text-sm font-medium">文档地址</Label>
          <Input
            id="doc"
            placeholder="https://docs.example.com"
            value={doc}
            onChange={(e) => setDoc(e.target.value)}
            className="h-10"
          />
        </div>
      </div>

      {/* 开源协议 + 操作系统 */}
      <div className="flex gap-4">
        <div className="flex-1 space-y-1.5">
          <Label htmlFor="licence" className="text-sm font-medium">开源协议</Label>
          <Select value={licence} onValueChange={setLicence}>
            <SelectTrigger className="h-10">
              <SelectValue placeholder="选择协议" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="MIT">MIT</SelectItem>
              <SelectItem value="Apache-2.0">Apache-2.0</SelectItem>
              <SelectItem value="GPL-3.0">GPL-3.0</SelectItem>
              <SelectItem value="BSD-3-Clause">BSD-3-Clause</SelectItem>
              <SelectItem value="MPL-2.0">MPL-2.0</SelectItem>
              <SelectItem value="其他">其他</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="flex-1 space-y-1.5">
          <Label htmlFor="os" className="text-sm font-medium">操作系统</Label>
          <Select value={os} onValueChange={setOs}>
            <SelectTrigger className="h-10">
              <SelectValue placeholder="选择操作系统" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="跨平台">跨平台</SelectItem>
              <SelectItem value="Linux">Linux</SelectItem>
              <SelectItem value="macOS">macOS</SelectItem>
              <SelectItem value="Windows">Windows</SelectItem>
              <SelectItem value="其他">其他</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>
    </>
  )
}
