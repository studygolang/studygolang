import { Page, Locator } from "@playwright/test"

/** 站点导航栏 Page Object */
export class SiteHeaderPage {
  readonly page: Page
  readonly logo: Locator
  readonly navLinks: Locator
  readonly searchInput: Locator
  readonly loginBtn: Locator
  readonly registerBtn: Locator
  readonly userMenu: Locator
  readonly mobileMenuToggle: Locator

  constructor(page: Page) {
    this.page = page
    this.logo = page.locator("header a[href='/']").first()
    this.navLinks = page.locator("header nav a")
    this.searchInput = page.locator("header input[placeholder*='搜索']")
    this.loginBtn = page.locator("header a[href='/account/login']")
    this.registerBtn = page.locator("header a[href='/account/register']")
    this.userMenu = page.locator("header button").filter({ hasText: /\w+/ }).first()
    this.mobileMenuToggle = page.locator("header button[aria-label*='菜单']")
  }

  async searchFor(query: string) {
    await this.searchInput.fill(query)
    await this.searchInput.press("Enter")
    await this.page.waitForLoadState("networkidle")
  }
}
