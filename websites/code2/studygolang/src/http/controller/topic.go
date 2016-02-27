package controller

import (
	"logic"
	"strconv"

	"github.com/labstack/echo"
)

type TopicController struct{}

func (*TopicController) Topics(ctx *echo.Context) error {
	// nodes := service.GenNodes()

	page, _ := strconv.Atoi(ctx.Query("p"))
	if page == 0 {
		page = 1
	}

	// order := ""
	// where := ""
	// view := ""
	// switch vars["view"] {
	// case "/no_reply":
	// 	view = "no_reply"
	// 	where = "lastreplyuid=0"
	// case "/last":
	// 	view = "last"
	// 	order = "ctime DESC"
	// }

	// topics, total := service.FindTopics(page, 0, where, order)
	// pageHtml := service.GetPageHtml(page, total, req.URL.Path)

	topicLogic := &logic.TopicLogic{}
	topics := topicLogic.FindAll(ctx)

	data := map[string]interface{}{
		"topics":       topics,
		"activeTopics": "active",
	}

	ctx.Set(TplFileKey, "topics/list.html")
	ctx.Set(DataKey, data)

	return render(ctx)
	// pageHtml := service.GetPageHtml(page, total, req.URL.Path)
	// req.Form.Set(filter.CONTENT_TPL_KEY, "/template/topics/list.html")
	// 设置模板数据
	// filter.SetData(req, map[string]interface{}{"activeTopics": "active", "topics": topics, "page": template.HTML(pageHtml), "nodes": nodes, "view": view})
	// return ctx.Render(http.StatusOK, "layout.html", data)
}
