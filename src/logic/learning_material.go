// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	. "db"

	"model"

	"golang.org/x/net/context"
)

type LearningMaterialLogic struct{}

var DefaultLearningMaterial = LearningMaterialLogic{}

func (LearningMaterialLogic) FindAll(ctx context.Context) []*model.LearningMaterial {
	objLog := GetLogger(ctx)

	learningMaterials := make([]*model.LearningMaterial, 0)
	err := MasterDB.Asc("type").Desc("seq").Find(&learningMaterials)
	if err != nil {
		objLog.Errorln("LearningMaterialLogic FindAll error:", err)
		return nil
	}

	return learningMaterials
}
