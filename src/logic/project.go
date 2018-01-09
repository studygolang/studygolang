// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"errors"
	"math/rand"
	"model"
	"net/url"
	"strconv"
	"strings"
	"time"
	"util"

	. "db"

	"github.com/PuerkitoBio/goquery"
	"github.com/lunny/html2md"
	"github.com/polaris1119/config"
	"github.com/polaris1119/logger"
	"github.com/polaris1119/set"
	"golang.org/x/net/context"
)

type ProjectLogic struct{}

var DefaultProject = ProjectLogic{}

func (self ProjectLogic) Publish(ctx context.Context, user *model.Me, form url.Values) (err error) {
	objLog := GetLogger(ctx)

	id := form.Get("id")
	isModify := id != ""

	project := &model.OpenProject{}

	if isModify {
		_, err = MasterDB.Id(id).Get(project)
		if err != nil {
			objLog.Errorln("Publish Project find error:", err)
			return
		}

		if !CanEdit(user, project) {
			err = NotModifyAuthorityErr
			return
		}

		err = schemaDecoder.Decode(project, form)
		if err != nil {
			objLog.Errorln("Publish Project schema decode error:", err)
			return
		}
	} else {
		err = schemaDecoder.Decode(project, form)
		if err != nil {
			objLog.Errorln("Publish Project schema decode error:", err)
			return
		}

		project.Username = user.Username
	}
	if project.Uri == "" {
		project.Uri = strings.Replace(project.Name, " ", "-", -1)
	}
	project.Uri = strings.ToLower(project.Uri)

	if !isModify && self.UriExists(ctx, form.Get("uri")) {
		err = errors.New("项目已存在")
		return
	}

	github := "github.com"
	pos := strings.Index(project.Src, github)
	if pos != -1 {
		project.Repo = project.Src[pos+len(github)+1:]
	}

	var affected int64
	if !isModify {
		affected, err = MasterDB.Insert(project)
	} else {
		affected, err = MasterDB.Id(id).Update(project)
	}

	if err != nil {
		objLog.Errorln("Publish Project error:", err)
		return
	}

	if affected == 0 {
		return
	}

	if isModify {
		go modifyObservable.NotifyObservers(user.Uid, model.TypeProject, project.Id)
	} else {
		go publishObservable.NotifyObservers(user.Uid, model.TypeProject, project.Id)
	}

	return
}

// UriExists 通过 uri 是否存在 project
func (ProjectLogic) UriExists(ctx context.Context, uri string) bool {
	total, err := MasterDB.Where("uri=?", uri).Count(new(model.OpenProject))
	if err != nil || total == 0 {
		return false
	}

	return true
}

// Total 开源项目总数
func (ProjectLogic) Total() int64 {
	total, err := MasterDB.Count(new(model.OpenProject))
	if err != nil {
		logger.Errorln("ProjectLogic Total error:", err)
	}
	return total
}

// FindBy 获取开源项目列表（分页）
func (ProjectLogic) FindBy(ctx context.Context, limit int, lastIds ...int) []*model.OpenProject {
	objLog := GetLogger(ctx)

	dbSession := MasterDB.Where("status IN(?,?)", model.ProjectStatusNew, model.ProjectStatusOnline)
	if len(lastIds) > 0 && lastIds[0] > 0 {
		dbSession.And("id<?", lastIds[0])
	}

	projectList := make([]*model.OpenProject, 0)
	err := dbSession.OrderBy("id DESC").Limit(limit).Find(&projectList)
	if err != nil {
		objLog.Errorln("ProjectLogic FindBy Error:", err)
		return nil
	}

	return projectList
}

// FindByIds 获取多个项目详细信息
func (ProjectLogic) FindByIds(ids []int) []*model.OpenProject {
	if len(ids) == 0 {
		return nil
	}

	projects := make([]*model.OpenProject, 0)
	err := MasterDB.In("id", ids).Find(&projects)
	if err != nil {
		logger.Errorln("ProjectLogic FindByIds error:", err)
		return nil
	}
	return projects
}

