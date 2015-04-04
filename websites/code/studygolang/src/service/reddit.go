// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com, http://golang.top
// Author：polaris	golang@studygolang.com

// 解析 http://www.reddit.com/r/golang 最新 Go 信息
package service

import (
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"logger"
	"model"
	"util"
)

const (
	Reddit       = "http://www.reddit.com"
	RedditGolang = "/r/golang/new/"
)

// 获取url对应的文章并根据规则进行解析
func ParseReddit(redditUrl string) error {
	redditUrl = strings.TrimSpace(redditUrl)
	if redditUrl == "" {
		redditUrl = Reddit + RedditGolang
	} else if !strings.HasPrefix(redditUrl, "http") {
		redditUrl = "http://" + redditUrl
	}

	var (
		doc *goquery.Document
		err error
	)

	if doc, err = goquery.NewDocument(redditUrl); err != nil {
		logger.Errorln("goquery reddit newdocument error:", err)
		return err
	}

	/*
		doc.Find("#siteTable .link").Each(func(i int, contentSelection *goquery.Selection) {

			err = dealRedditOneResource(contentSelection)

			if err != nil {
				logger.Errorln(err)
			}
		})
	*/

	// 最后面的先入库处理
	resourcesSelection := doc.Find("#siteTable .link")

	for i := resourcesSelection.Length() - 1; i >= 0; i-- {
		err = dealRedditOneResource(goquery.NewDocumentFromNode(resourcesSelection.Get(i)).Selection)

		if err != nil {
			logger.Errorln(err)
		}
	}

	return err
}

var PresetUids = []int{1, 1747, 1748, 1827}

var resourceRe = regexp.MustCompile(`\n\n`)

// 处理 Reddit 中的一条资源
func dealRedditOneResource(contentSelection *goquery.Selection) error {
	aSelection := contentSelection.Find(".title a.title")

	title := aSelection.Text()
	if title == "" {
		return errors.New("title is empty")
	}

	resourceUrl, ok := aSelection.Attr("href")
	if !ok || resourceUrl == "" {
		return errors.New("resource url is empty")
	}

	isReddit := false

	resource := model.NewResource()
	// Reddit 自身的内容
	if contentSelection.HasClass("self") {
		isReddit = true
		resourceUrl = Reddit + resourceUrl
	}

	err := resource.Where("url=?", resourceUrl).Find("id")
	// 已经存在
	if resource.Id != 0 {
		// 如果是 reddit 本身的，可以更新评论信息
		if !isReddit {
			return errors.New("url" + resourceUrl + "has exists!")
		}
	}

	if isReddit {

		resource.Form = model.ContentForm

		var doc *goquery.Document

		if doc, err = goquery.NewDocument(resourceUrl); err != nil {
			return errors.New("goquery reddit.com/r/golang self newdocument error:" + err.Error())
		}

		content, err := doc.Find("#siteTable .usertext .md").Html()
		if err != nil {
			return err
		}

		doc.Find(".commentarea .comment .usertext .md").Each(func(i int, contentSel *goquery.Selection) {
			if i == 0 {
				content += `<hr/>**评论：**<br/><br/>`
			}

			comment, err := contentSel.Html()
			if err != nil {
				return
			}

			comment = strings.TrimSpace(comment)
			comment = resourceRe.ReplaceAllLiteralString(comment, "\n")

			author := contentSel.ParentsFiltered(".usertext").Prev().Find(".author").Text()
			content += author + ": <pre>" + comment + "</pre>"
		})

		resource.Content = content

		// reddit 本身的，当做其他资源
		resource.Catid = 4
	} else {
		resource.Form = model.LinkForm

		// Github，是开源项目
		if contentSelection.Find(".title .domain a").Text() == "github.com" {
			resource.Catid = 2
		} else {
			resource.Catid = 1
		}
	}

	resource.Title = title
	resource.Url = resourceUrl
	resource.Uid = PresetUids[rand.Intn(4)]

	ctime := util.TimeNow()
	datetime, ok := contentSelection.Find(".tagline time").Attr("datetime")
	if ok {
		dtime, err := time.ParseInLocation(time.RFC3339, datetime, time.UTC)
		if err != nil {
			logger.Errorln("parse ctime error:", err)
		} else {
			ctime = dtime.Local().Format("2006-01-02 15:04:05")
		}
	}
	resource.Ctime = ctime

	if resource.Id == 0 {
		var id int64
		id, err = resource.Insert()

		if err != nil {
			return errors.New("insert into Resource error:" + err.Error())
		}

		// 存扩展信息
		resourceEx := model.NewResourceEx()
		resourceEx.Id = int(id)
		if _, err = resourceEx.Insert(); err != nil {
			return errors.New("insert into ResourceEx error:" + err.Error())
		}
	} else {
		if err = resource.Persist(resource); err != nil {
			return errors.New("persist resource:" + strconv.Itoa(resource.Id) + " error:" + err.Error())
		}
	}

	return nil
}
