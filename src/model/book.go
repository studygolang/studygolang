// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

const (
	IsFreeFalse = iota
	IsFreeTrue
)

type Book struct {
	Id          int       `json:"id" xorm:"pk autoincr"`
	Name        string    `json:"name"`
	Ename       string    `json:"ename"`
	Cover       string    `json:"cover"`
	Author      string    `json:"author"`
	Translator  string    `json:"translator"`
	Lang        int       `json:"lang"`
	PubDate     string    `json:"pub_date"`
	Desc        string    `json:"desc"`
	Catalogue   string    `json:"catalogue"`
	IsFree      bool      `json:"is_free"`
	OnlineUrl   string    `json:"online_url"`
	DownloadUrl string    `json:"download_url"`
	BuyUrl      string    `json:"buy_url"`
	Price       float32   `json:"price"`
	Viewnum     int       `json:"viewnum"`
	Cmtnum      int       `json:"cmtnum"`
	Likenum     int       `json:"likenum"`
	CreatedAt   OftenTime `json:"created_at" xorm:"created"`
	UpdatedAt   OftenTime `json:"updated_at" xorm:"<-"`
}

func (*Book) TableName() string {
	return "book"
}
