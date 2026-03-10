import { Page, Locator } from "@playwright/test"

export class RegisterPage {
  readonly page: Page
  readonly usernameInput: Locator
  readonly emailInput: Locator
  readonly passwordInput: Locator
  readonly confirmPasswordInput: Locator
  readonly submitBtn: Locator
  readonly errorMsg: Locator
  readonly loginLink: Locator

  constructor(page: Page) {
    this.page = page
    this.usernameInput = page.locator("input#username")
    this.emailInput = page.locator("input#email")
    this.passwordInput = page.locator("input#password")
    this.confirmPasswordInput = page.locator("input#confirmPassword")
    this.submitBtn = page.locator("button[type='submit']")
    this.errorMsg = page.locator(".bg-destructive\\/10")
    this.loginLink = page.locator("a[href='/account/login']")
  }

  async goto() {
    await this.page.goto("/account/register")
    await this.page.waitForLoadState("networkidle")
  }

  async fillForm(opts: {
    username: string
    email: string
    password: string
    confirmPassword?: string
  }) {
    await this.usernameInput.fill(opts.username)
    await this.emailInput.fill(opts.email)
    await this.passwordInput.fill(opts.password)
    await this.confirmPasswordInput.fill(opts.confirmPassword ?? opts.password)
  }
}
