// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package api

import (
	"fmt"
	"net/http"
	"service"
)

func AddArticleHandler(rw http.ResponseWriter, req *http.Request) {
	article, err := service.ParseArticle(req.FormValue("url"))
	if err != nil {
		fmt.Fprint(rw, err)
		return
	}

	fmt.Fprint(rw, article.Title)
}
