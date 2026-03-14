import { test, expect } from "@playwright/test"
import { MessagePage } from "../../pages/MessagePage"
import { mockLogin } from "../../fixtures/auth"

/**
 * 消息系统快速验证测试
 *
 * 注意：运行此测试前，请确保：
 * 1. 后端 API 认证响应格式已修复（返回 JSON 而不是 HTML）
 * 2. 前端认证状态检查逻辑已修复
 * 3. 后端服务运行在 http://localhost:8090
 * 4. 前端服务运行在 http://localhost:3000
 */

test.describe("消息系统快速验证", () => {
  let messagePage: MessagePage

  test.beforeEach(async ({ page }) => {
    await mockLogin(page)
    messagePage = new MessagePage(page)
  })

  test("访问消息页面不应超时", async ({ page }) => {
    // 设置较短的超时时间来快速发现问题
    page.setDefaultTimeout(10000)

    await messagePage.goto("system")

    // 页面应该加载完成，显示标签页
    await expect(messagePage.systemTab).toBeVisible({ timeout: 5000 })
    await expect(messagePage.inboxTab).toBeVisible()
    await expect(messagePage.outboxTab).toBeVisible()
  })

  test("API 应返回 JSON 格式的响应", async ({ page }) => {
    await page.goto("http://localhost:3000")

    // 监听 API 请求
    const apiResponse = page.waitForResponse(
      (response) =>
        response.url().includes("/api/v1/messages") &&
        response.request().method() === "GET"
    )

    await messagePage.goto("system")

    const response = await apiResponse
    const contentType = response.headers()["content-type"]

    // 验证返回的是 JSON 而不是 HTML
    expect(contentType).toContain("application/json")

    const body = await response.json()
    expect(body).toHaveProperty("code")
    expect(body).toHaveProperty("data")
  })

  test("未登录时应显示友好的错误提示", async ({ page }) => {
    // 清除登录状态
    await page.evaluate(() => {
      localStorage.clear()
    })

    await page.goto("/message/system")

    // 应该显示登录提示或重定向到登录页
    const hasLoginPrompt = await page.locator("text=请先登录").count()
    const isLoginPage = page.url().includes("/account/login")

    expect(hasLoginPrompt > 0 || isLoginPage).toBeTruthy()
  })

  test("空消息列表应显示提示", async ({ page }) => {
    await messagePage.goto("system")

    // 等待加载完成
    await page.waitForTimeout(2000)

    // 应该显示空状态或消息列表
    const hasEmptyState = (await messagePage.emptyState.count()) > 0
    const hasMessages = (await messagePage.getMessageCount()) > 0

    expect(hasEmptyState || hasMessages).toBeTruthy()
  })

  test("标签页切换应更新 URL", async ({ page }) => {
    await messagePage.goto("system")

    // 切换到收件箱
    await messagePage.inboxTab.click()
    await expect(page).toHaveURL(/\/message\/inbox/)

    // 切换到发件箱
    await messagePage.outboxTab.click()
    await expect(page).toHaveURL(/\/message\/outbox/)

    // 切换回系统消息
    await messagePage.systemTab.click()
    await expect(page).toHaveURL(/\/message\/system/)
  })
})
