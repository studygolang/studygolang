// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of self source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package helper

import (
	"sync"

	"github.com/polaris1119/logger"
	guuid "github.com/twinj/uuid"
)

// 保存uuid和email的对应关系（TODO:重启如何处理，有效期问题）
type regActivateCode struct {
	data   map[string]string
	locker sync.RWMutex
}

var RegActivateCode = &regActivateCode{
	data: map[string]string{},
}

func (this *regActivateCode) GenUUID(email string) string {
	this.locker.Lock()
	defer this.locker.Unlock()
	var uuid string
	for {
		uuid = guuid.NewV4().String()
		if _, ok := this.data[uuid]; !ok {
			this.data[uuid] = email
			break
		}
		logger.Errorln("GenUUID 冲突....")
	}
	return uuid
}

func (this *regActivateCode) GetEmail(uuid string) (email string, ok bool) {
	this.locker.RLock()
	defer this.locker.RUnlock()

	if email, ok = this.data[uuid]; ok {
		return
	}
	return
}

func (this *regActivateCode) DelUUID(uuid string) {
	this.locker.Lock()
	defer this.locker.Unlock()

	delete(this.data, uuid)
}
