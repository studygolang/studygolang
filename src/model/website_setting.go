// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
)

type DocMenu struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type FriendLogo struct {
	Image  string `json:"image"`
	Url    string `json:"url"`
	Name   string `json:"name"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

type FooterNav struct {
	Name      string `json:"name"`
	Url       string `json:"url"`
	OuterSite bool   `json:"outer_site"`
}

type WebsiteSetting struct {
	Id          int `xorm:"pk autoincr"`
	Name        string
	Domain      string
	TitleSuffix string
	Favicon     string
	Logo        string
	StartYear   int
	BlogUrl     string
	ReadingMenu string
	DocsMenu    string
	Slogan      string
	Beian       string
	FooterNav   string
	FriendsLogo string
	CreatedAt   time.Time `xorm:"created"`
	UpdatedAt   time.Time `xorm:"<-"`

	DocMenus    []*DocMenu    `xorm:"-"`
	FriendLogos []*FriendLogo `xorm:"-"`
	FooterNavs  []*FooterNav  `xorm:"-"`
}

func (this *WebsiteSetting) AfterSet(name string, cell xorm.Cell) {
	if name == "docs_menu" {
		this.DocMenus = this.unmarshalDocsMenu()
	} else if name == "friends_logo" {
		this.FriendLogos = this.unmarshalFriendsLogo()
	}
}

func (this *WebsiteSetting) unmarshalDocsMenu() []*DocMenu {
	if this.DocsMenu == "" {
		return nil
	}

	var docMenus = []*DocMenu{}
	err := json.Unmarshal([]byte(this.DocsMenu), &docMenus)
	if err != nil {
		fmt.Println("unmarshal docs menu error:", err)
		return nil
	}

	return docMenus
}

func (this *WebsiteSetting) unmarshalFriendsLogo() []*FriendLogo {
	if this.FriendsLogo == "" {
		return nil
	}

	var friendLogos = []*FriendLogo{}
	err := json.Unmarshal([]byte(this.FriendsLogo), &friendLogos)
	if err != nil {
		fmt.Println("unmarshal friends logo error:", err)
		return nil
	}

	return friendLogos
}

func (this *WebsiteSetting) unmarshalFooterNav() []*FooterNav {
	var footerNavs = []*FooterNav{}
	err := json.Unmarshal([]byte(this.FooterNav), &footerNavs)
	if err != nil {
		fmt.Println("unmarshal footer nav error:", err)
		return nil
	}

	for _, footerNav := range footerNavs {
		if strings.HasPrefix(footerNav.Url, "/") {
			footerNav.OuterSite = false
		} else {
			footerNav.OuterSite = true
		}
	}

	return footerNavs
}