// findByIds 获取多个项目详细信息 包内使用
func (ProjectLogic) findByIds(ids []int) map[int]*model.OpenProject {
	if len(ids) == 0 {
		return nil
	}

	projects := make(map[int]*model.OpenProject)
	err := MasterDB.In("id", ids).Find(&projects)
	if err != nil {
		logger.Errorln("ProjectLogic FindByIds error:", err)
		return nil
	}
	return projects
}

// FindOne 获取单个项目
func (ProjectLogic) FindOne(ctx context.Context, val interface{}) *model.OpenProject {
	objLog := GetLogger(ctx)

	field := "id"
	_, ok := val.(int)
	if !ok {
		val := val.(string)
		if _, err := strconv.Atoi(val); err != nil {
			field = "uri"
		}
	}

	project := &model.OpenProject{}
	_, err := MasterDB.Where(field+"=? AND status IN(?,?)", val, model.ProjectStatusNew, model.ProjectStatusOnline).Get(project)
	if err != nil {
		objLog.Errorln("project service FindProject error:", err)
		return nil
	}

	project.User = DefaultUser.FindOne(ctx, "username", project.Username)

	return project
}

// FindRecent 获得某个用户最近发布的开源项目
func (ProjectLogic) FindRecent(ctx context.Context, username string) []*model.OpenProject {
	projectList := make([]*model.OpenProject, 0)
	err := MasterDB.Where("username=?", username).Limit(5).OrderBy("id DESC").Find(&projectList)
	if err != nil {
		logger.Errorln("project logic FindRecent error:", err)
		return nil
	}
	return projectList
}

// FindAll 支持多页翻看
func (self ProjectLogic) FindAll(ctx context.Context, paginator *Paginator, orderBy string, querystring string, args ...interface{}) []*model.OpenProject {
	objLog := GetLogger(ctx)

	projects := make([]*model.OpenProject, 0)
	session := MasterDB.OrderBy(orderBy)
	if querystring != "" {
		session.Where(querystring, args...)
	}
	err := session.Limit(paginator.PerPage(), paginator.Offset()).Find(&projects)
	if err != nil {
		objLog.Errorln("ProjectLogic FindAll error:", err)
		return nil
	}

	self.fillUser(projects)

	return projects
}

func (ProjectLogic) Count(ctx context.Context, querystring string, args ...interface{}) int64 {
	objLog := GetLogger(ctx)

	var (
		total int64
		err   error
	)
	if querystring == "" {
		total, err = MasterDB.Count(new(model.OpenProject))
	} else {
		total, err = MasterDB.Where(querystring, args...).Count(new(model.OpenProject))
	}

	if err != nil {
		objLog.Errorln("ProjectLogic Count error:", err)
	}

	return total
}

func (ProjectLogic) fillUser(projects []*model.OpenProject) {
	usernameSet := set.New(set.NonThreadSafe)
	uidSet := set.New(set.NonThreadSafe)
	for _, project := range projects {
		usernameSet.Add(project.Username)

		if project.Lastreplyuid != 0 {
			uidSet.Add(project.Lastreplyuid)
		}
	}
	if !usernameSet.IsEmpty() {
		userMap := DefaultUser.FindUserInfos(nil, set.StringSlice(usernameSet))
		for _, project := range projects {
			for _, user := range userMap {
				if project.Username == user.Username {
					project.User = user
					break
				}
			}
		}
	}

	if !uidSet.IsEmpty() {
		replyUserMap := DefaultUser.FindUserInfos(nil, set.IntSlice(uidSet))
		for _, project := range projects {
			if project.Lastreplyuid == 0 {
				continue
			}

			project.LastReplyUser = replyUserMap[project.Lastreplyuid]
		}
	}
}

