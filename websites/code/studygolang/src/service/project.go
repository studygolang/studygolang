// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"errors"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"logger"
	"model"
	"util"

	"github.com/PuerkitoBio/goquery"
	"github.com/lunny/html2md"
)

func PublishProject(user map[string]interface{}, form url.Values) (err error) {
	id := form.Get("id")
	isModify := id != ""

	if !isModify && ProjectUriExists(form.Get("uri")) {
		err = errors.New("uri存在")
		return
	}

	username := user["username"].(string)

	project := model.NewOpenProject()

	if isModify {
		err = project.Where("id=?", id).Find()
		if err != nil {
			logger.Errorln("Publish Project find error:", err)
			return
		}
		isAdmin := false
		if _, ok := user["isadmin"]; ok {
			isAdmin = user["isadmin"].(bool)
		}
		if project.Username != username && !isAdmin {
			err = NotModifyAuthorityErr
			return
		}

		util.ConvertAssign(project, form)
	} else {
		util.ConvertAssign(project, form)

		project.Username = username
		project.Ctime = util.TimeNow()
	}

	project.Uri = strings.ToLower(project.Uri)

	github := "github.com"
	pos := strings.Index(project.Src, github)
	if pos != -1 {
		project.Repo = project.Src[pos+len(github)+1:]
	}

	if !isModify {
		_, err = project.Insert()
	} else {
		err = project.Persist(project)
	}

	if err != nil {
		logger.Errorln("Publish Project error:", err)
	}

	// 发布項目，活跃度+10
	if uid, ok := user["uid"].(int); ok {
		weight := 10
		if isModify {
			weight = 2
		}
		go IncUserWeight("uid="+strconv.Itoa(uid), weight)
	}

	return
}

func UpdateProjectStatus(id, status int, username string) error {
	if status < model.StatusNew || status > model.StatusOffline {
		return errors.New("status is illegal")
	}

	logger.Infoln("UpdateProjectStatus by username:", username)

	return model.NewOpenProject().Set("status=?", status).Where("id=?", id).Update()
}

func ProjectUriExists(uri string) bool {
	project := model.NewOpenProject()
	err := project.Where("uri=?", uri).Find("id")
	if err != nil {
		return false
	}

	if project.Id > 0 {
		return true
	}

	return false
}

// 获取开源项目列表（分页）
func FindProjects(lastId, limit string) []*model.OpenProject {
	project := model.NewOpenProject()

	query := "status IN(?,?)"
	args := []interface{}{model.StatusNew, model.StatusOnline}
	if lastId != "0" {
		query += " AND id<?"
		args = append(args, lastId)
	}

	projectList, err := project.Where(query, args...).Order("id DESC").Limit(limit).
		FindAll()
	if err != nil {
		logger.Errorln("project service FindProjects Error:", err)
		return nil
	}

	return projectList
}

// 获取单个项目
func FindProject(uniq string) *model.OpenProject {
	field := "id"
	_, err := strconv.Atoi(uniq)
	if err != nil {
		field = "uri"
	}

	project := model.NewOpenProject()
	err = project.Where(field+"=? AND status IN(?,?)", uniq, model.StatusNew, model.StatusOnline).Find()

	if err != nil {
		logger.Errorln("project service FindProject error:", err)
		return nil
	}

	if project.Id == 0 {
		return nil
	}

	return project
}

// 获得某个用户最近发布的开源项目
func FindUserRecentProjects(username string) []*model.OpenProject {
	projectList, err := model.NewOpenProject().Where("username=?", username).Limit("0,5").Order("ctime DESC").FindAll()
	if err != nil {
		logger.Errorln("project service FindUserRecentProjects error:", err)
		return nil
	}

	return projectList
}

