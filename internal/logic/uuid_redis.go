// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package logic

import (
	"time"
)

const (
	// Redis key 前缀
	resetPwdKeyPrefix = "resetpwd:"
	activateKeyPrefix = "activate:"
)

// SetResetPwdUUID 存储重置密码的 UUID 映射到 Redis
// TTL: 1 小时
func SetResetPwdUUID(uuid, email string) error {
	redisClient := GetRedisClient() // 使用全局连接池
	key := resetPwdKeyPrefix + uuid
	return redisClient.SET(key, email, int(time.Hour.Seconds()))
}

// GetResetPwdUUID 从 Redis 获取重置密码的邮箱
// 返回空字符串表示不存在或已过期
func GetResetPwdUUID(uuid string) string {
	redisClient := GetRedisClient() // 使用全局连接池
	key := resetPwdKeyPrefix + uuid
	email := redisClient.GET(key)
	return email
}

// DelResetPwdUUID 删除重置密码的 UUID 映射
func DelResetPwdUUID(uuid string) error {
	redisClient := GetRedisClient() // 使用全局连接池
	key := resetPwdKeyPrefix + uuid
	return redisClient.DEL(key)
}

// SetActivateUUID 存储注册激活的 UUID 映射到 Redis
// TTL: 24 小时
func SetActivateUUID(uuid, email string) error {
	redisClient := GetRedisClient() // 使用全局连接池
	key := activateKeyPrefix + uuid
	return redisClient.SET(key, email, int(time.Hour.Seconds())*24)
}

// GetActivateUUID 从 Redis 获取注册激活的邮箱
// 返回空字符串表示不存在或已过期
func GetActivateUUID(uuid string) string {
	redisClient := GetRedisClient() // 使用全局连接池
	key := activateKeyPrefix + uuid
	email := redisClient.GET(key)
	return email
}

// DelActivateUUID 删除注册激活的 UUID 映射
func DelActivateUUID(uuid string) error {
	redisClient := GetRedisClient() // 使用全局连接池
	key := activateKeyPrefix + uuid
	return redisClient.DEL(key)
}
