import { test, expect } from "@playwright/test"

const BASE_URL = process.env.BASE_URL || "http://localhost:3000"

// 冒烟测试 - 测试所有主要页面是否能正常渲染
test.describe("Smoke Tests", () => {
  // 设置更长的超时时间（SSR 页面需要等待后端 API）
  test.use({ navigationTimeout: 60000, actionTimeout: 30000 })

  // 测试首页
  test("homepage should render without errors", async ({ page }) => {
    const errors: string[] = []
    page.on("pageerror", (err) => errors.push(err.message))

    await page.goto(BASE_URL, { timeout: 60000 })
    await page.waitForLoadState("networkidle")

    // 检查页面是否有 React 错误
    expect(errors).toHaveLength(0)

    // 检查页面标题
    await expect(page).toHaveTitle(/Go语言中文网/)

    // 检查主要元素存在
    await expect(page.locator("header")).toBeVisible()
    await expect(page.locator("main")).toBeVisible()
    await expect(page.locator("footer")).toBeVisible()
  })

  // 测试话题列表页
  test("/topics should render without errors", async ({ page }) => {
    const errors: string[] = []
    page.on("pageerror", (err) => errors.push(err.message))

    await page.goto(`${BASE_URL}/topics`, { timeout: 60000 })
    await page.waitForLoadState("networkidle")

    expect(errors).toHaveLength(0)
    await expect(page).toHaveTitle(/主题讨论/)
    await expect(page.locator("h1, h2").first()).toBeVisible()
  })

  // 测试文章列表页
  test("/articles should render without errors", async ({ page }) => {
    const errors: string[] = []
    page.on("pageerror", (err) => errors.push(err.message))

    await page.goto(`${BASE_URL}/articles`)
    await page.waitForLoadState("networkidle")

    expect(errors).toHaveLength(0)
    await expect(page).toHaveTitle(/文章/)
  })

  // 测试项目列表页
  test("/projects should render without errors", async ({ page }) => {
    const errors: string[] = []
    page.on("pageerror", (err) => errors.push(err.message))

    await page.goto(`${BASE_URL}/projects`)
    await page.waitForLoadState("networkidle")

    expect(errors).toHaveLength(0)
    await expect(page).toHaveTitle(/项目/)
  })

  // 测试资源列表页
  test("/resources should render without errors", async ({ page }) => {
    const errors: string[] = []
    page.on("pageerror", (err) => errors.push(err.message))

    await page.goto(`${BASE_URL}/resources`)
    await page.waitForLoadState("networkidle")

    expect(errors).toHaveLength(0)
    await expect(page).toHaveTitle(/资源/)
  })

  // 测试图书列表页
  test("/books should render without errors", async ({ page }) => {
    const errors: string[] = []
    page.on("pageerror", (err) => errors.push(err.message))

    await page.goto(`${BASE_URL}/books`)
    await page.waitForLoadState("networkidle")

    expect(errors).toHaveLength(0)
    await expect(page).toHaveTitle(/图书|书籍/)
  })

  // 测试晨读列表页
  test("/readings should render without errors", async ({ page }) => {
    const errors: string[] = []
    page.on("pageerror", (err) => errors.push(err.message))

    await page.goto(`${BASE_URL}/readings`)
    await page.waitForLoadState("networkidle")

    expect(errors).toHaveLength(0)
    await expect(page).toHaveTitle(/晨读/)
  })

  // 测试招聘列表页
  test("/jobs should render without errors", async ({ page }) => {
    const errors: string[] = []
    page.on("pageerror", (err) => errors.push(err.message))

    await page.goto(`${BASE_URL}/jobs`)
    await page.waitForLoadState("networkidle")

    expect(errors).toHaveLength(0)
    await expect(page).toHaveTitle(/工作|招聘/)
  })

  // 测试 Wiki 页
  test("/wiki should render without errors", async ({ page }) => {
    const errors: string[] = []
    page.on("pageerror", (err) => errors.push(err.message))

    await page.goto(`${BASE_URL}/wiki`)
    await page.waitForLoadState("networkidle")

    expect(errors).toHaveLength(0)
    await expect(page).toHaveTitle(/Wiki/)
  })

  // 测试搜索页
  test("/search should render without errors", async ({ page }) => {
    const errors: string[] = []
    page.on("pageerror", (err) => errors.push(err.message))

    await page.goto(`${BASE_URL}/search`)
    await page.waitForLoadState("networkidle")

    expect(errors).toHaveLength(0)
    await expect(page).toHaveTitle(/搜索/)
  })

  // 测试发布页
  test("/publish should render without errors", async ({ page }) => {
    const errors: string[] = []
    page.on("pageerror", (err) => errors.push(err.message))

    await page.goto(`${BASE_URL}/publish`)
    await page.waitForLoadState("networkidle")

    expect(errors).toHaveLength(0)
    // 发布页应该显示表单或重定向到登录页
    const url = page.url()
    expect(url).toMatch(/\/(publish|login)/)
  })

  // 测试登录页
  test("/account/login should render without errors", async ({ page }) => {
    const errors: string[] = []
    page.on("pageerror", (err) => errors.push(err.message))

    await page.goto(`${BASE_URL}/account/login`)
    await page.waitForLoadState("networkidle")

    expect(errors).toHaveLength(0)
    await expect(page).toHaveTitle(/登录/)
    // 检查登录表单元素
    await expect(page.locator('input[type="text"], input[name="username"]')).toBeVisible()
    await expect(page.locator('input[type="password"]')).toBeVisible()
  })

  // 测试注册页
  test("/account/register should render without errors", async ({ page }) => {
    const errors: string[] = []
    page.on("pageerror", (err) => errors.push(err.message))

    await page.goto(`${BASE_URL}/account/register`)
    await page.waitForLoadState("networkidle")

    expect(errors).toHaveLength(0)
    await expect(page).toHaveTitle(/注册/)
    // 检查注册表单元素
    await expect(page.locator('input[name="username"]')).toBeVisible()
    await expect(page.locator('input[type="password"]')).toBeVisible()
  })

  // 测试导航链接
  test("navigation links should work", async ({ page }) => {
    await page.goto(BASE_URL)
    await page.waitForLoadState("networkidle")

    // 检查主导航链接存在
    const navLinks = ["/topics", "/articles", "/projects", "/resources", "/books", "/jobs"]

    for (const href of navLinks) {
      const link = page.locator(`a[href="${href}"]`).first()
      await expect(link).toBeVisible()
    }
  })

  // 测试搜索功能
  test("search should work", async ({ page }) => {
    await page.goto(BASE_URL)
    await page.waitForLoadState("networkidle")

    // 找到搜索框
    const searchInput = page.locator('input[placeholder*="搜索"]')
    await expect(searchInput).toBeVisible()

    // 输入搜索词
    await searchInput.fill("golang")
    await searchInput.press("Enter")

    // 等待导航到搜索页
    await page.waitForURL(/\/search/)
    await expect(page).toHaveTitle(/搜索/)
  })
})
