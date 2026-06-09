// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package logic

import (
	"encoding/json"
	"strconv"
)

const (
	// Redis key 前缀
	userInfoCachePrefix = "user:info:"
	// 用户缓存 TTL: 30 分钟
	userCacheTTL = 30 * 60
)

// UserInfoCache 用户信息缓存
type UserInfoCache struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsRoot   bool   `json:"is_root"`
	IsVip    bool   `json:"is_vip"`
	IsAdmin  bool   `json:"is_admin"`
	Avatar   string `json:"avatar"`
	Balance  int    `json:"balance"`
}

// GetCachedUserInfo 获取缓存的用户信息
// 返回 nil 表示缓存未命中
func GetCachedUserInfo(uid int) *UserInfoCache {
	redisClient := GetRedisClient()

	key := userInfoCachePrefix + strconv.Itoa(uid)
	data := redisClient.GET(key)
	if data == "" {
		return nil
	}

	var userInfo UserInfoCache
	if err := json.Unmarshal([]byte(data), &userInfo); err != nil {
		GetLogger(nil).Errorln("GetCachedUserInfo unmarshal error:", err)
		return nil
	}

	return &userInfo
}

// SetCachedUserInfo 缓存用户信息
func SetCachedUserInfo(userInfo *UserInfoCache) error {
	if userInfo == nil || userInfo.Uid == 0 {
		return nil
	}

	redisClient := GetRedisClient()

	key := userInfoCachePrefix + strconv.Itoa(userInfo.Uid)
	data, err := json.Marshal(userInfo)
	if err != nil {
		return err
	}

	return redisClient.SET(key, string(data), userCacheTTL)
}

// InvalidateUserCache 使用户缓存失效
func InvalidateUserCache(uid int) error {
	redisClient := GetRedisClient()
	key := userInfoCachePrefix + strconv.Itoa(uid)
	return redisClient.DEL(key)
}

// GetOrFetchUserInfo 获取用户信息，优先从缓存获取
func GetOrFetchUserInfo(uid int, fetchFunc func() *UserInfoCache) *UserInfoCache {
	// 1. 尝试从缓存获取
	cached := GetCachedUserInfo(uid)
	if cached != nil {
		return cached
	}

	// 2. 缓存未命中，从数据库获取
	userInfo := fetchFunc()
	if userInfo == nil {
		return nil
	}

	// 3. 存入缓存
	_ = SetCachedUserInfo(userInfo)

	return userInfo
}
