// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"logger"
	"model"
	"util"
)

var domainPatch = map[string]string{
	"iteye.com":      "iteye.com",
	"blog.51cto.com": "blog.51cto.com",
}

var articleRe = regexp.MustCompile("[\r　\n  \t\v]+")
var articleSpaceRe = regexp.MustCompile("[ ]+")

// 获取url对应的文章并根据规则进行解析
func ParseArticle(articleUrl string, auto bool) (*model.Article, error) {
	articleUrl = strings.TrimSpace(articleUrl)
	if !strings.HasPrefix(articleUrl, "http") {
		articleUrl = "http://" + articleUrl
	}

	tmpArticle := model.NewArticle()
	err := tmpArticle.Where("url=" + articleUrl).Find("id")
	if err != nil || tmpArticle.Id != 0 {
		logger.Errorln(articleUrl, "has exists:", err)
		return nil, errors.New("has exists!")
	}

	urlPaths := strings.SplitN(articleUrl, "/", 5)
	domain := urlPaths[2]

	for k, v := range domainPatch {
		if strings.Contains(domain, k) && !strings.Contains(domain, "www."+k) {
			domain = v
			break
		}
	}

	rule := model.NewCrawlRule()
	err = rule.Where("domain=" + domain).Find()
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
		if strings.HasPrefix(rule.Author, ".") || strings.HasPrefix(rule.Author, "#") {
			authorSelection := doc.Find(rule.Author)
			author, err = authorSelection.Html()
			if err != nil {
				logger.Errorln("goquery parse author error:", err)
				return nil, err
			}

			author = strings.TrimSpace(author)
			authorTxt = strings.TrimSpace(authorSelection.Text())
		} else {
			// 某些个人博客，页面中没有作者的信息，因此，规则中 author 即为 作者
			author = rule.Author
			authorTxt = rule.Author
		}
	}

	title := ""
	doc.Find(rule.Title).Each(func(i int, selection *goquery.Selection) {
		if title != "" {
			return
		}

		tmpTitle := strings.TrimSpace(strings.TrimPrefix(selection.Text(), "原"))
		tmpTitle = strings.TrimSpace(strings.TrimPrefix(tmpTitle, "荐"))
		tmpTitle = strings.TrimSpace(strings.TrimPrefix(tmpTitle, "转"))
		tmpTitle = strings.TrimSpace(strings.TrimPrefix(tmpTitle, "顶"))
		if tmpTitle != "" {
			title = tmpTitle
		}
	})

	if title == "" {
		logger.Errorln("url:", articleUrl, "parse title error:", err)
		return nil, err
	}

	replacer := strings.NewReplacer("[置顶]", "", "[原]", "", "[转]", "")
	title = strings.TrimSpace(replacer.Replace(title))

	contentSelection := doc.Find(rule.Content)
	content, err := contentSelection.Html()
	if err != nil {
		logger.Errorln("goquery parse content error:", err)
		return nil, err
	}
	content = strings.TrimSpace(content)
	txt := strings.TrimSpace(contentSelection.Text())
	txt = articleRe.ReplaceAllLiteralString(txt, " ")
	txt = articleSpaceRe.ReplaceAllLiteralString(txt, " ")

	// 自动抓取，内容长度不能少于 300 字
	if auto && len(txt) < 300 {
		logger.Infoln(articleUrl, "content is short")
		return nil, errors.New("content is short")
	}

	pubDate := util.TimeNow()
	if rule.PubDate != "" {
		pubDate = strings.TrimSpace(doc.Find(rule.PubDate).First().Text())

		// sochina patch
		re := regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}")
		submatches := re.FindStringSubmatch(pubDate)
		if len(submatches) > 0 {
			pubDate = submatches[0]
		}
	}

	if pubDate == "" {
		pubDate = util.TimeNow()
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
	article.Ctime = util.TimeNow()

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

	cond := "status IN(0,1)"
	if lastId != "0" {
		cond += " AND id<" + lastId
	}

	articleList, err := article.Where(cond).Order("id DESC").Limit(limit).
		FindAll()
	if err != nil {
		logger.Errorln("article service FindArticles Error:", err)
		return nil
	}

	return articleList
}

