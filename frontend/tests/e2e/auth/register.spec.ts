import { test, expect } from "@playwright/test"
import { RegisterPage } from "../../pages/RegisterPage"

test.describe("用户注册", () => {
  let registerPage: RegisterPage

  test.beforeEach(async ({ page }) => {
    registerPage = new RegisterPage(page)
    await registerPage.goto()
  })

  test("注册页渲染正确", async ({ page }) => {
    await expect(page).toHaveTitle(/注册|创建账号|Go语言中文网/)
    await expect(registerPage.usernameInput).toBeVisible()
    await expect(registerPage.emailInput).toBeVisible()
    await expect(registerPage.passwordInput).toBeVisible()
    await expect(registerPage.confirmPasswordInput).toBeVisible()
    await expect(registerPage.submitBtn).toContainText("注册")
  })

  test("用户名太短时显示校验错误", async () => {
    await registerPage.fillForm({
      username: "ab",
      email: "test@example.com",
      password: "password123",
    })
    await registerPage.submitBtn.click()

    await expect(registerPage.errorMsg).toBeVisible()
    await expect(registerPage.errorMsg).toContainText("3")
  })

  test("邮箱格式不正确时显示错误", async () => {
    await registerPage.fillForm({
      username: "validuser",
      email: "not-an-email",
      password: "password123",
    })
    await registerPage.submitBtn.click()

    await expect(registerPage.errorMsg).toBeVisible()
    await expect(registerPage.errorMsg).toContainText("邮箱")
  })

  test("密码太短时显示错误", async () => {
    await registerPage.fillForm({
      username: "validuser",
      email: "test@example.com",
      password: "123",
    })
    await registerPage.submitBtn.click()

    await expect(registerPage.errorMsg).toBeVisible()
    await expect(registerPage.errorMsg).toContainText("6")
  })

  test("两次密码不一致时显示错误", async () => {
    await registerPage.fillForm({
      username: "validuser",
      email: "test@example.com",
      password: "password123",
      confirmPassword: "different456",
    })
    await registerPage.submitBtn.click()

    await expect(registerPage.errorMsg).toBeVisible()
    await expect(registerPage.errorMsg).toContainText("不一致")
  })

  test("登录链接可以跳转到登录页", async ({ page }) => {
    await registerPage.loginLink.click()
    await expect(page).toHaveURL("/account/login")
  })

  test("注册表单有正确的 autocomplete 属性", async () => {
    expect(await registerPage.usernameInput.getAttribute("autocomplete")).toBe("username")
    expect(await registerPage.emailInput.getAttribute("autocomplete")).toBe("email")
    expect(await registerPage.passwordInput.getAttribute("autocomplete")).toBe("new-password")
    expect(await registerPage.confirmPasswordInput.getAttribute("autocomplete")).toBe("new-password")
  })
})
