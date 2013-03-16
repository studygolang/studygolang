package util

import (
	"github.com/gorilla/context"
	"net/http"
)

func Redirect(rw http.ResponseWriter, req *http.Request, uri string) {
	// 避免跳转，context中没有清除
	context.Clear(req)

	http.Redirect(rw, req, uri, http.StatusFound)
}
