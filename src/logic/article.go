// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	. "db"
	"errors"
	"model"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/polaris1119/logger"
	"github.com/polaris1119/times"
	"golang.org/x/net/context"
)

type ArticleLogic struct{}

var DefaultArticle = ArticleLogic{}

var domainPatch = map[string]string{
	"iteye.com":      "iteye.com",
	"blog.51cto.com": "blog.51cto.com",
}

var articleRe = regexp.MustCompile("[\r　\n  \t\v]+")
var articleSpaceRe = regexp.MustCompile("[ ]+")

// ParseArticle 获取 url 对应的文章并根据规则进行解析
func (ArticleLogic) ParseArticle(ctx context.Context, articleUrl string, auto bool) (*model.Article, error) {
	articleUrl = strings.TrimSpace(articleUrl)
	if !strings.HasPrefix(articleUrl, "http") {
		articleUrl = "http://" + articleUrl
	}

	tmpArticle := &model.Article{}
	_, err := MasterDB.Where("url=?", articleUrl).Get(tmpArticle)
	if err != nil || tmpArticle.Id != 0 {
		logger.Infoln(articleUrl, "has exists:", err)
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

	rule := &model.CrawlRule{}
	_, err = MasterDB.Where("domain=?", domain).Get(rule)
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

	// relative url -> abs url
	contentSelection.Find("img").Each(func(i int, s *goquery.Selection) {
		if v, ok := s.Attr("src"); ok {
			if !strings.HasPrefix(v, "http") {
				s.SetAttr("src", domain+v)
			}
		}
	})

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

	pubDate := times.Format("Y-m-d H:i:s")
	if rule.PubDate != "" {
		pubDate = strings.TrimSpace(doc.Find(rule.PubDate).First().Text())

		// oschina patch
		re := regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}")
		submatches := re.FindStringSubmatch(pubDate)
		if len(submatches) > 0 {
			pubDate = submatches[0]
		}
	}

	if pubDate == "" {
		pubDate = times.Format("Y-m-d H:i:s")
	} else {
		// YYYYY-MM-dd HH:mm
		if len(pubDate) == 16 && auto {
			// 三个月之前不入库
			pubTime, err := time.ParseInLocation("2006-01-02 15:04", pubDate, time.Local)
			if err == nil {
				if pubTime.Add(3 * 30 * 86400 * time.Second).Before(time.Now()) {
					return nil, errors.New("article is old!")
				}
			}
		}
	}

	article := &model.Article{
		Domain:    domain,
		Name:      rule.Name,
		Author:    author,
		AuthorTxt: authorTxt,
		Title:     title,
		Content:   content,
		Txt:       txt,
		PubDate:   pubDate,
		Url:       articleUrl,
		Lang:      rule.Lang,
	}

	_, err = MasterDB.Insert(article)
	if err != nil {
		logger.Errorln("insert article error:", err)
		return nil, err
	}

	return article, nil
}

func (ArticleLogic) FindLastList(beginTime string, limit int) ([]*model.Article, error) {
	articles := make([]*model.Article, 0)
	err := MasterDB.Where("ctime>? AND status!=?", beginTime, model.ArticleStatusOffline).
		OrderBy("cmtnum DESC, likenum DESC, viewnum DESC").Limit(limit).Find(&articles)

	return articles, err
}

// Total 博文总数
func (ArticleLogic) Total() int64 {
	total, err := MasterDB.Count(new(model.Article))
	if err != nil {
		logger.Errorln("ArticleLogic Total error:", err)
	}
	return total
}

// FindBy 获取抓取的文章列表（分页）
func (ArticleLogic) FindBy(ctx context.Context, limit int, lastIds ...int) []*model.Article {
	objLog := GetLogger(ctx)

	dbSession := MasterDB.Where("status IN(?,?)", model.ArticleStatusNew, model.ArticleStatusOnline)

	if len(lastIds) > 0 && lastIds[0] > 0 {
		dbSession.And("id<?", lastIds[0])
	}

	articles := make([]*model.Article, 0)
	err := dbSession.OrderBy("id DESC").Limit(limit).Find(&articles)
	if err != nil {
		objLog.Errorln("ArticleLogic FindBy Error:", err)
		return nil
	}

	topArticles := make([]*model.Article, 0)
	err = MasterDB.Where("top=?", 1).OrderBy("id DESC").Find(&topArticles)
	if err != nil {
		objLog.Errorln("ArticleLogic Find Top Articles Error:", err)
		return nil
	}
	if len(topArticles) > 0 {
		articles = append(topArticles, articles...)
	}

	return articles
}

