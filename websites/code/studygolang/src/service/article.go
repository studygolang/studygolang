// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"logger"
	"model"
	"net/url"
	"strconv"
	"strings"
	"util"
)

// 获取url对应的文章并根据规则进行解析
func ParseArticle(articleUrl string) (*model.Article, error) {
	if !strings.HasPrefix(articleUrl, "http") {
		articleUrl = "http://" + articleUrl
	}
	urlPaths := strings.SplitN(articleUrl, "/", 5)
	domain := urlPaths[2]

	rule := model.NewCrawlRule()
	err := rule.Where("domain=" + domain).Find()
	if err != nil {
		logger.Errorln("find rule by domain error:", err)
		return nil, err
	}

	if rule.Id == 0 {
		logger.Errorln("domain:", domain, "not exists!")
		return nil, errors.New("domain not exists")
	}

	var doc *goquery.Document
	if doc, err = goquery.NewDocument(articleUrl); err != nil {
		logger.Errorln("goquery newdocument error:", err)
		return nil, err
	}

	author, authorTxt := "", ""
	if rule.InUrl {
		index, err := strconv.Atoi(rule.Author)
		if err != nil {
			logger.Errorln("author rule is illegal:", rule.Author, "error:", err)
			return nil, err
		}
		author = urlPaths[index]
		authorTxt = author
	} else {
		authorSelection := doc.Find(rule.Author)
		author, err = authorSelection.Html()
		if err != nil {
			logger.Errorln("goquery parse author error:", err)
			return nil, err
		}

		author = strings.TrimSpace(author)
		authorTxt = strings.TrimSpace(authorSelection.Text())
	}

	title := strings.TrimSpace(doc.Find(rule.Title).Text())

	contentSelection := doc.Find(rule.Content)
	content, err := contentSelection.Html()
	if err != nil {
		logger.Errorln("goquery parse content error:", err)
		return nil, err
	}
	content = strings.TrimSpace(content)
	txt := strings.TrimSpace(contentSelection.Text())

	pubDate := util.TimeNow()
	if rule.PubDate != "" {
		pubDate = strings.TrimSpace(doc.Find(rule.PubDate).Text())
	}

	article := model.NewArticle()
	article.Domain = domain
	article.Name = rule.Name
	article.Author = author
	article.AuthorTxt = authorTxt
	article.Title = title
	article.Content = content
	article.Txt = txt
	article.PubDate = pubDate
	article.Url = articleUrl

	_, err = article.Insert()
	if err != nil {
		logger.Errorln("insert article error:", err)
		return nil, err
	}

	return article, nil
}

// 获取抓取的文章列表（分页）
func FindArticleByPage(conds map[string]string, curPage, limit int) ([]*model.Article, int) {
	conditions := make([]string, 0, len(conds))
	for k, v := range conds {
		conditions = append(conditions, k+"="+v)
	}

	article := model.NewArticle()

	limitStr := strconv.Itoa((curPage-1)*limit) + "," + strconv.Itoa(limit)
	articleList, err := article.Where(strings.Join(conditions, " AND ")).Order("id DESC").Limit(limitStr).
		FindAll()
	if err != nil {
		logger.Errorln("article service FindArticleByPage Error:", err)
		return nil, 0
	}

	total, err := article.Count()
	if err != nil {
		logger.Errorln("article service FindArticleByPage COUNT Error:", err)
		return nil, 0
	}

	return articleList, total
}

// 获取抓取的文章列表（分页）
func FindArticles(lastId, limit string) []*model.Article {
	article := model.NewArticle()

	articleList, err := article.Where("id>" + lastId).Order("id DESC").Limit(limit).
		FindAll()
	if err != nil {
		logger.Errorln("article service FindArticles Error:", err)
		return nil
	}

	return articleList
}

func FindArticleById(id string) (*model.Article, error) {
	article := model.NewArticle()
	err := article.Where("id=" + id).Find()
	if err != nil {
		logger.Errorln("article service FindArticleById Error:", err)
	}

	return article, err
}

// 修改文章信息
func ModifyArticle(user map[string]interface{}, form url.Values) (errMsg string, err error) {

	username := user["username"].(string)
	form.Set("op_user", username)

	fields := []string{
		"title", "url", "author", "author_txt",
		"lang", "pub_date", "content",
		"tags", "status", "op_user",
	}
	query, args := updateSetClause(form, fields)

	id := form.Get("id")

	err = model.NewArticle().Set(query, args...).Where("id=" + id).Update()
	if err != nil {
		logger.Errorf("更新文章 【%s】 信息失败：%s\n", id, err)
		errMsg = "对不起，服务器内部错误，请稍后再试！"
		return
	}

	return
}

func DelArticle(id string) error {
	return model.NewArticle().Where("id=" + id).Delete()
}

// 获取抓取规则列表（分页）
func FindRuleByPage(conds map[string]string, curPage, limit int) ([]*model.CrawlRule, int) {
	conditions := make([]string, 0, len(conds))
	for k, v := range conds {
		conditions = append(conditions, k+"="+v)
	}

	rule := model.NewCrawlRule()

	limitStr := strconv.Itoa((curPage-1)*limit) + "," + strconv.Itoa(limit)
	ruleList, err := rule.Where(strings.Join(conditions, " AND ")).Order("id DESC").Limit(limitStr).
		FindAll()
	if err != nil {
		logger.Errorln("rule service FindArticleByPage Error:", err)
		return nil, 0
	}

	total, err := rule.Count()
	if err != nil {
		logger.Errorln("rule service FindArticleByPage COUNT Error:", err)
		return nil, 0
	}

	return ruleList, total
}

func SaveRule(form url.Values, opUser string) (errMsg string, err error) {
	rule := model.NewCrawlRule()
	err = util.ConvertAssign(rule, form)
	if err != nil {
		logger.Errorln("rule ConvertAssign error", err)
		errMsg = err.Error()
		return
	}

	rule.OpUser = opUser

	if rule.Id != 0 {
		err = rule.Persist(rule)
	} else {
		_, err = rule.Insert()
	}

	if err != nil {
		errMsg = "内部服务器错误"
		logger.Errorln("rule save:", errMsg, ":", err)
		return
	}

	return
}
