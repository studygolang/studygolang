// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package logic

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	. "db"

	"github.com/polaris1119/config"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"

	"model"
)

type SearcherLogic struct {
	maxRows int

	engineUrl string
}

var DefaultSearcher = SearcherLogic{maxRows: 100, engineUrl: config.ConfigFile.MustValue("search", "engine_url")}

// 准备索引数据，post 给 solr
// isAll: 是否全量
func (self SearcherLogic) Indexing(isAll bool) {
	self.IndexingArticle(isAll)
	self.IndexingTopic(isAll)
	self.IndexingResource(isAll)
}

// IndexingArticle 索引博文
func (self SearcherLogic) IndexingArticle(isAll bool) {
	// solrClient := NewSolrClient()

	// articleObj := model.NewArticle()

	// limit := strconv.Itoa(self.maxRows)
	// if isAll {
	// 	id := 0
	// 	for {
	// 		articleList, err := articleObj.Where("id>? AND status!=?", id, model.StatusOffline).Limit(limit).FindAll()
	// 		if err != nil {
	// 			logger.Errorln("IndexingArticle error:", err)
	// 			break
	// 		}

	// 		if len(articleList) == 0 {
	// 			break
	// 		}

	// 		for _, article := range articleList {
	// 			if id < article.Id {
	// 				id = article.Id
	// 			}

	// 			document := model.NewDocument(article, nil)
	// 			addCommand := model.NewDefaultArgsAddCommand(document)

	// 			solrClient.Push(addCommand)
	// 		}

	// 		solrClient.Post()
	// 	}
	// }
}

// 索引帖子
func (SearcherLogic) IndexingTopic(isAll bool) {
	// solrClient := NewSolrClient()

	// topicObj := model.NewTopic()
	// topicExObj := model.NewTopicEx()

	// limit := strconv.Itoa(MaxRows)
	// if isAll {
	// 	id := 0
	// 	for {
	// 		topicList, err := topicObj.Where("tid>?", id).Limit(limit).FindAll()
	// 		if err != nil {
	// 			logger.Errorln("IndexingTopic error:", err)
	// 			break
	// 		}

	// 		if len(topicList) == 0 {
	// 			break
	// 		}

	// 		tids := util.Models2Intslice(topicList, "Tid")

	// 		tmpStr := strings.Repeat("?,", len(tids))
	// 		query := "tid in(" + tmpStr[:len(tmpStr)-1] + ")"
	// 		args := make([]interface{}, len(tids))
	// 		for i, tid := range tids {
	// 			args[i] = tid
	// 		}

	// 		topicExList, err := topicExObj.Where(query, args...).FindAll()
	// 		if err != nil {
	// 			logger.Errorln("IndexingTopic error:", err)
	// 			break
	// 		}

	// 		topicExMap := make(map[int]*model.TopicEx, len(topicExList))
	// 		for _, topicEx := range topicExList {
	// 			topicExMap[topicEx.Tid] = topicEx
	// 		}

	// 		for _, topic := range topicList {
	// 			if id < topic.Tid {
	// 				id = topic.Tid
	// 			}

	// 			topicEx, _ := topicExMap[topic.Tid]

	// 			document := model.NewDocument(topic, topicEx)
	// 			addCommand := model.NewDefaultArgsAddCommand(document)

	// 			solrClient.Push(addCommand)
	// 		}

	// 		solrClient.Post()
	// 	}
	// }
}

// 索引资源
func (SearcherLogic) IndexingResource(isAll bool) {
	// solrClient := NewSolrClient()

	// resourceObj := model.NewResource()
	// resourceExObj := model.NewResourceEx()

	// limit := strconv.Itoa(MaxRows)
	// if isAll {
	// 	id := 0
	// 	for {
	// 		resourceList, err := resourceObj.Where("id>?", id).Limit(limit).FindAll()
	// 		if err != nil {
	// 			logger.Errorln("IndexingResource error:", err)
	// 			break
	// 		}

	// 		if len(resourceList) == 0 {
	// 			break
	// 		}

	// 		ids := util.Models2Intslice(resourceList, "Id")

	// 		tmpStr := strings.Repeat("?,", len(ids))
	// 		query := "id in(" + tmpStr[:len(tmpStr)-1] + ")"
	// 		args := make([]interface{}, len(ids))
	// 		for i, rid := range ids {
	// 			args[i] = rid
	// 		}

	// 		resourceExList, err := resourceExObj.Where(query, args...).FindAll()
	// 		if err != nil {
	// 			logger.Errorln("IndexingResource error:", err)
	// 			break
	// 		}

	// 		resourceExMap := make(map[int]*model.ResourceEx, len(resourceExList))
	// 		for _, resourceEx := range resourceExList {
	// 			resourceExMap[resourceEx.Id] = resourceEx
	// 		}

	// 		for _, resource := range resourceList {
	// 			if id < resource.Id {
	// 				id = resource.Id
	// 			}

	// 			resourceEx, _ := resourceExMap[resource.Id]

	// 			document := model.NewDocument(resource, resourceEx)
	// 			addCommand := model.NewDefaultArgsAddCommand(document)

	// 			solrClient.Push(addCommand)
	// 		}

	// 		solrClient.Post()
	// 	}
	// }
}

