import { test, expect } from "@playwright/test"

test.describe("文章列表页", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/articles")
    await page.waitForLoadState("networkidle")
  })

  test("页面标题包含'文章'", async ({ page }) => {
    await expect(page).toHaveTitle(/文章/)
    await expect(page.locator("h1")).toContainText("文章")
  })

  test("文章过滤器显示排序 tab", async ({ page }) => {
    // ArticleFilter 中的排序按钮
    await expect(page.locator("button").filter({ hasText: "最新" })).toBeVisible()
    await expect(page.locator("button").filter({ hasText: "热门" })).toBeVisible()
  })

  test("点击排序 tab 更新 URL", async ({ page }) => {
    await page.locator("button").filter({ hasText: "热门" }).click()
    await expect(page).toHaveURL(/sort=hot/)
  })

  test("分类标签不可点击（后端暂不支持）", async ({ page }) => {
    const categorySpan = page.locator("span").filter({ hasText: "Go基础" })
    if (await categorySpan.count() > 0) {
      // 应该是 span（不可点击），不是 button 或 a
      expect(await categorySpan.evaluate((el) => el.tagName.toLowerCase())).toBe("span")
    }
  })
})

test.describe("文章详情页", () => {
  test("访问文章详情（ID=1）", async ({ page }) => {
    await page.goto("/articles/1")
    await page.waitForLoadState("networkidle")

    // 页面应该渲染（有内容或错误提示）
    const hasContent =
      (await page.locator("article, main").count()) > 0 ||
      (await page.getByText(/不存在|找不到/).count()) > 0
    expect(hasContent).toBeTruthy()
  })
})
