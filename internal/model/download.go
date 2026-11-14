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
	Id          int       `json:"id" xorm:"pk autoincr"`
	Version     string    `json:"version"`
	Filename    string    `json:"filename"`
	Kind        string    `json:"kind"`
	OS          string    `json:"os" xorm:"os"`
	Arch        string    `json:"arch"`
	Size        int       `json:"size"`
	Checksum    string    `json:"checksum"`
	Category    int       `json:"category"`
	IsRecommend bool      `json:"is_recommend"`
	Seq         int       `json:"seq"`
	CreatedAt   time.Time `json:"created_at" xorm:"created"`
}
