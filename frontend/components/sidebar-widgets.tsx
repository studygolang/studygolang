"use client"

import Link from "next/link"
import {
  Github,
  TrendingUp,
  Users,
  MessageSquare,
  BookOpen,
  Trophy,
  ArrowRight,
  Flame,
  Calendar,
  BarChart3,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"

/* ---------- Login Card ---------- */
export function LoginCard() {
  return (
    <Card>
      <CardContent className="p-4">
        <div className="mb-3 text-center">
          <div className="mx-auto mb-2 flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
            <Users className="h-6 w-6 text-primary" />
          </div>
          <p className="text-sm font-medium text-foreground">
            {"\u52a0\u5165 Go \u4e2d\u6587\u793e\u533a"}
          </p>
          <p className="mt-1 text-xs text-muted-foreground">
            {"\u4e0e\u5168\u56fd Gopher \u4e00\u8d77\u5b66\u4e60\u4ea4\u6d41"}
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm" className="flex-1 gap-1.5 text-xs">
            <Github className="h-3.5 w-3.5" />
            GitHub
          </Button>
          <Button size="sm" className="flex-1 bg-primary text-xs text-primary-foreground hover:bg-primary/90">
            {"\u6ce8\u518c\u8d26\u53f7"}
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}

/* ---------- Daily Interview Question ---------- */
export function DailyQuestion() {
  return (
    <Card>
      <CardHeader className="p-4 pb-2">
        <CardTitle className="flex items-center gap-2 text-sm font-semibold">
          <Flame className="h-4 w-4 text-primary" />
          {"Go \u4eca\u65e5\u9762\u8bd5\u9898"}
        </CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-4">
        <Link
          href="/interview/today"
          className="group block rounded-md bg-secondary/50 p-3 transition-colors hover:bg-secondary"
        >
          <p className="text-sm font-medium leading-relaxed text-foreground group-hover:text-primary">
            {"Go \u4e2d slice \u548c array \u7684\u533a\u522b\u662f\u4ec0\u4e48\uff1f\u5982\u4f55\u907f\u514d slice \u7684\u5185\u5b58\u6cc4\u6f0f\uff1f"}
          </p>
          <span className="mt-2 inline-flex items-center gap-1 text-xs font-medium text-primary">
            {"\u67e5\u770b\u89e3\u7b54"}
            <ArrowRight className="h-3 w-3" />
          </span>
        </Link>
      </CardContent>
    </Card>
  )
}

/* ---------- Morning Reading ---------- */
export function MorningReading() {
  const articles = [
    {
      title: "Go 1.26 \u4e2d\u7684\u65b0\u7279\u6027\u5168\u89e3\u6790",
      href: "/articles/1",
    },
    {
      title: "\u5982\u4f55\u7528 Go \u6784\u5efa\u9ad8\u6027\u80fd\u7f51\u5173",
      href: "/articles/2",
    },
    {
      title: "Go \u5e76\u53d1\u7f16\u7a0b\u6700\u4f73\u5b9e\u8df5 2026 \u7248",
      href: "/articles/3",
    },
  ]

  return (
    <Card>
      <CardHeader className="p-4 pb-2">
        <CardTitle className="flex items-center gap-2 text-sm font-semibold">
          <Calendar className="h-4 w-4 text-primary" />
          {"\u4eca\u65e5\u6668\u8bfb"}
        </CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-4">
        <div className="space-y-2.5">
          {articles.map((article, i) => (
            <Link
              key={i}
              href={article.href}
              className="flex items-start gap-2 text-sm text-muted-foreground transition-colors hover:text-primary"
            >
              <span className="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded bg-secondary text-xs font-semibold text-secondary-foreground">
                {i + 1}
              </span>
              <span className="line-clamp-2 leading-snug">{article.title}</span>
            </Link>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}

/* ---------- Trending Topics ---------- */
export function TrendingTopics() {
  const trending = [
    { title: "Go \u6cdb\u578b\u5b9e\u6218\u6280\u5de7\u603b\u7ed3", views: "3.2k", hot: true },
    { title: "gRPC \u4e0e RESTful API \u5bf9\u6bd4\u5206\u6790", views: "2.8k", hot: true },
    { title: "Docker \u591a\u9636\u6bb5\u6784\u5efa\u4f18\u5316 Go \u955c\u50cf", views: "2.1k", hot: false },
    { title: "\u5fae\u670d\u52a1\u67b6\u6784\u4e0b\u7684\u670d\u52a1\u53d1\u73b0\u4e0e\u6cbb\u7406", views: "1.9k", hot: false },
    { title: "Go \u8bed\u8a00\u5e74\u85aa 50W \u5fc5\u987b\u638c\u63e1\u7684 10 \u4e2a\u5b9e\u6218\u6280\u5de7", views: "1.7k", hot: false },
  ]

  return (
    <Card>
      <CardHeader className="p-4 pb-2">
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2 text-sm font-semibold">
            <TrendingUp className="h-4 w-4 text-primary" />
            {"\u70ed\u95e8\u8bdd\u9898"}
          </CardTitle>
          <Link href="/trending" className="text-xs text-muted-foreground hover:text-primary">
            {"\u66f4\u591a"}
          </Link>
        </div>
      </CardHeader>
      <CardContent className="px-4 pb-4">
        <div className="space-y-2.5">
          {trending.map((item, i) => (
            <Link
              key={i}
              href={`/topics/${i}`}
              className="flex items-start gap-2 text-sm transition-colors hover:text-primary"
            >
              <span
                className={`mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded text-xs font-bold ${
                  i < 2
                    ? "bg-primary text-primary-foreground"
                    : "bg-secondary text-secondary-foreground"
                }`}
              >
                {i + 1}
              </span>
              <span className="line-clamp-2 flex-1 leading-snug text-foreground">
                {item.title}
              </span>
              <span className="shrink-0 text-xs text-muted-foreground">
                {item.views}
              </span>
            </Link>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}

/* ---------- Latest Comments ---------- */
export function LatestComments() {
  const comments = [
    {
      user: "alice_go",
      initial: "A",
      content: "\u8fd9\u7bc7\u6587\u7ae0\u5199\u5f97\u771f\u597d\uff0c\u5b66\u5230\u4e86\u5f88\u591a\uff01",
      topic: "Go \u5e76\u53d1\u7f16\u7a0b\u6700\u4f73\u5b9e\u8df5",
      time: "10 \u5206\u949f\u524d",
    },
    {
      user: "bob_dev",
      initial: "B",
      content: "\u8bf7\u95ee sync.Pool \u548c\u5bf9\u8c61\u6c60\u6709\u4ec0\u4e48\u533a\u522b\uff1f",
      topic: "\u6df1\u5165\u7406\u89e3 sync.Pool",
      time: "30 \u5206\u949f\u524d",
    },
    {
      user: "charlie",
      initial: "C",
      content: "\u5df2\u661f\u6807\uff0c\u975e\u5e38\u5b9e\u7528\u7684\u603b\u7ed3\uff01",
      topic: "Go \u6cdb\u578b\u5b9e\u6218\u6280\u5de7",
      time: "1 \u5c0f\u65f6\u524d",
    },
  ]

  return (
    <Card>
      <CardHeader className="p-4 pb-2">
        <CardTitle className="flex items-center gap-2 text-sm font-semibold">
          <MessageSquare className="h-4 w-4 text-primary" />
          {"\u6700\u65b0\u8bc4\u8bba"}
        </CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-4">
        <div className="space-y-3">
          {comments.map((comment, i) => (
            <div key={i} className="space-y-1">
              <div className="flex items-center gap-2">
                <Avatar className="h-5 w-5">
                  <AvatarFallback className="bg-primary/10 text-[10px] font-semibold text-primary">
                    {comment.initial}
                  </AvatarFallback>
                </Avatar>
                <span className="text-xs font-medium text-foreground">
                  {comment.user}
                </span>
                <span className="text-[10px] text-muted-foreground">
                  {comment.time}
                </span>
              </div>
              <p className="line-clamp-1 pl-7 text-xs text-muted-foreground">
                {comment.content}
              </p>
              <Link
                href="#"
                className="block pl-7 text-[10px] text-primary hover:underline"
              >
                {comment.topic}
              </Link>
              {i < comments.length - 1 && <Separator className="mt-2" />}
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}

/* ---------- Active Members ---------- */
export function ActiveMembers() {
  const members = [
    { name: "polaris", initial: "P" },
    { name: "guo_hongzhi", initial: "G" },
    { name: "Minnie", initial: "M" },
    { name: "Share111", initial: "S" },
    { name: "chowyu12", initial: "C" },
    { name: "KPbl", initial: "K" },
    { name: "umansyds", initial: "U" },
    { name: "rainyun", initial: "R" },
  ]

  return (
    <Card>
      <CardHeader className="p-4 pb-2">
        <CardTitle className="flex items-center gap-2 text-sm font-semibold">
          <Trophy className="h-4 w-4 text-primary" />
          {"\u6d3b\u8dc3\u4f1a\u5458"}
        </CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-4">
        <div className="flex flex-wrap gap-2">
          {members.map((member) => (
            <Link
              key={member.name}
              href={`/user/${member.name}`}
              className="group flex items-center gap-1.5 rounded-full bg-secondary py-1 pl-1 pr-2.5 transition-colors hover:bg-primary/10"
            >
              <Avatar className="h-5 w-5">
                <AvatarFallback className="bg-primary/10 text-[10px] font-semibold text-primary">
                  {member.initial}
                </AvatarFallback>
              </Avatar>
              <span className="text-xs text-secondary-foreground group-hover:text-primary">
                {member.name}
              </span>
            </Link>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}

/* ---------- Stats Card ---------- */
export function StatsCard() {
  const stats = [
    { label: "\u4f1a\u5458\u6570", value: "189,432", icon: Users },
    { label: "\u4e3b\u9898\u6570", value: "52,871", icon: BookOpen },
    { label: "\u8bc4\u8bba\u6570", value: "328,109", icon: MessageSquare },
    { label: "\u6587\u7ae0\u6570", value: "41,256", icon: BarChart3 },
  ]

  return (
    <Card>
      <CardHeader className="p-4 pb-2">
        <CardTitle className="text-sm font-semibold">{"\u793e\u533a\u7edf\u8ba1"}</CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-4">
        <div className="grid grid-cols-2 gap-3">
          {stats.map((stat) => (
            <div key={stat.label} className="rounded-md bg-secondary/50 p-2.5 text-center">
              <stat.icon className="mx-auto mb-1 h-4 w-4 text-primary" />
              <p className="text-sm font-bold text-foreground">{stat.value}</p>
              <p className="text-[10px] text-muted-foreground">{stat.label}</p>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}

/* ---------- Friend Links ---------- */
export function FriendLinks() {
  const links = [
    { name: "Go \u5b98\u65b9\u7f51\u7ad9", href: "https://go.dev" },
    { name: "Go \u5305\u67e5\u8be2", href: "https://pkg.go.dev" },
    { name: "GitHub Trending", href: "https://github.com/trending/go" },
    { name: "Awesome Go", href: "https://awesome-go.com" },
  ]

  return (
    <Card>
      <CardHeader className="p-4 pb-2">
        <CardTitle className="text-sm font-semibold">{"\u53cb\u60c5\u94fe\u63a5"}</CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-4">
        <div className="flex flex-wrap gap-2">
          {links.map((link) => (
            <a
              key={link.name}
              href={link.href}
              target="_blank"
              rel="noopener noreferrer"
              className="text-xs text-muted-foreground transition-colors hover:text-primary"
            >
              {link.name}
            </a>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}
