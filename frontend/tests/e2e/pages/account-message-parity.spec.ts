import { test, expect } from "@playwright/test"
import { mockLogin } from "../../fixtures/auth"

function json(data: unknown) {
  return {
    status: 200,
    contentType: "application/json",
    body: JSON.stringify({ code: 0, data }),
  }
}

test("资料保存保留微博和签名，纯 OAuth 用户可首次设置密码", async ({ page }) => {
  await mockLogin(page)
  let profileUpdate: Record<string, unknown> | null = null
  let passwordUpdate: Record<string, unknown> | null = null

  await page.route("**/api/v1/account/bind_users", async (route) => {
    await route.fulfill(json({ bind_users: [] }))
  })
  await page.route("**/api/v1/user/profile", async (route) => {
    if (route.request().method() === "PUT") {
      profileUpdate = route.request().postDataJSON()
      await route.fulfill(json({ need_activate: false }))
      return
    }
    await route.fulfill(json({
      has_passwd: false,
      user: {
        uid: 999,
        username: "testuser",
        name: "测试用户",
        city: "上海",
        company: "StudyGolang",
        github: "studygolang",
        weibo: "studygolang",
        website: "https://studygolang.com",
        monlog: "保持好奇",
        introduce: "个人介绍",
        open: 1,
      },
    }))
  })
  await page.route("**/api/v1/user/password", async (route) => {
    passwordUpdate = route.request().postDataJSON()
    await route.fulfill(json(null))
  })

  await page.goto("/user/testuser/settings")
  await expect(page.getByLabel("微博")).toHaveValue("studygolang")
  await expect(page.getByLabel("签名", { exact: true })).toHaveValue("保持好奇")
  await page.getByRole("button", { name: "保存修改" }).click()
  await expect(page.getByText("资料更新成功")).toBeVisible()
  expect(profileUpdate).toMatchObject({ weibo: "studygolang", monlog: "保持好奇" })

  await page.getByRole("button", { name: "设置密码", exact: true }).click()
  await expect(page.getByLabel("当前密码")).toHaveCount(0)
  await page.getByLabel("新密码").fill("new-password")
  await page.getByRole("button", { name: "确认设置" }).click()
  await expect(page.getByText("密码设置成功")).toBeVisible()
  expect(passwordUpdate).toEqual({ cur_passwd: "", new_passwd: "new-password" })
})

test("回复私信时预填收件人并展示匹配的原消息", async ({ page }) => {
  await mockLogin(page)
  await page.route("**/api/v1/user/alice", async (route) => {
    await route.fulfill(json({ user: { uid: 42, username: "alice" } }))
  })
  await page.route("**/api/v1/messages/7", async (route) => {
    await route.fulfill(json({
      message: { id: 7, content: "原私信内容", from: 42, ctime: "2026-01-01" },
    }))
  })

  await page.goto("/messages/send?username=alice&id=7")
  await expect(page.getByLabel(/收件人/)).toHaveValue("alice")
  await expect(page.getByText("回复 alice 的私信")).toBeVisible()
  await expect(page.getByText("原私信内容", { exact: true })).toBeVisible()
})

test("收件箱提供带用户名和消息 ID 的回复入口", async ({ page }) => {
  await mockLogin(page)
  await page.route("**/api/v1/messages?type=inbox*", async (route) => {
    await route.fulfill(json({
      messages: [{
        id: 7,
        content: "原私信内容",
        hasread: "未读",
        ctime: "2026-01-01",
        user: { uid: 42, username: "alice", avatar: "" },
      }],
      page: 1,
      total: 1,
      has_more: false,
    }))
  })

  await page.goto("/message/inbox")
  await expect(page.getByRole("link", { name: "回复", exact: true })).toHaveAttribute(
    "href",
    "/messages/send?username=alice&id=7",
  )
})
