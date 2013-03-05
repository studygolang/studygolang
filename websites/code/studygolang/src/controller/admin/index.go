package admin

import (
	"config"
	"fmt"
	"html/template"
	"net/http"
)

var ROOT = config.ROOT

func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles(ROOT+"/template/admin/common.html", ROOT+"/template/admin/index.html")
	if err != nil {
		fmt.Fprintln(rw, err)
		return
	}
	tpl.Execute(rw, nil)
}
