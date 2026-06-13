// Copyright 2024 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package config 提供配置管理功能，支持环境变量覆盖
package config

import (
	"os"
	"strconv"

	"github.com/Unknwon/goconfig"
	config "github.com/polaris1119/config"
)

// EnvConfig 提供配置读取功能，支持环境变量覆盖
type EnvConfig struct {
	configFile *goconfig.ConfigFile
}

// NewEnvConfig 创建支持环境变量的配置读取器
func NewEnvConfig() *EnvConfig {
	return &EnvConfig{
		configFile: config.ConfigFile,
	}
}

// MustValue 获取配置值，优先从环境变量读取
// envKey: 环境变量名（如 DB_PASSWORD）
// section: 配置文件的 section
// key: 配置文件的 key
// defaultValue: 默认值
func (ec *EnvConfig) MustValue(envKey, section, key, defaultValue string) string {
	// 1. 优先从环境变量读取
	if envValue := os.Getenv(envKey); envValue != "" {
		return envValue
	}

	// 2. 从配置文件读取
	return ec.configFile.MustValue(section, key, defaultValue)
}

// MustInt 获取整数配置值，优先从环境变量读取
func (ec *EnvConfig) MustInt(envKey, section, key string, defaultValue int) int {
	// 1. 优先从环境变量读取
	if envValue := os.Getenv(envKey); envValue != "" {
		if intValue, err := strconv.Atoi(envValue); err == nil {
			return intValue
		}
	}

	// 2. 从配置文件读取
	return ec.configFile.MustInt(section, key, defaultValue)
}

// MustBool 获取布尔配置值，优先从环境变量读取
func (ec *EnvConfig) MustBool(envKey, section, key string, defaultValue bool) bool {
	// 1. 优先从环境变量读取
	if envValue := os.Getenv(envKey); envValue != "" {
		if boolValue, err := strconv.ParseBool(envValue); err == nil {
			return boolValue
		}
	}

	// 2. 从配置文件读取
	return ec.configFile.MustBool(section, key, defaultValue)
}

// 全局配置实例
var EnvConfigInstance = NewEnvConfig()

// 便捷函数

// MustValue 获取配置值（支持环境变量）
func MustValue(envKey, section, key, defaultValue string) string {
	return EnvConfigInstance.MustValue(envKey, section, key, defaultValue)
}

// MustInt 获取整数配置值（支持环境变量）
func MustInt(envKey, section, key string, defaultValue int) int {
	return EnvConfigInstance.MustInt(envKey, section, key, defaultValue)
}

// MustBool 获取布尔配置值（支持环境变量）
func MustBool(envKey, section, key string, defaultValue bool) bool {
	return EnvConfigInstance.MustBool(envKey, section, key, defaultValue)
}

// GetEnvWithFallback 获取环境变量，如果不存在则返回 fallback
func GetEnvWithFallback(envKey, fallback string) string {
	if value := os.Getenv(envKey); value != "" {
		return value
	}
	return fallback
}
