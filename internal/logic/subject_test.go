// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic_test

import (
	"reflect"
	"testing"

	"github.com/studygolang/studygolang/internal/logic"
	"github.com/studygolang/studygolang/internal/model"

	"golang.org/x/net/context"
)

func TestFindArticles(t *testing.T) {
	type args struct {
		ctx context.Context
		sid int
	}
	tests := []struct {
		name string
		self logic.SubjectLogic
		args args
		want []*model.Article
	}{
		{
			name: "subject1",
			args: args{
				nil,
				1,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := logic.SubjectLogic{}
			if got := self.FindArticles(tt.args.ctx, tt.args.sid, nil, ""); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SubjectLogic.FindArticles() = %v, want %v", got, tt.want)
			}
		})
	}
}
