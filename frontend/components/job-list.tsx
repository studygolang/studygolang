/**
 * JobList - 招聘列表组件（Client Component）
 *
 * 注意：此组件使用的是 Mock 数据。
 * 后端目前没有招聘（jobs）相关的 API，无法从服务端获取真实数据。
 * Mock 数据仅用于 UI 展示和交互原型验证。
 *
 * TODO: 后端实现 GET /api/v1/jobs 接口后，将 mock 数据替换为真实 API 调用。
 * TODO: 分页按钮目前为静态展示，接入 API 后需实现真实分页逻辑。
 * TODO: 职位详情链接 /jobs/[id] 对应的页面尚未实现，接入 API 后同步创建。
 */
"use client"

import { useState } from "react"
import {
  MapPin,
  Clock,
  Building2,
  Briefcase,
  Users,
  Flame,
} from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { cn } from "@/lib/utils"

const cityFilters = ["全部", "北京", "上海", "深圳", "杭州", "成都", "远程"]

const expFilters = ["不限", "1-3年", "3-5年", "5-10年", "10年以上"]

interface Job {
  id: number
  title: string
  company: string
  companySize: string
  city: string
  salary: string
  experience: string
  tags: string[]
  time: string
  hot?: boolean
  urgent?: boolean
}

const jobs: Job[] = [
  {
    id: 1,
    title: "Go 高级后端工程师",
    company: "字节跳动",
    companySize: "10000人以上",
    city: "北京",
    salary: "40k-70k",
    experience: "3-5年",
    tags: ["Go", "微服务", "Kubernetes"],
    time: "今天",
    hot: true,
  },
  {
    id: 2,
    title: "Go 基础架构工程师",
    company: "阿里巴巴",
    companySize: "10000人以上",
    city: "杭州",
    salary: "35k-65k",
    experience: "3-5年",
    tags: ["Go", "分布式", "中间件"],
    time: "今天",
    hot: true,
  },
  {
    id: 3,
    title: "Go 后端开发 - 云原生方向",
    company: "腾讯",
    companySize: "10000人以上",
    city: "深圳",
    salary: "30k-60k",
    experience: "3-5年",
    tags: ["Go", "Docker", "云原生"],
    time: "1天前",
  },
  {
    id: 4,
    title: "Go 语言开发工程师",
    company: "PingCAP",
    companySize: "500-999人",
    city: "北京",
    salary: "35k-60k",
    experience: "3-5年",
    tags: ["Go", "TiDB", "数据库"],
    time: "1天前",
    urgent: true,
  },
  {
    id: 5,
    title: "Go 微服务架构师",
    company: "美团",
    companySize: "10000人以上",
    city: "北京",
    salary: "45k-80k",
    experience: "5-10年",
    tags: ["Go", "微服务", "架构设计"],
    time: "2天前",
  },
  {
    id: 6,
    title: "全栈工程师 (Go + React)",
    company: "蚂蚁集团",
    companySize: "10000人以上",
    city: "杭州",
    salary: "30k-55k",
    experience: "3-5年",
    tags: ["Go", "React", "全栈"],
    time: "2天前",
  },
  {
    id: 7,
    title: "Go 后端工程师 - 区块链方向",
    company: "Web3 创业公司",
    companySize: "50-150人",
    city: "远程",
    salary: "40k-70k",
    experience: "3-5年",
    tags: ["Go", "区块链", "Web3"],
    time: "3天前",
    urgent: true,
  },
  {
    id: 8,
    title: "Go SRE 工程师",
    company: "滴滴出行",
    companySize: "10000人以上",
    city: "北京",
    salary: "30k-50k",
    experience: "3-5年",
    tags: ["Go", "SRE", "监控"],
    time: "3天前",
  },
  {
    id: 9,
    title: "Go 初级开发工程师",
    company: "七牛云",
    companySize: "500-999人",
    city: "上海",
    salary: "15k-25k",
    experience: "1-3年",
    tags: ["Go", "对象存储", "CDN"],
    time: "4天前",
  },
  {
    id: 10,
    title: "Go 高级开发 - AI 基础设施",
    company: "商汤科技",
    companySize: "1000-9999人",
    city: "上海",
    salary: "40k-65k",
    experience: "5-10年",
    tags: ["Go", "AI", "GPU调度"],
    time: "4天前",
    hot: true,
  },
  {
    id: 11,
    title: "Go 后端开发 (可远程)",
    company: "Milvus 开源社区",
    companySize: "150-500人",
    city: "远程",
    salary: "25k-45k",
    experience: "1-3年",
    tags: ["Go", "向量数据库", "开源"],
    time: "5天前",
  },
  {
    id: 12,
    title: "Go DevOps 工程师",
    company: "华为云",
    companySize: "10000人以上",
    city: "成都",
    salary: "25k-45k",
    experience: "3-5年",
    tags: ["Go", "DevOps", "CI/CD"],
    time: "5天前",
  },
]

