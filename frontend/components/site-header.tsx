"use client"

import { useState, useRef, useEffect, type KeyboardEvent } from "react"
import { useRouter } from "next/navigation"
import Link from "next/link"
import {
  Search,
  Bell,
  PenSquare,
  Menu,
  X,
  ChevronDown,
  BookOpen,
  FileText,
  FolderGit2,
  Package,
  BookMarked,
  Download,
  ExternalLink,
  Compass,
  LogIn,
  UserPlus,
  Briefcase,
  User,
  LogOut,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

const navItems = [
  { label: "\u4e3b\u9898", href: "/topics", icon: BookOpen },
  { label: "\u6587\u7ae0", href: "/articles", icon: FileText },
  { label: "\u9879\u76ee", href: "/projects", icon: FolderGit2 },
  { label: "\u8d44\u6e90", href: "/resources", icon: Package },
  { label: "\u56fe\u4e66", href: "/books", icon: BookMarked },
  { label: "\u9177\u5de5\u4f5c", href: "/jobs", icon: Briefcase },
]

interface DocItem {
  label: string
  href: string
  /** 是否外部链接，true 时用 <a> 并 target="_blank" */
  external?: boolean
}

// 中文文档指向外部资源，前端暂无 /docs/* 页面
const docItems: DocItem[] = [
  { label: "\u82f1\u6587\u6587\u6863", href: "https://go.dev/doc/", external: true },
  { label: "\u4e2d\u6587\u6587\u6863", href: "https://go-zh.org/doc/", external: true },
  {
    label: "\u6807\u51c6\u5e93\u4e2d\u6587\u7248",
    href: "https://books.studygolang.com/The-Golang-Standard-Library-by-Example/",
    external: true,
  },
  // Go 指南指向 wiki 内部页面
  { label: "Go \u6307\u5357", href: "/wiki" },
]

interface AuthState {
  username: string
  uid: string
}

export function SiteHeader() {
  const router = useRouter()
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const [searchFocused, setSearchFocused] = useState(false)
  const [auth, setAuth] = useState<AuthState | null>(null)
  const searchInputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    const token = localStorage.getItem("token")
    const username = localStorage.getItem("username")
    const uid = localStorage.getItem("uid")
    if (token && username && uid) {
      setAuth({ username, uid })
    }
  }, [])

  function handleLogout() {
    localStorage.removeItem("token")
    localStorage.removeItem("uid")
    localStorage.removeItem("username")
    setAuth(null)
    router.push("/")
  }

  function handleSearchKeyDown(e: KeyboardEvent<HTMLInputElement>) {
    if (e.key !== "Enter") return
    const q = searchInputRef.current?.value.trim()
    if (!q) return
    router.push(`/search?q=${encodeURIComponent(q)}`)
  }

  return (
    <header className="sticky top-0 z-50 border-b border-border bg-card/95 backdrop-blur-sm">
      <div className="mx-auto flex h-16 max-w-7xl items-center gap-4 px-4 lg:px-6">
        {/* Logo */}
        <Link href="/" className="flex shrink-0 items-center gap-2.5">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
            <span className="text-sm font-bold text-primary-foreground">Go</span>
          </div>
          <span className="hidden text-lg font-bold tracking-tight text-foreground sm:inline-block">
            Go<span className="text-primary">{"\u4e2d\u6587\u7f51"}</span>
          </span>
        </Link>

        {/* Desktop Nav */}
        <nav className="hidden items-center gap-1 lg:flex">
          {navItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className="flex items-center gap-1.5 rounded-md px-3 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground"
            >
              {item.label}
            </Link>
          ))}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <button className="flex cursor-pointer items-center gap-1 rounded-md px-3 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground">
                <Compass className="h-4 w-4" />
                {"\u5bfc\u822a"}
                <ChevronDown className="h-3 w-3" />
              </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" className="w-48">
              {/* /sites 前端暂无此页面，临时指向 # */}
              <DropdownMenuItem asChild>
                <a href="#" className="flex items-center gap-2">
                  <ExternalLink className="h-4 w-4" />
                  {"Go\u7f51\u5740\u5bfc\u822a"}
                </a>
              </DropdownMenuItem>
              {/* /dl 前端无此页面，直接跳转 go.dev 官方下载 */}
              <DropdownMenuItem asChild>
                <a href="https://go.dev/dl/" target="_blank" rel="noopener noreferrer" className="flex items-center gap-2">
                  <Download className="h-4 w-4" />
                  {"\u4e0b\u8f7d Go"}
                </a>
              </DropdownMenuItem>
              {docItems.map((doc) => (
                <DropdownMenuItem key={doc.href} asChild>
                  {doc.external ? (
                    <a href={doc.href} target="_blank" rel="noopener noreferrer" className="flex items-center gap-2">
                      <FileText className="h-4 w-4" />
                      {doc.label}
                    </a>
                  ) : (
                    <Link href={doc.href} className="flex items-center gap-2">
                      <FileText className="h-4 w-4" />
                      {doc.label}
                    </Link>
                  )}
                </DropdownMenuItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
        </nav>

        {/* Search */}
        <div className="relative ml-auto flex-1 lg:max-w-xs">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            ref={searchInputRef}
            placeholder={"\u641c\u7d22\u4e3b\u9898\u3001\u6587\u7ae0\u3001\u8d44\u6e90..."}
            className={`h-9 w-full pl-9 text-sm transition-all ${
              searchFocused ? "ring-2 ring-primary/30" : ""
            }`}
            onFocus={() => setSearchFocused(true)}
            onBlur={() => setSearchFocused(false)}
            onKeyDown={handleSearchKeyDown}
          />
        </div>

        {/* Actions */}
        <div className="hidden items-center gap-2 md:flex">
          <Button variant="ghost" size="icon" className="h-9 w-9 text-muted-foreground" aria-label={"\u901a\u77e5"}>
            <Bell className="h-4 w-4" />
          </Button>
          {/* /topics/create 页面暂未实现，临时跳转 /topics */}
          <Button asChild size="sm" className="h-9 gap-1.5 bg-primary text-primary-foreground hover:bg-primary/90">
            <Link href="/topics">
              <PenSquare className="h-3.5 w-3.5" />
              {"\u53d1\u5e03"}
            </Link>
          </Button>
          <div className="mx-1 h-5 w-px bg-border" />
          {auth ? (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="sm" className="h-9 gap-1.5 text-muted-foreground">
                  <User className="h-4 w-4" />
                  {auth.username}
                  <ChevronDown className="h-3 w-3" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-40">
                <DropdownMenuItem asChild>
                  <Link href={`/user/${auth.username}`} className="flex items-center gap-2">
                    <User className="h-4 w-4" />
                    {"个人主页"}
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={handleLogout} className="flex items-center gap-2 text-destructive focus:text-destructive">
                  <LogOut className="h-4 w-4" />
                  {"退出登录"}
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          ) : (
            <>
              <Button asChild variant="ghost" size="sm" className="h-9 gap-1.5 text-muted-foreground">
                <Link href="/account/login">
                  <LogIn className="h-4 w-4" />
                  {"\u767b\u5f55"}
                </Link>
              </Button>
              <Button asChild variant="outline" size="sm" className="h-9 gap-1.5">
                <Link href="/account/register">
                  <UserPlus className="h-4 w-4" />
                  {"\u6ce8\u518c"}
                </Link>
              </Button>
            </>
          )}
        </div>

        {/* Mobile menu toggle */}
        <Button
          variant="ghost"
          size="icon"
          className="h-9 w-9 lg:hidden"
          onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
          aria-label={mobileMenuOpen ? "\u5173\u95ed\u83dc\u5355" : "\u6253\u5f00\u83dc\u5355"}
        >
          {mobileMenuOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
        </Button>
      </div>

      {/* Mobile Menu */}
      {mobileMenuOpen && (
        <div className="border-t border-border bg-card lg:hidden">
          <nav className="mx-auto max-w-7xl space-y-1 px-4 py-3">
            {navItems.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className="flex items-center gap-3 rounded-md px-3 py-2.5 text-sm font-medium text-foreground hover:bg-secondary"
                onClick={() => setMobileMenuOpen(false)}
              >
                <item.icon className="h-4 w-4 text-muted-foreground" />
                {item.label}
              </Link>
            ))}
            <div className="my-2 h-px bg-border" />
            {auth ? (
              <div className="flex flex-col gap-1 px-3 pt-1">
                <Button asChild variant="ghost" size="sm" className="justify-start gap-2">
                  <Link href={`/user/${auth.username}`} onClick={() => setMobileMenuOpen(false)}>
                    <User className="h-4 w-4" />
                    {auth.username}
                  </Link>
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  className="justify-start gap-2 text-destructive hover:text-destructive"
                  onClick={() => { setMobileMenuOpen(false); handleLogout() }}
                >
                  <LogOut className="h-4 w-4" />
                  {"退出登录"}
                </Button>
              </div>
            ) : (
              <div className="flex gap-2 px-3 pt-1">
                <Button asChild variant="outline" size="sm" className="flex-1 gap-1.5">
                  <Link href="/account/login" onClick={() => setMobileMenuOpen(false)}>
                    <LogIn className="h-4 w-4" />
                    {"\u767b\u5f55"}
                  </Link>
                </Button>
                <Button asChild size="sm" className="flex-1 bg-primary text-primary-foreground hover:bg-primary/90">
                  <Link href="/account/register" onClick={() => setMobileMenuOpen(false)}>
                    {"\u6ce8\u518c"}
                  </Link>
                </Button>
              </div>
            )}
          </nav>
        </div>
      )}
    </header>
  )
}
