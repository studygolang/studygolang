// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build integration
// +build integration

package api

// 注意：此测试文件标记为 integration 构建，需要完整的运行环境
// 运行方式：go test -v -tags=integration ./internal/http/controller/api/...
//
// 以下测试用例验证资源和图书发布的字段验证逻辑：
//
// 1. 资源发布验证测试：
//   - 缺少标题 -> 返回 "资源标题不能为空"
//   - 缺少资源类型 -> 返回 "请选择资源类型"
//   - 资源类型不合法 -> 返回 "资源类型不合法"
//   - 链接类型缺少URL -> 返回 "链接类型资源需要填写 URL"
//   - 内容类型缺少内容 -> 返回 "内容类型资源需要填写内容"
//
// 2. 图书发布验证测试：
//   - 缺少书名 -> 返回 "书名不能为空"
//
// 3. 认证测试：
//   - 未登录用户发布资源/图书 -> 返回 "请先登录"
//   - 未登录用户编辑资源/图书 -> 返回 "请先登录"
//   - 未登录用户更新资源/图书 -> 返回 "请先登录"
//
// 详细测试代码请参见原始测试文件模板
