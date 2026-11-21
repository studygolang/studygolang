import { test, expect } from "@playwright/test"
import { mockLogin } from "../../fixtures/auth"

test.describe("发布页面 /publish", () => {
  test.beforeEach(async ({ page }) => {
    await mockLogin(page)
    await page.route("**/api/v1/nodes", async (route) => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          code: 0,
          data: [
            { nid: 1, name: "社区", pid: 0 },
            { nid: 2, name: "Go 语言", pid: 1 },
          ],
        }),
      })
    })
    await page.route("**/api/v1/captcha/required", async (route) => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ code: 0, data: { required: false } }),
      })
    })
    await page.route("**/api/v1/resources/categories", async (route) => {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          code: 0,
          data: { categories: [{ catid: 7, name: "Go 语言" }] },
        }),
      })
    })
    await page.goto("/publish")
  })

  test("应该正确加载发布页面", async ({ page }) => {
    await expect(page).toHaveTitle(/发布/)
    await expect(page.getByRole("heading", { name: "发布内容", exact: true })).toBeVisible()
  })

  test("移动端不应出现横向溢出", async ({ page }) => {
    await page.setViewportSize({ width: 390, height: 844 })
    await expect(page.locator(".bn-editor")).toBeVisible({ timeout: 10000 })
    const { scrollWidth, viewportWidth } = await page.evaluate(() => ({
      scrollWidth: document.documentElement.scrollWidth,
      viewportWidth: window.innerWidth,
    }))
    expect(scrollWidth).toBeLessThanOrEqual(viewportWidth + 1)
  })

  test("应该显示三种内容类型选项", async ({ page }) => {
    await expect(page.getByRole("button", { name: "主题", exact: true })).toBeVisible()
    await expect(page.getByRole("button", { name: "文章", exact: true })).toBeVisible()
    await expect(page.getByRole("button", { name: "项目", exact: true })).toBeVisible()
  })

  test("应该能够切换内容类型", async ({ page }) => {
    const topicButton = page.getByRole("button", { name: "主题", exact: true })
    await expect(topicButton).toHaveClass(/bg-background/)

    const articleButton = page.getByRole("button", { name: "文章", exact: true })
    await articleButton.click()
    await expect(articleButton).toHaveClass(/bg-background/)
  })

  test("应该加载 BlockNote 编辑器", async ({ page }) => {
    await page.waitForTimeout(2000) // 等待动态导入
    const editor = page.locator(".bn-editor")
    await expect(editor).toBeVisible({ timeout: 10000 })
  })

  test("应该显示标题输入框", async ({ page }) => {
    const titleInput = page.getByRole("textbox", { name: /标题/ })
    await expect(titleInput).toBeVisible()
  })

  test("应该显示节点选择器（主题类型）", async ({ page }) => {
    const nodeSelector = page.getByText("选择节点")
    await expect(nodeSelector).toBeVisible()
  })

  test("应该显示标签输入框", async ({ page }) => {
    const tagInput = page.getByRole("textbox", { name: /标签/ })
    await expect(tagInput).toBeVisible()
  })

  test("应该能够输入标题", async ({ page }) => {
    const titleInput = page.getByRole("textbox", { name: /标题/ })
    await titleInput.fill("测试标题")
    await expect(titleInput).toHaveValue("测试标题")
  })

  test("应该能够添加标签", async ({ page }) => {
    const tagInput = page.getByRole("textbox", { name: /标签/ })
    await tagInput.fill("测试标签")
    await tagInput.press("Enter")
    await expect(page.getByText("测试标签")).toBeVisible()
  })

  test("应该能够删除标签", async ({ page }) => {
    const tagInput = page.getByRole("textbox", { name: /标签/ })
    await tagInput.fill("测试标签")
    await tagInput.press("Enter")

    const removeButton = page.locator('button[aria-label*="删除"]').first()
    await removeButton.click()
    await expect(page.getByText("测试标签")).not.toBeVisible()
  })

  test("应该显示发布和保存草稿按钮", async ({ page }) => {
    await expect(page.getByRole("main").locator("[data-publish-btn]")).toBeVisible()
    await expect(page.getByRole("button", { name: /保存草稿/ })).toBeVisible()
  })

  test("应该显示取消按钮", async ({ page }) => {
    await expect(page.getByRole("button", { name: "取消" })).toBeVisible()
  })

  test("标题为空时应该显示验证错误", async ({ page }) => {
    const publishButton = page.getByRole("main").locator("[data-publish-btn]")
    await publishButton.click()
    await expect(page.getByText("请填写标题", { exact: true })).toBeVisible()
  })

  test("内容为空时应该显示验证错误", async ({ page }) => {
    await page.getByRole("button", { name: "文章", exact: true }).click()
    const titleInput = page.getByRole("textbox", { name: /标题/ })
    await titleInput.fill("测试标题")

    const publishButton = page.getByRole("main").locator("[data-publish-btn]")
    await publishButton.click()
    await expect(page.getByText("请填写内容", { exact: true })).toBeVisible()
  })

  test("主题类型未选择节点时应该显示验证错误", async ({ page }) => {
    const titleInput = page.getByRole("textbox", { name: /标题/ })
    await titleInput.fill("测试标题")

    // 等待编辑器加载
    await page.waitForTimeout(2000)

    const publishButton = page.getByRole("main").locator("[data-publish-btn]")
    await publishButton.click()
    await expect(page.getByText(/请选择节点/)).toBeVisible()
  })

  test("资源发布应该选择真实分类", async ({ page }) => {
    await page.getByRole("button", { name: "资源", exact: true }).click()
    await page.locator("#title").fill("测试资源")
    await page.locator("#resourceUrl").fill("https://example.com/resource")
    await page.getByRole("main").locator("[data-publish-btn]").click()
    await expect(page.getByText("请选择资源分类")).toBeVisible()

    await page.locator("#resourceCatid").click()
    await page.getByRole("option", { name: "Go 语言" }).click()
    await expect(page.locator("#resourceCatid")).toContainText("Go 语言")
  })

  test("项目发布应该包含完整资料字段", async ({ page }) => {
    await page.getByRole("button", { name: "项目", exact: true }).click()
    for (const id of ["uri", "download", "logo", "author"]) {
      await expect(page.locator(`#${id}`)).toBeVisible()
    }
    await page.locator("#uri").fill("studygolang")
    await page.locator("#download").fill("https://example.com/download")
    await page.locator("#logo").fill("https://example.com/logo.png")
    await page.locator("#author").fill("StudyGolang")
    await expect(page.locator("#author")).toHaveValue("StudyGolang")
  })

  test("恢复草稿后应该同步标题和编辑器内容", async ({ page }) => {
    await page.addInitScript(() => {
      localStorage.setItem("draft:topic:new", JSON.stringify({
        title: "草稿标题",
        content: "**草稿正文**",
        savedAt: Date.now(),
      }))
    })
    await page.reload()
    await page.getByRole("button", { name: "恢复草稿" }).click()
    await expect(page.locator("#title")).toHaveValue("草稿标题")
    await expect(page.locator(".bn-editor")).toContainText("草稿正文", { timeout: 10000 })
  })

  test("点击取消应该返回上一页", async ({ page }) => {
    await page.goto("/topics")
    await page.goto("/publish")
    const cancelButton = page.getByRole("button", { name: "取消" })
    await cancelButton.click()
    await expect(page).toHaveURL(/\/topics$/)
  })
})
