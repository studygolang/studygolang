import { test, expect } from "@playwright/test"
import { LoginPage } from "../../pages/LoginPage"

test.describe("用户登录", () => {
  let loginPage: LoginPage

  test.beforeEach(async ({ page }) => {
    loginPage = new LoginPage(page)
    await loginPage.goto()
  })

  test("登录页渲染正确", async ({ page }) => {
    await expect(page).toHaveTitle(/登录|Go语言中文网/)
    await expect(loginPage.usernameInput).toBeVisible()
    await expect(loginPage.passwordInput).toBeVisible()
    await expect(loginPage.submitBtn).toBeVisible()
    await expect(loginPage.submitBtn).toContainText("登录")

    // 有指向注册页的链接
    await expect(loginPage.registerLink).toBeVisible()
  })

  test("用户名为空时显示错误提示", async ({ page }) => {
    await loginPage.passwordInput.fill("somepassword")
    await loginPage.submitBtn.click()

    await expect(loginPage.errorMsg).toBeVisible()
    await expect(loginPage.errorMsg).toContainText("用户名")
  })

  test("密码为空时显示错误提示", async ({ page }) => {
    await loginPage.usernameInput.fill("someuser")
    await loginPage.submitBtn.click()

    await expect(loginPage.errorMsg).toBeVisible()
    await expect(loginPage.errorMsg).toContainText("密码")
  })

  test("错误凭证时显示后端错误", async ({ page }) => {
    await loginPage.login("wronguser", "wrongpassword")

    // 等待 API 响应
    await page.waitForLoadState("networkidle")

    // 应显示错误信息（来自后端或网络）
    await expect(loginPage.errorMsg).toBeVisible()
  })

  test("登录页有 autocomplete 属性", async ({ page }) => {
    const usernameAutoComplete = await loginPage.usernameInput.getAttribute("autocomplete")
    const passwordAutoComplete = await loginPage.passwordInput.getAttribute("autocomplete")

    expect(usernameAutoComplete).toBe("username")
    expect(passwordAutoComplete).toBe("current-password")
  })

  test("注册链接可以跳转到注册页", async ({ page }) => {
    await loginPage.registerLink.click()
    await expect(page).toHaveURL("/account/register")
  })

  test("点击忘记密码链接", async ({ page }) => {
    const forgotLink = page.locator("a[href='/account/forgot-password']")
    await expect(forgotLink).toBeVisible()
    await expect(forgotLink).toContainText("忘记密码")
  })
})
