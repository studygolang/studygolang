// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package logic

import (
	. "db"
	"model"
)

type ArticleLogic struct{}

var DefaultArticle = ArticleLogic{}

func (ArticleLogic) FindLastList(beginTime string, limit int) ([]*model.Article, error) {
	articles := make([]*model.Article, 0)
	err := MasterDB.Where("ctime>? AND status!=?", beginTime, model.ArticleStatusOffline).
		OrderBy("cmtnum DESC, likenum DESC, viewnum DESC").Limit(limit).Find(&articles)

	return articles, err
}
