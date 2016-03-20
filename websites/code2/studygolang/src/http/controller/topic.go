package controller

import (
	"html/template"
	"logic"

	. "http"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
)

type TopicController struct{}

// 注册路由
func (this *TopicController) RegisterRoute(e *echo.Echo) {
	e.Get("/topics", echo.HandlerFunc(this.Topics))
	e.Get("/topics/no_reply", echo.HandlerFunc(this.TopicsNoReply))
	e.Get("/topics/last", echo.HandlerFunc(this.TopicsLast))
}

func (self TopicController) Topics(ctx echo.Context) error {
	return self.topicList(ctx, "", "topics.mtime DESC", "")
}

func (self TopicController) TopicsNoReply(ctx echo.Context) error {
	return self.topicList(ctx, "no_reply", "topics.mtime DESC", "lastreplyuid=?", 0)
}

func (self TopicController) TopicsLast(ctx echo.Context) error {
	return self.topicList(ctx, "last", "ctime DESC", "")
}

func (TopicController) topicList(ctx echo.Context, view, orderBy, querystring string, args ...interface{}) error {
	curPage := goutils.MustInt(ctx.Query("p"), 1)
	paginator := logic.NewPaginator(curPage)

	topics := logic.DefaultTopic.FindAll(ctx, paginator, orderBy, querystring, args...)
	total := logic.DefaultTopic.Count(ctx, querystring, args...)
	pageHtml := paginator.SetTotal(total).GetPageHtml(Request(ctx).URL.Path)

	data := map[string]interface{}{
		"topics":       topics,
		"activeTopics": "active",
		"nodes":        logic.GenNodes(),
		"view":         view,
		"page":         template.HTML(pageHtml),
	}

	return render(ctx, "topics/list.html", data)
}
