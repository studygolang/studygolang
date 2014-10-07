// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"logger"
	"model"
)

// 准备索引数据，post 给 solr
// isAll: 是否全量
func Indexing(isAll bool) {

}

// 索引博文
func IndexingArticle(isAll bool) {
	article := model.NewArticle()

	if isAll {
		id := 0
		for {
			articleList, err := article.Where("id>?", id).FindAll()
			if err != nil {
				logger.Errorln("IndexingArticle error:", err)
				break
			}

			if len(articleList) == 0 {
				break
			}

		}
	}
}
