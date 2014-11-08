// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"logger"
	"model"
	"util"
)

func PublishProject(user map[string]interface{}, form url.Values) (err error) {
	isModify := form.Get("id") != ""

	if !isModify && ProjectUriExists(form.Get("uri")) {
		err = errors.New("uri存在")
		return
	}

	username := user["username"].(string)

	project := model.NewOpenProject()
	util.ConvertAssign(project, form)

	if isModify {
		isAdmin := user["isadmin"].(bool)
		if project.Username != username && !isAdmin {
			err = NotModifyAuthorityErr
			return
		}
	} else {
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

	cond := "status=" + strconv.Itoa(model.StatusOnline)
	if lastId != "0" {
		cond += " AND id<" + lastId
	}

	projectList, err := project.Where(cond).Order("id DESC").Limit(limit).
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
	err = project.Where(field+"=?", uniq).Find()

	if err != nil {
		logger.Errorln("project service FindProject error:", err)
		return nil
	}

	if project.Id == 0 {
		return nil
	}

	return project
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
	id := strconv.Itoa(objid)

	// 更新喜欢数（TODO：暂时每次都更新表）
	err := model.NewOpenProject().Where("id="+id).Increment("likenum", num)
	if err != nil {
		logger.Errorln("更新项目喜欢数失败：", err)
	}
}

func (self ProjectLike) String() string {
	return "project"
}
