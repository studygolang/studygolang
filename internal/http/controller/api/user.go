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

// pwdResetEntry 重置密码条目，带过期时间
type pwdResetEntry struct {
	email    string
	expireAt time.Time
}

// resetPwdMap 存储重置密码 token -> *pwdResetEntry 的映射（内存存储，重启失效）
var resetPwdMap = sync.Map{}

func init() {
	// 后台定期清理过期 token（每 10 分钟）
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			now := time.Now()
			resetPwdMap.Range(func(key, value interface{}) bool {
				if entry, ok := value.(*pwdResetEntry); ok && now.After(entry.expireAt) {
					resetPwdMap.Delete(key)
				}
				return true
			})
		}
	}()
}

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