// 获取单条博文
func FindArticleById(id string) (*model.Article, error) {
	article := model.NewArticle()
	err := article.Where("id=" + id).Find()
	if err != nil {
		logger.Errorln("article service FindArticleById Error:", err)
	}

	return article, err
}

// 获取当前(id)博文以及前后博文
func FindArticlesById(idstr string) (curArticle *model.Article, prevNext []*model.Article, err error) {

	id := util.MustInt(idstr)
	cond := "id BETWEEN ? AND ? AND status!=2"

	articles, err := model.NewArticle().Where(cond, id-5, id+5).FindAll()
	if err != nil {
		logger.Errorln("article service FindArticlesById Error:", err)
		return
	}

	if len(articles) == 0 {
		return
	}

	prevNext = make([]*model.Article, 2)
	prevId, nextId := articles[0], id
	for _, article := range articles {
		if article.Id < id && article.Id > prevId {
			prevId = article.Id
			prevNext[0] = article
		} else if article.Id > id {
			nextId = article.Id
			prevNext[1] = article
		} else {
			curArticle = article
		}
	}

	if prevId == id {
		prevNext[0] = nil
	}

	if nextId == id {
		prevNext[1] = nil
	}

	return
}

// 获取多个文章详细信息
func FindArticlesByIds(ids []int) []*model.Article {
	if len(ids) == 0 {
		return nil
	}
	inIds := util.Join(ids, ",")
	articles, err := model.NewArticle().Where("id in(" + inIds + ")").FindAll()
	if err != nil {
		logger.Errorln("article service FindArticlesByIds error:", err)
		return nil
	}
	return articles
}

// 修改文章信息
func ModifyArticle(user map[string]interface{}, form url.Values) (errMsg string, err error) {

	username := user["username"].(string)
	form.Set("op_user", username)

	fields := []string{
		"title", "url", "cover", "author", "author_txt",
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

// 博文总数
func ArticlesTotal() (total int) {
	total, err := model.NewArticle().Count()
	if err != nil {
		logger.Errorln("article service ArticlesTotal error:", err)
	}
	return
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

// 博文评论
type ArticleComment struct{}

// 更新该文章的评论信息
// cid：评论id；objid：被评论对象id；uid：评论者；cmttime：评论时间
func (self ArticleComment) UpdateComment(cid, objid, uid int, cmttime string) {
	id := strconv.Itoa(objid)

	// 更新评论数（TODO：暂时每次都更新表）
	err := model.NewArticle().Where("id="+id).Increment("cmtnum", 1)
	if err != nil {
		logger.Errorln("更新文章评论数失败：", err)
	}
}

func (self ArticleComment) String() string {
	return "article"
}

// 实现 CommentObjecter 接口
func (self ArticleComment) SetObjinfo(ids []int, commentMap map[int][]*model.Comment) {
	articles := FindArticlesByIds(ids)
	if len(articles) == 0 {
		return
	}

	for _, article := range articles {
		objinfo := make(map[string]interface{})
		objinfo["title"] = article.Title

		for _, comment := range commentMap[article.Id] {
			comment.Objinfo = objinfo
		}
	}
}

// 博文喜欢
type ArticleLike struct{}

// 更新该文章的喜欢数
// objid：被喜欢对象id；num: 喜欢数(负数表示取消喜欢)
func (self ArticleLike) UpdateLike(objid, num int) {
	id := strconv.Itoa(objid)

	// 更新喜欢数（TODO：暂时每次都更新表）
	err := model.NewArticle().Where("id="+id).Increment("likenum", num)
	if err != nil {
		logger.Errorln("更新文章喜欢数失败：", err)
	}
}

func (self ArticleLike) String() string {
	return "article"
}
