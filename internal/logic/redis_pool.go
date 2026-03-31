// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package logic

import (
	"sync"

	"github.com/polaris1119/nosql"
)

var (
	// globalRedisClient 全局 Redis 客户端单例
	globalRedisClient *nosql.RedisClient
	redisClientOnce   sync.Once
)

// GetRedisClient 获取全局 Redis 客户端（单例模式）
// 避免每次操作都创建新连接，提升性能
func GetRedisClient() *nosql.RedisClient {
	redisClientOnce.Do(func() {
		globalRedisClient = nosql.NewRedisClient()
	})
	return globalRedisClient
}

// CloseRedisClient 关闭全局 Redis 客户端（程序退出时调用）
func CloseRedisClient() {
	if globalRedisClient != nil {
		globalRedisClient.Close()
	}
}
