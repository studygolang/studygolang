import { test, expect } from "@playwright/test"

test.describe("搜索功能", () => {
  test("通过 URL 参数直接访问搜索页", async ({ page }) => {
    await page.goto("/search?q=golang")
    await page.waitForLoadState("networkidle")

    await expect(page.locator("h1, [class*='title']")).toContainText(/golang|搜索/i)
  })

  test("搜索结果页展示搜索词", async ({ page }) => {
    await page.goto("/search?q=goroutine")
    await page.waitForLoadState("networkidle")

    // 页面上应该显示搜索词
    await expect(page.locator("body")).toContainText(/goroutine/i)
  })

  test("空搜索词处理正常", async ({ page }) => {
    await page.goto("/search")
    await page.waitForLoadState("networkidle")

    // 不应崩溃
    await expect(page.locator("body")).toBeVisible()
  })

  test("从导航栏搜索跳转到搜索页", async ({ page }) => {
    await page.goto("/")
    await page.waitForLoadState("networkidle")

    const searchInput = page.locator("header input[placeholder*='搜索']")
    await searchInput.fill("channel")
    await searchInput.press("Enter")

    await expect(page).toHaveURL(/\/search\?q=channel/)
  })
})
