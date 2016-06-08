package controller

import (
	"logic"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
)

type SearchController struct{}

// 注册路由
func (self SearchController) RegisterRoute(g *echo.Group) {
	g.GET("/search", self.Search)
}

// Search
func (SearchController) Search(ctx echo.Context) error {
	q := ctx.QueryParam("q")
	field := ctx.QueryParam("f")
	p := goutils.MustInt(ctx.QueryParam("p"), 1)

	rows := 20

	respBody, err := logic.DefaultSearcher.DoSearch(q, field, (p-1)*rows, rows)

	data := map[string]interface{}{
		"respBody": respBody,
		"q":        q,
		"f":        field,
	}
	if err == nil {
		uri := "/search?q=" + q + "&f=" + field + "&"
		paginator := logic.NewPaginatorWithPerPage(p, rows)
		data["pageHtml"] = paginator.SetTotal(int64(respBody.NumFound)).GetPageHtml(uri)
	}

	return render(ctx, "search.html", data)
}
