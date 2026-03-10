import { Page, Locator } from "@playwright/test"

export class LoginPage {
  readonly page: Page
  readonly usernameInput: Locator
  readonly passwordInput: Locator
  readonly submitBtn: Locator
  readonly errorMsg: Locator
  readonly registerLink: Locator

  constructor(page: Page) {
    this.page = page
    this.usernameInput = page.locator("input#username")
    this.passwordInput = page.locator("input#password")
    this.submitBtn = page.locator("button[type='submit']")
    this.errorMsg = page.locator(".bg-destructive\\/10")
    this.registerLink = page.locator("a[href='/account/register']")
  }

  async goto() {
    await this.page.goto("/account/login")
    await this.page.waitForLoadState("networkidle")
  }

  async login(username: string, password: string) {
    await this.usernameInput.fill(username)
    await this.passwordInput.fill(password)
    await this.submitBtn.click()
  }
}
