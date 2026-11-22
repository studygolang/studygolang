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

  test("首页频道和子频道与后端配置一致", async ({ page, request }) => {
    const response = await request.get("/api/v1/home")
    expect(response.ok()).toBeTruthy()
    const home = (await response.json()).data
    const navs = home.index_navs ?? []
    test.skip(navs.length === 0, "需要至少一个首页频道")

    const nav = page.getByRole("navigation", { name: "首页频道" })
    await expect(nav.getByRole("link")).toHaveCount(navs.length)
    for (const item of navs) {
      await expect(nav.getByRole("link", { name: item.name, exact: true })).toBeVisible()
    }

    const children = home.cur_nav?.children ?? []
    const childNav = page.getByRole("navigation", { name: "首页子频道" })
    await expect(childNav.getByRole("link")).toHaveCount(children.length)
    for (const child of children) {
      await expect(childNav.getByRole("link", { name: child.name, exact: true })).toHaveAttribute("href", child.uri)
    }
  })

  test("切换首页频道更新 URL 和选中状态", async ({ page, request }) => {
    const response = await request.get("/api/v1/home")
    const home = (await response.json()).data
    const navs = home.index_navs ?? []
    test.skip(navs.length < 2, "需要至少两个首页频道")

    const target = navs[1]
    const nav = page.getByRole("navigation", { name: "首页频道" })
    await nav.getByRole("link", { name: target.name, exact: true }).click()
    await expect(page).toHaveURL(/[?&]tab=/)
    expect(new URL(page.url()).searchParams.get("tab")).toBe(target.tab)
    await expect(nav.getByRole("link", { name: target.name, exact: true })).toHaveAttribute("aria-current", "page")
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
