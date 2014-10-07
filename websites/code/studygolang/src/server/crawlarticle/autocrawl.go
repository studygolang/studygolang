// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package main

import (
	"fmt"
	"regexp"
	"strings"

	"config"
	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron"
	"logger"
	"service"
	"util"
)

var websites = map[string]map[string]string{
	"cnblogs": {
		"all_url":        "http://zzk.cnblogs.com/s?t=b&w=",
		"incr_url":       "http://zzk.cnblogs.com/s?t=b&dateRange=One-Week&w=",
		"keywords":       "golang,go语言",
		"listselector":   "#searchResult .searchItem", // 搜索结果一项的选择器
		"resultselector": "h3 a",
		"page_field":     "p",
		"max_page":       "30",
	},
	"csdn": {
		"all_url":        "http://so.csdn.net/so/search/s.do?t=blog&q=",
		"incr_url":       "http://so.csdn.net/so/search/s.do?t=blog&q=",
		"keywords":       "go,golang,golang语言,go语言",
		"listselector":   ".search-list",
		"resultselector": "dt a",
		"page_field":     "p",
		"max_page":       "13",
	},
	"oschina": {
		"all_url":        "http://www.oschina.net/search?scope=blog&q=",
		"incr_url":       "http://www.oschina.net/search?scope=blog&sort_by_time=1&q=",
		"keywords":       "go,golang",
		"listselector":   "#results li",
		"resultselector": "h3 a",
		"page_field":     "p",
		"max_page":       "50",
	},
	"oschina_translate": {
		"all_url":        "http://www.oschina.net/search?scope=translate&q=",
		"incr_url":       "http://www.oschina.net/search?scope=translate&sort_by_time=1&q=",
		"keywords":       "go,golang",
		"listselector":   "#results li",
		"resultselector": "h3 a",
		"page_field":     "p",
		"max_page":       "50",
	},
	"iteye": {
		"all_url":        "http://www.iteye.com/search?type=blog&query=",
		"incr_url":       "http://www.iteye.com/search?type=blog&sort=created_at&query=",
		"keywords":       "go语言,golang",
		"listselector":   "#search_result .topic",
		"resultselector": ".content h4 a",
		"page_field":     "page",
		"max_page":       "20",
	},
	"iteye_news": {
		"all_url":        "http://www.iteye.com/search?type=news&query=",
		"incr_url":       "http://www.iteye.com/search?type=news&sort=created_at&query=",
		"keywords":       "go语言,golang",
		"listselector":   "#search_result .topic",
		"resultselector": ".content h4 a",
		"page_field":     "page",
		"max_page":       "20",
	},
}

const pattern = "go|golang|goroutine|channel/i"

func autocrawl(needAll bool) {

	if needAll {
		// 全量
		for website, wbconf := range websites {
			logger.Infoln("all crawl", website)
			go doCrawl(wbconf, true)
		}
	}

	// 定时增量
	c := cron.New()
	c.AddFunc(config.Config["crawl_spec"], startCrawl)
	c.Start()
}

func startCrawl() {

	for website, wbconf := range websites {
		logger.Infoln("do crawl", website)
		go doCrawl(wbconf, false)
	}
}

func doCrawl(wbconf map[string]string, isAll bool) {
	crawlUrl := wbconf["incr_url"]
	if isAll {
		crawlUrl = wbconf["all_url"]
	}

	keywords := strings.Split(wbconf["keywords"], ",")
	listselector := wbconf["listselector"]
	resultselector := wbconf["resultselector"]
	pageField := wbconf["page_field"]

	maxPage := 1
	if isAll {
		maxPage = util.MustInt(wbconf["max_page"])
	}

	var (
		doc *goquery.Document
		err error
	)

	for _, keyword := range keywords {
		for p := 1; p <= maxPage; p++ {

			page := fmt.Sprintf("&%s=%d", pageField, p)
			logger.Infoln("parse url:", crawlUrl+keyword+page)
			if doc, err = goquery.NewDocument(crawlUrl + keyword + page); err != nil {
				break
			}

			doc.Find(listselector).Each(func(i int, contentSelection *goquery.Selection) {

				aSelection := contentSelection.Find(resultselector)
				title := aSelection.Text()
				matched, err := regexp.MatchString(pattern, title)
				if err != nil {
					logger.Errorln(err)
					return
				}

				if !matched {
					return
				}

				articleUrl, ok := aSelection.Attr("href")
				if ok {
					service.ParseArticle(articleUrl, true)
				}
			})
		}
	}
}