// 获取开源项目列表（分页，后台用）
func FindProjectByPage(conds map[string]string, curPage, limit int) ([]*model.OpenProject, int) {
	conditions := make([]string, 0, len(conds))
	for k, v := range conds {
		conditions = append(conditions, k+"="+v)
	}

	project := model.NewOpenProject()

	limitStr := strconv.Itoa((curPage-1)*limit) + "," + strconv.Itoa(limit)
	projectList, err := project.Where(strings.Join(conditions, " AND ")).Order("id DESC").Limit(limitStr).
		FindAll()
	if err != nil {
		logger.Errorln("project service FindProjectByPage Error:", err)
		return nil, 0
	}

	total, err := project.Count()
	if err != nil {
		logger.Errorln("project service FindProjectByPage COUNT Error:", err)
		return nil, 0
	}

	return projectList, total
}

// 获取多个项目详细信息
func FindProjectsByIds(ids []int) []*model.OpenProject {
	if len(ids) == 0 {
		return nil
	}
	inIds := util.Join(ids, ",")
	projects, err := model.NewOpenProject().Where("id in(" + inIds + ")").FindAll()
	if err != nil {
		logger.Errorln("project service FindProjectsByIds error:", err)
		return nil
	}
	return projects
}

// 开源项目总数
func ProjectsTotal() (total int) {
	total, err := model.NewOpenProject().Count()
	if err != nil {
		logger.Errorln("project service ProjectsTotal error:", err)
	}
	return
}

// ParseProjectList 解析其他网站的开源项目
func ParseProjectList(pUrl string) error {
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
	projectsSelection := doc.Find(".ProjectList .List li")

	for i := projectsSelection.Length() - 1; i >= 0; i-- {

		contentSelection := goquery.NewDocumentFromNode(projectsSelection.Get(i)).Selection
		projectUrl, ok := contentSelection.Find("h3 a").Attr("href")

		if !ok || projectUrl == "" {
			continue
		}
		err = ParseOneProject(projectUrl)

		if err != nil {
			logger.Errorln(err)
		}
	}

	return err
}

const OsChinaDomain = "http://www.oschina.net"

var PresetUsernames = []string{"polaris", "blov", "agolangf", "xuanbao"}

// ParseOneProject 处理单个 project
func ParseOneProject(projectUrl string) error {
	if !strings.HasPrefix(projectUrl, "http") {
		projectUrl = OsChinaDomain + projectUrl
	}

	var (
		doc *goquery.Document
		err error
	)

	if doc, err = goquery.NewDocument(projectUrl); err != nil {
		return errors.New("goquery fetch " + projectUrl + " error:" + err.Error())
	}

	// 标题
	category := strings.TrimSpace(doc.Find(".Project .name").Text())
	name := strings.TrimSpace(doc.Find(".Project .name u").Text())
	tmpIndex := strings.LastIndex(category, name)
	if tmpIndex != -1 {
		category = category[:tmpIndex]
	}

	// uri
	uri := projectUrl[strings.LastIndex(projectUrl, "/")+1:]

	project := model.NewOpenProject()

	err = project.Where("uri=?", uri).Find("id")
	// 已经存在
	if project.Id != 0 {
		return errors.New("url" + projectUrl + "has exists!")
	}

	// 获取项目相关链接
	doc.Find("#Body .urls li").Each(func(i int, liSelection *goquery.Selection) {
		aSelection := liSelection.Find("a")
		uri := util.FetchRealUrl(OsChinaDomain + aSelection.AttrOr("href", ""))
		switch aSelection.Text() {
		case "软件首页":
			project.Home = uri
		case "软件文档":
			project.Doc = uri
		case "软件下载":
			project.Download = uri
		}
	})

	ctime := util.TimeNow()

	doc.Find("#Body .attrs li").Each(func(i int, liSelection *goquery.Selection) {
		aSelection := liSelection.Find("a")
		txt := aSelection.Text()
		if i == 0 {
			project.Licence = txt
			if txt == "未知" {
				project.Licence = "其他"
			}
		} else if i == 1 {
			project.Lang = txt
		} else if i == 2 {
			project.Os = txt
		} else if i == 3 {
			dtime, err := time.ParseInLocation("2006年01月02日", aSelection.Last().Text(), time.Local)
			if err != nil {
				logger.Errorln("parse ctime error:", err)
			} else {
				ctime = dtime.Local().Format("2006-01-02 15:04:05")
			}
		}
	})

	project.Name = name
	project.Category = category
	project.Uri = uri
	project.Repo = strings.TrimSpace(doc.Find("#Body .github-widget").AttrOr("data-repo", ""))
	project.Src = "https://github.com/" + project.Repo

	pos := strings.Index(project.Repo, "/")
	if pos > -1 {
		project.Author = project.Repo[:pos]
	} else {
		project.Author = "网友"
	}

	// TODO: logo

	if project.Doc == "" {
		// TODO：暂时认为一定是 Go 语言
		project.Doc = "https://godoc.org/" + project.Src
	}

	desc := ""
	doc.Find("#Body .detail").Find("p").NextAll().Each(func(i int, domSelection *goquery.Selection) {
		doc.FindSelection(domSelection).WrapHtml(`<div id="tmp` + strconv.Itoa(i) + `"></div>`)
		domHtml, _ := doc.Find("#tmp" + strconv.Itoa(i)).Html()
		if domSelection.Is("pre") {
			desc += domHtml + "\n\n"
		} else {
			desc += html2md.Convert(domHtml) + "\n\n"
		}
	})

	project.Desc = strings.TrimSpace(desc)
	project.Username = PresetUsernames[rand.Intn(4)]
	project.Status = model.ProjectStatusOnline
	project.Ctime = ctime

	_, err = project.Insert()
	if err != nil {
		return errors.New("insert into open project error:" + err.Error())
	}

	return nil
}

