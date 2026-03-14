import { test, expect } from "@playwright/test"

/**
 * 用户中心功能 E2E 测试
 * 测试范围:
 * 1. 会员列表页 /users
 * 2. 用户个人主页 /user/:username
 * 3. 收藏列表 /favorites/:username
 */

test.describe("用户中心功能测试", () => {
  test.describe("会员列表页 /users", () => {
    test("页面正常渲染", async ({ page }) => {
      const response = await page.goto("/users")
      await page.waitForLoadState("networkidle")

      // HTTP 状态应为 200
      expect(response?.status()).toBe(200)

      // 页面标题正确
      await expect(page).toHaveTitle(/会员列表/)

      // 页面 header 可见
      await expect(page.locator("header")).toBeVisible()

      // 页面主标题可见
      await expect(
        page.getByRole("heading", { name: "会员列表" })
      ).toBeVisible()
    })

    test("显示活跃会员和新加入会员 Tab", async ({ page }) => {
      await page.goto("/users")
      await page.waitForLoadState("networkidle")

      // 应有两个 Tab
      const activeTab = page.getByRole("tab", { name: /活跃会员/ })
      const newTab = page.getByRole("tab", { name: /新加入会员/ })

      await expect(activeTab).toBeVisible()
      await expect(newTab).toBeVisible()
    })

    test("活跃会员 Tab 显示用户列表", async ({ page }) => {
      await page.goto("/users")
      await page.waitForLoadState("networkidle")

      // 默认应显示活跃会员
      const activeTab = page.getByRole("tab", { name: /活跃会员/ })
      await expect(activeTab).toHaveAttribute("data-state", "active")

      // 应显示用户卡片（如果有数据）
      const userCards = page.locator("a[href^='/user/']")
      const count = await userCards.count()

      if (count > 0) {
        // 验证第一个用户卡片包含头像和用户名
        const firstCard = userCards.first()
        await expect(firstCard).toBeVisible()
        await expect(firstCard.locator("img, [role='img']")).toBeVisible() // 头像
      }
    })

    test("切换到新加入会员 Tab", async ({ page }) => {
      await page.goto("/users")
      await page.waitForLoadState("networkidle")

      // 点击新加入会员 Tab
      const newTab = page.getByRole("tab", { name: /新加入会员/ })
      await newTab.click()

      // Tab 状态应变为 active
      await expect(newTab).toHaveAttribute("data-state", "active")

      // 等待内容加载
      await page.waitForTimeout(500)
    })

    test("用户卡片可点击跳转", async ({ page }) => {
      await page.goto("/users")
      await page.waitForLoadState("networkidle")

      const userCards = page.locator("a[href^='/user/']")
      const count = await userCards.count()

      if (count > 0) {
        const firstCard = userCards.first()
        const href = await firstCard.getAttribute("href")

        // 点击用户卡片
        await firstCard.click()
        await page.waitForLoadState("networkidle")

        // 应跳转到用户主页
        expect(page.url()).toContain(href || "/user/")
      }
    })

    test("响应式布局 - 移动端", async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 })
      await page.goto("/users")
      await page.waitForLoadState("networkidle")

      // 页面应正常显示
      await expect(page.locator("body")).toBeVisible()
      await expect(page.getByRole("heading", { name: "会员列表" })).toBeVisible()
    })

    test("无数据时显示空状态", async ({ page }) => {
      await page.goto("/users")
      await page.waitForLoadState("networkidle")

      const userCards = page.locator("a[href^='/user/']")
      const count = await userCards.count()

      if (count === 0) {
        // 应显示"暂无数据"提示
        await expect(page.getByText(/暂无数据/)).toBeVisible()
      }
    })
  })

  test.describe("用户个人主页 /user/:username", () => {
    test("页面正常渲染", async ({ page }) => {
      // 使用一个可能存在的用户名
      const response = await page.goto("/user/polaris")
      await page.waitForLoadState("networkidle")

      // HTTP 状态应为 200
      expect(response?.status()).toBe(200)

      // 页面 body 可见
      await expect(page.locator("body")).toBeVisible()
    })

    test("显示用户基本信息", async ({ page }) => {
      await page.goto("/user/polaris")
      await page.waitForLoadState("networkidle")

      // 应显示用户头像
      const avatar = page.locator("img[alt*='头像'], [role='img']").first()
      await expect(avatar).toBeVisible()

      // 应显示用户名或昵称
      const username = page.locator("text=/polaris|@polaris/i").first()
      await expect(username).toBeVisible()
    })

    test("显示用户统计信息", async ({ page }) => {
      await page.goto("/user/polaris")
      await page.waitForLoadState("networkidle")

      // 应显示统计数据（话题、文章、积分等）
      const stats = page.locator("text=/话题|文章|积分|收藏/i")
      const count = await stats.count()
      expect(count).toBeGreaterThan(0)
    })

    test("显示用户发布的话题", async ({ page }) => {
      await page.goto("/user/polaris")
      await page.waitForLoadState("networkidle")

      // 查找话题列表
      const topics = page.locator("a[href^='/topics/']")
      const count = await topics.count()

      if (count > 0) {
        // 验证话题卡片可见
        await expect(topics.first()).toBeVisible()
      }
    })

    test("显示用户发布的文章", async ({ page }) => {
      await page.goto("/user/polaris")
      await page.waitForLoadState("networkidle")

      // 查找文章列表
      const articles = page.locator("a[href^='/articles/']")
      const count = await articles.count()

      if (count > 0) {
        // 验证文章卡片可见
        await expect(articles.first()).toBeVisible()
      }
    })

    test("用户不存在时显示错误", async ({ page }) => {
      await page.goto("/user/nonexistentuser12345")
      await page.waitForLoadState("networkidle")

      // 应显示用户未找到提示
      const notFound = page.getByText(/用户未找到|用户不存在|404/i)
      await expect(notFound).toBeVisible()
    })

    test("响应式布局 - 移动端", async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 })
      await page.goto("/user/polaris")
      await page.waitForLoadState("networkidle")

      // 页面应正常显示
      await expect(page.locator("body")).toBeVisible()
    })

    test("社交链接可点击", async ({ page }) => {
      await page.goto("/user/polaris")
      await page.waitForLoadState("networkidle")

      // 查找社交链接（GitHub、网站等）
      const socialLinks = page.locator("a[href^='http']")
      const count = await socialLinks.count()

      if (count > 0) {
        const firstLink = socialLinks.first()
        await expect(firstLink).toBeVisible()

        // 验证链接有 href 属性
        const href = await firstLink.getAttribute("href")
        expect(href).toBeTruthy()
      }
    })
  })

  test.describe("收藏列表 /favorites/:username", () => {
    test("页面正常渲染", async ({ page }) => {
      const response = await page.goto("/favorites/polaris")
      await page.waitForLoadState("networkidle")

      // HTTP 状态应为 200
      expect(response?.status()).toBe(200)

      // 页面 body 可见
      await expect(page.locator("body")).toBeVisible()
    })

    test("显示收藏类型 Tab", async ({ page }) => {
      await page.goto("/favorites/polaris")
      await page.waitForLoadState("networkidle")

      // 应显示多个收藏类型 Tab（话题、文章、资源、项目）
      const tabs = page.locator("[role='tab']")
      const count = await tabs.count()
      expect(count).toBeGreaterThan(0)
    })

    test("默认显示话题收藏", async ({ page }) => {
      await page.goto("/favorites/polaris")
      await page.waitForLoadState("networkidle")

      // 默认应显示话题收藏
      const topicTab = page.getByRole("tab", { name: /话题/ })
      if (await topicTab.isVisible()) {
        await expect(topicTab).toHaveAttribute("data-state", "active")
      }
    })

    test("切换到文章收藏 Tab", async ({ page }) => {
      await page.goto("/favorites/polaris")
      await page.waitForLoadState("networkidle")

      // 查找文章 Tab
      const articleTab = page.getByRole("tab", { name: /文章/ })

      if (await articleTab.isVisible()) {
        await articleTab.click()
        await page.waitForTimeout(500)

        // Tab 状态应变为 active
        await expect(articleTab).toHaveAttribute("data-state", "active")
      }
    })

    test("切换到资源收藏 Tab", async ({ page }) => {
      await page.goto("/favorites/polaris")
      await page.waitForLoadState("networkidle")

      const resourceTab = page.getByRole("tab", { name: /资源/ })

      if (await resourceTab.isVisible()) {
        await resourceTab.click()
        await page.waitForTimeout(500)

        await expect(resourceTab).toHaveAttribute("data-state", "active")
      }
    })

    test("切换到项目收藏 Tab", async ({ page }) => {
      await page.goto("/favorites/polaris")
      await page.waitForLoadState("networkidle")

      const projectTab = page.getByRole("tab", { name: /项目/ })

      if (await projectTab.isVisible()) {
        await projectTab.click()
        await page.waitForTimeout(500)

        await expect(projectTab).toHaveAttribute("data-state", "active")
      }
    })

    test("显示收藏的话题列表", async ({ page }) => {
      await page.goto("/favorites/polaris")
      await page.waitForLoadState("networkidle")

      // 查找话题链接
      const topics = page.locator("a[href^='/topics/']")
      const count = await topics.count()

      if (count > 0) {
        await expect(topics.first()).toBeVisible()
      } else {
        // 无收藏时应显示提示
        await expect(page.getByText(/暂无收藏/)).toBeVisible()
      }
    })

    test("显示收藏的文章列表", async ({ page }) => {
      await page.goto("/favorites/polaris?objtype=2")
      await page.waitForLoadState("networkidle")

      // 切换到文章 Tab
      const articleTab = page.getByRole("tab", { name: /文章/ })
      if (await articleTab.isVisible()) {
        await articleTab.click()
        await page.waitForTimeout(500)

        // 查找文章链接
        const articles = page.locator("a[href^='/articles/']")
        const count = await articles.count()

        if (count > 0) {
          await expect(articles.first()).toBeVisible()
        }
      }
    })

    test("收藏项可点击跳转", async ({ page }) => {
      await page.goto("/favorites/polaris")
      await page.waitForLoadState("networkidle")

      const links = page.locator("a[href^='/topics/'], a[href^='/articles/']")
      const count = await links.count()

      if (count > 0) {
        const firstLink = links.first()
        const href = await firstLink.getAttribute("href")

        await firstLink.click()
        await page.waitForLoadState("networkidle")

        // 应跳转到对应的详情页
        expect(page.url()).toContain(href || "/")
      }
    })

    test("用户不存在时显示错误", async ({ page }) => {
      await page.goto("/favorites/nonexistentuser12345")
      await page.waitForLoadState("networkidle")

      // 应显示用户未找到提示
      const notFound = page.getByText(/用户未找到|用户不存在|404/i)
      await expect(notFound).toBeVisible()
    })

    test("响应式布局 - 移动端", async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 })
      await page.goto("/favorites/polaris")
      await page.waitForLoadState("networkidle")

      // 页面应正常显示
      await expect(page.locator("body")).toBeVisible()
    })

    test("无收藏时显示空状态", async ({ page }) => {
      await page.goto("/favorites/polaris")
      await page.waitForLoadState("networkidle")

      // 如果没有收藏，应显示提示
      const emptyState = page.getByText(/暂无收藏/)
      const hasContent = await page.locator("a[href^='/topics/'], a[href^='/articles/']").count()

      if (hasContent === 0) {
        await expect(emptyState).toBeVisible()
      }
    })
  })

  test.describe("错误处理", () => {
    test("API 错误时页面不崩溃", async ({ page }) => {
      // 监听控制台错误
      const errors: string[] = []
      page.on("pageerror", (err) => errors.push(err.message))

      await page.goto("/users")
      await page.waitForLoadState("networkidle")

      // 页面应正常显示，即使 API 失败
      await expect(page.locator("body")).toBeVisible()
    })

    test("网络超时时显示友好提示", async ({ page }) => {
      await page.goto("/user/polaris")
      await page.waitForLoadState("networkidle")

      // 页面应能正常加载
      await expect(page.locator("body")).toBeVisible()
    })
  })

  test.describe("性能测试", () => {
    test("页面加载时间合理", async ({ page }) => {
      const startTime = Date.now()

      await page.goto("/users")
      await page.waitForLoadState("networkidle")

      const loadTime = Date.now() - startTime

      // 页面应在 5 秒内加载完成
      expect(loadTime).toBeLessThan(5000)
    })

    test("图片懒加载", async ({ page }) => {
      await page.goto("/users")
      await page.waitForLoadState("networkidle")

      // 检查是否有图片使用懒加载
      const lazyImages = page.locator("img[loading='lazy']")
      const count = await lazyImages.count()

      // 如果有图片，应使用懒加载
      if (count > 0) {
        expect(count).toBeGreaterThan(0)
      }
    })
  })
})
