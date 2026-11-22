import { test, expect } from "@playwright/test"

test("文章排序透传 SSR 且切换时重置页码", async ({ page, request }) => {
  const response = await request.get("/api/v1/articles?p=1&sort=hot")
  const articles = (await response.json()).data.list as Array<{ title: string }>
  test.skip(articles.length === 0, "需要至少一篇文章")

  await page.goto("/articles?p=2&sort=latest")
  await page.getByRole("button", { name: "热门", exact: true }).click()
  await expect(page).toHaveURL(/sort=hot/)
  expect(new URL(page.url()).searchParams.has("p")).toBe(false)
  await expect(page.locator("article").first().getByRole("link").first()).toHaveText(articles[0].title)
})

test("搜索兼容旧 f/p 参数，重新搜索时清理旧参数", async ({ page }) => {
  await page.goto("/search?q=golang&f=title&p=2")
  await expect(page.getByLabel("搜索字段")).toHaveValue("title")

  await page.locator("form").filter({ has: page.getByLabel("搜索字段") }).getByRole("button", { name: "搜索", exact: true }).click()
  await expect(page).toHaveURL(/type=title/)
  const params = new URL(page.url()).searchParams
  expect(params.has("f")).toBe(false)
  expect(params.has("p")).toBe(false)
  expect(params.has("page")).toBe(false)
})

test("话题默认按最近活跃排序", async ({ page }) => {
  await page.goto("/topics")
  await expect(page.getByRole("combobox")).toContainText("最近活跃")
  await expect(page.getByRole("link", { name: "Go语言", exact: true })).toHaveAttribute(
    "href",
    "/topics?tab=go&sort=active"
  )
})

test("项目排序后渲染新的服务端数据", async ({ page, request }) => {
  const latestResponse = await request.get("/api/v1/projects?p=1&sort=latest")
  const hotResponse = await request.get("/api/v1/projects?p=1&sort=hot")
  const latest = (await latestResponse.json()).data.projects.map((project: { name: string }) => project.name)
  const hot = (await hotResponse.json()).data.projects.map((project: { name: string }) => project.name)
  test.skip(latest.length === 0 || latest[0] === hot[0], "需要两种排序首项不同的项目数据")

  await page.goto("/projects?sort=latest")
  await expect(page.locator("h3").first()).toHaveText(latest[0])

  await page.getByRole("combobox").click()
  await page.getByRole("option", { name: "最热", exact: true }).click()
  await expect(page).toHaveURL(/sort=hot/)
  await expect(page.locator("h3").first()).toHaveText(hot[0])
})

test("资源页展示后端分类并进入对应分类", async ({ page, request }) => {
  const response = await request.get("/api/v1/resources/categories")
  const categories = (await response.json()).data.categories as Array<{ catid: number; name: string }>
  test.skip(categories.length === 0, "需要至少一个资源分类")
  const first = categories[0]

  await page.goto("/resources")
  const category = page.getByRole("navigation", { name: "资源分类" }).getByRole("link", { name: first.name })
  await expect(category).toHaveAttribute("href", `/resources/cat/${first.catid}`)
  await category.click()
  await expect(page).toHaveURL(new RegExp(`/resources/cat/${first.catid}$`))
  await expect(page.getByRole("navigation", { name: "资源分类" }).getByRole("link", { name: first.name })).toHaveAttribute("aria-current", "page")
})
