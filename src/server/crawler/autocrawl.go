// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"logic"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/polaris1119/config"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"

	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron"
	"golang.org/x/net/context"
)

var websites = make(map[string]map[string]string)

const pattern = "(?i)go|golang|goroutine|channel"

func autocrawl(needAll bool, crawlConfFile string, whichSite string) {
	content, err := ioutil.ReadFile(crawlConfFile)
	if err != nil {
		log.Fatalln("parse crawl config read file error:", err)
	}

	err = json.Unmarshal(content, &websites)
	if err != nil {
		log.Fatalln("parse crawl config json parse error:", err)
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
	c.AddFunc(config.ConfigFile.MustValue("crawl", "spec", "0 0 */1 * * ?"), func() {
		// 抓取 reddit
		go logic.DefaultReddit.Parse("")

		// 抓取 www.oschina.net/project
		go logic.DefaultProject.ParseProjectList("http://www.oschina.net/project/lang/358/go?tag=0&os=0&sort=time")

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
		maxPage = goutils.MustInt(wbconf["max_page"])
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

			curUrl := ""

			page := fmt.Sprintf("&%s=%d", pageField, p)
			if strings.Contains(crawlUrl, "%s") {
				curUrl = fmt.Sprintf(crawlUrl, keyword) + page
			} else {
				curUrl = crawlUrl + keyword + page
			}

			if err := parseArticleList(curUrl, listselector, resultselector, true); err != nil {
				logger.Errorln("parse article url error:", err)
				break
			}
		}
	}
}

func parseArticleList(url, listselector, resultselector string, isAuto bool) (err error) {

	logger.Infoln("parse url:", url)

	var doc *goquery.Document

	if strings.Contains(url, "oschina.net") {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		req.Header.Add("Referer", "http://www.oschina.net/search?q=go&scope=blog&onlytitle=1&sort_by_time=1")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		doc, err = goquery.NewDocumentFromResponse(resp)
	} else {
		doc, err = goquery.NewDocument(url)
	}

	if err != nil {
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
			pos := strings.LastIndex(articleUrl, "?")
			if pos != -1 {
				articleUrl = articleUrl[:pos]
			}
			logic.DefaultArticle.ParseArticle(context.Background(), articleUrl, isAuto)
		}
	})

	return
}
