// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/internal/logic"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type NodeController struct{}

// RegisterRoute 注册节点相关路由
func (self NodeController) RegisterRoute(g *echo.Group) {
	// 节点别名路由，通过 slug 或名称查询节点及话题列表
	g.GET("/go/:node", self.NodeTopics)
}

// NodeTopics 通过节点别名获取节点信息和话题列表
// GET /api/v1/go/:node?p=1
// node 参数可以是节点的 ename (slug) 或名称
func (NodeController) NodeTopics(ctx echo.Context) error {
	nodeParam := ctx.Param("node")
	if nodeParam == "" {
		return fail(ctx, "节点参数不能为空")
	}

	// 尝试通过 ename 查找节点
	node := logic.GetNodeByEname(nodeParam)
	if node == nil {
		// 如果 ename 找不到，尝试通过名称查找（兼容性）
		// GenNodes 返回 []map[string][]map[string]interface{}
		// 外层 slice 是分组，每个分组的 value 是子节点列表
		nodes := logic.GenNodes()
		found := false
		for _, group := range nodes {
			for _, children := range group {
				for _, n := range children {
					if name, ok := n["name"].(string); ok && name == nodeParam {
						node = n
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if found {
				break
			}
		}
		if node == nil {
			return fail(ctx, "节点不存在")
		}
	}

	// 获取节点 ID
	nid, ok := node["nid"].(int)
	if !ok {
		// 尝试从 map[string]interface{} 中获取
		if nidFloat, ok2 := node["nid"].(float64); ok2 {
			nid = int(nidFloat)
		} else if nidInt64, ok3 := node["nid"].(int64); ok3 {
			nid = int(nidInt64)
		} else {
			return fail(ctx, "节点 ID 格式错误")
		}
	}

	// 分页获取话题列表
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	if curPage < 1 {
		curPage = 1
	}
	paginator := logic.NewPaginatorWithPerPage(curPage, perPage)

	querystring := "nid=?"
	topics := logic.DefaultTopic.FindAll(context.EchoContext(ctx), paginator, "topics.mtime DESC", querystring, nid)
	total := logic.DefaultTopic.Count(context.EchoContext(ctx), querystring, nid)
	hasMore := paginator.SetTotal(total).HasMorePage()

	return success(ctx, map[string]interface{}{
		"node":     node,
		"list":     topics,
		"total":    total,
		"page":     curPage,
		"has_more": hasMore,
	})
}
