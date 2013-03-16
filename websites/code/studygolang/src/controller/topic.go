package controller

import (
	"filter"
	"fmt"
	"github.com/studygolang/mux"
	"html/template"
	"logger"
	"model"
	"net/http"
	"service"
	"strconv"
	"util"
)

// 在需要评论且要回调的地方注册评论对象
func init() {
	// 注册评论对象
	service.RegisterCommentObject("topic", service.TopicComment{})
}

// 社区帖子列表页
// uri: /topics{view:(|/popular|/no_reply|/last)}
func TopicsHandler(rw http.ResponseWriter, req *http.Request) {
	nodes := genNodes()
	// 设置内容模板
	page, _ := strconv.Atoi(req.FormValue("p"))
	if page == 0 {
		page = 1
	}
	vars := mux.Vars(req)
	order := ""
	where := ""
	switch vars["view"] {
	case "/no_reply":
		where = "lastreplyuid=0"
	case "/last":
		order = "ctime DESC"
	}
	topics, total := service.FindTopics(page, 0, where, order)
	pageHtml := service.GetPageHtml(page, total)
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/topics/list.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"activeTopics": "active", "topics": topics, "page": template.HTML(pageHtml), "nodes": nodes})
}

// 某节点下的帖子列表
// uri: /topics/node{nid:[0-9]+}
func NodesHandler(rw http.ResponseWriter, req *http.Request) {
	page, _ := strconv.Atoi(req.FormValue("p"))
	if page == 0 {
		page = 1
	}
	vars := mux.Vars(req)
	topics, total := service.FindTopics(page, 0, "nid="+vars["nid"])
	pageHtml := service.GetPageHtml(page, total)
	// 当前节点信息
	node := model.GetNode(util.MustInt(vars["nid"]))
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/topics/node.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"activeTopics": "active", "topics": topics, "page": template.HTML(pageHtml), "total": total, "node": node})
}

// 社区帖子详细页
// uri: /topics/{tid:[0-9]+}
func TopicDetailHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	topic, replies, err := service.FindTopicByTid(vars["tid"])
	if err != nil {
		// TODO:
	}
	uid := 0
	user, ok := filter.CurrentUser(req)
	if ok {
		uid = user["uid"].(int)
	}
	// TODO:刷屏暂时不处理
	// 增加浏览量
	service.IncrTopicView(vars["tid"], uid)
	// 设置内容模板
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/topics/detail.html")
	// 设置模板数据
	filter.SetData(req, map[string]interface{}{"activeTopics": "active", "topic": topic, "replies": replies})
}

// 新建帖子
// uri: /topics/new
func NewTopicHandler(rw http.ResponseWriter, req *http.Request) {
	nodes := genNodes()
	vars := mux.Vars(req)
	title := req.FormValue("title")
	// 请求新建帖子页面
	if title == "" || req.Method != "POST" || vars["json"] == "" {
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/topics/new.html")
		filter.SetData(req, map[string]interface{}{"nodes": nodes})
		return
	}

	user, _ := filter.CurrentUser(req)
	// 入库
	topic := model.NewTopic()
	topic.Uid = user["uid"].(int)
	topic.Nid = util.MustInt(req.FormValue("nid"))
	topic.Title = req.FormValue("title")
	topic.Content = req.FormValue("content")
	errMsg, err := service.PublishTopic(topic)
	if err != nil {
		fmt.Fprint(rw, `{"errno": 1, "error":"`, errMsg, `"}`)
		return
	}
	fmt.Fprint(rw, `{"errno": 0, "error":""}`)
}

// 将node组织成一定结构，方便前端展示
func genNodes() []map[string][]map[string]interface{} {
	sameParent := make(map[string][]map[string]interface{})
	allParentNodes := make([]string, 0)
	for _, node := range model.AllNode {
		if node["pid"].(int) != 0 {
			if len(sameParent[node["parent"].(string)]) == 0 {
				sameParent[node["parent"].(string)] = []map[string]interface{}{node}
			} else {
				sameParent[node["parent"].(string)] = append(sameParent[node["parent"].(string)], node)
			}
		} else {
			allParentNodes = append(allParentNodes, node["name"].(string))
		}
	}
	nodes := make([]map[string][]map[string]interface{}, 0)
	for _, parent := range allParentNodes {
		tmpMap := make(map[string][]map[string]interface{})
		tmpMap[parent] = sameParent[parent]
		nodes = append(nodes, tmpMap)
	}
	logger.Debugf("%v\n", nodes)
	return nodes
}
