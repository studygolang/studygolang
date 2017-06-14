// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"model"
	"strconv"
	"strings"

	"github.com/polaris1119/times"
	"golang.org/x/net/context"
)

type IndexLogic struct{}

var DefaultIndex = IndexLogic{}

func (IndexLogic) FindData(ctx context.Context, tab string) map[string]interface{} {
	indexNav := GetCurIndexNav(tab)
	data := map[string]interface{}{
		"tab":        tab,
		"index_navs": WebsiteSetting.IndexNavs,
		"cur_nav":    indexNav,
	}

	isNid := false
	nid, err := strconv.Atoi(indexNav.DataSource)
	if err == nil {
		isNid = true
	}

	switch {
	case indexNav.DataSource == "feed":
		topFeeds := DefaultFeed.FindTop(ctx)
		feeds := DefaultFeed.FindRecent(ctx, 50)
		data["feeds"] = append(topFeeds, feeds...)
	case isNid:
		paginator := NewPaginator(1)

		node := GetNode(nid)
		if node["pid"].(int) == 0 {
			nids := GetChildrenNode(nid, 10)
			questions := strings.TrimSuffix(strings.Repeat("?,", len(nids)), ",")
			querystring := "nid in(" + questions + ")"

			data["topics"] = DefaultTopic.FindAll(ctx, paginator, "topics.mtime DESC", querystring, nids...)
		} else {
			querystring := "nid=?"
			data["topics"] = DefaultTopic.FindAll(ctx, paginator, "topics.mtime DESC", querystring, nid)
		}
	case strings.Contains(indexNav.DataSource, ","):
		dsSlice := strings.Split(indexNav.DataSource, ",")
		nids := make([]interface{}, 0, len(dsSlice))
		tags := make([]string, 0, len(dsSlice))
		for _, d := range dsSlice {
			if nid, err := strconv.Atoi(d); err == nil {
				nids = append(nids, nid)
			} else {
				// 是 tag
				tags = append(tags, d)
			}
		}

		hasData := false
		if len(nids) > 0 {
			questions := strings.TrimSuffix(strings.Repeat("?,", len(nids)), ",")
			querystring := "nid in(" + questions + ")"
			paginator := NewPaginator(1)
			topics := DefaultTopic.FindAll(ctx, paginator, "topics.mtime DESC", querystring, nids...)
			if len(topics) > 0 {
				hasData = true
			}
			data["topics"] = topics
		}

		if !hasData && len(tags) > 0 {
			respBody, err := DefaultSearcher.SearchByField("title", strings.Join(tags, " "), 0, 50)
			if err != nil {
				break
			}
			users, nodes := DefaultSearcher.FillNodeAndUser(ctx, respBody)
			if respBody.NumFound == 0 {
				break
			}

			data["docs"] = respBody.Docs
			data["users"] = users
			data["nodes"] = nodes
		}
	case indexNav.DataSource == "rank":
		articles := DefaultRank.FindDayRank(ctx, model.TypeArticle, times.Format("ymd"), 25)
		articleNum := 0
		if articles != nil {
			articleNum = len(articles.([]*model.Article))
		}
		data["articles"] = articles
		data["topics"] = DefaultRank.FindDayRank(ctx, model.TypeTopic, times.Format("ymd"), 50-articleNum, true)

		newIndexNav := &model.IndexNav{
			Tab:        indexNav.Tab,
			Name:       indexNav.Name,
			DataSource: indexNav.DataSource,
		}

		hotNodes := DefaultTopic.FindHotNodes(ctx)
		newIndexNav.Children = make([]*model.IndexNavChild, len(hotNodes))
		for i, hotNode := range hotNodes {
			newIndexNav.Children[i] = &model.IndexNavChild{
				Uri:  "/go/" + hotNode["ename"].(string),
				Name: hotNode["name"].(string),
			}
		}

		data["cur_nav"] = newIndexNav
	case indexNav.DataSource == "article":
		data["articles"] = DefaultArticle.FindBy(ctx, 50)
	}

	return data
}
