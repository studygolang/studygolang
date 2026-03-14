import { Page, Locator } from "@playwright/test"

export type MessageType = "system" | "inbox" | "outbox"

export class MessagePage {
  readonly page: Page
  readonly systemTab: Locator
  readonly inboxTab: Locator
  readonly outboxTab: Locator
  readonly messageList: Locator
  readonly emptyState: Locator
  readonly loadingState: Locator
  readonly pagination: Locator
  readonly prevPageBtn: Locator
  readonly nextPageBtn: Locator

  constructor(page: Page) {
    this.page = page
    this.systemTab = page.locator('button[data-state]:has-text("系统消息")')
    this.inboxTab = page.locator('button[data-state]:has-text("收件箱")')
    this.outboxTab = page.locator('button[data-state]:has-text("发件箱")')
    this.messageList = page.locator('[data-testid="message-list"]').first()
    this.emptyState = page.locator("text=暂无消息").first()
    this.loadingState = page.locator("text=加载中...").first()
    this.pagination = page.locator(".flex.items-center.justify-center.gap-4").first()
    this.prevPageBtn = page.locator("button:has-text('上一页')").first()
    this.nextPageBtn = page.locator("button:has-text('下一页')").first()
  }

  async goto(msgtype: MessageType = "system") {
    await this.page.goto(`/message/${msgtype}`)
    await this.page.waitForLoadState("networkidle")
  }

  async switchTab(type: MessageType) {
    const tabMap = {
      system: this.systemTab,
      inbox: this.inboxTab,
      outbox: this.outboxTab,
    }
    await tabMap[type].click()
    await this.page.waitForLoadState("networkidle")
  }

  async getMessageItems() {
    return this.page.locator(".rounded-lg.border.bg-card").all()
  }

  async getMessageById(id: number) {
    return this.page.locator(`[data-message-id="${id}"]`).first()
  }

  async deleteMessage(index: number) {
    const messages = await this.getMessageItems()
    if (messages.length > index) {
      const deleteBtn = messages[index].locator("button:has-text('删除')").first()
      await deleteBtn.click()

      // 等待确认对话框
      const confirmBtn = this.page.locator("button:has-text('删除')").last()
      await confirmBtn.click()

      // 等待删除完成
      await this.page.waitForTimeout(500)
    }
  }

  async cancelDelete(index: number) {
    const messages = await this.getMessageItems()
    if (messages.length > index) {
      const deleteBtn = messages[index].locator("button:has-text('删除')").first()
      await deleteBtn.click()

      // 点击取消
      const cancelBtn = this.page.locator("button:has-text('取消')").last()
      await cancelBtn.click()
    }
  }

  async goToNextPage() {
    await this.nextPageBtn.click()
    await this.page.waitForLoadState("networkidle")
  }

  async goToPrevPage() {
    await this.prevPageBtn.click()
    await this.page.waitForLoadState("networkidle")
  }

  async getMessageCount() {
    const messages = await this.getMessageItems()
    return messages.length
  }

  async isMessageRead(index: number) {
    const messages = await this.getMessageItems()
    if (messages.length > index) {
      const badge = messages[index].locator(".bg-primary").first()
      return (await badge.count()) === 0
    }
    return false
  }
}
