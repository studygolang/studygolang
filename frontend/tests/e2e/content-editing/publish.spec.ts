import { test, expect } from "@playwright/test"
import { mockLogin } from "../../fixtures/auth"

test.describe("发布页面 /publish", () => {
  test.beforeEach(async ({ page }) => {
    await mockLogin(page)
    await page.goto("/publish")
  })

  test("应该正确加载发布页面", async ({ page }) => {
    await expect(page).toHaveTitle(/发布/)
    await expect(page.locator("h1")).toContainText("发布内容")
  })

  test("应该显示三种内容类型选项", async ({ page }) => {
    await expect(page.getByText("主题")).toBeVisible()
    await expect(page.getByText("文章")).toBeVisible()
    await expect(page.getByText("项目")).toBeVisible()
  })

  test("应该能够切换内容类型", async ({ page }) => {
    // 默认选中主题
    const topicCard = page.locator('[data-type="topic"]')
    await expect(topicCard).toHaveClass(/border-primary/)

    // 切换到文章
    await page.getByText("文章").click()
    const articleCard = page.locator('[data-type="article"]')
    await expect(articleCard).toHaveClass(/border-primary/)
  })

  test("应该加载 BlockNote 编辑器", async ({ page }) => {
    await page.waitForTimeout(2000) // 等待动态导入
    const editor = page.locator(".bn-editor")
    await expect(editor).toBeVisible({ timeout: 10000 })
  })

  test("应该显示标题输入框", async ({ page }) => {
    const titleInput = page.locator('input[placeholder*="标题"]')
    await expect(titleInput).toBeVisible()
  })

  test("应该显示节点选择器（主题类型）", async ({ page }) => {
    const nodeSelector = page.getByText("选择节点")
    await expect(nodeSelector).toBeVisible()
  })

  test("应该显示标签输入框", async ({ page }) => {
    const tagInput = page.locator('input[placeholder*="标签"]')
    await expect(tagInput).toBeVisible()
  })

  test("应该能够输入标题", async ({ page }) => {
    const titleInput = page.locator('input[placeholder*="标题"]')
    await titleInput.fill("测试标题")
    await expect(titleInput).toHaveValue("测试标题")
  })

  test("应该能够添加标签", async ({ page }) => {
    const tagInput = page.locator('input[placeholder*="标签"]')
    await tagInput.fill("测试标签")
    await tagInput.press("Enter")
    await expect(page.getByText("测试标签")).toBeVisible()
  })

  test("应该能够删除标签", async ({ page }) => {
    const tagInput = page.locator('input[placeholder*="标签"]')
    await tagInput.fill("测试标签")
    await tagInput.press("Enter")

    const removeButton = page.locator('button[aria-label*="删除"]').first()
    await removeButton.click()
    await expect(page.getByText("测试标签")).not.toBeVisible()
  })

  test("应该显示发布和保存草稿按钮", async ({ page }) => {
    await expect(page.getByRole("button", { name: /发布/ })).toBeVisible()
    await expect(page.getByRole("button", { name: /保存草稿/ })).toBeVisible()
  })

  test("应该显示取消按钮", async ({ page }) => {
    await expect(page.getByRole("button", { name: "取消" })).toBeVisible()
  })

  test("标题为空时应该显示验证错误", async ({ page }) => {
    const publishButton = page.getByRole("button", { name: /发布/ })
    await publishButton.click()
    await expect(page.getByText(/请输入标题/)).toBeVisible()
  })

  test("内容为空时应该显示验证错误", async ({ page }) => {
    const titleInput = page.locator('input[placeholder*="标题"]')
    await titleInput.fill("测试标题")

    const publishButton = page.getByRole("button", { name: /发布/ })
    await publishButton.click()
    await expect(page.getByText(/请输入内容/)).toBeVisible()
  })

  test("主题类型未选择节点时应该显示验证错误", async ({ page }) => {
    const titleInput = page.locator('input[placeholder*="标题"]')
    await titleInput.fill("测试标题")

    // 等待编辑器加载
    await page.waitForTimeout(2000)

    const publishButton = page.getByRole("button", { name: /发布/ })
    await publishButton.click()
    await expect(page.getByText(/请选择节点/)).toBeVisible()
  })

  test("点击取消应该返回上一页", async ({ page }) => {
    const cancelButton = page.getByRole("button", { name: "取消" })
    await cancelButton.click()
    // 应该导航到首页或上一页
    await expect(page).toHaveURL(/\/$|\/topics|\/articles/)
  })
})
