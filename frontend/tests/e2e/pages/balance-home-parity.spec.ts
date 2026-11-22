import { test, expect } from "@playwright/test"

test("首页和通用侧栏展示最新注册会员", async ({ page, request }) => {
  const response = await request.get("/api/v1/sidebar/users/newest")
  const users = (await response.json()).data.users as Array<{ username: string }>
  test.skip(users.length === 0, "需要至少一个用户")

  await page.goto("/")
  await expect(page.getByText("新会员", { exact: true })).toBeVisible()
  await expect(page.locator(`a[href="/user/${users[0].username}"]`).first()).toBeVisible()

  await page.goto("/gctt")
  await expect(page.getByText("新会员", { exact: true })).toBeVisible()
  await expect(page.locator(`a[href="/user/${users[0].username}"]`).first()).toBeVisible()
})

test("匿名访问充值记录保留原页并跳转登录", async ({ page }) => {
  await page.goto("/balance/add?p=2")
  await expect(page).toHaveURL(/\/account\/login\?redirect=/)
  expect(decodeURIComponent(new URL(page.url()).searchParams.get("redirect") ?? "")).toBe("/balance/add?p=2")
})

test("余额翻页同步旧 p 参数", async ({ page }) => {
  await page.route("**/api/v1/balance?p=*", async route => {
    const currentPage = Number(new URL(route.request().url()).searchParams.get("p")) || 1
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        code: 0,
        msg: "ok",
        data: {
          balance: 200,
          details: [],
          incomes: [],
          expenses: [],
          total: 100,
          page: currentPage,
          has_more: currentPage < 3,
          initial_mission: null,
        },
      }),
    })
  })

  await page.goto("/balance?p=2")
  await expect(page.getByRole("button", { name: "上一页" })).toBeVisible()
  await expect(page.getByRole("button", { name: "下一页" })).toBeVisible()

  await page.getByRole("button", { name: "上一页" }).click()
  await expect(page).toHaveURL(/\/balance$/)
  await expect(page.getByRole("button", { name: "上一页" })).toHaveCount(0)
  await expect(page.getByRole("button", { name: "下一页" })).toBeVisible()
})
