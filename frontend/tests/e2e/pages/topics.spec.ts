import { test, expect } from "@playwright/test"

test.describe("话题列表页", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/topics")
    await page.waitForLoadState("networkidle")
  })

  test("页面标题包含'主题'", async ({ page }) => {
    await expect(page).toHaveTitle(/主题|话题/)
    await expect(page.locator("h1")).toContainText("主题")
  })

  test("面包屑导航存在", async ({ page }) => {
    const breadcrumb = page.locator("nav[aria-label='breadcrumb']")
    await expect(breadcrumb).toBeVisible()
    // 首页图标
    await expect(breadcrumb.locator("a[href='/']")).toBeVisible()
  })

  test("Tab 过滤器存在（全部/Go语言/问与答等）", async ({ page }) => {
    // 话题列表有 tab 导航
    await expect(page.locator("a").filter({ hasText: "全部" })).toBeVisible()
    await expect(page.locator("a").filter({ hasText: "Go语言" })).toBeVisible()
  })

  test("切换 tab 更新 URL 参数", async ({ page }) => {
    await page.locator("a[href*='tab=go']").click()
    await expect(page).toHaveURL(/tab=go/)
    await page.waitForLoadState("networkidle")
  })

  test("页面有分页或话题列表", async ({ page }) => {
    // 要么有话题列表，要么有空状态提示
    const hasTopics = await page.locator("[class*='topic'], .topic-item, article").count()
    const hasEmpty = await page.locator("text=暂无话题, text=没有数据").count()
    expect(hasTopics + hasEmpty).toBeGreaterThan(0)
  })
})

test.describe("话题详情页", () => {
  test("访问话题详情（如 ID=1）", async ({ page }) => {
    await page.goto("/topics/1")
    await page.waitForLoadState("networkidle")

    // 要么有内容，要么有错误提示（ID 1 可能不存在）
    const statusOk =
      (await page.locator("h1, article").count()) > 0 ||
      (await page.locator("text=不存在, text=404, text=找不到").count()) > 0
    expect(statusOk).toBeTruthy()
  })
})
