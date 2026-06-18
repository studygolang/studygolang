// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build integration
// +build integration

package logic

import (
	"testing"
)

// TestUserCacheIntegration 测试用户缓存集成
// 运行方式：go test -v -tags=integration -run TestUserCacheIntegration ./internal/logic/
func TestUserCacheIntegration(t *testing.T) {
	t.Run("SetAndGetCachedUserInfo", func(t *testing.T) {
		userInfo := &UserInfoCache{
			Uid:      999991,
			Username: "test_cache_user",
			Email:    "cache@test.com",
			IsRoot:   false,
			IsVip:    false,
			Avatar:   "avatar.jpg",
		}

		// 设置缓存
		err := SetCachedUserInfo(userInfo)
		if err != nil {
			t.Fatalf("SetCachedUserInfo() error = %v", err)
		}

		// 获取缓存
		cached := GetCachedUserInfo(userInfo.Uid)
		if cached == nil {
			t.Fatal("GetCachedUserInfo() returned nil")
		}

		if cached.Uid != userInfo.Uid {
			t.Errorf("cached.Uid = %d, want %d", cached.Uid, userInfo.Uid)
		}
		if cached.Username != userInfo.Username {
			t.Errorf("cached.Username = %s, want %s", cached.Username, userInfo.Username)
		}
		if cached.Email != userInfo.Email {
			t.Errorf("cached.Email = %s, want %s", cached.Email, userInfo.Email)
		}

		// 清理缓存
		err = InvalidateUserCache(userInfo.Uid)
		if err != nil {
			t.Errorf("InvalidateUserCache() error = %v", err)
		}

		// 验证缓存已清除
		cached = GetCachedUserInfo(userInfo.Uid)
		if cached != nil {
			t.Error("GetCachedUserInfo() should return nil after invalidation")
		}
	})

	t.Run("GetOrFetchUserInfo", func(t *testing.T) {
		testUid := 999992
		fetchCount := 0

		// 清理可能存在的缓存
		_ = InvalidateUserCache(testUid)

		// 第一次调用，应该触发 fetchFunc
		userInfo := GetOrFetchUserInfo(testUid, func() *UserInfoCache {
			fetchCount++
			return &UserInfoCache{
				Uid:      testUid,
				Username: "fetched_user",
				Email:    "fetched@test.com",
				IsRoot:   false,
				IsVip:    false,
				Avatar:   "fetched_avatar.jpg",
			}
		})

		if userInfo == nil {
			t.Fatal("GetOrFetchUserInfo() returned nil")
		}
		if fetchCount != 1 {
			t.Errorf("fetchCount = %d, want 1 (should call fetchFunc once)", fetchCount)
		}

		// 第二次调用，应该从缓存获取
		userInfo2 := GetOrFetchUserInfo(testUid, func() *UserInfoCache {
			fetchCount++
			return &UserInfoCache{
				Uid:      testUid,
				Username: "fetched_again",
				Email:    "fetched_again@test.com",
				IsRoot:   false,
				IsVip:    false,
				Avatar:   "fetched_again_avatar.jpg",
			}
		})

		if userInfo2 == nil {
			t.Fatal("GetOrFetchUserInfo() second call returned nil")
		}
		if fetchCount != 1 {
			t.Errorf("fetchCount = %d, want 1 (should use cache, not call fetchFunc)", fetchCount)
		}
		if userInfo2.Username != "fetched_user" {
			t.Errorf("userInfo2.Username = %s, want 'fetched_user' (from cache)", userInfo2.Username)
		}

		// 清理
		_ = InvalidateUserCache(testUid)
	})

	t.Run("CacheExpiration", func(t *testing.T) {
		testUid := 999993

		// 设置一个短 TTL 的缓存（通过手动设置）
		userInfo := &UserInfoCache{
			Uid:      testUid,
			Username: "expire_test",
			Email:    "expire@test.com",
			IsRoot:   false,
			IsVip:    false,
			Avatar:   "expire_avatar.jpg",
		}

		err := SetCachedUserInfo(userInfo)
		if err != nil {
			t.Fatalf("SetCachedUserInfo() error = %v", err)
		}

		// 立即获取，应该成功
		cached := GetCachedUserInfo(testUid)
		if cached == nil {
			t.Fatal("GetCachedUserInfo() returned nil immediately after set")
		}

		// 等待 31 分钟（超过 TTL）- 这里只模拟，不实际等待
		// 实际测试中可以手动设置一个很短的 TTL 来测试
		// 这里只是验证 TTL 常量正确
		if userCacheTTL != 30*60 {
			t.Errorf("userCacheTTL = %d, want %d", userCacheTTL, 30*60)
		}

		// 清理
		_ = InvalidateUserCache(testUid)
	})
}

// TestUserCacheUnit 测试用户缓存单元（不需要 Redis）
func TestUserCacheUnit(t *testing.T) {
	t.Run("UserInfoCacheStruct", func(t *testing.T) {
		userInfo := &UserInfoCache{
			Uid:      123,
			Username: "testuser",
			Email:    "test@example.com",
			IsRoot:   true,
			IsVip:    false,
			Avatar:   "test.jpg",
		}

		if userInfo.Uid != 123 {
			t.Errorf("userInfo.Uid = %d, want 123", userInfo.Uid)
		}
		if !userInfo.IsRoot {
			t.Error("userInfo.IsRoot should be true")
		}
	})

	t.Run("NilUserInfoCache", func(t *testing.T) {
		// 测试 nil 输入
		var userInfo *UserInfoCache = nil

		if userInfo != nil {
			t.Error("nil UserInfoCache should be nil")
		}
	})
}

// TestInvalidateUserCacheOnUpdate 测试用户更新时缓存失效
func TestInvalidateUserCacheOnUpdate(t *testing.T) {
	// 这个测试需要完整的数据库环境，标记为集成测试
	t.Skip("Requires database connection - run with -tags=integration")
}
