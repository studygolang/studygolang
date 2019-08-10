package controller

import (
	"github.com/studygolang/studygolang/context"
	"github.com/studygolang/studygolang/logic"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

type SearchController struct{}

// 注册路由
func (self SearchController) RegisterRoute(g *echo.Group) {
	g.GET("/search", self.Search)
	g.GET("/tag/:name", self.TagList)
}

// Search
func (SearchController) Search(ctx echo.Context) error {
	q := ctx.QueryParam("q")
	field := ctx.QueryParam("f")
	p := goutils.MustInt(ctx.QueryParam("p"), 1)

	rows := 50

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

// TagList
func (SearchController) TagList(ctx echo.Context) error {
	field := "tag"
	p := goutils.MustInt(ctx.QueryParam("p"), 1)
	q := ctx.Param("name")
	if q == "" {
		return render(ctx, "notfound", nil)
	}

	rows := 50

	respBody, err := logic.DefaultSearcher.DoSearch(q, field, (p-1)*rows, rows)
	users, nodes := logic.DefaultSearcher.FillNodeAndUser(context.EchoContext(ctx), respBody)

	data := map[string]interface{}{
		"respBody": respBody,
		"name":     q,
		"users":    users,
		"nodes":    nodes,
	}
	if err == nil {
		uri := "/tag/" + q + "?"
		paginator := logic.NewPaginatorWithPerPage(p, rows)
		data["pageHtml"] = paginator.SetTotal(int64(respBody.NumFound)).GetPageHtml(uri)
	}

	return render(ctx, "feed/tag.html", data)
}