// 获取抓取的文章列表（分页）：后台用
func (ArticleLogic) FindArticleByPage(ctx context.Context, conds map[string]string, curPage, limit int) ([]*model.Article, int) {
	objLog := GetLogger(ctx)

	session := MasterDB.NewSession()
	session.IsAutoClose = true

	for k, v := range conds {
		session.And(k+"=?", v)
	}

	totalSession := session.Clone()

	offset := (curPage - 1) * limit
	articleList := make([]*model.Article, 0)
	err := session.OrderBy("id DESC").Limit(limit, offset).Find(&articleList)
	if err != nil {
		objLog.Errorln("find error:", err)
		return nil, 0
	}

	total, err := totalSession.Count(new(model.Article))
	if err != nil {
		objLog.Errorln("find count error:", err)
		return nil, 0
	}

	return articleList, int(total)
}

// FindByIds 获取多个文章详细信息
func (ArticleLogic) FindByIds(ids []int) []*model.Article {
	if len(ids) == 0 {
		return nil
	}
	articles := make([]*model.Article, 0)
	err := MasterDB.In("id", ids).Find(&articles)
	if err != nil {
		logger.Errorln("ArticleLogic FindByIds error:", err)
		return nil
	}
	return articles
}

// findByIds 获取多个文章详细信息 包内使用
func (ArticleLogic) findByIds(ids []int) map[int]*model.Article {
	if len(ids) == 0 {
		return nil
	}
	articles := make(map[int]*model.Article)
	err := MasterDB.In("id", ids).Find(&articles)
	if err != nil {
		logger.Errorln("ArticleLogic findByIds error:", err)
		return nil
	}
	return articles
}

// FindByIdAndPreNext 获取当前(id)博文以及前后博文
func (ArticleLogic) FindByIdAndPreNext(ctx context.Context, id int) (curArticle *model.Article, prevNext []*model.Article, err error) {
	objLog := GetLogger(ctx)

	articles := make([]*model.Article, 0)

	err = MasterDB.Where("id BETWEEN ? AND ? AND status!=?", id-5, id+5, model.ArticleStatusOffline).Find(&articles)
	if err != nil {
		objLog.Errorln("ArticleLogic FindByIdAndPreNext Error:", err)
		return
	}

	if len(articles) == 0 {
		objLog.Errorln("ArticleLogic FindByIdAndPreNext not find articles, id:", id)
		return
	}

	prevNext = make([]*model.Article, 2)
	prevId, nextId := articles[0].Id, articles[len(articles)-1].Id
	for _, article := range articles {
		if article.Id < id && article.Id > prevId {
			prevId = article.Id
			prevNext[0] = article
		} else if article.Id > id && article.Id < nextId {
			nextId = article.Id
			prevNext[1] = article
		} else if article.Id == id {
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

// Modify 修改文章信息
func (ArticleLogic) Modify(ctx context.Context, user *model.Me, form url.Values) (errMsg string, err error) {
	form.Set("op_user", user.Username)

	fields := []string{
		"title", "url", "cover", "author", "author_txt",
		"lang", "pub_date", "content",
		"tags", "status", "op_user",
	}
	change := make(map[string]string)

	for _, field := range fields {
		change[field] = form.Get(field)
	}

	id := form.Get("id")
	_, err = MasterDB.Table(new(model.Article)).Id(id).Update(change)
	if err != nil {
		logger.Errorf("更新文章 【%s】 信息失败：%s\n", id, err)
		errMsg = "对不起，服务器内部错误，请稍后再试！"
		return
	}

	return
}

// FindById 获取单条博文
func (ArticleLogic) FindById(ctx context.Context, id string) (*model.Article, error) {
	article := &model.Article{}
	_, err := MasterDB.Id(id).Get(article)
	if err != nil {
		logger.Errorln("article logic FindById Error:", err)
	}

	return article, err
}

// 博文评论
type ArticleComment struct{}

// UpdateComment 更新该文章的评论信息
// cid：评论id；objid：被评论对象id；uid：评论者；cmttime：评论时间
func (self ArticleComment) UpdateComment(cid, objid, uid int, cmttime time.Time) {
	// 更新评论数（TODO：暂时每次都更新表）
	_, err := MasterDB.Id(objid).Incr("cmtnum", 1).Update(new(model.Article))
	if err != nil {
		logger.Errorln("更新文章评论数失败：", err)
	}
}

func (self ArticleComment) String() string {
	return "article"
}

// SetObjinfo 实现 CommentObjecter 接口
func (self ArticleComment) SetObjinfo(ids []int, commentMap map[int][]*model.Comment) {
	articles := DefaultArticle.FindByIds(ids)
	if len(articles) == 0 {
		return
	}

	for _, article := range articles {
		objinfo := make(map[string]interface{})
		objinfo["title"] = article.Title
		objinfo["uri"] = model.PathUrlMap[model.TypeArticle]
		objinfo["type_name"] = model.TypeNameMap[model.TypeArticle]

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
	// 更新喜欢数（TODO：暂时每次都更新表）
	_, err := MasterDB.Where("id=?", objid).Incr("likenum", num).Update(new(model.Article))
	if err != nil {
		logger.Errorln("更新文章喜欢数失败：", err)
	}
}

func (self ArticleLike) String() string {
	return "article"
}
