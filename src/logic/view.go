// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"net/http"
	"strconv"
	"sync"

	. "db"
	"model"

	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
)

// 话题/文章/资源的浏览数
// 避免每次写库，同时避免刷屏
type view struct {
	objtype int // 对象类型（model/comment 中的 type 常量）
	objid   int // 对象id（相应的表中的id）

	num    int // 当前浏览数
	locker sync.Mutex
}

func newView(objtype, objid int) *view {
	return &view{objtype: objtype, objid: objid}
}

func (this *view) incr() {
	this.locker.Lock()
	defer this.locker.Unlock()

	this.num++
}

// flush 将浏览数刷入数据库中
func (this *view) flush() {
	this.locker.Lock()
	defer this.locker.Unlock()

	session := MasterDB.Id(this.objid)
	switch this.objtype {
	case model.TypeTopic:
		session.Incr("view", this.num).Update(new(model.TopicEx))
	case model.TypeArticle:
		session.Incr("viewnum", this.num).Update(new(model.Article))
	case model.TypeResource:
		session.Incr("viewnum", this.num).Update(new(model.ResourceEx))
	case model.TypeProject:
		session.Incr("viewnum", this.num).Update(new(model.OpenProject))
	case model.TypeWiki:
		session.Incr("viewnum", this.num).Update(new(model.Wiki))
	}

	this.num = 0
}

// 保存所有对象的浏览数
type views struct {
	data map[string]*view
	// 记录用户是否已经看过（防止刷屏）
	users map[string]bool

	locker sync.Mutex
}

func newViews() *views {
	return &views{data: make(map[string]*view), users: make(map[string]bool)}
}

// TODO: 用户登录了，应该用用户标识，而不是IP
func (this *views) Incr(req *http.Request, objtype, objid int) {
	user := goutils.Ip2long(goutils.RemoteIp(req))

	key := strconv.Itoa(objtype) + strconv.Itoa(objid)

	this.locker.Lock()
	defer this.locker.Unlock()

	if user != 0 {
		userKey := key + strconv.FormatUint(uint64(user), 10)

		if _, ok := this.users[userKey]; ok {
			return
		} else {
			this.users[userKey] = true
		}
	}

	if _, ok := this.data[key]; !ok {
		this.data[key] = newView(objtype, objid)
	}

	this.data[key].incr()
}

func (this *views) Flush() {
	logger.Debugln("start views flush")
	this.locker.Lock()
	defer this.locker.Unlock()

	// TODO：量大时，考虑copy一份，然后异步 入库，以免堵塞 锁 太久
	for _, view := range this.data {
		view.flush()
	}

	this.data = make(map[string]*view)
	this.users = make(map[string]bool)

	logger.Debugln("end views flush")
}

var Views = newViews()
