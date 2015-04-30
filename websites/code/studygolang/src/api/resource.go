// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com http://golang.top
// Author：polaris	polaris@studygolang.com

package api

import (
	"fmt"
	"net/http"
	"service"
)

func AddRedditResourceHandler(rw http.ResponseWriter, req *http.Request) {
	err := service.ParseReddit(req.FormValue("url"))
	if err != nil {
		fmt.Fprint(rw, err)
		return
	}

	fmt.Fprint(rw, "success")
}
