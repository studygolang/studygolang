// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package logic

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"util"

	. "db"

	"github.com/polaris1119/config"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
	"github.com/polaris1119/set"

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
	go self.IndexingOpenProject(isAll)
	go self.IndexingTopic(isAll)
	go self.IndexingResource(isAll)
	self.IndexingArticle(isAll)
}

// IndexingArticle 索引博文
func (self SearcherLogic) IndexingArticle(isAll bool) {
	solrClient := NewSolrClient()

	var (
		articleList []*model.Article
		err         error
	)

	id := 0
	for {
		articleList = make([]*model.Article, 0)
		if isAll {
			err = MasterDB.Where("id>?", id).Limit(self.maxRows).OrderBy("id ASC").Find(&articleList)
		} else {
			timeAgo := time.Now().Add(-2 * time.Minute).Format("2006-01-02 15:04:05")
			err = MasterDB.Where("mtime>?", timeAgo).Find(&articleList)
		}
		if err != nil {
			logger.Errorln("IndexingArticle error:", err)
			break
		}

		if len(articleList) == 0 {
			break
		}

		for _, article := range articleList {
			logger.Infoln("deal article_id:", article.Id)

			if id < article.Id {
				id = article.Id
			}

			if article.Tags == "" {
				// 自动生成
				article.Tags = model.AutoTag(article.Title, article.Txt, 4)
				if article.Tags != "" {
					MasterDB.Id(article.Id).Cols("tags").Update(article)
				}
			}

			document := model.NewDocument(article, nil)
			if article.Status != model.ArticleStatusOffline {
				solrClient.PushAdd(model.NewDefaultArgsAddCommand(document))
			} else {
				solrClient.PushDel(model.NewDelCommand(document))
			}
		}

		solrClient.Post()

		if !isAll {
			break
		}
	}
}

