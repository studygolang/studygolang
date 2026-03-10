import { test, expect } from "@playwright/test"

test.describe("首页", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/")
    await page.waitForLoadState("networkidle")
  })

  test("页面标题正确", async ({ page }) => {
    await expect(page).toHaveTitle(/Go语言中文网/)
  })

  test("HeroBanner 渲染，公告可见", async ({ page }) => {
    // Hero banner 区域存在
    const heroBanner = page.locator(".border-b.border-border.bg-card").first()
    await expect(heroBanner).toBeVisible()

    // 公告文字存在
    const announcement = page.locator("a").filter({ hasText: /Go/ }).first()
    await expect(announcement).toBeVisible()
  })

  test("侧边栏 LoginCard 在未登录时显示注册入口", async ({ page }) => {
    const loginCard = page.locator("text=加入 Go 中文社区")
    await expect(loginCard).toBeVisible()

    const registerBtn = page.locator("a[href='/account/register']")
    await expect(registerBtn.first()).toBeVisible()
  })

  test("话题 Tab 切换通过 URL 参数驱动", async ({ page }) => {
    // 点击"最热" tab
    await page.locator("button").filter({ hasText: "最热" }).click()
    await expect(page).toHaveURL(/[?&]tab=hot/)

    // 点击"最新" tab
    await page.locator("button").filter({ hasText: "最新" }).click()
    await expect(page).toHaveURL(/[?&]tab=latest/)
  })

  test("页面底部有 SiteFooter", async ({ page }) => {
    const footer = page.locator("footer")
    await expect(footer).toBeVisible()
    await expect(footer).toContainText("Go语言中文网")
  })

  test("页面包含 SEO meta 标签", async ({ page }) => {
    const description = page.locator('meta[name="description"]')
    await expect(description).toHaveAttribute("content", /Go语言/)

    const keywords = page.locator('meta[name="keywords"]')
    await expect(keywords).toHaveAttribute("content", /Go/)
  })
})
