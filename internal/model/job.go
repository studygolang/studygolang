// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import "time"

// Job 招聘职位
type Job struct {
	Id          int       `json:"id" xorm:"pk autoincr"`
	Title       string    `json:"title"`                            // 职位名称
	Company     string    `json:"company"`                          // 公司名称
	CompanySize string    `json:"company_size" xorm:"company_size"` // 公司规模
	City        string    `json:"city"`                             // 城市
	SalaryMin   int       `json:"salary_min" xorm:"salary_min"`     // 最低薪资（k）
	SalaryMax   int       `json:"salary_max" xorm:"salary_max"`     // 最高薪资（k）
	Experience  string    `json:"experience"`                       // 经验要求
	Tags        string    `json:"tags"`                             // 标签，逗号分隔
	Uid         int       `json:"uid"`                              // 发布者 uid
	Status      int       `json:"status"`                           // 状态：0=正常，1=关闭
	Ctime       time.Time `json:"ctime" xorm:"created"`
	Mtime       time.Time `json:"mtime" xorm:"updated"`
}
