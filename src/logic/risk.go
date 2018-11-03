// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	. "db"
	"model"

	"github.com/polaris1119/nosql"
)

type RiskLogic struct{}

var DefaultRisk = RiskLogic{}

// AddBlackIP 加入 IP 黑名单
func (RiskLogic) AddBlackIP(ip string) error {
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()

	key := "black:ip"
	return redisClient.HSET(key, ip, "1")
}

// AddBlackIPByUID 通过用户 UID 将最后一次登录 IP 加入黑名单
func (self RiskLogic) AddBlackIPByUID(uid int) error {
	userLogin := &model.UserLogin{}
	_, err := MasterDB.Where("uid=?", uid).Get(userLogin)
	if err != nil {
		return err
	}

	if userLogin.LoginIp != "" {
		return self.AddBlackIP(userLogin.LoginIp)
	}

	return nil
}

// IsBlackIP 是否是 IP 黑名单
func (RiskLogic) IsBlackIP(ip string) bool {
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()

	key := "black:ip"
	val, err := redisClient.HGET(key, ip)
	if err != nil {
		return false
	}

	return val == "1"
}
