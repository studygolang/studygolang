// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"config"
	"logger"
	"model"
	"util"
)

const MaxRows = 100

// 准备索引数据，post 给 solr
// isAll: 是否全量
func Indexing(isAll bool) {
	IndexingArticle(isAll)
	IndexingTopic(isAll)
	IndexingResource(isAll)
}

// 索引博文
func IndexingArticle(isAll bool) {
	solrClient := NewSolrClient()

	articleObj := model.NewArticle()

	limit := strconv.Itoa(MaxRows)
	if isAll {
		id := 0
		for {
			articleList, err := articleObj.Where("id>? AND status!=?", id, model.StatusOffline).Limit(limit).FindAll()
			if err != nil {
				logger.Errorln("IndexingArticle error:", err)
				break
			}

			if len(articleList) == 0 {
				break
			}

			for _, article := range articleList {
				if id < article.Id {
					id = article.Id
				}

				document := model.NewDocument(article, nil)
				addCommand := model.NewDefaultArgsAddCommand(document)

				solrClient.Push(addCommand)
			}

			solrClient.Post()
		}
	}
}

// 索引帖子
func IndexingTopic(isAll bool) {
	solrClient := NewSolrClient()

	topicObj := model.NewTopic()
	topicExObj := model.NewTopicEx()

	limit := strconv.Itoa(MaxRows)
	if isAll {
		id := 0
		for {
			topicList, err := topicObj.Where("tid>?", id).Limit(limit).FindAll()
			if err != nil {
				logger.Errorln("IndexingTopic error:", err)
				break
			}

			if len(topicList) == 0 {
				break
			}

			tids := util.Models2Intslice(topicList, "Tid")

			tmpStr := strings.Repeat("?,", len(tids))
			query := "tid in(" + tmpStr[:len(tmpStr)-1] + ")"
			args := make([]interface{}, len(tids))
			for i, tid := range tids {
				args[i] = tid
			}

			topicExList, err := topicExObj.Where(query, args...).FindAll()
			if err != nil {
				logger.Errorln("IndexingTopic error:", err)
				break
			}

			topicExMap := make(map[int]*model.TopicEx, len(topicExList))
			for _, topicEx := range topicExList {
				topicExMap[topicEx.Tid] = topicEx
			}

			for _, topic := range topicList {
				if id < topic.Tid {
					id = topic.Tid
				}

				topicEx, _ := topicExMap[topic.Tid]

				document := model.NewDocument(topic, topicEx)
				addCommand := model.NewDefaultArgsAddCommand(document)

				solrClient.Push(addCommand)
			}

			solrClient.Post()
		}
	}
}

// 索引资源
func IndexingResource(isAll bool) {
	solrClient := NewSolrClient()

	resourceObj := model.NewResource()
	resourceExObj := model.NewResourceEx()

	limit := strconv.Itoa(MaxRows)
	if isAll {
		id := 0
		for {
			resourceList, err := resourceObj.Where("id>?", id).Limit(limit).FindAll()
			if err != nil {
				logger.Errorln("IndexingResource error:", err)
				break
			}

			if len(resourceList) == 0 {
				break
			}

			ids := util.Models2Intslice(resourceList, "Id")

			tmpStr := strings.Repeat("?,", len(ids))
			query := "id in(" + tmpStr[:len(tmpStr)-1] + ")"
			args := make([]interface{}, len(ids))
			for i, rid := range ids {
				args[i] = rid
			}

			resourceExList, err := resourceExObj.Where(query, args...).FindAll()
			if err != nil {
				logger.Errorln("IndexingResource error:", err)
				break
			}

			resourceExMap := make(map[int]*model.ResourceEx, len(resourceExList))
			for _, resourceEx := range resourceExList {
				resourceExMap[resourceEx.Id] = resourceEx
			}

			for _, resource := range resourceList {
				if id < resource.Id {
					id = resource.Id
				}

				resourceEx, _ := resourceExMap[resource.Id]

				document := model.NewDocument(resource, resourceEx)
				addCommand := model.NewDefaultArgsAddCommand(document)

				solrClient.Push(addCommand)
			}

			solrClient.Post()
		}
	}
}

type SolrClient struct {
	addCommands []*model.AddCommand
}

func NewSolrClient() *SolrClient {
	return &SolrClient{
		addCommands: make([]*model.AddCommand, 0, MaxRows),
	}
}

func (this *SolrClient) Push(addCommand *model.AddCommand) {
	this.addCommands = append(this.addCommands, addCommand)
}

func (this *SolrClient) Post() error {
	stringBuilder := util.NewBuffer().Append("{")

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

		stringBuilder.Append(`"add":`).AppendBytes(commandJson)
	}

	if stringBuilder.Len() == 1 {
		logger.Errorln("post docs:no right addcommand")
		return errors.New("no right addcommand")
	}

	stringBuilder.Append("}")

	resp, err := http.Post(config.Config["engine_url"]+"/update?wt=json&commit=true", "application/json", stringBuilder)
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

const ContentLen = 350

func DoSearch(q, field string, start, rows int) (*model.ResponseBody, error) {
	selectUrl := config.Config["engine_url"] + "/select?"

	var values = url.Values{
		"wt":             []string{"json"},
		"hl":             []string{"true"},
		"hl.fl":          []string{"title,content"},
		"hl.simple.pre":  []string{"<em>"},
		"hl.simple.post": []string{"</em>"},
		"hl.fragsize":    []string{strconv.Itoa(ContentLen)},
		"start":          []string{strconv.Itoa(start)},
		"rows":           []string{strconv.Itoa(rows)},
	}

	if q == "" {
		values.Add("q", "*:*")
	} else {
		searchStat := model.NewSearchStat()
		searchStat.Where("keyword=?", q).Find()
		if searchStat.Id > 0 {
			searchStat.Where("keyword=?", q).Increment("times", 1)
		} else {
			searchStat.Keyword = q
			searchStat.Times = 1
			_, err := searchStat.Insert()
			if err != nil {
				searchStat.Where("keyword=?", q).Increment("times", 1)
			}
		}
	}

	if field == "text" {
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
			values.Add("q", "title:"+q+"^2"+" OR content:"+q+"^0.2")
		}
	}

	resp, err := http.Get(selectUrl + values.Encode())
	if err != nil {
		logger.Errorln("search error:", err)
		return nil, err
	}

	defer resp.Body.Close()

	var searchResponse model.SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResponse)
	if err != nil {
		logger.Errorln("parse response error:", err)
		return nil, err
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
				if maxLen > ContentLen {
					maxLen = ContentLen
				}
				doc.HlContent = doc.Content[:maxLen]
			}

			doc.HlContent += "..."
		}
	}

	return searchResponse.RespBody, nil
}
