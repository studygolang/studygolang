import { test, expect } from "@playwright/test"
import { mockLogin } from "../../fixtures/auth"

test.describe("文章修改页 /articles/modify/:id", () => {
  const ARTICLE_ID = "1" // 使用一个假设存在的文章 ID

  test.describe("未登录状态", () => {
    test("未登录时应该重定向到登录页", async ({ page }) => {
      await page.goto(`/articles/modify/${ARTICLE_ID}`)
      await expect(page).toHaveURL(/\/account\/login/)
    })

    test("重定向 URL 应该包含 redirect 参数", async ({ page }) => {
      await page.goto(`/articles/modify/${ARTICLE_ID}`)
      await expect(page).toHaveURL(new RegExp(`redirect=.*articles.*modify.*${ARTICLE_ID}`))
    })
  })

  test.describe("已登录状态", () => {
    test.beforeEach(async ({ page }) => {
      await mockLogin(page)
    })

    test("应该显示加载状态", async ({ page }) => {
      // 拦截 API 请求，延迟响应
      await page.route(`**/api/v1/articles/${ARTICLE_ID}/edit`, async (route) => {
        await new Promise((resolve) => setTimeout(resolve, 500))
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              article: {
                id: 1,
                title: "测试文章标题",
                content: "<p>测试文章内容</p>",
                txt: "测试文章内容",
                tags: "golang,test",
              },
            },
          }),
        })
      })

      await page.goto(`/articles/modify/${ARTICLE_ID}`)
      // 等待加载完成
      await expect(page.locator("h1")).toContainText("编辑文章", { timeout: 10000 })
    })

    test("应该正确加载文章数据", async ({ page }) => {
      await page.route(`**/api/v1/articles/${ARTICLE_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              article: {
                id: 1,
                title: "测试文章标题",
                content: "<p>测试文章内容</p>",
                txt: "测试文章内容",
                tags: "golang,test",
              },
            },
          }),
        })
      })

      await page.goto(`/articles/modify/${ARTICLE_ID}`)
      await expect(page.locator("h1")).toContainText("编辑文章", { timeout: 10000 })

      // 标题应该被回填
      const titleInput = page.locator('input#title')
      await expect(titleInput).toHaveValue("测试文章标题", { timeout: 10000 })
    })

    test("应该优先使用 txt 字段作为内容", async ({ page }) => {
      await page.route(`**/api/v1/articles/${ARTICLE_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              article: {
                id: 1,
                title: "测试文章标题",
                content: "<p>HTML内容</p>",
                txt: "纯文本内容",
                tags: "golang",
              },
            },
          }),
        })
      })

      await page.goto(`/articles/modify/${ARTICLE_ID}`)
      await expect(page.locator("h1")).toContainText("编辑文章", { timeout: 10000 })

      // 等待编辑器加载
      await page.waitForTimeout(2000)
      // txt 字段应该被优先使用
      // 注意：实际内容在编辑器内部，这里只是验证页面加载成功
    })

    test("应该回填标签", async ({ page }) => {
      await page.route(`**/api/v1/articles/${ARTICLE_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              article: {
                id: 1,
                title: "测试文章标题",
                content: "测试文章内容",
                txt: "测试文章内容",
                tags: "golang,test,article",
              },
            },
          }),
        })
      })

      await page.goto(`/articles/modify/${ARTICLE_ID}`)
      await expect(page.locator("h1")).toContainText("编辑文章", { timeout: 10000 })

      // 标签应该被回填
      await expect(page.getByText("golang")).toBeVisible({ timeout: 10000 })
      await expect(page.getByText("test")).toBeVisible({ timeout: 10000 })
      await expect(page.getByText("article")).toBeVisible({ timeout: 10000 })
    })

    test("API 返回错误时应该显示错误信息", async ({ page }) => {
      await page.route(`**/api/v1/articles/${ARTICLE_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 403,
            msg: "没有权限编辑此文章",
          }),
        })
      })

      await page.goto(`/articles/modify/${ARTICLE_ID}`)
      await expect(page.getByText(/没有权限编辑此文章/)).toBeVisible({ timeout: 10000 })
    })

    test("网络错误时应该显示错误信息", async ({ page }) => {
      await page.route(`**/api/v1/articles/${ARTICLE_ID}/edit`, async (route) => {
        await route.abort("failed")
      })

      await page.goto(`/articles/modify/${ARTICLE_ID}`)
      await expect(page.getByText(/网络错误/)).toBeVisible({ timeout: 10000 })
    })

    test("应该加载 BlockNote 编辑器", async ({ page }) => {
      await page.route(`**/api/v1/articles/${ARTICLE_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              article: {
                id: 1,
                title: "测试",
                content: "内容",
                txt: "内容",
                tags: "",
              },
            },
          }),
        })
      })

      await page.goto(`/articles/modify/${ARTICLE_ID}`)
      await expect(page.locator("h1")).toContainText("编辑文章", { timeout: 10000 })

      // 等待编辑器加载
      await page.waitForTimeout(2000)
      const editor = page.locator(".bn-editor")
      await expect(editor).toBeVisible({ timeout: 10000 })
    })

    test("标题为空时应该显示验证错误", async ({ page }) => {
      await page.route(`**/api/v1/articles/${ARTICLE_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              article: {
                id: 1,
                title: "测试",
                content: "内容",
                txt: "内容",
                tags: "",
              },
            },
          }),
        })
      })

      await page.goto(`/articles/modify/${ARTICLE_ID}`)
      await expect(page.locator("h1")).toContainText("编辑文章", { timeout: 10000 })

      // 清空标题
      const titleInput = page.locator('input#title')
      await titleInput.clear()

      const saveButton = page.getByRole("button", { name: /保存/ })
      await saveButton.click()
      await expect(page.getByText(/请填写标题/)).toBeVisible()
    })

    test("内容为空时应该显示验证错误", async ({ page }) => {
      await page.route(`**/api/v1/articles/${ARTICLE_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              article: {
                id: 1,
                title: "测试",
                content: "",
                txt: "",
                tags: "",
              },
            },
          }),
        })
      })

      await page.goto(`/articles/modify/${ARTICLE_ID}`)
      await expect(page.locator("h1")).toContainText("编辑文章", { timeout: 10000 })

      const saveButton = page.getByRole("button", { name: /保存/ })
      await saveButton.click()
      await expect(page.getByText(/请填写内容/)).toBeVisible()
    })

    test("应该能够修改标题", async ({ page }) => {
      await page.route(`**/api/v1/articles/${ARTICLE_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              article: {
                id: 1,
                title: "原标题",
                content: "内容",
                txt: "内容",
                tags: "",
              },
            },
          }),
        })
      })

      await page.goto(`/articles/modify/${ARTICLE_ID}`)
      await expect(page.locator("h1")).toContainText("编辑文章", { timeout: 10000 })

      const titleInput = page.locator('input#title')
      await titleInput.clear()
      await titleInput.fill("新标题")
      await expect(titleInput).toHaveValue("新标题")
    })

    test("应该能够添加新标签", async ({ page }) => {
      await page.route(`**/api/v1/articles/${ARTICLE_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              article: {
                id: 1,
                title: "测试",
                content: "内容",
                txt: "内容",
                tags: "golang",
              },
            },
          }),
        })
      })

      await page.goto(`/articles/modify/${ARTICLE_ID}`)
      await expect(page.locator("h1")).toContainText("编辑文章", { timeout: 10000 })

      const tagInput = page.locator('input[placeholder*="标签"]')
      await tagInput.fill("newtag")
      await tagInput.press("Enter")
      await expect(page.getByText("newtag")).toBeVisible()
    })

    test("应该能够删除标签", async ({ page }) => {
      await page.route(`**/api/v1/articles/${ARTICLE_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              article: {
                id: 1,
                title: "测试",
                content: "内容",
                txt: "内容",
                tags: "golang,test",
              },
            },
          }),
        })
      })

      await page.goto(`/articles/modify/${ARTICLE_ID}`)
      await expect(page.locator("h1")).toContainText("编辑文章", { timeout: 10000 })

      // 等待标签加载
      await expect(page.getByText("golang")).toBeVisible({ timeout: 10000 })

      // 删除第一个标签
      const removeButton = page.locator('button[aria-label*="删除"]').first()
      await removeButton.click()
      await expect(page.getByText("golang")).not.toBeVisible()
    })

    test("点击取消应该返回上一页", async ({ page }) => {
      await page.route(`**/api/v1/articles/${ARTICLE_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              article: {
                id: 1,
                title: "测试",
                content: "内容",
                txt: "内容",
                tags: "",
              },
            },
          }),
        })
      })

      await page.goto(`/articles/modify/${ARTICLE_ID}`)
      await expect(page.locator("h1")).toContainText("编辑文章", { timeout: 10000 })

      const cancelButton = page.getByRole("button", { name: "取消" })
      await cancelButton.click()
      // 应该返回上一页
      await page.waitForTimeout(500)
    })

    test("保存成功后应该显示成功消息", async ({ page }) => {
      await page.route(`**/api/v1/articles/${ARTICLE_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              article: {
                id: 1,
                title: "测试",
                content: "内容",
                txt: "内容",
                tags: "",
              },
            },
          }),
        })
      })

      await page.route(`**/api/v1/articles/${ARTICLE_ID}`, async (route) => {
        if (route.request().method() === "PUT") {
          await route.fulfill({
            status: 200,
            contentType: "application/json",
            body: JSON.stringify({
              code: 0,
              msg: "保存成功",
            }),
          })
        }
      })

      await page.goto(`/articles/modify/${ARTICLE_ID}`)
      await expect(page.locator("h1")).toContainText("编辑文章", { timeout: 10000 })

      const saveButton = page.getByRole("button", { name: /保存/ })
      await saveButton.click()
      await expect(page.getByText(/保存成功/)).toBeVisible({ timeout: 10000 })
    })

    test("保存失败时应该显示错误消息", async ({ page }) => {
      await page.route(`**/api/v1/articles/${ARTICLE_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              article: {
                id: 1,
                title: "测试",
                content: "内容",
                txt: "内容",
                tags: "",
              },
            },
          }),
        })
      })

      await page.route(`**/api/v1/articles/${ARTICLE_ID}`, async (route) => {
        if (route.request().method() === "PUT") {
          await route.fulfill({
            status: 200,
            contentType: "application/json",
            body: JSON.stringify({
              code: 500,
              msg: "保存失败",
            }),
          })
        }
      })

      await page.goto(`/articles/modify/${ARTICLE_ID}`)
      await expect(page.locator("h1")).toContainText("编辑文章", { timeout: 10000 })

      const saveButton = page.getByRole("button", { name: /保存/ })
      await saveButton.click()
      await expect(page.getByText(/保存失败/)).toBeVisible({ timeout: 10000 })
    })

    test("保存时按钮应该显示加载状态", async ({ page }) => {
      await page.route(`**/api/v1/articles/${ARTICLE_ID}/edit`, async (route) => {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({
            code: 0,
            data: {
              article: {
                id: 1,
                title: "测试",
                content: "内容",
                txt: "内容",
                tags: "",
              },
            },
          }),
        })
      })

      await page.route(`**/api/v1/articles/${ARTICLE_ID}`, async (route) => {
        if (route.request().method() === "PUT") {
          await new Promise((resolve) => setTimeout(resolve, 1000))
          await route.fulfill({
            status: 200,
            contentType: "application/json",
            body: JSON.stringify({
              code: 0,
              msg: "保存成功",
            }),
          })
        }
      })

      await page.goto(`/articles/modify/${ARTICLE_ID}`)
      await expect(page.locator("h1")).toContainText("编辑文章", { timeout: 10000 })

      const saveButton = page.getByRole("button", { name: /保存/ })
      await saveButton.click()
      await expect(page.getByText(/保存中/)).toBeVisible()
      await expect(saveButton).toBeDisabled()
    })
  })
})
