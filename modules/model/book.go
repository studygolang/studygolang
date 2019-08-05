// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import "time"

const (
	IsFreeFalse = iota
	IsFreeTrue
)

type Book struct {
	Id            int       `json:"id" xorm:"pk autoincr"`
	Name          string    `json:"name"`
	Ename         string    `json:"ename"`
	Cover         string    `json:"cover"`
	Author        string    `json:"author"`
	Translator    string    `json:"translator"`
	Lang          int       `json:"lang"`
	PubDate       string    `json:"pub_date"`
	Desc          string    `json:"desc"`
	Tags          string    `json:"tags"`
	Catalogue     string    `json:"catalogue"`
	IsFree        bool      `json:"is_free"`
	OnlineUrl     string    `json:"online_url"`
	DownloadUrl   string    `json:"download_url"`
	BuyUrl        string    `json:"buy_url"`
	Price         float32   `json:"price"`
	Lastreplyuid  int       `json:"lastreplyuid"`
	Lastreplytime OftenTime `json:"lastreplytime"`
	Viewnum       int       `json:"viewnum"`
	Cmtnum        int       `json:"cmtnum"`
	Likenum       int       `json:"likenum"`
	Uid           int       `json:"uid"`
	CreatedAt     OftenTime `json:"created_at" xorm:"created"`
	UpdatedAt     OftenTime `json:"updated_at" xorm:"<-"`

	// 排行榜阅读量
	RankView int `json:"rank_view" xorm:"-"`
}

func (this *Book) AfterInsert() {
	go func() {
		// AfterInsert 时，自增 ID 还未赋值，这里 sleep 一会，确保自增 ID 有值
		for {
			if this.Id > 0 {
				PublishFeed(this, nil)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}
