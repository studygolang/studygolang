package admin

import (
	"net/http"
	"strconv"
)

func parsePage(req *http.Request) (curPage, limit int) {

	curPage, err := strconv.Atoi(req.FormValue("page"))
	if err != nil {
		curPage = 0
	}

	limit, err = strconv.Atoi(req.FormValue("limit"))
	if err != nil {
		limit = 20
	}

	return
}

func parseConds(req *http.Request, fields []string) map[string]string {
	conds := make(map[string]string)

	for _, field := range fields {
		if value := req.PostFormValue(field); value != "" {
			conds[field] = value
		}
	}

	return conds
}
