// Copyright 2015 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com  http://golang.top
// Author：polaris	polaris@studygolang.com

// golang.org/x/... 系列的包，go get 不下来，原因你懂的
// 该程序解决此问题，要求使用 go get 的客户端，配置如下 host:
// 		101.251.196.90	golang.org
package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"text/template"

	"github.com/studygolang/mux"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	router := initRouter()
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":9999", nil))
}

// TODO:多人同时请求不同的包，验证可能会失败
var project = ""

func initRouter() *mux.Router {
	router := mux.NewRouter()

	router.PathPrefix("/x").HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		paths := strings.SplitN(path, "/", 4)

		if len(paths) < 3 {
			return parseTmpl(rw)
		}
		project = paths[2]

		parseTmpl(rw)
	})

	return router
}

func parseTmpl(rw http.ResponseWriter) {
	tmpl, err := template.New("gox").Parse(tpl)
	if err != nil {
		fmt.Fprintln(rw, "error:", err)
		return
	}

	err = tmpl.Execute(rw, project)
	if err != nil {
		fmt.Fprintln(rw, "error:", err)
		return
	}
}

var tpl = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
	<meta charset="utf-8">
	<meta name="go-import" content="golang.org/x/{{.}} git https://github.com/golang/{{.}}">
</head>
<body>
</body>
</html>
`
