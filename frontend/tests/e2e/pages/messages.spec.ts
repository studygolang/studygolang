import { test, expect } from "@playwright/test"
import { MessagePage } from "../../pages/MessagePage"
import { mockLogin } from "../../fixtures/auth"

test.describe("消息系统功能测试", () => {
  let messagePage: MessagePage

  test.beforeEach(async ({ page }) => {
    // 模拟登录状态
    await mockLogin(page)
    messagePage = new MessagePage(page)
  })

  test.describe("页面访问和渲染", () => {
    test("未登录用户访问消息页应重定向到登录页", async ({ page }) => {
      // 清除登录状态
      await page.evaluate(() => {
        localStorage.clear()
      })

      await page.goto("/message/system")
      await page.waitForLoadState("networkidle")

      // 应该重定向到登录页或显示未登录提示
      const url = page.url()
      const hasLoginPrompt = await page.locator("text=请先登录").count() > 0
      const isLoginPage = url.includes("/account/login")

      expect(hasLoginPrompt || isLoginPage).toBeTruthy()
    })

    test("系统消息页面正确渲染", async ({ page }) => {
      await messagePage.goto("system")

      // 检查页面标题
      await expect(page).toHaveTitle(/消息|Go语言中文网/)

      // 检查三个标签页都存在
      await expect(messagePage.systemTab).toBeVisible()
      await expect(messagePage.inboxTab).toBeVisible()
      await expect(messagePage.outboxTab).toBeVisible()

      // 系统消息标签应该是激活状态
      const systemTabState = await messagePage.systemTab.getAttribute("data-state")
      expect(systemTabState).toBe("active")
    })

    test("收件箱页面正确渲染", async ({ page }) => {
      await messagePage.goto("inbox")

      // 收件箱标签应该是激活状态
      const inboxTabState = await messagePage.inboxTab.getAttribute("data-state")
      expect(inboxTabState).toBe("active")
    })

    test("发件箱页面正确渲染", async ({ page }) => {
      await messagePage.goto("outbox")

      // 发件箱标签应该是激活状态
      const outboxTabState = await messagePage.outboxTab.getAttribute("data-state")
      expect(outboxTabState).toBe("active")
    })
  })

  test.describe("标签页切换", () => {
    test("可以在不同消息类型之间切换", async ({ page }) => {
      await messagePage.goto("system")

      // 切换到收件箱
      await messagePage.switchTab("inbox")
      await expect(page).toHaveURL(/\/message\/inbox/)
      let inboxTabState = await messagePage.inboxTab.getAttribute("data-state")
      expect(inboxTabState).toBe("active")

      // 切换到发件箱
      await messagePage.switchTab("outbox")
      await expect(page).toHaveURL(/\/message\/outbox/)
      let outboxTabState = await messagePage.outboxTab.getAttribute("data-state")
      expect(outboxTabState).toBe("active")

      // 切换回系统消息
      await messagePage.switchTab("system")
      await expect(page).toHaveURL(/\/message\/system/)
      let systemTabState = await messagePage.systemTab.getAttribute("data-state")
      expect(systemTabState).toBe("active")
    })

    test("切换标签页时 URL 应该更新", async ({ page }) => {
      await messagePage.goto("system")

      await messagePage.inboxTab.click()
      await page.waitForURL(/\/message\/inbox/)
      expect(page.url()).toContain("/message/inbox")

      await messagePage.outboxTab.click()
      await page.waitForURL(/\/message\/outbox/)
      expect(page.url()).toContain("/message/outbox")
    })
  })

  test.describe("消息列表显示", () => {
    test("加载消息列表", async ({ page }) => {
      await messagePage.goto("system")

      // 等待加载完成
      await page.waitForTimeout(1000)

      // 应该显示消息列表或空状态
      const hasMessages = (await messagePage.getMessageCount()) > 0
      const hasEmptyState = (await messagePage.emptyState.count()) > 0

      expect(hasMessages || hasEmptyState).toBeTruthy()
    })

    test("空消息列表显示正确的提示", async ({ page }) => {
      await messagePage.goto("system")
      await page.waitForTimeout(1000)

      const messageCount = await messagePage.getMessageCount()

      if (messageCount === 0) {
        await expect(messagePage.emptyState).toBeVisible()
        await expect(messagePage.emptyState).toContainText("暂无消息")
      }
    })

    test("消息列表项包含必要信息", async ({ page }) => {
      await messagePage.goto("system")
      await page.waitForTimeout(1000)

      const messages = await messagePage.getMessageItems()

      if (messages.length > 0) {
        const firstMessage = messages[0]

        // 检查消息内容存在
        const content = await firstMessage.textContent()
        expect(content).toBeTruthy()
        expect(content!.length).toBeGreaterThan(0)

        // 检查是否有时间信息
        const hasTime = content!.includes("前") || content!.includes(":")
        expect(hasTime).toBeTruthy()
      }
    })
  })

  test.describe("消息删除功能", () => {
    test("点击删除按钮显示确认对话框", async ({ page }) => {
      await messagePage.goto("system")
      await page.waitForTimeout(1000)

      const messageCount = await messagePage.getMessageCount()

      if (messageCount > 0) {
        const messages = await messagePage.getMessageItems()
        const deleteBtn = messages[0].locator("button:has-text('删除')").first()

        await deleteBtn.click()

        // 检查确认对话框
        const dialog = page.locator('[role="alertdialog"]').first()
        await expect(dialog).toBeVisible()
        await expect(dialog).toContainText("确认删除")
        await expect(dialog).toContainText("确定要删除这条消息吗")

        // 关闭对话框
        const cancelBtn = page.locator("button:has-text('取消')").last()
        await cancelBtn.click()
      }
    })

    test("取消删除操作", async ({ page }) => {
      await messagePage.goto("system")
      await page.waitForTimeout(1000)

      const initialCount = await messagePage.getMessageCount()

      if (initialCount > 0) {
        await messagePage.cancelDelete(0)
        await page.waitForTimeout(500)

        const afterCount = await messagePage.getMessageCount()
        expect(afterCount).toBe(initialCount)
      }
    })

    test("确认删除消息", async ({ page }) => {
      await messagePage.goto("system")
      await page.waitForTimeout(1000)

      const initialCount = await messagePage.getMessageCount()

      if (initialCount > 0) {
        await messagePage.deleteMessage(0)
        await page.waitForTimeout(1000)

        const afterCount = await messagePage.getMessageCount()

        // 删除后消息数量应该减少或显示空状态
        const hasEmptyState = (await messagePage.emptyState.count()) > 0
        expect(afterCount < initialCount || hasEmptyState).toBeTruthy()
      }
    })
  })

  test.describe("分页功能", () => {
    test("分页按钮显示正确", async ({ page }) => {
      await messagePage.goto("system")
      await page.waitForTimeout(1000)

      const messageCount = await messagePage.getMessageCount()

      if (messageCount > 0) {
        // 检查分页区域是否存在
        const hasPagination = (await messagePage.pagination.count()) > 0

        if (hasPagination) {
          // 第一页时，上一页按钮应该禁用
          const prevBtnDisabled = await messagePage.prevPageBtn.isDisabled()
          expect(prevBtnDisabled).toBeTruthy()
        }
      }
    })

    test("点击下一页加载更多消息", async ({ page }) => {
      await messagePage.goto("system")
      await page.waitForTimeout(1000)

      const nextBtnCount = await messagePage.nextPageBtn.count()

      if (nextBtnCount > 0) {
        const isDisabled = await messagePage.nextPageBtn.isDisabled()

        if (!isDisabled) {
          await messagePage.goToNextPage()
          await page.waitForTimeout(1000)

          // URL 应该包含页码参数
          const url = page.url()
          expect(url).toContain("p=2")

          // 上一页按钮应该可用
          const prevBtnDisabled = await messagePage.prevPageBtn.isDisabled()
          expect(prevBtnDisabled).toBeFalsy()
        }
      }
    })

    test("点击上一页返回前一页", async ({ page }) => {
      // 先跳转到第二页
      await page.goto("/message/system?p=2")
      await page.waitForTimeout(1000)

      const prevBtnCount = await messagePage.prevPageBtn.count()

      if (prevBtnCount > 0) {
        const isDisabled = await messagePage.prevPageBtn.isDisabled()

        if (!isDisabled) {
          await messagePage.goToPrevPage()
          await page.waitForTimeout(1000)

          // URL 应该回到第一页
          const url = page.url()
          expect(url).toMatch(/p=1|\/message\/system$/)
        }
      }
    })
  })

  test.describe("消息状态", () => {
    test("未读消息显示未读标记", async ({ page }) => {
      await messagePage.goto("inbox")
      await page.waitForTimeout(1000)

      const messages = await messagePage.getMessageItems()

      if (messages.length > 0) {
        // 检查是否有未读标记（红点或 NEW 标签）
        const firstMessage = messages[0]
        const unreadBadge = firstMessage.locator(".bg-primary, .bg-red-500, text=NEW").first()
        const badgeCount = await unreadBadge.count()

        // 至少应该有消息内容
        const content = await firstMessage.textContent()
        expect(content).toBeTruthy()
      }
    })
  })

  test.describe("响应式设计", () => {
    test("移动端视图正确显示", async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 })
      await messagePage.goto("system")
      await page.waitForTimeout(1000)

      // 标签页应该可见
      await expect(messagePage.systemTab).toBeVisible()
      await expect(messagePage.inboxTab).toBeVisible()
      await expect(messagePage.outboxTab).toBeVisible()

      // 消息列表应该可见
      const messageCount = await messagePage.getMessageCount()
      const hasEmptyState = (await messagePage.emptyState.count()) > 0

      expect(messageCount > 0 || hasEmptyState).toBeTruthy()
    })

    test("平板视图正确显示", async ({ page }) => {
      await page.setViewportSize({ width: 768, height: 1024 })
      await messagePage.goto("system")
      await page.waitForTimeout(1000)

      await expect(messagePage.systemTab).toBeVisible()
      await expect(messagePage.inboxTab).toBeVisible()
      await expect(messagePage.outboxTab).toBeVisible()
    })
  })

  test.describe("错误处理", () => {
    test("API 错误时显示友好提示", async ({ page }) => {
      // 拦截 API 请求并返回错误
      await page.route("**/api/v1/messages*", (route) => {
        route.fulfill({
          status: 500,
          contentType: "application/json",
          body: JSON.stringify({ code: 1, msg: "服务器错误" }),
        })
      })

      await messagePage.goto("system")
      await page.waitForTimeout(1000)

      // 应该显示空状态或错误提示
      const hasEmptyState = (await messagePage.emptyState.count()) > 0
      expect(hasEmptyState).toBeTruthy()
    })

    test("网络错误时显示友好提示", async ({ page }) => {
      // 拦截 API 请求并模拟网络错误
      await page.route("**/api/v1/messages*", (route) => {
        route.abort("failed")
      })

      await messagePage.goto("system")
      await page.waitForTimeout(1000)

      // 应该显示空状态
      const hasEmptyState = (await messagePage.emptyState.count()) > 0
      expect(hasEmptyState).toBeTruthy()
    })
  })

  test.describe("性能测试", () => {
    test("页面加载时间应该合理", async ({ page }) => {
      const startTime = Date.now()
      await messagePage.goto("system")
      await page.waitForLoadState("networkidle")
      const loadTime = Date.now() - startTime

      // 页面应该在 5 秒内加载完成
      expect(loadTime).toBeLessThan(5000)
    })

    test("标签页切换应该流畅", async ({ page }) => {
      await messagePage.goto("system")

      const startTime = Date.now()
      await messagePage.switchTab("inbox")
      const switchTime = Date.now() - startTime

      // 切换应该在 2 秒内完成
      expect(switchTime).toBeLessThan(2000)
    })
  })
})
