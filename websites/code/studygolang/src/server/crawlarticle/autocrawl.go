// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"config"
	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron"
	"logger"
	"service"
	"util"
)

var websites = make(map[string]map[string]string)

const pattern = "(?i)go|golang|goroutine|channel"

func autocrawl(needAll bool, crawlConfFile string, whichSite string) {

	_, err := config.ParseConfig(crawlConfFile, &websites)
	if err != nil {
		log.Fatalln("parse crawl config error:", err)
	}

	if needAll {
		// 全量
		for website, wbconf := range websites {
			if whichSite != "" && whichSite != website {
				continue
			}

			logger.Infoln("all crawl", website)
			go doCrawl(wbconf, true)
		}
	}

	// 定时增量
	c := cron.New()
	c.AddFunc(config.Config["crawl_spec"], func() {
		for website, wbconf := range websites {
			if whichSite != "" && whichSite != website {
				continue
			}

			logger.Infoln("do crawl", website)
			go doCrawl(wbconf, false)
		}
	})
	c.Start()
}

func doCrawl(wbconf map[string]string, isAll bool) {
	crawlUrl := wbconf["incr_url"]
	if isAll {
		crawlUrl = wbconf["all_url"]
	}

	listselector := wbconf["listselector"]
	resultselector := wbconf["resultselector"]
	pageField := wbconf["page_field"]

	maxPage := 1
	if isAll {
		maxPage = util.MustInt(wbconf["max_page"])
	}

	// 个人博客，一般通过 tag 方式获取，这种处理方式和搜索不一样
	if wbconf["keywords"] == "" {
		for p := maxPage; p >= 1; p-- {
			if pageField == "" {

				// 标题不包含 go 等关键词的，也入库
				if err := parseArticleList(crawlUrl+strconv.Itoa(p), listselector, resultselector, false); err != nil {
					break
				}
			}
		}

		return
	}

	keywords := strings.Split(wbconf["keywords"], ",")

	for _, keyword := range keywords {
		for p := 1; p <= maxPage; p++ {

			page := fmt.Sprintf("&%s=%d", pageField, p)
			if err := parseArticleList(crawlUrl+keyword+page, listselector, resultselector, true); err != nil {
				logger.Errorln("parse article url error:", err)
				break
			}
		}
	}
}

func parseArticleList(url, listselector, resultselector string, isAuto bool) (err error) {

	logger.Infoln("parse url:", url)

	var doc *goquery.Document

	if doc, err = goquery.NewDocument(url); err != nil {
		return
	}

	doc.Find(listselector).Each(func(i int, contentSelection *goquery.Selection) {

		aSelection := contentSelection.Find(resultselector)

		if isAuto {
			title := aSelection.Text()

			matched, err := regexp.MatchString(pattern, title)
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
			service.ParseArticle(articleUrl, isAuto)
		}
	})

	return
}
