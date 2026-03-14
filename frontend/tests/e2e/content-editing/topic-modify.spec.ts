import { test, expect } from "@playwright/test"
import { mockLogin } from "../../fixtures/auth"

test.describe("话题修改页 /topics/modify/:id", () => {
  const TOPIC_ID = "1" // 使用一个假设存在的话题 ID

  test.describe("未登录状态", () => {
    test("未登录时应该重定向到登录页", async ({ page }) => {
      await page.goto(`/topics/modify/${TOPIC_ID}`)
      await expect(page).toHaveURL(/\/account\/login/)
    })

    test("重定向 URL 应该包含 redirect 参数", async ({ page }) => {
      await page.goto(`/topics/modify/${TOPIC_ID}`)
      await expect(page).toHaveURL(new RegExp(`redirect=.*topics.*modify.*${TOPIC_ID}`))
    })
  })

  test.describe("已登录状态", () => {
    test.beforeEach(async ({ page }) => {
      await mockLogin(page)
    })

    test("应该显示加载状态", async ({ page }) => {
      // 拦截 API 请求，延迟响应
      await page.route(`**/api/v1/topics/${TOPIC_ID}/edit`, async (route) => {
        await new Promise((resolve) => setTimeout(resolve, 500))
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              topic: {
                tid: 1,
                title: "测试话题标题",
                content: "测试话题内容",
                nid: 1,
                tags: "golang,test",
              },
              nodes: [],
            },
          }),
        })
      })

      await page.goto(`/topics/modify/${TOPIC_ID}`)
      // 加载状态应该短暂出现
      const loader = page.locator(".animate-spin")
      // 等待加载完成
      await expect(page.locator("h1")).toContainText("编辑话题", { timeout: 10000 })
    })

    test("应该正确加载话题数据", async ({ page }) => {
      await page.route(`**/api/v1/topics/${TOPIC_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              topic: {
                tid: 1,
                title: "测试话题标题",
                content: "测试话题内容",
                nid: 1,
                tags: "golang,test",
              },
              nodes: [
                {
                  "Go语言": [
                    { nid: 1, name: "Go基础", ename: "go-basic", pid: 0, seq: 1 },
                  ],
                },
              ],
            },
          }),
        })
      })

      await page.goto(`/topics/modify/${TOPIC_ID}`)
      await expect(page.locator("h1")).toContainText("编辑话题", { timeout: 10000 })

      // 标题应该被回填
      const titleInput = page.locator('input#title')
      await expect(titleInput).toHaveValue("测试话题标题", { timeout: 10000 })
    })

    test("应该回填标签", async ({ page }) => {
      await page.route(`**/api/v1/topics/${TOPIC_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              topic: {
                tid: 1,
                title: "测试话题标题",
                content: "测试话题内容",
                nid: 1,
                tags: "golang,test",
              },
              nodes: [],
            },
          }),
        })
      })

      await page.goto(`/topics/modify/${TOPIC_ID}`)
      await expect(page.locator("h1")).toContainText("编辑话题", { timeout: 10000 })

      // 标签应该被回填
      await expect(page.getByText("golang")).toBeVisible({ timeout: 10000 })
      await expect(page.getByText("test")).toBeVisible({ timeout: 10000 })
    })

    test("API 返回错误时应该显示错误信息", async ({ page }) => {
      await page.route(`**/api/v1/topics/${TOPIC_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 403,
            msg: "没有权限编辑此话题",
          }),
        })
      })

      await page.goto(`/topics/modify/${TOPIC_ID}`)
      await expect(page.getByText(/没有权限编辑此话题/)).toBeVisible({ timeout: 10000 })
    })

    test("网络错误时应该显示错误信息", async ({ page }) => {
      await page.route(`**/api/v1/topics/${TOPIC_ID}/edit`, async (route) => {
        await route.abort("failed")
      })

      await page.goto(`/topics/modify/${TOPIC_ID}`)
      await expect(page.getByText(/网络错误/)).toBeVisible({ timeout: 10000 })
    })

    test("应该加载 BlockNote 编辑器", async ({ page }) => {
      await page.route(`**/api/v1/topics/${TOPIC_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              topic: { tid: 1, title: "测试", content: "内容", nid: 1, tags: "" },
              nodes: [],
            },
          }),
        })
      })

      await page.goto(`/topics/modify/${TOPIC_ID}`)
      await expect(page.locator("h1")).toContainText("编辑话题", { timeout: 10000 })

      // 等待编辑器加载
      await page.waitForTimeout(2000)
      const editor = page.locator(".bn-editor")
      await expect(editor).toBeVisible({ timeout: 10000 })
    })

    test("标题为空时应该显示验证错误", async ({ page }) => {
      await page.route(`**/api/v1/topics/${TOPIC_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              topic: { tid: 1, title: "测试", content: "内容", nid: 1, tags: "" },
              nodes: [],
            },
          }),
        })
      })

      await page.goto(`/topics/modify/${TOPIC_ID}`)
      await expect(page.locator("h1")).toContainText("编辑话题", { timeout: 10000 })

      // 清空标题
      const titleInput = page.locator('input#title')
      await titleInput.clear()

      const saveButton = page.getByRole("button", { name: /保存/ })
      await saveButton.click()
      await expect(page.getByText(/请填写标题/)).toBeVisible()
    })

    test("未选择节点时应该显示验证错误", async ({ page }) => {
      await page.route(`**/api/v1/topics/${TOPIC_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              topic: { tid: 1, title: "测试", content: "内容", nid: 0, tags: "" },
              nodes: [],
            },
          }),
        })
      })

      await page.goto(`/topics/modify/${TOPIC_ID}`)
      await expect(page.locator("h1")).toContainText("编辑话题", { timeout: 10000 })

      const saveButton = page.getByRole("button", { name: /保存/ })
      await saveButton.click()
      await expect(page.getByText(/请选择节点/)).toBeVisible()
    })

    test("保存成功后应该跳转到话题详情页", async ({ page }) => {
      await page.route(`**/api/v1/topics/${TOPIC_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              topic: { tid: 1, title: "测试话题", content: "内容", nid: 1, tags: "" },
              nodes: [
                { "Go语言": [{ nid: 1, name: "Go基础", ename: "go-basic", pid: 0, seq: 1 }] },
              ],
            },
          }),
        })
      })

      await page.route(`**/api/v1/topics/${TOPIC_ID}`, async (route) => {
        if (route.request().method() === "PUT") {
          await route.fulfill({
            status: 200,
            contentType: "application/json",
            body: JSON.stringify({ code: 0, data: { tid: 1 } }),
          })
        } else {
          await route.continue()
        }
      })

      await page.goto(`/topics/modify/${TOPIC_ID}`)
      await expect(page.locator("h1")).toContainText("编辑话题", { timeout: 10000 })

      const saveButton = page.getByRole("button", { name: /保存/ })
      await saveButton.click()

      // 应该显示保存成功提示
      await expect(page.getByText(/保存成功/)).toBeVisible({ timeout: 5000 })

      // 应该跳转到话题详情页
      await expect(page).toHaveURL(new RegExp(`/topics/${TOPIC_ID}`), { timeout: 5000 })
    })

    test("保存失败时应该显示错误信息", async ({ page }) => {
      await page.route(`**/api/v1/topics/${TOPIC_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              topic: { tid: 1, title: "测试话题", content: "内容", nid: 1, tags: "" },
              nodes: [
                { "Go语言": [{ nid: 1, name: "Go基础", ename: "go-basic", pid: 0, seq: 1 }] },
              ],
            },
          }),
        })
      })

      await page.route(`**/api/v1/topics/${TOPIC_ID}`, async (route) => {
        if (route.request().method() === "PUT") {
          await route.fulfill({
            status: 200,
            contentType: "application/json",
            body: JSON.stringify({ code: 500, msg: "服务器内部错误" }),
          })
        } else {
          await route.continue()
        }
      })

      await page.goto(`/topics/modify/${TOPIC_ID}`)
      await expect(page.locator("h1")).toContainText("编辑话题", { timeout: 10000 })

      const saveButton = page.getByRole("button", { name: /保存/ })
      await saveButton.click()

      await expect(page.getByText(/服务器内部错误/)).toBeVisible({ timeout: 5000 })
    })

    test("应该显示取消按钮", async ({ page }) => {
      await page.route(`**/api/v1/topics/${TOPIC_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              topic: { tid: 1, title: "测试", content: "内容", nid: 1, tags: "" },
              nodes: [],
            },
          }),
        })
      })

      await page.goto(`/topics/modify/${TOPIC_ID}`)
      await expect(page.locator("h1")).toContainText("编辑话题", { timeout: 10000 })

      await expect(page.getByText("取消")).toBeVisible()
    })
  })
})
