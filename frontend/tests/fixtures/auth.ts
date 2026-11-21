import { Page } from "@playwright/test"

export const TEST_USER = {
  username: "testuser",
  password: "test123456",
  email: "testuser@example.com",
}

/** 模拟当前 Cookie 会话对应的登录用户（无需真实网络请求） */
export async function mockLogin(page: Page, username = TEST_USER.username) {
  await page.addInitScript(
    ({ username }) => {
      localStorage.setItem("token", "mock-token-for-e2e-testing")
      localStorage.setItem("uid", "999")
      localStorage.setItem("username", username)
    },
    { username }
  )
  await page.route("**/api/v1/user/me", async (route) => {
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        code: 0,
        data: {
          user: {
            uid: 999,
            username,
            name: username,
            is_root: false,
            is_vip: true,
            balance: 100,
          },
        },
      }),
    })
  })
  await page.route("**/api/v1/interviews/today", async (route) => {
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ code: 0, data: { question: null } }),
    })
  })
}

/** 清除登录状态 */
export async function logout(page: Page) {
  await page.evaluate(() => {
    localStorage.removeItem("token")
    localStorage.removeItem("uid")
    localStorage.removeItem("username")
  })
}
