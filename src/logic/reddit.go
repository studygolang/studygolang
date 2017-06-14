// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com, http://golang.top
// Author: polaris	polaris@studygolang.com

// 解析 http://www.reddit.com/r/golang 最新 Go 信息
package logic

import (
	"errors"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	. "db"
	"model"

	"github.com/PuerkitoBio/goquery"
	"github.com/polaris1119/config"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
)

type RedditLogic struct {
	domain string
	path   string
}

var DefaultReddit = newRedditLogic()

func newRedditLogic() *RedditLogic {
	return &RedditLogic{
		domain: "https://www.reddit.com",
		path:   config.ConfigFile.MustValue("crawl", "reddit_path"),
	}
}

// Parse 获取url对应的资源并根据规则进行解析
func (this *RedditLogic) Parse(redditUrl string) error {
	redditUrl = strings.TrimSpace(redditUrl)
	if redditUrl == "" {
		if this.path == "" {
			return nil
		}
		redditUrl = this.domain + this.path
	} else if !strings.HasPrefix(redditUrl, "https") {
		redditUrl = "https://" + redditUrl
	}

	var (
		doc *goquery.Document
		err error
	)

	// if doc, err = goquery.NewDocument(redditUrl); err != nil {
	if doc, err = this.newDocumentFromResp(redditUrl); err != nil {
		logger.Errorln("goquery reddit newdocument error:", err)
		return err
	}

	// 最后面的先入库处理
	resourcesSelection := doc.Find("#siteTable .link")

	for i := resourcesSelection.Length() - 1; i >= 0; i-- {
		err = this.dealRedditOneResource(goquery.NewDocumentFromNode(resourcesSelection.Get(i)).Selection)

		if err != nil {
			logger.Errorln(err)
		}
	}

	return err
}

func (this *RedditLogic) newDocumentFromResp(url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromResponse(resp)
}

var PresetUids = config.ConfigFile.MustValueArray("crawl", "preset_uids", ",")

var resourceRe = regexp.MustCompile(`\n\n`)

// 处理 Reddit 中的一条资源
func (this *RedditLogic) dealRedditOneResource(contentSelection *goquery.Selection) error {
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

	resource := &model.Resource{}
	// Reddit 自身的内容
	if contentSelection.HasClass("self") {
		isReddit = true
		resourceUrl = this.domain + resourceUrl
	}

	_, err := MasterDB.Where("url=?", resourceUrl).Get(resource)
	if err != nil {
		return err
	}
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
			return errors.New("goquery reddit.com" + this.path + " self newdocument error:" + err.Error())
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

		if strings.TrimSpace(content) == "" {
			return errors.New("goquery reddit.com" + this.path + " self newdocument(" + resourceUrl + ") error: content is empty")
		}

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
	resource.Uid = goutils.MustInt(PresetUids[rand.Intn(len(PresetUids))])

	ctime := time.Now()
	datetime, ok := contentSelection.Find(".tagline time").Attr("datetime")
	if ok {
		dtime, err := time.ParseInLocation(time.RFC3339, datetime, time.UTC)
		if err != nil {
			logger.Errorln("parse ctime error:", err)
		} else {
			ctime = dtime.Local()
		}
	}
	resource.Ctime = model.OftenTime(ctime)

	if resource.Id == 0 {
		session := MasterDB.NewSession()
		defer session.Close()
		session.Begin()

		_, err = session.Insert(resource)
		if err != nil {
			session.Rollback()
			return errors.New("insert into Resource error:" + err.Error())
		}

		// 存扩展信息
		resourceEx := &model.ResourceEx{}
		resourceEx.Id = resource.Id
		if _, err = session.Insert(resourceEx); err != nil {
			session.Rollback()
			return errors.New("insert into ResourceEx error:" + err.Error())
		}
		session.Commit()

		DefaultFeed.publish(resource, resourceEx)
	} else {
		if _, err = MasterDB.Id(resource.Id).Update(resource); err != nil {
			return errors.New("update resource:" + strconv.Itoa(resource.Id) + " error:" + err.Error())
		}
	}

	return nil
}
