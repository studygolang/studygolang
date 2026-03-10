import { test, expect } from "@playwright/test"

/** 验证各页面能正常加载（200 状态，无崩溃） */
const pages = [
  { name: "项目列表", url: "/projects" },
  { name: "资源列表", url: "/resources" },
  { name: "晨读列表", url: "/readings" },
  { name: "书籍列表", url: "/books" },
  { name: "招聘列表", url: "/jobs" },
  { name: "节点列表", url: "/nodes" },
  { name: "Wiki 列表", url: "/wiki" },
]

for (const { name, url } of pages) {
  test(`${name}（${url}）页面正常渲染`, async ({ page }) => {
    const response = await page.goto(url)
    await page.waitForLoadState("networkidle")

    // HTTP 状态应为 200
    expect(response?.status()).toBe(200)

    // 页面 body 可见（无崩溃）
    await expect(page.locator("body")).toBeVisible()

    // 页面应有 header 导航
    await expect(page.locator("header")).toBeVisible()

    // 页面不应有未处理的 JS 错误
    const errors: string[] = []
    page.on("pageerror", (err) => errors.push(err.message))
    expect(errors).toHaveLength(0)
  })
}

test("用户主页渲染正常", async ({ page }) => {
  // 访问一个可能存在的用户主页
  const response = await page.goto("/user/polaris")
  await page.waitForLoadState("networkidle")

  expect(response?.status()).toBe(200)
  await expect(page.locator("body")).toBeVisible()
})

test("所有页面底部都有 footer", async ({ page }) => {
  const pagesToCheck = ["/", "/topics", "/articles", "/projects"]

  for (const url of pagesToCheck) {
    await page.goto(url)
    await page.waitForLoadState("networkidle")

    await expect(page.locator("footer"), `${url} 应有 footer`).toBeVisible()
  }
})
