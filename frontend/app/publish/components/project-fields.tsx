"use client"

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
  uri: string
  setUri: (uri: string) => void
  category: string
  setCategory: (category: string) => void
  lang: string
  setLang: (lang: string) => void
  home: string
  setHome: (home: string) => void
  doc: string
  setDoc: (doc: string) => void
  download: string
  setDownload: (download: string) => void
  logo: string
  setLogo: (logo: string) => void
  author: string
  setAuthor: (author: string) => void
  licence: string
  setLicence: (licence: string) => void
  os: string
  setOs: (os: string) => void
}

export function ProjectFields({
  src,
  setSrc,
  uri,
  setUri,
  category,
  setCategory,
  lang,
  setLang,
  home,
  setHome,
  doc,
  setDoc,
  download,
  setDownload,
  logo,
  setLogo,
  author,
  setAuthor,
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

      <div className="space-y-1.5">
        <Label htmlFor="uri" className="text-sm font-medium">个性 URI</Label>
        <Input
          id="uri"
          placeholder="留空则根据项目名称生成"
          value={uri}
          onChange={(e) => setUri(e.target.value)}
          className="h-10"
        />
      </div>

      {/* 项目分类 + 开发语言 */}
      <div className="flex flex-col gap-4 sm:flex-row">
        <div className="min-w-0 flex-1 space-y-1.5">
          <Label htmlFor="category" className="text-sm font-medium">项目分类 <span className="text-destructive">*</span></Label>
          <Input
            id="category"
            placeholder="例如：Web框架、CLI工具"
            value={category}
            onChange={(e) => setCategory(e.target.value)}
            className="h-10"
          />
        </div>
        <div className="min-w-0 flex-1 space-y-1.5">
          <Label htmlFor="lang" className="text-sm font-medium">开发语言 <span className="text-destructive">*</span></Label>
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
      <div className="flex flex-col gap-4 sm:flex-row">
        <div className="min-w-0 flex-1 space-y-1.5">
          <Label htmlFor="home" className="text-sm font-medium">项目主页 <span className="text-destructive">*</span></Label>
          <Input
            id="home"
            placeholder="https://example.com"
            value={home}
            onChange={(e) => setHome(e.target.value)}
            className="h-10"
          />
        </div>
        <div className="min-w-0 flex-1 space-y-1.5">
          <Label htmlFor="doc" className="text-sm font-medium">文档地址 <span className="text-destructive">*</span></Label>
          <Input
            id="doc"
            placeholder="https://docs.example.com"
            value={doc}
            onChange={(e) => setDoc(e.target.value)}
            className="h-10"
          />
        </div>
      </div>

      <div className="flex flex-col gap-4 sm:flex-row">
        <div className="min-w-0 flex-1 space-y-1.5">
          <Label htmlFor="download" className="text-sm font-medium">下载地址</Label>
          <Input
            id="download"
            placeholder="https://example.com/download"
            value={download}
            onChange={(e) => setDownload(e.target.value)}
            className="h-10"
          />
        </div>
        <div className="min-w-0 flex-1 space-y-1.5">
          <Label htmlFor="logo" className="text-sm font-medium">Logo URL</Label>
          <Input
            id="logo"
            placeholder="https://example.com/logo.png"
            value={logo}
            onChange={(e) => setLogo(e.target.value)}
            className="h-10"
          />
        </div>
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="author" className="text-sm font-medium">作者</Label>
        <Input
          id="author"
          placeholder="项目作者或组织"
          value={author}
          onChange={(e) => setAuthor(e.target.value)}
          className="h-10"
        />
      </div>

      {/* 开源协议 + 操作系统 */}
      <div className="flex flex-col gap-4 sm:flex-row">
        <div className="min-w-0 flex-1 space-y-1.5">
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
        <div className="min-w-0 flex-1 space-y-1.5">
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
