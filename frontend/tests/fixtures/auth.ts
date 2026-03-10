import { Page } from "@playwright/test"

export const TEST_USER = {
  username: "testuser",
  password: "test123456",
  email: "testuser@example.com",
}

/** 通过 localStorage 模拟已登录状态（无需真实网络请求） */
export async function mockLogin(page: Page, username = TEST_USER.username) {
  await page.addInitScript(
    ({ username }) => {
      localStorage.setItem("token", "mock-token-for-e2e-testing")
      localStorage.setItem("uid", "999")
      localStorage.setItem("username", username)
    },
    { username }
  )
}

/** 清除登录状态 */
export async function logout(page: Page) {
  await page.evaluate(() => {
    localStorage.removeItem("token")
    localStorage.removeItem("uid")
    localStorage.removeItem("username")
  })
}
