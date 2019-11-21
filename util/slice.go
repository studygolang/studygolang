// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package util

// InSlice 判断 i 是否在 slice 中
func InSlice(i int, slice []int) bool {
	for _, val := range slice {
		if val == i {
			return true
		}
	}

	return false
}