// getOwner 通过objid获得 project 的所有者
func (ProjectLogic) getOwner(ctx context.Context, id int) int {
	project := &model.OpenProject{}
	_, err := MasterDB.Id(id).Get(project)
	if err != nil {
		logger.Errorln("project logic getOwner Error:", err)
		return 0
	}

	user := DefaultUser.FindOne(ctx, "username", project.Username)
	return user.Uid
}

// ParseProjectList 解析其他网站的开源项目
func (self ProjectLogic) ParseProjectList(pUrl string) error {
	pUrl = strings.TrimSpace(pUrl)
	if !strings.HasPrefix(pUrl, "http") {
		pUrl = "http://" + pUrl
	}

	var (
		doc *goquery.Document
		err error
	)

	if doc, err = goquery.NewDocument(pUrl); err != nil {
		logger.Errorln("goquery opensource project newdocument error:", err)
		return err
	}

	// 最后面的先入库处理
	projectsSelection := doc.Find(".news-list").Children()

	for i := projectsSelection.Length() - 1; i >= 0; i-- {

		contentSelection := goquery.NewDocumentFromNode(projectsSelection.Get(i)).Selection
		projectUrl, ok := contentSelection.Find(".box-aw a").First().Attr("href")

		if !ok || projectUrl == "" {
			logger.Errorln("project url is empty")
			continue
		}
		go func(projectUrl string) {
			err := self.ParseOneProject(projectUrl)

			if err != nil {
				logger.Errorln(err)
			}
		}(projectUrl)
	}

	return err
}

const OsChinaDomain = "http://www.oschina.net"

// ProjectLogoPrefix 开源项目 logo 前缀
const ProjectLogoPrefix = "plogo"

var PresetUsernames = config.ConfigFile.MustValueArray("crawl", "preset_users", ",")

// ParseOneProject 处理单个 project
func (ProjectLogic) ParseOneProject(projectUrl string) error {
	if !strings.HasPrefix(projectUrl, "http") {
		projectUrl = OsChinaDomain + projectUrl
	}

	var (
		doc *goquery.Document
		err error
	)

	// 加上 ?fromerr=xfwefs，否则页面有 js 重定向
	if doc, err = goquery.NewDocument(projectUrl + "?fromerr=xfwefs"); err != nil {
		return errors.New("goquery fetch " + projectUrl + " error:" + err.Error())
	}

	// 标题
	category := strings.TrimSpace(doc.Find("#v-header header .box-aw h1").Text())
	name := strings.TrimSpace(doc.Find("#v-header header .box-aw h1 span").Text())
	if category == "" && name == "" {
		return errors.New("projectUrl:" + projectUrl + " category and name are empty")
	}

	tmpIndex := strings.LastIndex(category, name)
	if tmpIndex != -1 {
		category = category[:tmpIndex]
	}

	// uri
	uri := projectUrl[strings.LastIndex(projectUrl, "/")+1:]

	project := &model.OpenProject{}

	_, err = MasterDB.Where("uri=?", uri).Get(project)
	// 已经存在
	if project.Id != 0 {
		logger.Infoln("url", projectUrl, "has exists!")
		return nil
	}

	logoSelection := doc.Find("#v-header header .logo img")
	if logoSelection.AttrOr("title", "") != "" {
		project.Logo = logoSelection.AttrOr("src", "")

		if !strings.HasPrefix(project.Logo, "http") {
			// project.Logo = OsChinaDomain + project.Logo
			// 默认图
			project.Logo = ""
		} else {
			project.Logo, err = DefaultUploader.TransferUrl(nil, project.Logo, ProjectLogoPrefix)
			if err != nil {
				logger.Errorln("project logo upload error:", err)
			}
		}
	}

	// 获取项目相关链接
	doc.Find("#v-details .urls a").Each(func(i int, aSelection *goquery.Selection) {
		uri := util.FetchRealUrl(OsChinaDomain + aSelection.AttrOr("href", ""))
		switch aSelection.Find("span").Text() {
		case "软件首页":
			project.Home = uri
		case "软件文档":
			project.Doc = uri
		case "软件下载":
			project.Download = uri
		}
	})

	doc.Find("#v-basic .list .box").Each(func(i int, liSelection *goquery.Selection) {
		aSelection := liSelection.Find("span")
		txt := strings.TrimSpace(aSelection.Text())
		if i == 0 {
			project.Licence = txt
			if txt == "未知" {
				project.Licence = "其他"
			}
		} else if i == 1 {
			txt = liSelection.Find("span a:first-child").Text()
			project.Lang = txt
		} else if i == 2 {
			project.Os = txt
		}
	})

	project.Name = name
	project.Category = strings.TrimSpace(category)
	project.Uri = uri
	project.Repo = strings.TrimSpace(doc.Find("#v-details .github-widget").AttrOr("data-repo", ""))
	project.Src = "https://github.com/" + project.Repo

	pos := strings.Index(project.Repo, "/")
	if pos > -1 {
		project.Author = project.Repo[:pos]
	} else {
		project.Author = "网友"
	}

	if project.Doc == "" {
		// TODO：暂时认为一定是 Go 语言
		project.Doc = "https://godoc.org/" + project.Src[8:]
	}

	desc := ""
	doc.Find("#v-details .detail").Find("p").Each(func(i int, domSelection *goquery.Selection) {
		doc.FindSelection(domSelection).WrapHtml(`<div id="tmp` + strconv.Itoa(i) + `"></div>`)
		domHtml, _ := doc.Find("#tmp" + strconv.Itoa(i)).Html()
		if domSelection.Is("pre") {
			desc += domHtml + "\n\n"
		} else {
			desc += html2md.Convert(domHtml) + "\n\n"
		}
	})

	project.Desc = strings.TrimSpace(desc)
	project.Username = PresetUsernames[rand.Intn(len(PresetUsernames))]
	project.Status = model.ProjectStatusOnline
	project.Ctime = model.OftenTime(time.Now())

	_, err = MasterDB.Insert(project)
	if err != nil {
		return errors.New("insert into open project error:" + err.Error())
	}

	return nil
}

