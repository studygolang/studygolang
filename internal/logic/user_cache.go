// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package logic

import (
	"encoding/json"
	"strconv"

	"github.com/polaris1119/logger"
	"golang.org/x/sync/singleflight"
)

// userInfoSF 用于在缓存击穿时合并同一 uid 的并发请求
var userInfoSF singleflight.Group

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
	Status   int    `json:"status"`
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
// 使用 singleflight 防止缓存击穿：同一 uid 的并发请求只会触发一次 DB 查询
func GetOrFetchUserInfo(uid int, fetchFunc func() *UserInfoCache) *UserInfoCache {
	// 1. 尝试从缓存获取
	cached := GetCachedUserInfo(uid)
	if cached != nil {
		return cached
	}

	// 2. 缓存未命中，使用 singleflight 合并并发请求
	v, _, _ := userInfoSF.Do(strconv.Itoa(uid), func() (interface{}, error) {
		// 双重检查：进入 singleflight 后再次查缓存，避免重复 DB
		if cached := GetCachedUserInfo(uid); cached != nil {
			return cached, nil
		}
		userInfo := fetchFunc()
		if userInfo == nil {
			return nil, nil
		}
		// 存入缓存
		if err := SetCachedUserInfo(userInfo); err != nil {
			logger.Errorln("GetOrFetchUserInfo: cache set failed for uid", uid, ":", err)
		}
		return userInfo, nil
	})
	if v == nil {
		return nil
	}
	return v.(*UserInfoCache)
}