// 索引主题
func (self SearcherLogic) IndexingTopic(isAll bool) {
	solrClient := NewSolrClient()

	var (
		topicList   []*model.Topic
		topicExList map[int]*model.TopicUpEx

		err error
	)

	id := 0
	for {
		topicList = make([]*model.Topic, 0)
		topicExList = make(map[int]*model.TopicUpEx)

		if isAll {
			err = MasterDB.Where("tid>?", id).OrderBy("tid ASC").Limit(self.maxRows).Find(&topicList)
		} else {
			timeAgo := time.Now().Add(-2 * time.Minute).Format("2006-01-02 15:04:05")
			err = MasterDB.Where("mtime>?", timeAgo).Find(&topicList)
		}
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
			logger.Infoln("deal topic_id:", topic.Tid)

			if id < topic.Tid {
				id = topic.Tid
			}

			if topic.Tags == "" {
				// 自动生成
				topic.Tags = model.AutoTag(topic.Title, topic.Content, 4)
				if topic.Tags != "" {
					MasterDB.Id(topic.Tid).Cols("tags").Update(topic)
				}
			}

			topicEx := topicExList[topic.Tid]

			document := model.NewDocument(topic, topicEx)
			addCommand := model.NewDefaultArgsAddCommand(document)

			solrClient.PushAdd(addCommand)
		}

		solrClient.Post()

		if !isAll {
			break
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

	id := 0
	for {
		resourceList = make([]*model.Resource, 0)
		resourceExList = make(map[int]*model.ResourceEx)

		if isAll {
			err = MasterDB.Where("id>?", id).OrderBy("id ASC").Limit(self.maxRows).Find(&resourceList)
		} else {
			timeAgo := time.Now().Add(-2 * time.Minute).Format("2006-01-02 15:04:05")
			err = MasterDB.Where("mtime>?", timeAgo).Find(&resourceList)
		}
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
			logger.Infoln("deal resource_id:", resource.Id)

			if id < resource.Id {
				id = resource.Id
			}

			if resource.Tags == "" {
				// 自动生成
				resource.Tags = model.AutoTag(resource.Title+resource.CatName, resource.Content, 4)
				if resource.Tags != "" {
					MasterDB.Id(resource.Id).Cols("tags").Update(resource)
				}
			}

			resourceEx := resourceExList[resource.Id]

			document := model.NewDocument(resource, resourceEx)
			addCommand := model.NewDefaultArgsAddCommand(document)

			solrClient.PushAdd(addCommand)
		}

		solrClient.Post()

		if !isAll {
			break
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

	id := 0
	for {
		projectList = make([]*model.OpenProject, 0)

		if isAll {
			err = MasterDB.Where("id>?", id).OrderBy("id ASC").Limit(self.maxRows).Find(&projectList)
		} else {
			timeAgo := time.Now().Add(-2 * time.Minute).Format("2006-01-02 15:04:05")
			err = MasterDB.Where("mtime>?", timeAgo).Find(&projectList)
		}
		if err != nil {
			logger.Errorln("IndexingArticle error:", err)
			break
		}

		if len(projectList) == 0 {
			break
		}

		for _, project := range projectList {
			logger.Infoln("deal project_id:", project.Id)

			if id < project.Id {
				id = project.Id
			}

			if project.Tags == "" {
				// 自动生成
				project.Tags = model.AutoTag(project.Name+project.Category, project.Desc, 4)
				if project.Tags != "" {
					MasterDB.Id(project.Id).Cols("tags").Update(project)
				}
			}

			document := model.NewDocument(project, nil)
			if project.Status != model.ProjectStatusOffline {
				solrClient.PushAdd(model.NewDefaultArgsAddCommand(document))
			} else {
				solrClient.PushDel(model.NewDelCommand(document))
			}
		}

		solrClient.Post()

		if !isAll {
			break
		}
	}

}

const searchContentLen = 350

// DoSearch 搜索
func (this *SearcherLogic) DoSearch(q, field string, start, rows int) (*model.ResponseBody, error) {
	selectUrl := this.engineUrl + "/select?"

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
	} else if field == "tag" {
		values.Add("q", "*:*")
		values.Add("fq", "tags:"+q)
		values.Add("sort", "viewnum desc")
		q = ""
		field = ""
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
	logger.Infoln(selectUrl + values.Encode())
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

// DoSearch 搜索
func (this *SearcherLogic) SearchByField(field, value string, start, rows int, sorts ...string) (*model.ResponseBody, error) {
	selectUrl := this.engineUrl + "/select?"

	sort := "sort_time desc,cmtnum desc,viewnum desc"
	if len(sorts) > 0 {
		sort = sorts[0]
	}
	var values = url.Values{
		"wt":    []string{"json"},
		"start": []string{strconv.Itoa(start)},
		"rows":  []string{strconv.Itoa(rows)},
		"sort":  []string{sort},
		"fl":    []string{"objid,objtype,title,author,uid,pub_time,tags,viewnum,cmtnum,likenum,lastreplyuid,lastreplytime,updated_at,top,nid"},
	}

	values.Add("q", value)
	values.Add("df", field)

	logger.Infoln(selectUrl + values.Encode())

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

	if searchResponse.RespBody == nil {
		searchResponse.RespBody = &model.ResponseBody{}
	}

	return searchResponse.RespBody, nil
}

func (this *SearcherLogic) FindAtomFeeds(rows int) (*model.ResponseBody, error) {
	selectUrl := this.engineUrl + "/select?"

	var values = url.Values{
		"q":     []string{"*:*"},
		"sort":  []string{"sort_time desc"},
		"wt":    []string{"json"},
		"start": []string{"0"},
		"rows":  []string{strconv.Itoa(rows)},
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

	if searchResponse.RespBody == nil {
		searchResponse.RespBody = &model.ResponseBody{}
	}

	return searchResponse.RespBody, nil
}

func (this *SearcherLogic) FillNodeAndUser(ctx context.Context, respBody *model.ResponseBody) (map[int]*model.User, map[int]*model.TopicNode) {
	if respBody.NumFound == 0 {
		return nil, nil
	}

	uidSet := set.New(set.NonThreadSafe)
	nidSet := set.New(set.NonThreadSafe)

	for _, doc := range respBody.Docs {
		if doc.Uid > 0 {
			uidSet.Add(doc.Uid)
		}
		if doc.Lastreplyuid > 0 {
			uidSet.Add(doc.Lastreplyuid)
		}
		if doc.Nid > 0 {
			nidSet.Add(doc.Nid)
		}
	}

	users := DefaultUser.FindUserInfos(nil, set.IntSlice(uidSet))
	// 获取节点信息
	nodes := GetNodesByNids(set.IntSlice(nidSet))

	return users, nodes
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
