// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"context"
	. "db"
	"model"
	"net/url"

	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
)

type TopicNodeLogic struct{}

var DefaultNode = TopicNodeLogic{}

func (self TopicNodeLogic) FindOne(nid int) *model.TopicNode {
	topicNode := &model.TopicNode{}
	_, err := MasterDB.Id(nid).Get(topicNode)
	if err != nil {
		logger.Errorln("TopicNodeLogic FindOne error:", err, "nid:", nid)
	}

	return topicNode
}

func (self TopicNodeLogic) FindByEname(ename string) *model.TopicNode {
	topicNode := &model.TopicNode{}
	_, err := MasterDB.Where("ename=?", ename).Get(topicNode)
	if err != nil {
		logger.Errorln("TopicNodeLogic FindByEname error:", err, "ename:", ename)
	}

	return topicNode
}

func (self TopicNodeLogic) FindByNids(nids []int) map[int]*model.TopicNode {
	nodeList := make(map[int]*model.TopicNode, 0)
	err := MasterDB.In("nid", nids).Find(&nodeList)
	if err != nil {
		logger.Errorln("TopicNodeLogic FindByNids error:", err, "nids:", nids)
	}

	return nodeList
}

func (self TopicNodeLogic) FindByParent(pid, num int) []*model.TopicNode {
	nodeList := make([]*model.TopicNode, 0)
	err := MasterDB.Where("parent=?", pid).Limit(num).Find(&nodeList)
	if err != nil {
		logger.Errorln("TopicNodeLogic FindByParent error:", err, "parent:", pid)
	}

	return nodeList
}

func (self TopicNodeLogic) FindAll(ctx context.Context) []*model.TopicNode {
	nodeList := make([]*model.TopicNode, 0)
	err := MasterDB.Asc("seq").Find(&nodeList)
	if err != nil {
		logger.Errorln("TopicNodeLogic FindAll error:", err)
	}

	return nodeList
}

func (self TopicNodeLogic) Modify(ctx context.Context, form url.Values) error {
	objLog := GetLogger(ctx)

	node := &model.TopicNode{}
	err := schemaDecoder.Decode(node, form)
	if err != nil {
		objLog.Errorln("TopicNodeLogic Modify decode error:", err)
		return err
	}

	nid := goutils.MustInt(form.Get("nid"))
	if nid == 0 {
		// 新增
		_, err = MasterDB.Insert(node)
		if err != nil {
			objLog.Errorln("TopicNodeLogic Modify insert error:", err)
		}
		return err
	}

	change := make(map[string]interface{})

	fields := []string{"parent", "logo", "name", "ename", "intro", "seq", "show_index"}
	for _, field := range fields {
		change[field] = form.Get(field)
	}

	_, err = MasterDB.Table(new(model.TopicNode)).Id(nid).Update(change)
	if err != nil {
		objLog.Errorln("TopicNodeLogic Modify update error:", err)
	}
	return err
}

func (self TopicNodeLogic) ModifySeq(ctx context.Context, nid, seq int) error {
	_, err := MasterDB.Table(new(model.TopicNode)).Id(nid).Update(map[string]interface{}{"seq": seq})
	return err
}

func (self TopicNodeLogic) FindParallelTree(ctx context.Context) []*model.TopicNode {
	nodeList := make([]*model.TopicNode, 0)
	err := MasterDB.Asc("parent").Asc("seq").Find(&nodeList)
	if err != nil {
		logger.Errorln("TopicNodeLogic FindTreeList error:", err)

		return nil
	}

	showNodeList := make([]*model.TopicNode, 0, len(nodeList))
	self.tileNodes(&showNodeList, nodeList, 0, 1, 3, 0)

	return showNodeList
}

func (self TopicNodeLogic) tileNodes(showNodeList *[]*model.TopicNode, nodeList []*model.TopicNode, parentId, curLevel, showLevel, pos int) {
	for num := len(nodeList); pos < num; pos++ {
		node := nodeList[pos]

		if node.Parent == parentId {
			*showNodeList = append(*showNodeList, node)

			if node.Level == 0 {
				node.Level = curLevel
			}

			if curLevel <= showLevel {
				self.tileNodes(showNodeList, nodeList, node.Nid, curLevel+1, showLevel, pos+1)
			}
		}

		if node.Parent > parentId {
			break
		}
	}
}
