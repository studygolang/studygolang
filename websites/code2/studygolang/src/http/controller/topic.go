package controller

import (
	"html/template"
	"http/middleware"
	"logic"
	"model"
	"net/http"

	. "http"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
)

// 在需要评论（喜欢）且要回调的地方注册评论（喜欢）对象
func init() {
	// 注册评论（喜欢）对象
	logic.RegisterCommentObject(model.TypeTopic, logic.TopicComment{})
	logic.RegisterLikeObject(model.TypeTopic, logic.TopicLike{})
}

type TopicController struct{}

// 注册路由
func (self TopicController) RegisterRoute(e *echo.Echo) {
	e.Get("/topics", echo.HandlerFunc(self.Topics))
	e.Get("/topics/no_reply", echo.HandlerFunc(self.TopicsNoReply))
	e.Get("/topics/last", echo.HandlerFunc(self.TopicsLast))
	e.Get("/topics/:tid", echo.HandlerFunc(self.Detail))
	e.Get("/topics/node/:nid", echo.HandlerFunc(self.NodeTopics))

	e.Any("/topics/new", echo.HandlerFunc(self.Create), middleware.NeedLogin())
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
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
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

// NodeTopics 某节点下的主题列表
func (TopicController) NodeTopics(ctx echo.Context) error {
	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginator(curPage)

	querystring, nid := "nid=?", goutils.MustInt(ctx.Param("nid"))
	topics := logic.DefaultTopic.FindAll(ctx, paginator, "topics.mtime DESC", querystring, nid)
	total := logic.DefaultTopic.Count(ctx, querystring, nid)
	pageHtml := paginator.SetTotal(total).GetPageHtml(Request(ctx).URL.Path)

	// 当前节点信息
	node := logic.GetNode(nid)

	return render(ctx, "topics/node.html", map[string]interface{}{"activeTopics": "active", "topics": topics, "page": template.HTML(pageHtml), "total": total, "node": node})
}

// Detail 社区主题详细页
func (TopicController) Detail(ctx echo.Context) error {
	tid := goutils.MustInt(ctx.Param("tid"))
	if tid == 0 {
		return ctx.Redirect(http.StatusSeeOther, "/topics")
	}

	topic, replies, err := logic.DefaultTopic.FindByTid(ctx, tid)
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/topics")
	}

	likeFlag := 0
	hadCollect := 0
	// me, ok := ctx.Get("user").(*model.Me)
	// if ok {
	// 	tid := topic["tid"].(int)
	// 	likeFlag = service.HadLike(me.Uid, tid, model.TypeTopic)
	// 	hadCollect = service.HadFavorite(me.Uid, tid, model.TypeTopic)
	// }

	// service.Views.Incr(req, model.TypeTopic, util.MustInt(vars["tid"]))

	return render(ctx, "topics/detail.html,common/comment.html", map[string]interface{}{"activeTopics": "active", "topic": topic, "replies": replies, "likeflag": likeFlag, "hadcollect": hadCollect})
}

// Create 新建主题
func (TopicController) Create(ctx echo.Context) error {
	nodes := logic.GenNodes()

	title := ctx.FormValue("title")
	// 请求新建主题页面
	if title == "" || Request(ctx).Method != "POST" {
		return render(ctx, "topics/new.html", map[string]interface{}{"nodes": nodes, "activeTopics": "active"})
	}

	me := ctx.Get("user").(*model.Me)
	err := logic.DefaultTopic.Publish(ctx, me, Request(ctx).PostForm)
	if err != nil {
		return fail(ctx, 1, "内部服务错误")
	}

	return success(ctx, nil)
}

// 修改主题
// uri: /topics/modify{json:(|.json)}
// func ModifyTopicHandler(rw http.ResponseWriter, req *http.Request) {
// 	tid := req.FormValue("tid")
// 	if tid == "" {
// 		util.Redirect(rw, req, "/topics")
// 		return
// 	}

// 	nodes := service.GenNodes()

// 	vars := mux.Vars(req)
// 	// 请求编辑主题页面
// 	if req.Method != "POST" || vars["json"] == "" {
// 		topic := service.FindTopic(tid)
// 		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/topics/new.html")
// 		filter.SetData(req, map[string]interface{}{"nodes": nodes, "topic": topic, "activeTopics": "active"})
// 		return
// 	}

// 	user, _ := filter.CurrentUser(req)
// 	err := service.PublishTopic(user, req.PostForm)
// 	if err != nil {
// 		if err == service.NotModifyAuthorityErr {
// 			rw.WriteHeader(http.StatusForbidden)
// 			return
// 		}
// 		fmt.Fprint(rw, `{"ok": 0, "error":"内部服务错误！"}`)
// 		return
// 	}
// 	fmt.Fprint(rw, `{"ok": 1, "data":""}`)
// }
