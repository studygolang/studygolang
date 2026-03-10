import { test, expect } from "@playwright/test"
import { SiteHeaderPage } from "../../pages/SiteHeaderPage"
import { mockLogin } from "../../fixtures/auth"

test.describe("站点导航栏", () => {
  test("首页可以正常访问，导航栏渲染完整", async ({ page }) => {
    await page.goto("/")
    await page.waitForLoadState("networkidle")

    const header = new SiteHeaderPage(page)

    // Logo 存在
    await expect(header.logo).toBeVisible()

    // 主导航链接存在（desktop + mobile 各一个，取第一个）
    await expect(page.locator("header a[href='/topics']").first()).toBeVisible()
    await expect(page.locator("header a[href='/articles']")).toBeVisible()
    await expect(page.locator("header a[href='/projects']")).toBeVisible()
    await expect(page.locator("header a[href='/resources']")).toBeVisible()

    // 搜索框存在
    await expect(header.searchInput).toBeVisible()
  })

  test("未登录状态显示登录/注册按钮", async ({ page }) => {
    await page.goto("/")
    await page.waitForLoadState("networkidle")

    const header = new SiteHeaderPage(page)
    await expect(header.loginBtn).toBeVisible()
    await expect(header.registerBtn).toBeVisible()
  })

  test("已登录状态显示用户名，不显示登录/注册按钮", async ({ page }) => {
    await mockLogin(page, "polarisxu")
    await page.goto("/")
    await page.waitForLoadState("networkidle")

    // 应显示用户名
    await expect(page.locator("header").getByText("polarisxu")).toBeVisible()

    // 不应显示登录按钮
    await expect(page.locator("header a[href='/account/login']")).toBeHidden()
    await expect(page.locator("header a[href='/account/register']")).toBeHidden()
  })

  test("搜索框回车跳转到搜索页", async ({ page }) => {
    await page.goto("/")
    await page.waitForLoadState("networkidle")

    const header = new SiteHeaderPage(page)
    await header.searchFor("golang")

    await expect(page).toHaveURL(/\/search\?q=golang/)
  })

  test("导航链接可以正确跳转", async ({ page }) => {
    await page.goto("/")
    await page.waitForLoadState("networkidle")

    // 点击话题链接
    await page.locator("header a[href='/topics']").first().click()
    await expect(page).toHaveURL("/topics")
    await page.waitForLoadState("networkidle")

    // 页面应该有话题相关内容
    await expect(page.locator("h1")).toContainText("主题")
  })

  test("导航菜单中的下载链接指向 go.dev", async ({ page }) => {
    await page.goto("/")

    // 找到导航下拉菜单并打开
    await page.locator("header").getByText("导航").click()

    // 下载链接应该指向外部
    const downloadLink = page.locator("a[href='https://go.dev/dl/']")
    await expect(downloadLink).toBeVisible()
  })

  test("所有导航按钮鼠标悬浮显示 pointer 样式", async ({ page }) => {
    await page.goto("/")
    await page.waitForLoadState("networkidle")

    // 检查 Button 组件有 cursor-pointer
    const loginBtn = page.locator("header a[href='/account/login']")
    const cursor = await loginBtn.evaluate(
      (el) => window.getComputedStyle(el).cursor
    )
    expect(cursor).toBe("pointer")
  })
})
