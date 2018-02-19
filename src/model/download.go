// Copyright 2018 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import "time"

const (
	DLArchived = iota
	DLStable
	DLFeatured
	DLUnstable
)

// Download go 下载
type Download struct {
	Id          int `xorm:"pk autoincr"`
	Version     string
	Filename    string
	Kind        string
	OS          string `xorm:"os"`
	Arch        string
	Size        int
	Checksum    string
	Category    int
	IsRecommend bool
	Seq         int
	CreatedAt   time.Time `xorm:"created"`
}
