// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"fmt"
	"os"
	"strconv"
	"text/template"
	"time"

	"config"
	"logger"
	"model"
	"util"
)

var sitemapTpl = template.Must(template.ParseFiles(config.ROOT + "/template/sitemap.xml"))
var sitemapIndexTpl = template.Must(template.ParseFiles(config.ROOT + "/template/sitemapindex.xml"))

var sitemapPath = config.ROOT + "/sitemap/"

func init() {
	if !util.Exist(sitemapPath) {
		err := os.MkdirAll(sitemapPath, 0777)
		if err != nil {
			panic(err)
		}
	}
}

func GenSitemap() {
	sitemapFiles := []string{}

	// 首页
	home := map[string]string{
		"loc":      "http://" + config.Config["domain"],
		"lastmode": time.Now().Format(time.RFC3339),
	}

	var (
		little = 1
		step   = 4999
		large  = little + step
	)

	// 文章
	article := model.NewArticle()
	for {
		sitemapFile := "sitemap_article_" + strconv.Itoa(large) + ".xml"

		articles, err := article.Where("id BETWEEN ? AND ? AND status!=?", little, large, model.StatusOffline).FindAll("id", "mtime")
		little = large + 1
		large = little + step

		if err != nil {
			continue
		}

		if len(articles) == 0 {
			break
		}

		data := map[string]interface{}{
			"home":     home,
			"articles": articles,
		}

		if err = output(sitemapFile, data); err == nil {
			sitemapFiles = append(sitemapFiles, sitemapFile)
		}
	}

	little = 1
	large = little + step

	// 主题（帖子）
	topic := model.NewTopic()
	for {
		sitemapFile := "sitemap_topic_" + strconv.Itoa(large) + ".xml"

		topics, err := topic.Where("tid BETWEEN ? AND ? AND flag IN(?,?)", little, large, 0, 1).FindAll("tid", "mtime")
		little, large = large+1, little+step

		if err != nil {
			continue
		}

		if len(topics) == 0 {
			break
		}

		data := map[string]interface{}{
			"home":   home,
			"topics": topics,
		}

		if err = output(sitemapFile, data); err == nil {
			sitemapFiles = append(sitemapFiles, sitemapFile)
		}
	}

	little = 1
	large = little + step

	// 资源
	resource := model.NewResource()
	for {
		sitemapFile := "sitemap_resource_" + strconv.Itoa(large) + ".xml"

		resources, err := resource.Where("id BETWEEN ? AND ?", little, large).FindAll("id", "mtime")
		little, large = large+1, little+step

		if err != nil {
			continue
		}

		if len(resources) == 0 {
			break
		}

		data := map[string]interface{}{
			"home":      home,
			"resources": resources,
		}

		if err = output(sitemapFile, data); err == nil {
			sitemapFiles = append(sitemapFiles, sitemapFile)
		}
	}

	little = 1
	large = little + step

	// 项目
	project := model.NewOpenProject()
	for {
		sitemapFile := "sitemap_project_" + strconv.Itoa(large) + ".xml"

		projects, err := project.Where("id BETWEEN ? AND ? AND status=?", little, large, model.StatusOnline).FindAll("id", "uri", "mtime")
		little, large = large+1, little+step

		if err != nil {
			continue
		}

		if len(projects) == 0 {
			break
		}

		data := map[string]interface{}{
			"home":     home,
			"projects": projects,
		}

		if err = output(sitemapFile, data); err == nil {
			sitemapFiles = append(sitemapFiles, sitemapFile)
		}
	}

	little = 1
	large = little + step

	// wiki
	wiki := model.NewWiki()
	for {
		sitemapFile := "sitemap_wiki_" + strconv.Itoa(large) + ".xml"

		wikis, err := wiki.Where("id BETWEEN ? AND ?", little, large).FindAll("id", "uri", "mtime")
		little, large = large+1, little+step

		if err != nil {
			continue
		}

		if len(wikis) == 0 {
			break
		}

		data := map[string]interface{}{
			"home":  home,
			"wikis": wikis,
		}

		if err = output(sitemapFile, data); err == nil {
			sitemapFiles = append(sitemapFiles, sitemapFile)
		}
	}

	file, err := os.Create(sitemapPath + "sitemapindex.xml")
	if err != nil {
		logger.Errorln("gen sitemap index file error:", err)
		return
	}
	defer file.Close()

	err = sitemapIndexTpl.Execute(file, map[string]interface{}{
		"home":         home,
		"sitemapFiles": sitemapFiles,
	})
	if err != nil {
		logger.Errorln("execute sitemap index template error:", err)
	}

	fmt.Println("finish")
}

func output(filename string, data map[string]interface{}) (err error) {
	var file *os.File
	file, err = os.Create(sitemapPath + filename)
	if err != nil {
		logger.Errorln("open file error:", err)
		return
	}
	defer file.Close()
	if err = sitemapTpl.Execute(file, data); err != nil {
		logger.Errorln("execute template error:", err)
	}

	return
}
