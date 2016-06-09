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
	"util"

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
	self.IndexingOpenProject(isAll)
}

// IndexingArticle 索引博文
func (self SearcherLogic) IndexingArticle(isAll bool) {
	solrClient := NewSolrClient()

	var (
		articleList []*model.Article
		err         error
	)

	if isAll {
		id := 0
		for {
			articleList = make([]*model.Article, 0)
			err = MasterDB.Where("id>?", id).Limit(self.maxRows).OrderBy("id ASC").Find(&articleList)
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
				if article.Status != model.ArticleStatusOffline {
					solrClient.PushAdd(model.NewDefaultArgsAddCommand(document))
				} else {
					solrClient.PushDel(model.NewDelCommand(document))
				}
			}

			solrClient.Post()
		}
	}
}

// 索引主题
func (self SearcherLogic) IndexingTopic(isAll bool) {
	solrClient := NewSolrClient()

	var (
		topicList   []*model.Topic
		topicExList map[int]*model.TopicEx

		err error
	)

	if isAll {
		id := 0
		for {
			topicList = make([]*model.Topic, 0)
			topicExList = make(map[int]*model.TopicEx)

			err = MasterDB.Where("tid>?", id).OrderBy("tid ASC").Limit(self.maxRows).Find(&topicList)
			if err != nil {
				logger.Errorln("IndexingTopic error:", err)
				break
			}

			if len(topicList) == 0 {
				break
			}

			tids := util.Models2Intslice(topicList, "Tid")

			err = MasterDB.In("tid", tids).Find(&topicExList)
			if err != nil {
				logger.Errorln("IndexingTopic error:", err)
				break
			}

			for _, topic := range topicList {
				if id < topic.Tid {
					id = topic.Tid
				}

				topicEx := topicExList[topic.Tid]

				document := model.NewDocument(topic, topicEx)
				addCommand := model.NewDefaultArgsAddCommand(document)

				solrClient.PushAdd(addCommand)
			}

			solrClient.Post()
		}
	}
}

// 索引资源
func (self SearcherLogic) IndexingResource(isAll bool) {
	solrClient := NewSolrClient()

	var (
		resourceList   []*model.Resource
		resourceExList map[int]*model.ResourceEx
		err            error
	)

	if isAll {
		id := 0
		for {
			resourceList = make([]*model.Resource, 0)
			resourceExList = make(map[int]*model.ResourceEx)

			err = MasterDB.Where("id>?", id).OrderBy("id ASC").Limit(self.maxRows).Find(&resourceList)
			if err != nil {
				logger.Errorln("IndexingResource error:", err)
				break
			}

			if len(resourceList) == 0 {
				break
			}

			ids := util.Models2Intslice(resourceList, "Id")

			err = MasterDB.In("id", ids).Find(&resourceExList)
			if err != nil {
				logger.Errorln("IndexingResource error:", err)
				break
			}

			for _, resource := range resourceList {
				if id < resource.Id {
					id = resource.Id
				}

				resourceEx := resourceExList[resource.Id]

				document := model.NewDocument(resource, resourceEx)
				addCommand := model.NewDefaultArgsAddCommand(document)

				solrClient.PushAdd(addCommand)
			}

			solrClient.Post()
		}
	}
}

// IndexingOpenProject 索引博文
func (self SearcherLogic) IndexingOpenProject(isAll bool) {
	solrClient := NewSolrClient()

	var (
		projectList []*model.OpenProject
		err         error
	)

	if isAll {
		id := 0
		for {
			projectList = make([]*model.OpenProject, 0)
			err = MasterDB.Where("id>?", id).OrderBy("id ASC").Limit(self.maxRows).Find(&projectList)
			if err != nil {
				logger.Errorln("IndexingArticle error:", err)
				break
			}

			if len(projectList) == 0 {
				break
			}

			for _, project := range projectList {
				if id < project.Id {
					id = project.Id
				}

				document := model.NewDocument(project, nil)
				if project.Status != model.ProjectStatusOffline {
					solrClient.PushAdd(model.NewDefaultArgsAddCommand(document))
				} else {
					solrClient.PushDel(model.NewDelCommand(document))
				}
			}

			solrClient.Post()
		}
	}
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
				utf8string := util.NewString(doc.Content)
				maxLen := utf8string.RuneCount() - 1
				if maxLen > searchContentLen {
					maxLen = searchContentLen
				}
				doc.HlContent = util.NewString(doc.Content).Slice(0, maxLen)
			}

			doc.HlContent += "..."
		}
	}

	if searchResponse.RespBody == nil {
		searchResponse.RespBody = &model.ResponseBody{}
	}

	return searchResponse.RespBody, nil
}

type SolrClient struct {
	addCommands []*model.AddCommand
	delCommands []*model.DelCommand
}

func NewSolrClient() *SolrClient {
	return &SolrClient{
		addCommands: make([]*model.AddCommand, 0, 100),
		delCommands: make([]*model.DelCommand, 0, 100),
	}
}

func (this *SolrClient) PushAdd(addCommand *model.AddCommand) {
	this.addCommands = append(this.addCommands, addCommand)
}

func (this *SolrClient) PushDel(delCommand *model.DelCommand) {
	this.delCommands = append(this.delCommands, delCommand)
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

	for _, delCommand := range this.delCommands {
		commandJson, err := json.Marshal(delCommand)
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

		stringBuilder.Append(`"delete":`).Append(commandJson)
	}

	if stringBuilder.Len() == 1 {
		logger.Errorln("post docs:no right addcommand")
		return errors.New("no right addcommand")
	}

	stringBuilder.Append("}")

	logger.Infoln("start post data to solr...")

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

	logger.Infoln("post data result:", result)

	return nil
}