export function JobList() {
  const [activeCity, setActiveCity] = useState("全部")
  const [activeExp, setActiveExp] = useState("不限")

  let filtered = jobs
  if (activeCity !== "全部") {
    filtered = filtered.filter((j) => j.city === activeCity)
  }
  if (activeExp !== "不限") {
    filtered = filtered.filter((j) => j.experience === activeExp)
  }

  return (
    <div>
      {/* Filters */}
      <div className="mb-5 space-y-3 rounded-lg border border-border bg-card p-4">
        {/* City */}
        <div className="flex flex-wrap items-center gap-2">
          <span className="mr-1 text-xs font-semibold text-muted-foreground">
            {"城市:"}
          </span>
          {cityFilters.map((city) => (
            <button
              key={city}
              onClick={() => setActiveCity(city)}
              className={cn(
                "cursor-pointer rounded-full px-3 py-1 text-xs font-medium transition-all",
                activeCity === city
                  ? "bg-primary text-primary-foreground"
                  : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
              )}
            >
              {city}
            </button>
          ))}
        </div>
        {/* Experience */}
        <div className="flex flex-wrap items-center gap-2">
          <span className="mr-1 text-xs font-semibold text-muted-foreground">
            {"经验:"}
          </span>
          {expFilters.map((exp) => (
            <button
              key={exp}
              onClick={() => setActiveExp(exp)}
              className={cn(
                "cursor-pointer rounded-full px-3 py-1 text-xs font-medium transition-all",
                activeExp === exp
                  ? "bg-primary text-primary-foreground"
                  : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
              )}
            >
              {exp}
            </button>
          ))}
        </div>
      </div>

      {/* Result count */}
      <div className="mb-4 flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          {"共 "}
          <span className="font-semibold text-foreground">{filtered.length}</span>
          {" 个职位"}
        </p>
        <div className="flex items-center gap-2">
          {["最新", "薪资最高"].map((sort, i) => (
            <button
              key={sort}
              className={cn(
                "rounded-md px-2.5 py-1 text-xs font-medium transition-colors",
                i === 0
                  ? "bg-primary text-primary-foreground"
                  : "text-muted-foreground hover:bg-secondary"
              )}
            >
              {sort}
            </button>
          ))}
        </div>
      </div>

      {/* Job cards */}
      <div className="space-y-3">
        {filtered.map((job) => (
          <Card
            key={job.id}
            className="group transition-all hover:border-primary/20 hover:shadow-sm"
          >
            <CardContent className="p-4 sm:p-5">
              <div className="flex items-start justify-between gap-4">
                <div className="min-w-0 flex-1">
                  {/* Title row */}
                  {/* TODO: /jobs/[id] 详情页尚未实现，接入后端 API 后改为 <Link href={`/jobs/${job.id}`}> */}
                  <div className="flex flex-wrap items-center gap-2">
                    <span
                      className="text-base font-semibold text-foreground sm:text-[17px]"
                    >
                      {job.title}
                    </span>
                    {job.hot && (
                      <Badge className="gap-0.5 bg-destructive/10 text-[10px] font-semibold text-destructive">
                        <Flame className="h-2.5 w-2.5" />
                        {"热招"}
                      </Badge>
                    )}
                    {job.urgent && (
                      <Badge className="bg-chart-3/10 text-[10px] font-semibold text-chart-3">
                        {"急聘"}
                      </Badge>
                    )}
                  </div>

                  {/* Company + meta */}
                  <div className="mt-2 flex flex-wrap items-center gap-x-4 gap-y-1">
                    <span className="flex items-center gap-1.5 text-sm font-medium text-foreground">
                      <Building2 className="h-3.5 w-3.5 text-muted-foreground" />
                      {job.company}
                    </span>
                    <span className="flex items-center gap-1 text-xs text-muted-foreground">
                      <Users className="h-3 w-3" />
                      {job.companySize}
                    </span>
                  </div>

                  {/* Tags */}
                  <div className="mt-2.5 flex flex-wrap items-center gap-x-3 gap-y-1.5">
                    {job.tags.map((tag) => (
                      <Badge
                        key={tag}
                        variant="secondary"
                        className="text-[10px] font-medium"
                      >
                        {tag}
                      </Badge>
                    ))}
                    <span className="flex items-center gap-1 text-xs text-muted-foreground">
                      <Briefcase className="h-3 w-3" />
                      {job.experience}
                    </span>
                    <span className="flex items-center gap-1 text-xs text-muted-foreground">
                      <Clock className="h-3 w-3" />
                      {job.time}
                    </span>
                  </div>
                </div>

                {/* Salary + Location */}
                <div className="shrink-0 text-right">
                  <p className="text-lg font-bold text-primary">{job.salary}</p>
                  <p className="mt-1 flex items-center justify-end gap-1 text-xs text-muted-foreground">
                    <MapPin className="h-3 w-3" />
                    {job.city}
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* 分页 - 当前为静态 mock，接入真实 API 后需实现分页逻辑 */}
      <div className="mt-6 flex items-center justify-center gap-2">
        <button
          className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-muted-foreground"
          disabled
        >
          {"上一页"}
        </button>
        {[1, 2, 3].map((p) => (
          <button
            key={p}
            className={cn(
              "rounded-md px-3 py-1.5 text-sm font-medium transition-colors",
              p === 1
                ? "bg-primary text-primary-foreground"
                : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
            )}
            disabled
          >
            {p}
          </button>
        ))}
        <button
          className="rounded-md bg-secondary px-3 py-1.5 text-sm font-medium text-muted-foreground"
          disabled
        >
          {"下一页"}
        </button>
      </div>
    </div>
  )
}
