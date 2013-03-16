package admin

import (
	"config"
	"filter"
	"net/http"
)

var ROOT = config.ROOT

func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	user, _ := filter.CurrentUser(req)
	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/admin/index.html")
	filter.SetData(req, map[string]interface{}{"user": user})
}
