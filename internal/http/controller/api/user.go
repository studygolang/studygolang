// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package api

import (
	"sync"
	"time"
)

type UserController struct{}

// ======================== 速率限制器 ========================

// rateLimiter 简单的内存速率限制器，按 IP 维度计数
type rateLimiter struct {
	mu       sync.Mutex
	attempts map[string]*attemptInfo
}

// attemptInfo 记录某个 key 在时间窗口内的尝试次数
type attemptInfo struct {
	count    int
	expireAt time.Time
}

// check 检查是否超过限制，返回 true 表示允许，false 表示被限流
func (rl *rateLimiter) check(key string, maxAttempts int, window time.Duration) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	info, exists := rl.attempts[key]
	if !exists || now.After(info.expireAt) {
		rl.attempts[key] = &attemptInfo{count: 1, expireAt: now.Add(window)}
		return true
	}

	info.count++
	if info.count > maxAttempts {
		return false
	}
	return true
}

var (
	loginLimiter    = &rateLimiter{attempts: make(map[string]*attemptInfo)}
	registerLimiter = &rateLimiter{attempts: make(map[string]*attemptInfo)}
)

func init() {
	// 后台定期清理过期速率限制条目，防止内存泄漏
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			loginLimiter.cleanup()
			registerLimiter.cleanup()
		}
	}()
}

// cleanup 清理过期的速率限制条目
func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	for key, info := range rl.attempts {
		if now.After(info.expireAt) {
			delete(rl.attempts, key)
		}
	}
}
