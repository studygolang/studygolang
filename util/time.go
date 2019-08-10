// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package util

import "time"
import "fmt"

// MonthDayNum t 所在时间的月份总天数
func MonthDayNum(t time.Time) int {
	isLeapYear := isLeap(t.Year())

	month := t.Month()
	switch month {
	case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
		return 31
	case time.February:
		if isLeapYear {
			return 29
		}

		return 28
	default:
		return 30
	}
}

func TimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)
	hours := diff.Hours()
	if hours < 1.0 {
		return fmt.Sprintf("约 %.0f 分钟前", diff.Minutes())
	}

	if hours < 24.0 {
		return fmt.Sprintf("约 %.0f 小时前", hours)
	}

	if hours < 72.0 {
		return fmt.Sprintf("约 %.0f 天前", hours/24.0)
	}

	// 同一年，不用年份
	if now.Year() == t.Year() {
		return t.Format("01-02 15:04")
	}

	return t.Format("2006-01-02")
}

// 是否闰年
func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
