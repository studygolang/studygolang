import { test, expect } from "@playwright/test"
import { mockLogin } from "../../fixtures/auth"

test.describe("资源和图书发布功能", () => {
  test.beforeEach(async ({ page }) => {
    await mockLogin(page)
    await page.goto("/publish")
  })

  // 辅助函数：点击发布页面中的内容类型按钮（避免与导航栏冲突）
  async function selectContentType(page: import("@playwright/test").Page, type: string) {
    // 内容类型选择器在页面右上角的标签页样式中，使用 role button 精确定位
    await page.getByRole("button", { name: type, exact: true }).click()
  }

  test.describe("资源发布", () => {
    test.beforeEach(async ({ page }) => {
      // 切换到资源类型
      await selectContentType(page, "资源")
    })

    test("应该显示资源类型选项", async ({ page }) => {
      await expect(page.getByText("只是链接")).toBeVisible()
      await expect(page.getByText("包括内容")).toBeVisible()
    })

    test("应该显示资源表单字段", async ({ page }) => {
      // 检查标题输入框
      const titleInput = page.locator('input[placeholder*="资源标题"]')
      await expect(titleInput).toBeVisible()

      // 检查资源类型选择器
      await expect(page.getByText("资源类型")).toBeVisible()
    })

    test("链接类型应显示链接地址输入框", async ({ page }) => {
      // 默认选中"只是链接"
      const linkRadio = page.locator('input[value="link"]')
      await expect(linkRadio).toBeChecked()

      // 应该显示资源链接输入框
      const urlInput = page.locator('input[placeholder*="example.com/resource"]')
      await expect(urlInput).toBeVisible()
    })

    test("内容类型应显示编辑器而不是链接输入框", async ({ page }) => {
      // 选择"包括内容"
      await page.getByText("包括内容").click()

      // 等待编辑器加载
      await page.waitForTimeout(2000)

      // 应该显示编辑器
      const editor = page.locator(".bn-editor")
      await expect(editor).toBeVisible({ timeout: 10000 })

      // 不应该显示链接输入框
      const urlInput = page.locator('input[placeholder*="example.com/resource"]')
      await expect(urlInput).not.toBeVisible()
    })

    test("链接类型必填验证 - 缺少URL", async ({ page }) => {
      // 填写标题但不填写链接
      const titleInput = page.locator('input[placeholder*="资源标题"]')
      await titleInput.fill("测试资源")

      // 点击发布（使用 main 区域内的按钮，避免匹配导航栏）
      const publishButton = page.locator("main").getByRole("button", { name: /^发布$/ })
      await publishButton.click()

      // 应该显示错误提示
      await expect(page.getByText(/请填写资源链接/)).toBeVisible()
    })

    test("内容类型必填验证 - 缺少内容", async ({ page }) => {
      // 选择"包括内容"
      await page.getByText("包括内容").click()

      // 填写标题
      const titleInput = page.locator('input[placeholder*="资源标题"]')
      await titleInput.fill("测试资源")

      // 点击发布（使用 main 区域内的按钮，避免匹配导航栏）
      const publishButton = page.locator("main").getByRole("button", { name: /^发布$/ })
      await publishButton.click()

      // 应该显示错误提示
      await expect(page.getByText(/请填写资源内容/)).toBeVisible()
    })

    test("应该能够填写完整的链接类型资源表单", async ({ page }) => {
      // 填写标题
      const titleInput = page.locator('input[placeholder*="资源标题"]')
      await titleInput.fill("Go语言学习资源")

      // 填写链接
      const urlInput = page.locator('input[placeholder*="example.com/resource"]')
      await urlInput.fill("https://golang.org/doc/")

      // 添加标签
      const tagInput = page.locator("main").locator('input[id="tags"]')
      await tagInput.fill("Go")
      await tagInput.press("Enter")

      // 验证表单数据
      await expect(titleInput).toHaveValue("Go语言学习资源")
      await expect(urlInput).toHaveValue("https://golang.org/doc/")
      // 验证标签已添加（通过 Badge 文本 "Go ×" 验证）
      const addedTag = page.locator("main").locator("span", { hasText: /^Go\s*×$/ })
      await expect(addedTag).toBeVisible()
    })
  })

  test.describe("图书发布", () => {
    test.beforeEach(async ({ page }) => {
      // 切换到图书类型
      await selectContentType(page, "图书")
    })

    test("应该显示图书类型选项", async ({ page }) => {
      // 验证图书选项被选中（使用 button role）
      const bookButton = page.getByRole("button", { name: "图书", exact: true })
      await expect(bookButton).toBeVisible()
      // 验证按钮处于选中状态（通过类名判断）
      await expect(bookButton).toHaveClass(/bg-background/)
    })

    test("应该显示所有图书表单字段", async ({ page }) => {
      // 必填字段：书名、作者
      const titleInput = page.locator('input[placeholder*="图书名称"]')
      await expect(titleInput).toBeVisible()

      const authorInput = page.locator('input[placeholder*="作者姓名"]')
      await expect(authorInput).toBeVisible()

      // 可选字段：译者
      const translatorInput = page.locator('input[placeholder*="译者姓名"]')
      await expect(translatorInput).toBeVisible()

      // 可选字段：封面
      const coverInput = page.locator('input[placeholder*="封面图片URL"]')
      await expect(coverInput).toBeVisible()

      // 可选字段：出版日期
      const pubDateInput = page.locator('input[placeholder*="2024-01"]')
      await expect(pubDateInput).toBeVisible()

      // 语言选择器
      await expect(page.locator("label").filter({ hasText: "语言" })).toBeVisible()

      // 是否免费复选框（使用更精确的选择器）
      const freeLabel = page.locator("label").filter({ hasText: "免费资源" })
      await expect(freeLabel).toBeVisible()

      // 在线阅读地址
      const onlineUrlInput = page.locator('input[placeholder*="在线阅读地址"]')
      await expect(onlineUrlInput).toBeVisible()

      // 下载地址
      const downloadUrlInput = page.locator('input[placeholder*="电子版下载地址"]')
      await expect(downloadUrlInput).toBeVisible()

      // 购买地址
      const buyUrlInput = page.locator('input[placeholder*="购买链接"]')
      await expect(buyUrlInput).toBeVisible()

      // 价格
      const priceInput = page.locator('input[placeholder*="99.00"]')
      await expect(priceInput).toBeVisible()
    })

    test("语言选项应正确显示", async ({ page }) => {
      // 点击语言选择器
      await page.locator('button:has-text("中文")').click()

      // 验证选项
      await expect(page.getByRole("option", { name: "中文" })).toBeVisible()
      await expect(page.getByRole("option", { name: "英文" })).toBeVisible()
      await expect(page.getByRole("option", { name: "中英双语" })).toBeVisible()
      await expect(page.getByRole("option", { name: "其他" })).toBeVisible()
    })

    test("是否免费复选框应可切换", async ({ page }) => {
      // 图书页面的复选框
      const freeCheckbox = page.locator("main").locator('input[type="checkbox"]')

      // 默认未选中
      await expect(freeCheckbox).not.toBeChecked()

      // 点击选中
      const freeLabel = page.locator("main").locator("label").filter({ hasText: "免费资源" })
      await freeLabel.click()
      await expect(freeCheckbox).toBeChecked()
    })

    test("图书必填验证 - 缺少作者", async ({ page }) => {
      // 只填写书名
      const titleInput = page.locator('input[placeholder*="图书名称"]')
      await titleInput.fill("Go语言圣经")

      // 点击发布（使用 main 区域内的按钮，避免匹配导航栏）
      const publishButton = page.locator("main").getByRole("button", { name: /^发布$/ })
      await publishButton.click()

      // 应该显示错误提示
      await expect(page.getByText(/请填写作者/)).toBeVisible()
    })

    test("应该能够填写完整的图书表单", async ({ page }) => {
      // 填写书名
      const titleInput = page.locator('input[placeholder*="图书名称"]')
      await titleInput.fill("Go语言圣经")

      // 填写作者
      const authorInput = page.locator('input[placeholder*="作者姓名"]')
      await authorInput.fill("Alan A. A. Donovan")

      // 填写译者
      const translatorInput = page.locator('input[placeholder*="译者姓名"]')
      await translatorInput.fill("李笑")

      // 填写封面
      const coverInput = page.locator('input[placeholder*="封面图片URL"]')
      await coverInput.fill("https://example.com/cover.jpg")

      // 填写出版日期
      const pubDateInput = page.locator('input[placeholder*="2024-01"]')
      await pubDateInput.fill("2016-01")

      // 选择语言
      await page.locator('button:has-text("中文")').click()
      await page.getByRole("option", { name: "中文" }).click()

      // 勾选免费
      const freeLabel = page.locator("main").locator("label").filter({ hasText: "免费资源" })
      await freeLabel.click()

      // 填写在线阅读地址
      const onlineUrlInput = page.locator('input[placeholder*="在线阅读地址"]')
      await onlineUrlInput.fill("https://example.com/read")

      // 填写下载地址
      const downloadUrlInput = page.locator('input[placeholder*="电子版下载地址"]')
      await downloadUrlInput.fill("https://example.com/download")

      // 填写购买地址
      const buyUrlInput = page.locator('input[placeholder*="购买链接"]')
      await buyUrlInput.fill("https://example.com/buy")

      // 填写价格
      const priceInput = page.locator('input[placeholder*="99.00"]')
      await priceInput.fill("79.00")

      // 添加标签（使用 main 区域内的标签输入框）
      const tagInput = page.locator("main").locator('input[id="tags"]')
      await tagInput.fill("Go")
      await tagInput.press("Enter")

      // 验证表单数据
      await expect(titleInput).toHaveValue("Go语言圣经")
      await expect(authorInput).toHaveValue("Alan A. A. Donovan")
      await expect(translatorInput).toHaveValue("李笑")
      // 验证标签已添加（通过 Badge 文本 "Go ×" 验证）
      const addedTag = page.locator("main").locator("span", { hasText: /^Go\s*×$/ })
      await expect(addedTag).toBeVisible()
    })
  })

  test.describe("内容类型切换", () => {
    test("应该正确显示所有五种内容类型选项", async ({ page }) => {
      // 使用 button role 精确定位内容类型按钮
      await expect(page.getByRole("button", { name: "主题", exact: true })).toBeVisible()
      await expect(page.getByRole("button", { name: "文章", exact: true })).toBeVisible()
      await expect(page.getByRole("button", { name: "项目", exact: true })).toBeVisible()
      await expect(page.getByRole("button", { name: "资源", exact: true })).toBeVisible()
      await expect(page.getByRole("button", { name: "图书", exact: true })).toBeVisible()
    })

    test("切换到资源类型应显示正确的字段", async ({ page }) => {
      await selectContentType(page, "资源")

      // 应该显示资源类型选择
      await expect(page.getByText("只是链接")).toBeVisible()
      await expect(page.getByText("包括内容")).toBeVisible()

      // 不应该显示节点选择器
      await expect(page.getByText("选择节点")).not.toBeVisible()
    })

    test("切换到图书类型应显示正确的字段", async ({ page }) => {
      await selectContentType(page, "图书")

      // 应该显示图书特有字段
      await expect(page.getByText("作者")).toBeVisible()
      await expect(page.getByText("译者")).toBeVisible()

      // 不应该显示节点选择器
      await expect(page.getByText("选择节点")).not.toBeVisible()
    })
  })
})
