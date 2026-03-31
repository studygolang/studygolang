// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build integration
// +build integration

package logic

import (
	"fmt"
	"sync"
	"testing"

	"github.com/polaris1119/nosql"
)

// TestGetNextCommentFloorIntegration 测试评论楼层 Redis 原子计数器
// 运行方式：go test -v -tags=integration -run TestGetNextCommentFloorIntegration ./internal/logic/
func TestGetNextCommentFloorIntegration(t *testing.T) {
	t.Run("SequentialIncrements", func(t *testing.T) {
		objid := 9999991
		objtype := 99

		// 获取第一个楼层号
		floor1, err := GetNextCommentFloor(objid, objtype)
		if err != nil {
			t.Fatalf("GetNextCommentFloor() error = %v", err)
		}
		if floor1 < 1 {
			t.Errorf("floor1 = %d, want >= 1", floor1)
		}

		// 获取第二个楼层号，应该比第一个大 1
		floor2, err := GetNextCommentFloor(objid, objtype)
		if err != nil {
			t.Fatalf("GetNextCommentFloor() second call error = %v", err)
		}
		if floor2 != floor1+1 {
			t.Errorf("floor2 = %d, want %d", floor2, floor1+1)
		}

		// 获取第三个楼层号
		floor3, err := GetNextCommentFloor(objid, objtype)
		if err != nil {
			t.Fatalf("GetNextCommentFloor() third call error = %v", err)
		}
		if floor3 != floor2+1 {
			t.Errorf("floor3 = %d, want %d", floor3, floor2+1)
		}
	})

	t.Run("DifferentObjects", func(t *testing.T) {
		objid1 := 9999992
		objid2 := 9999993
		objtype := 99

		floor1, err := GetNextCommentFloor(objid1, objtype)
		if err != nil {
			t.Fatalf("GetNextCommentFloor() obj1 error = %v", err)
		}

		floor2, err := GetNextCommentFloor(objid2, objtype)
		if err != nil {
			t.Fatalf("GetNextCommentFloor() obj2 error = %v", err)
		}

		// 不同对象的楼层号应该是独立的
		// 两者都应该是有效的楼层号（>=1）
		if floor1 < 1 {
			t.Errorf("floor1 = %d, want >= 1", floor1)
		}
		if floor2 < 1 {
			t.Errorf("floor2 = %d, want >= 1", floor2)
		}
	})

	t.Run("DifferentTypes", func(t *testing.T) {
		objid := 9999994
		objtype1 := 98
		objtype2 := 97

		floor1, err := GetNextCommentFloor(objid, objtype1)
		if err != nil {
			t.Fatalf("GetNextCommentFloor() type1 error = %v", err)
		}

		floor2, err := GetNextCommentFloor(objid, objtype2)
		if err != nil {
			t.Fatalf("GetNextCommentFloor() type2 error = %v", err)
		}

		// 不同类型的楼层号应该是独立的
		if floor1 < 1 {
			t.Errorf("floor1 = %d, want >= 1", floor1)
		}
		if floor2 < 1 {
			t.Errorf("floor2 = %d, want >= 1", floor2)
		}
	})
}

// TestCommentFloorKeyFormat 测试 Redis key 格式
func TestCommentFloorKeyFormat(t *testing.T) {
	// 验证 key 格式正确性
	expectedPrefix := "comment:floor:"
	if commentFloorKeyPrefix != expectedPrefix {
		t.Errorf("commentFloorKeyPrefix = %s, want %s", commentFloorKeyPrefix, expectedPrefix)
	}
}

// TestGetNextCommentFloorConcurrency 测试并发获取楼层号（竞态条件）
// 运行方式：go test -v -tags=integration -run TestGetNextCommentFloorConcurrency ./internal/logic/
func TestGetNextCommentFloorConcurrency(t *testing.T) {
	objid := 999
	objtype := 1
	goroutines := 100

	// 清理 Redis
	redisClient := nosql.NewRedisClient()
	key := fmt.Sprintf("%s%d:%d", commentFloorKeyPrefix, objtype, objid)
	redisClient.DEL(key)
	redisClient.Close()

	// 并发获取楼层号
	floors := make([]int, goroutines)
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(index int) {
			defer wg.Done()
			floor, err := GetNextCommentFloor(objid, objtype)
			if err != nil {
				t.Errorf("GetNextCommentFloor failed: %v", err)
				return
			}
			floors[index] = floor
		}(i)
	}

	wg.Wait()

	// 验证：所有楼层号应该唯一且连续
	floorSet := make(map[int]bool)
	for _, floor := range floors {
		if floorSet[floor] {
			t.Errorf("Duplicate floor number: %d", floor)
		}
		floorSet[floor] = true
	}

	// 验证范围：1 到 goroutines
	for i := 1; i <= goroutines; i++ {
		if !floorSet[i] {
			t.Errorf("Missing floor number: %d", i)
		}
	}

	t.Logf("Concurrency test passed: %d unique floors generated", goroutines)
}
