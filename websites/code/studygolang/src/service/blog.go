// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"logger"
	"math/rand"
	"model"
	"time"
)

// 保存一段时间
var (
	articles     []*model.Article
	startTime    time.Time
	keepDuration time.Duration = 5 * time.Minute
)

// 获取最新的10篇博文，随机展示3篇
func FindNewBlogs() []*model.Article {
	if len(articles) != 0 && startTime.Sub(time.Now()) < keepDuration {
		rnd := rand.Intn(len(articles) - 3)
		return articles[rnd : rnd+3]
	}
	startTime = time.Now()
	articleList, err := model.NewArticle().Where("post_status=publish and post_type=post").Order("post_date DESC").Limit("0, 10").FindAll()
	if err != nil {
		logger.Errorln("获取博客文章失败")
		return nil
	}
	// 内容截取一部分
	for _, article := range articleList {
		t, _ := time.Parse("2006-01-02 15:04:05", article.PostDate)
		article.PostUri = t.Format("2006/01") + "/" + article.PostName
	}
	articles = articleList
	rnd := rand.Intn(len(articles) - 3)
	return articles[rnd : rnd+3]
}
