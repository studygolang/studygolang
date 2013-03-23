// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package model

import (
	"logger"
	"sort"
	"sync"
)

// 常驻内存数据

var (
	// 由于角色数据很少，为了方便，让其常驻内存。当数据有变动时，应该调用UpdateAllRole更新该数据
	AllRole map[int]*Role
	// 保存所有Roleid，升序。需要保证角色权限越来越小
	AllRoleId []int

	nodeRWMutex sync.RWMutex
	// 节点信息
	AllNode []map[string]interface{}

	catRWMutex sync.RWMutex
	// 资源分类
	AllCategory []*ResourceCat
)

func init() {
	UpdateAllRole()
	UpdateAllNode()
	UpdateAllCategory()
}

// 更新 AllRole 数据
func UpdateAllRole() {
	roleList, err := NewRole().FindAll()
	if err != nil {
		logger.Errorln("获取角色数据失败：", err)
		return
	}
	roleNum := len(roleList)
	AllRole = make(map[int]*Role, roleNum)
	AllRoleId = make([]int, roleNum)
	for i, role := range roleList {
		AllRole[role.Roleid] = role
		AllRoleId[i] = role.Roleid
	}

	// 对AllRoleId进行排序
	sort.Sort(sort.IntSlice(AllRoleId))
}

// 更新 AllNode 数据
func UpdateAllNode() {
	nodeList, err := NewTopicNode().FindAll()
	if err != nil {
		logger.Errorln("获取节点数据失败：", err)
		return
	}
	nodeNum := len(nodeList)
	tmpNodeList := make(map[int]*TopicNode, nodeNum)
	for _, node := range nodeList {
		tmpNodeList[node.Nid] = node
	}
	AllNode = make([]map[string]interface{}, nodeNum)
	for i, node := range nodeList {
		nodeMap := make(map[string]interface{}, 5)
		nodeMap["pid"] = node.Parent
		if node.Parent == 0 {
			nodeMap["parent"] = "根节点"
		} else {
			nodeMap["parent"] = tmpNodeList[node.Parent].Name
		}
		nodeMap["nid"] = node.Nid
		nodeMap["name"] = node.Name
		nodeMap["intro"] = node.Intro
		nodeMap["ctime"] = node.Ctime
		AllNode[i] = nodeMap
	}
}

// 获得单个节点名
func GetNodeName(nid int) string {
	nodeRWMutex.RLock()
	defer nodeRWMutex.RUnlock()
	for _, node := range AllNode {
		if node["nid"].(int) == nid {
			return node["name"].(string)
		}
	}
	return ""
}

// 获得单个节点信息
func GetNode(nid int) map[string]interface{} {
	nodeRWMutex.RLock()
	defer nodeRWMutex.RUnlock()
	for _, node := range AllNode {
		if node["nid"].(int) == nid {
			return node
		}
	}
	return nil
}

// 获得多个节点名
func GetNodesName(nids []int) map[int]string {
	nodes := make(map[int]string, len(nids))
	nodeRWMutex.RLock()
	defer nodeRWMutex.RUnlock()
	for _, nid := range nids {
		for _, node := range AllNode {
			if node["nid"].(int) == nid {
				nodes[nid] = node["name"].(string)
			}
		}
	}
	return nodes
}

// 更新 AllCategory 数据
func UpdateAllCategory() {
	var err error
	AllCategory, err = NewResourceCat().FindAll()
	if err != nil {
		logger.Errorln("获取资源分类数据失败：", err)
		return
	}
}

// 获得分类名
func GetCategoryName(catid int) string {
	catRWMutex.RLock()
	defer catRWMutex.RUnlock()
	for _, cat := range AllCategory {
		if cat.Catid == catid {
			return cat.Name
		}
	}
	return ""
}
