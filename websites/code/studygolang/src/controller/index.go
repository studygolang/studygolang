package controller

import (
	"config"
	"fmt"
	"html/template"
	"net/http"
)

var ROOT = config.ROOT

func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles(ROOT+"/template/index.html", ROOT+"/template/common/header.html", ROOT+"/template/common/footer.html")
	if err != nil {
		fmt.Fprintln(rw, err)
		return
	}
	tpl.Execute(rw, nil)
}