const searchContentLen = 350

// DoSearch 搜索
func (self SearcherLogic) DoSearch(q, field string, start, rows int) (*model.ResponseBody, error) {
	selectUrl := self.engineUrl + "/select?"

	var values = url.Values{
		"wt":             []string{"json"},
		"hl":             []string{"true"},
		"hl.fl":          []string{"title,content"},
		"hl.simple.pre":  []string{"<em>"},
		"hl.simple.post": []string{"</em>"},
		"hl.fragsize":    []string{strconv.Itoa(searchContentLen)},
		"start":          []string{strconv.Itoa(start)},
		"rows":           []string{strconv.Itoa(rows)},
	}

	if q == "" {
		values.Add("q", "*:*")
	} else {
		searchStat := &model.SearchStat{}
		MasterDB.Where("keyword=?", q).Get(searchStat)
		if searchStat.Id > 0 {
			MasterDB.Where("keyword=?", q).Incr("times", 1).Update(new(model.SearchStat))
		} else {
			searchStat.Keyword = q
			searchStat.Times = 1
			_, err := MasterDB.Insert(searchStat)
			if err != nil {
				MasterDB.Where("keyword=?", q).Incr("times", 1).Update(new(model.SearchStat))
			}
		}
	}

	isTag := false
	// TODO: 目前大部分都没有tag，因此，对tag特殊处理
	if field == "text" || field == "tag" {
		if field == "tag" {
			isTag = true
		}
		field = ""
	}

	if field != "" {
		values.Add("df", field)
		if q != "" {
			values.Add("q", q)
		}
	} else {
		// 全文检索
		if q != "" {
			if isTag {
				values.Add("q", "title:"+q+"^2"+" OR tags:"+q+"^4 OR content:"+q+"^0.2")
			} else {
				values.Add("q", "title:"+q+"^2"+" OR content:"+q+"^0.2")
			}
		}
	}

	resp, err := http.Get(selectUrl + values.Encode())
	if err != nil {
		logger.Errorln("search error:", err)
		return &model.ResponseBody{}, err
	}

	defer resp.Body.Close()

	var searchResponse model.SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResponse)
	if err != nil {
		logger.Errorln("parse response error:", err)
		return &model.ResponseBody{}, err
	}

	if len(searchResponse.Highlight) > 0 {
		for _, doc := range searchResponse.RespBody.Docs {
			highlighting, ok := searchResponse.Highlight[doc.Id]
			if ok {
				if len(highlighting.Title) > 0 {
					doc.HlTitle = highlighting.Title[0]
				}

				if len(highlighting.Content) > 0 {
					doc.HlContent = highlighting.Content[0]
				}
			}

			if doc.HlTitle == "" {
				doc.HlTitle = doc.Title
			}

			if doc.HlContent == "" && doc.Content != "" {
				maxLen := len(doc.Content) - 1
				if maxLen > searchContentLen {
					maxLen = searchContentLen
				}
				doc.HlContent = doc.Content[:maxLen]
			}

			doc.HlContent += "..."
		}
	}

	return searchResponse.RespBody, nil
}

type SolrClient struct {
	addCommands []*model.AddCommand
}

func NewSolrClient() *SolrClient {
	return &SolrClient{
		addCommands: make([]*model.AddCommand, 0, 100),
	}
}

func (this *SolrClient) Push(addCommand *model.AddCommand) {
	this.addCommands = append(this.addCommands, addCommand)
}

func (this *SolrClient) Post() error {
	stringBuilder := goutils.NewBuffer().Append("{")

	needComma := false
	for _, addCommand := range this.addCommands {
		commandJson, err := json.Marshal(addCommand)
		if err != nil {
			continue
		}

		if stringBuilder.Len() == 1 {
			needComma = false
		} else {
			needComma = true
		}

		if needComma {
			stringBuilder.Append(",")
		}

		stringBuilder.Append(`"add":`).Append(commandJson)
	}

	if stringBuilder.Len() == 1 {
		logger.Errorln("post docs:no right addcommand")
		return errors.New("no right addcommand")
	}

	stringBuilder.Append("}")

	resp, err := http.Post(config.ConfigFile.MustValue("search", "engine_url")+"/update?wt=json&commit=true", "application/json", stringBuilder)
	if err != nil {
		logger.Errorln("post error:", err)
		return err
	}

	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		logger.Errorln("parse response error:", err)
		return err
	}

	return nil
}
