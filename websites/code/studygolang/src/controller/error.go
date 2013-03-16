package controller

import (
	"filter"
	"net/http"
)

func NoAuthorizeHandler(rw http.ResponseWriter, req *http.Request) {
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/noauthorize.html")
}

func NotFoundHandler(rw http.ResponseWriter, req *http.Request) {
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/404.html")
}
