// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package util_test

import (
	. "github.com/studygolang/studygolang/util"
	"testing"
)

type model struct {
	Id   int
	Name string
}

func TestModels2Intslice(t *testing.T) {
	models := []*model{
		{12, "polaris"},
		{13, "xuxinhua"},
	}

	actualResult := Models2Intslice(models, "Id")
	expectResult := []int{12, 13}

	if !sliceIsEqual(actualResult, expectResult) {
		t.Fatalf("expect:%v, actual:%v", expectResult, actualResult)
	}
}

func sliceIsEqual(slice1, slice2 []int) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i, v := range slice1 {
		if v != slice2[i] {
			return false
		}
	}

	return true
}
