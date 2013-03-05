package controller

import (
	"html/template"
	"fmt"
	"net/http"
)

func NewTopicHandler(rw http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles(ROOT + "/template/topics/new.html", ROOT + "/template/common/header.html", ROOT + "/template/common/footer.html")
	if err != nil {
		fmt.Fprintln(rw, err)
		return
	}
	tpl.Execute(rw, nil)
}
