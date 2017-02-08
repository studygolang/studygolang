// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	. "db"
	"errors"
	"fmt"
	"model"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/polaris1119/logger"
	"golang.org/x/net/context"
)

const titlePattern = "(?i)go|golang|goroutine|channel"

type AutoCrawlLogic struct{}

var DefaultAutoCrawl = AutoCrawlLogic{}

func (self AutoCrawlLogic) DoCrawl(isAll bool) error {
	autoCrawlConfList := make([]*model.AutoCrawlRule, 0)
	err := MasterDB.Where("status=?", model.AutoCrawlOn).Find(&autoCrawlConfList)
	if err != nil {
		logger.Errorln("ArticleLogic FindBy Error:", err)
		return err
	}

	for _, autoCrawlConf := range autoCrawlConfList {
		self.crawlOneWebsite(autoCrawlConf, isAll)
	}

	return nil
}

// 通过网站标识抓取
func (self AutoCrawlLogic) CrawlWebsite(website string, isAll bool) error {
	autoCrawlConf := &model.AutoCrawlRule{}
	_, err := MasterDB.Where("website=?", website).Get(autoCrawlConf)
	if err != nil {
		return err
	}

	if autoCrawlConf.Id == 0 {
		return errors.New("the website is not exists in auto crawl rule.")
	}

	go self.crawlOneWebsite(autoCrawlConf, isAll)

	return nil
}

func (self AutoCrawlLogic) crawlOneWebsite(autoCrawlConf *model.AutoCrawlRule, isAll bool) {
	maxPage := 1
	crawlUrl := autoCrawlConf.IncrUrl
	if isAll {
		crawlUrl = autoCrawlConf.AllUrl
		maxPage = autoCrawlConf.MaxPage
	}

	pageField := autoCrawlConf.PageField

	// 个人博客，一般通过 tag 方式获取，这种处理方式和搜索不一样
	if autoCrawlConf.Keywords == "" {
		for p := maxPage; p >= 1; p-- {
			curUrl := ""

			if pageField == "" {
				if p > 1 {
					curUrl += crawlUrl + "page/" + strconv.Itoa(p)
				}
			} else {
				page := fmt.Sprintf("?%s=%d", pageField, p)
				curUrl += crawlUrl + page
			}

			// 标题不包含 go 等关键词的，也入库
			if err := self.parseArticleList(curUrl, autoCrawlConf, false); err != nil {
				logger.Errorln("parse article url", curUrl, "error:", err)
				break
			}
		}
		return
	}

	keywords := strings.Split(autoCrawlConf.Keywords, ",")
	for _, keyword := range keywords {
		for p := 1; p <= maxPage; p++ {

			curUrl := ""
			page := fmt.Sprintf("&%s=%d", pageField, p)
			if strings.Contains(crawlUrl, "%s") {
				curUrl = fmt.Sprintf(crawlUrl, keyword) + page
			} else {
				curUrl = crawlUrl + keyword + page
			}

			if err := self.parseArticleList(curUrl, autoCrawlConf, true); err != nil {
				logger.Errorln("parse article url", curUrl, "error:", err)
				break
			}
		}
	}
}

func (self AutoCrawlLogic) parseArticleList(strUrl string, autoCrawlConf *model.AutoCrawlRule, isSearch bool) (err error) {

	logger.Infoln("parse url:", strUrl)

	var doc *goquery.Document

	extMap := autoCrawlConf.ParseExt()
	if extMap == nil {
		doc, err = goquery.NewDocument(strUrl)
	} else {
		req, err := http.NewRequest("GET", strUrl, nil)
		if err != nil {
			return err
		}
		if referer, ok := extMap["referer"]; ok {
			req.Header.Add("Referer", referer)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		doc, err = goquery.NewDocumentFromResponse(resp)
	}

	if err != nil {
		return
	}

	listSelector := autoCrawlConf.ListSelector
	resultSelector := autoCrawlConf.ResultSelector

	u, err := url.Parse(autoCrawlConf.IncrUrl)
	if err != nil {
		logger.Errorln("parse incr_url error:", err)
		return
	}
	host := u.Scheme + "://" + u.Host

	doc.Find(listSelector).Each(func(i int, contentSelection *goquery.Selection) {

		aSelection := contentSelection.Find(resultSelector)

		// 搜索时，避免搜到垃圾，对标题进一步判断
		if isSearch {
			title := aSelection.Text()

			matched, err := regexp.MatchString(titlePattern, title)
			if err != nil {
				logger.Errorln(err)
				return
			}

			if !matched {
				return
			}
		}

		articleUrl, ok := aSelection.Attr("href")
		if ok {
			pos := strings.LastIndex(articleUrl, "?")
			if pos != -1 {
				articleUrl = articleUrl[:pos]
			}

			if !strings.HasPrefix(articleUrl, "http") {
				articleUrl = host + articleUrl
			}
			DefaultArticle.ParseArticle(context.Background(), articleUrl, isSearch)
		}
	})

	return
}
