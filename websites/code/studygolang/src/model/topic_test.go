// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package model_test

import (
	. "model"
	"testing"
)

func TestNewTopic(t *testing.T) {
	// err := NewTopic().Set("lastreplyuid=1").Where("uid=").Update()
	err := NewTopicEx().Where("tid=1").Increment("reply", 1)
	if err != nil {
		t.Fatal(err)
	}
}
