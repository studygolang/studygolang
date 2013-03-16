package controller

import (
	"filter"
	"net/http"
)

func ResourcesHandler(rw http.ResponseWriter, req *http.Request) {
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/resources.html")
	filter.SetData(req, map[string]interface{}{"activeResources": "active"})
}