// 通过objid获得 project 的所有者
func getProjectOwner(id int) int {
	project := model.NewOpenProject()
	err := project.Where("id=" + strconv.Itoa(id)).Find()
	if err != nil {
		logger.Errorln("project service getProjectOwner Error:", err)
		return 0
	}

	user := FindUserByUsername(project.Username)
	return user.Uid
}

// 提供给其他service调用（包内）
func getProjects(ids map[int]int) map[int]*model.OpenProject {
	projects := FindProjectsByIds(util.MapIntKeys(ids))
	projectMap := make(map[int]*model.OpenProject, len(projects))
	for _, project := range projects {
		projectMap[project.Id] = project
	}
	return projectMap
}

// 项目评论
type ProjectComment struct{}

// 更新该项目的评论信息
// cid：评论id；objid：被评论对象id；uid：评论者；cmttime：评论时间
func (self ProjectComment) UpdateComment(cid, objid, uid int, cmttime string) {
	id := strconv.Itoa(objid)

	// 更新评论数（TODO：暂时每次都更新表）
	err := model.NewOpenProject().Where("id="+id).Increment("cmtnum", 1)
	if err != nil {
		logger.Errorln("更新项目评论数失败：", err)
	}
}

func (self ProjectComment) String() string {
	return "project"
}

// 实现 CommentObjecter 接口
func (self ProjectComment) SetObjinfo(ids []int, commentMap map[int][]*model.Comment) {
	projects := FindProjectsByIds(ids)
	if len(projects) == 0 {
		return
	}

	for _, project := range projects {
		objinfo := make(map[string]interface{})
		objinfo["title"] = project.Category + project.Name
		objinfo["uri"] = model.PathUrlMap[model.TYPE_PROJECT]
		objinfo["type_name"] = model.TypeNameMap[model.TYPE_PROJECT]

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
	err := model.NewOpenProject().Where("id=?", objid).Increment("likenum", num)
	if err != nil {
		logger.Errorln("更新项目喜欢数失败：", err)
	}
}

func (self ProjectLike) String() string {
	return "project"
}