// 项目评论
type ProjectComment struct{}

// 更新该项目的评论信息
// cid：评论id；objid：被评论对象id；uid：评论者；cmttime：评论时间
func (self ProjectComment) UpdateComment(cid, objid, uid int, cmttime time.Time) {
	// 更新评论数（TODO：暂时每次都更新表）
	_, err := MasterDB.Table(new(model.OpenProject)).Id(objid).Incr("cmtnum", 1).Update(map[string]interface{}{
		"lastreplyuid":  uid,
		"lastreplytime": cmttime,
	})
	if err != nil {
		logger.Errorln("更新项目评论数失败：", err)
		return
	}
}

func (self ProjectComment) String() string {
	return "project"
}

// 实现 CommentObjecter 接口
func (self ProjectComment) SetObjinfo(ids []int, commentMap map[int][]*model.Comment) {
	projects := DefaultProject.FindByIds(ids)
	if len(projects) == 0 {
		return
	}

	for _, project := range projects {
		objinfo := make(map[string]interface{})
		objinfo["title"] = project.Category + project.Name
		objinfo["uri"] = model.PathUrlMap[model.TypeProject]
		objinfo["type_name"] = model.TypeNameMap[model.TypeProject]

		for _, comment := range commentMap[project.Id] {
			comment.Objinfo = objinfo
		}
	}
}

// 项目喜欢
type ProjectLike struct{}

// 更新该项目的喜欢数
// objid：被喜欢对象id；num: 喜欢数(负数表示取消喜欢)
func (self ProjectLike) UpdateLike(objid, num int) {
	// 更新喜欢数（TODO：暂时每次都更新表）
	_, err := MasterDB.Id(objid).Incr("likenum", num).Update(new(model.OpenProject))
	if err != nil {
		logger.Errorln("更新项目喜欢数失败：", err)
	}
}

func (self ProjectLike) String() string {
	return "project"
}
