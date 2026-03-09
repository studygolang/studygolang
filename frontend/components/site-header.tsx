"use client"

import { useState, useRef, type KeyboardEvent } from "react"
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
  Github,
  LogIn,
  UserPlus,
  Briefcase,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
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

const docItems = [
  { label: "\u82f1\u6587\u6587\u6863", href: "https://go.dev/doc/" },
  { label: "\u4e2d\u6587\u6587\u6863", href: "/docs/zh" },
  { label: "\u6807\u51c6\u5e93\u4e2d\u6587\u7248", href: "/docs/stdlib" },
  { label: "Go \u6307\u5357", href: "/docs/guide" },
]

export function SiteHeader() {
  const router = useRouter()
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const [searchFocused, setSearchFocused] = useState(false)
  const searchInputRef = useRef<HTMLInputElement>(null)

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
              <button className="flex items-center gap-1 rounded-md px-3 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground">
                <Compass className="h-4 w-4" />
                {"\u5bfc\u822a"}
                <ChevronDown className="h-3 w-3" />
              </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" className="w-48">
              <DropdownMenuItem asChild>
                <Link href="/sites" className="flex items-center gap-2">
                  <ExternalLink className="h-4 w-4" />
                  {"Go\u7f51\u5740\u5bfc\u822a"}
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link href="/dl" className="flex items-center gap-2">
                  <Download className="h-4 w-4" />
                  {"\u4e0b\u8f7d"}
                </Link>
              </DropdownMenuItem>
              {docItems.map((doc) => (
                <DropdownMenuItem key={doc.href} asChild>
                  <Link href={doc.href} className="flex items-center gap-2">
                    <FileText className="h-4 w-4" />
                    {doc.label}
                  </Link>
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
          <Button size="sm" className="h-9 gap-1.5 bg-primary text-primary-foreground hover:bg-primary/90">
            <PenSquare className="h-3.5 w-3.5" />
            {"\u53d1\u5e03"}
          </Button>
          <div className="mx-1 h-5 w-px bg-border" />
          <Button variant="ghost" size="sm" className="h-9 gap-1.5 text-muted-foreground">
            <LogIn className="h-4 w-4" />
            {"\u767b\u5f55"}
          </Button>
          <Button variant="outline" size="sm" className="h-9 gap-1.5">
            <UserPlus className="h-4 w-4" />
            {"\u6ce8\u518c"}
          </Button>
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
            <div className="flex gap-2 px-3 pt-1">
              <Button variant="outline" size="sm" className="flex-1 gap-1.5">
                <Github className="h-4 w-4" />
                GitHub {"\u767b\u5f55"}
              </Button>
              <Button size="sm" className="flex-1 bg-primary text-primary-foreground hover:bg-primary/90">
                {"\u6ce8\u518c"}
              </Button>
            </div>
          </nav>
        </div>
      )}
    </header>
  )
}
