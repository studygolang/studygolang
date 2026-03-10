"use client"

import { useState } from "react"
import Link from "next/link"
import {
  MapPin,
  Clock,
  Building2,
  Briefcase,
  Users,
  Flame,
  Banknote,
} from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { cn } from "@/lib/utils"
import type { Job } from "@/lib/types"

const cityFilters = ["全部", "北京", "上海", "深圳", "杭州", "成都", "远程"]
const expFilters = ["不限", "1-3年", "3-5年", "5-10年", "10年以上"]

function formatSalary(min: number, max: number): string {
  if (min === 0 && max === 0) return "面议"
  if (min === 0) return `${max}k 以下`
  if (max === 0) return `${min}k 以上`
  return `${min}k-${max}k`
}

interface JobListProps {
  jobs: Job[]
}

export function JobList({ jobs }: JobListProps) {
  const [activeCity, setActiveCity] = useState("全部")
  const [activeExp, setActiveExp] = useState("不限")

  let filtered = jobs
  if (activeCity !== "全部") {
    filtered = filtered.filter((j) => j.city === activeCity)
  }
  if (activeExp !== "不限") {
    filtered = filtered.filter((j) => j.experience === activeExp)
  }

  if (jobs.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-20 text-center">
        <Briefcase className="mb-3 h-10 w-10 text-muted-foreground/40" />
        <p className="text-sm font-medium text-foreground">暂无职位信息</p>
        <p className="mt-1 text-sm text-muted-foreground">期待更多 Go 职位上线</p>
      </div>
    )
  }

  return (
    <div>
      {/* Filters */}
      <div className="mb-5 space-y-3 rounded-lg border border-border bg-card p-4">
        {/* City */}
        <div className="flex flex-wrap items-center gap-2">
          <span className="mr-1 text-xs font-semibold text-muted-foreground">{"城市:"}</span>
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
          <span className="mr-1 text-xs font-semibold text-muted-foreground">{"经验:"}</span>
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
      </div>

      {/* Job cards */}
      {filtered.length === 0 ? (
        <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-12 text-center">
          <Briefcase className="mb-3 h-8 w-8 text-muted-foreground/40" />
          <p className="text-sm text-muted-foreground">当前筛选条件下暂无职位</p>
        </div>
      ) : (
        <div className="space-y-3">
          {filtered.map((job) => {
            const tagList = job.tags ? job.tags.split(",").filter(Boolean) : []
            return (
              <Card
                key={job.id}
                className="group transition-all hover:border-primary/20 hover:shadow-sm"
              >
                <CardContent className="p-4 sm:p-5">
                  <div className="flex items-start justify-between gap-4">
                    <div className="min-w-0 flex-1">
                      {/* Title */}
                      <div className="flex flex-wrap items-center gap-2">
                        <Link
                          href={`/jobs/${job.id}`}
                          className="text-base font-semibold text-foreground group-hover:text-primary sm:text-[17px]"
                        >
                          {job.title}
                        </Link>
                        {job.experience === "1-3年" && (
                          <Badge className="gap-0.5 bg-chart-4/10 text-[10px] font-semibold text-chart-4">
                            <Flame className="h-2.5 w-2.5" />
                            {"热招"}
                          </Badge>
                        )}
                      </div>

                      {/* Company */}
                      <div className="mt-2 flex flex-wrap items-center gap-x-4 gap-y-1">
                        <span className="flex items-center gap-1.5 text-sm font-medium text-foreground">
                          <Building2 className="h-3.5 w-3.5 text-muted-foreground" />
                          {job.company}
                        </span>
                        {job.company_size && (
                          <span className="flex items-center gap-1 text-xs text-muted-foreground">
                            <Users className="h-3 w-3" />
                            {job.company_size}
                          </span>
                        )}
                      </div>

                      {/* Tags + meta */}
                      <div className="mt-2.5 flex flex-wrap items-center gap-x-3 gap-y-1.5">
                        {tagList.slice(0, 3).map((tag) => (
                          <Badge key={tag} variant="secondary" className="text-[10px] font-medium">
                            {tag}
                          </Badge>
                        ))}
                        {job.experience && (
                          <span className="flex items-center gap-1 text-xs text-muted-foreground">
                            <Briefcase className="h-3 w-3" />
                            {job.experience}
                          </span>
                        )}
                        <span className="flex items-center gap-1 text-xs text-muted-foreground">
                          <Clock className="h-3 w-3" />
                          {job.ctime}
                        </span>
                      </div>
                    </div>

                    {/* Salary + location */}
                    <div className="shrink-0 text-right">
                      <p className="flex items-center justify-end gap-1 text-lg font-bold text-primary">
                        <Banknote className="h-4 w-4" />
                        {formatSalary(job.salary_min, job.salary_max)}
                      </p>
                      <p className="mt-1 flex items-center justify-end gap-1 text-xs text-muted-foreground">
                        <MapPin className="h-3 w-3" />
                        {job.city}
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            )
          })}
        </div>
      )}
    </div>
  )
}
