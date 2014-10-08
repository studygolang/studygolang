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
	"strings"

	"config"
	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron"
	"logger"
	"service"
	"util"
)

var websites = make(map[string]map[string]string)

const pattern = "go|golang|goroutine|channel/i"

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
