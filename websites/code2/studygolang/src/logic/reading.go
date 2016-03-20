// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package logic

import (
	. "db"
	"model"
)

type ReadingLogic struct{}

var DefaultReading = ReadingLogic{}

func (ReadingLogic) FindLastList(beginTime string) ([]*model.MorningReading, error) {
	readings := make([]*model.MorningReading, 0)
	err := MasterDB.Where("ctime>? AND rtype=0", beginTime).OrderBy("id DESC").Find(&readings)

	return readings, err
}
