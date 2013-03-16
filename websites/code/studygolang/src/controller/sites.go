package controller

import (
	"filter"
	"net/http"
)

func SitesHandler(rw http.ResponseWriter, req *http.Request) {
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/sites.html")
	filter.SetData(req, map[string]interface{}{"activeSites": "active"})
}
