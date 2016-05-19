// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"errors"
	"sync"

	"github.com/polaris1119/logger"

	. "db"
	"model"
)

// 常驻内存数据（多实例部署时，数据同步会有问题）

var (
	authLocker  sync.RWMutex
	Authorities []*model.Authority

	roleAuthLocker  sync.RWMutex
	RoleAuthorities map[int][]int

	roleLocker sync.RWMutex
	Roles      []*model.Role // 相应的 roleid-1 为索引

	nodeRWMutex sync.RWMutex
	// 节点信息
	AllNode []map[string]interface{}

	catRWMutex sync.RWMutex
	// 资源分类
	AllCategory []*model.ResourceCat
)

// 将所有 权限 加载到内存中；后台修改权限时，重新加载一次
func LoadAuthorities() error {
	authorities := make([]*model.Authority, 0)
	err := MasterDB.Find(&authorities)
	if err != nil {
		logger.Errorln("LoadAuthorities authority read fail:", err)
		return err
	}

	authLocker.Lock()
	defer authLocker.Unlock()

	Authorities = authorities

	logger.Infoln("LoadAuthorities successfully!")

	return nil
}

// 将所有 角色拥有的权限 加载到内存中；后台修改时，重新加载一次
func LoadRoleAuthorities() error {
	roleAuthorities := make([]*model.RoleAuthority, 0)
	err := MasterDB.Find(&roleAuthorities)
	if err != nil {
		logger.Errorln("LoadRoleAuthorities role_authority read fail:", err)
		return err
	}

	roleAuthLocker.Lock()
	defer roleAuthLocker.Unlock()

	RoleAuthorities = make(map[int][]int)

	for _, roleAuth := range roleAuthorities {
		roleId := roleAuth.Roleid

		if authorities, ok := RoleAuthorities[roleId]; ok {
			RoleAuthorities[roleId] = append(authorities, roleAuth.Aid)
		} else {
			RoleAuthorities[roleId] = []int{roleAuth.Aid}
		}
	}

	logger.Infoln("LoadRoleAuthorities successfully!")

	return nil
}

// 将所有 角色 加载到内存中；后台修改角色时，重新加载一次
func LoadRoles() error {
	roles := make([]*model.Role, 0)
	err := MasterDB.Find(&roles)
	if err != nil {
		logger.Errorln("LoadRoles role read fail:", err)
		return err
	}

	if len(roles) == 0 {
		logger.Errorln("LoadRoles role read fail: num is 0")
		return errors.New("no role")
	}

	roleLocker.Lock()
	defer roleLocker.Unlock()

	maxRoleid := roles[len(roles)-1].Roleid
	Roles = make([]*model.Role, maxRoleid)

	// 由于角色不多，而且一般角色id是连续自增的，因此这里以角色id当slice的index
	for _, role := range roles {
		Roles[role.Roleid-1] = role
	}

	logger.Infoln("LoadRoles successfully!")

	return nil
}

// 将所有 节点信息 加载到内存中：后台修改节点时，重新加载一次
func LoadNodes() error {
	nodeList := make([]*model.TopicNode, 0)
	err := MasterDB.Find(&nodeList)
	if err != nil {
		logger.Errorln("LoadNodes node read fail:", err)
		return err
	}
	nodeNum := len(nodeList)
	tmpNodeList := make(map[int]*model.TopicNode, nodeNum)
	for _, node := range nodeList {
		tmpNodeList[node.Nid] = node
	}

	nodeRWMutex.Lock()
	defer nodeRWMutex.Unlock()

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

	logger.Infoln("LoadNodes successfully!")

	return nil
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

// 将 node 组织成一定结构，方便前端展示
func GenNodes() []map[string][]map[string]interface{} {
	sameParent := make(map[string][]map[string]interface{})
	allParentNodes := make([]string, 0, 8)
	for _, node := range AllNode {
		if node["pid"].(int) != 0 {
			if len(sameParent[node["parent"].(string)]) == 0 {
				sameParent[node["parent"].(string)] = []map[string]interface{}{node}
			} else {
				sameParent[node["parent"].(string)] = append(sameParent[node["parent"].(string)], node)
			}
		} else {
			allParentNodes = append(allParentNodes, node["name"].(string))
		}
	}
	nodes := make([]map[string][]map[string]interface{}, 0, len(allParentNodes))
	for _, parent := range allParentNodes {
		tmpMap := make(map[string][]map[string]interface{})
		tmpMap[parent] = sameParent[parent]
		nodes = append(nodes, tmpMap)
	}
	logger.Debugf("%v\n", nodes)
	return nodes
}

// 将所有 资源分类信息 加载到内存中：后台修改节点时，重新加载一次
func LoadCategories() (err error) {
	categories := make([]*model.ResourceCat, 0)
	err = MasterDB.Find(&categories)
	if err != nil {
		logger.Errorln("LoadCategories category read fail:", err)
		return
	}

	catRWMutex.Lock()
	defer catRWMutex.Unlock()

	AllCategory = categories

	logger.Infoln("LoadCategories successfully!")

	return
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
