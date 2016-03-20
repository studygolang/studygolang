package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"logic"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"util"

	"vendor/github.com/polaris1119/goutils"

	"github.com/labstack/echo"
	"github.com/polaris1119/logger"
)

// 侧边栏的内容通过异步请求获取
type SidebarController struct{}

func (this *SidebarController) RegisterRoutes(e *echo.Echo) {
	e.Get("/readings/recent", echo.HandlerFunc(this.RecentReading))
	e.Get("/topics/:nid/others", echo.HandlerFunc(this.OtherTopics))
	e.Get("/websites/stat", echo.HandlerFunc(this.WebsiteStat))
	e.Get("/dymanics/recent", echo.HandlerFunc(this.RecentDymanic))
	e.Get("/topics/recent", echo.HandlerFunc(this.RecentTopic))
	e.Get("/articles/recent", echo.HandlerFunc(this.RecentArticle))
	e.Get("/projects/recent", echo.HandlerFunc(this.RecentProject))
	e.Get("/resources/recent", echo.HandlerFunc(this.RecentResource))
}

// RecentReading 技术晨读
func (SidebarController) RecentReading(ctx echo.Context) error {
	// limit := goutils.MustInt(ctx.Query("limit"), 7)

	// readings := service.FindReadings("0", limit, model.RtypeGo)
	// buf, err := json.Marshal(readings)
	// if err != nil {
	// 	logger.Errorln("[RecentReadingHandler] json.marshal error:", err)
	// 	fmt.Fprint(rw, `{"ok": 0, "error":"解析json出错"}`)
	// 	return
	// }
	// fmt.Fprint(rw, `{"ok": 1, "data":`+string(buf)+`}`)
	return nil
}

// OtherTopics 某节点下其他帖子
func (SidebarController) OtherTopics(ctx echo.Context) error {
	topics := logic.DefaultTopic.FindByNid(ctx, ctx.Param("nid"), ctx.Query("tid"))
	topics = logic.DefaultTopic.JSEscape(topics)
	return success(ctx, topics)
}

// WebsiteStat 网站统计信息
func (SidebarController) WebsiteStat(ctx echo.Context) error {
	articleTotal := logic.DefaultArticle.Total()
	projectTotal := logic.DefaultProject.Total()
	topicTotal := logic.DefaultTopic.Total()
	cmtTotal := logic.DefaultComment.Total()
	resourceTotal := logic.DefaultResource.Total()
	userTotal := logic.DefaultUser.Total()

	data := map[string]interface{}{
		"article":  articleTotal,
		"project":  projectTotal,
		"topic":    topicTotal,
		"resource": resourceTotal,
		"comment":  cmtTotal,
		"user":     userTotal,
	}

	return success(ctx, data)
}

// RecentDymanic 社区最新公告或go最新动态
func (SidebarController) RecentDymanic(ctx echo.Context) error {
	dynamics := logic.DefaultDynamic.FindBy(ctx, 0, 3)
	return success(ctx, dynamics)
}

// RecentTopic 最新帖子
func (SidebarController) RecentTopic(ctx echo.Context) error {
	limit := goutils.MustInt(ctx.Query("limit"), 10)
	topicList := logic.DefaultTopic.FindRecent(limit)
	return success(ctx, topicList)
}

// RecentArticle 最新博文
func (SidebarController) RecentArticle(ctx echo.Context) error {
	limit := goutils.MustInt(ctx.Query("limit"), 10)
	recentArticles := logic.DefaultArticle.FindBy(limit)
	return success(ctx, recentArticles)
}

// RecentProject 最新开源项目
func (SidebarController) RecentProject(ctx echo.Context) error {
	limit := goutils.MustInt(ctx.Query("limit"), 10)
	recentProjects := logic.DefaultProject.FindBy(limit)
	return success(ctx, recentProjects)
}

// RecentResource 最新资源
func (SidebarController) RecentResource(ctx echo.Context) error {
	limit := goutils.MustInt(ctx.Query("limit"), 10)
	recentResources := logic.DefaultResource.FindBy(limit)
	return success(ctx, recentResources)
}

// 最新评论
// uri: /comments/recent.json
func RecentCommentHandler(rw http.ResponseWriter, req *http.Request) {
	limit := req.FormValue("limit")
	if limit == "" {
		limit = "10"
	}
	recentComments := service.FindRecentComments(0, -1, limit)

	uids := util.Models2Intslice(recentComments, "Uid")
	users := service.GetUserInfos(uids)

	result := map[string]interface{}{
		"comments": recentComments,
	}

	// json encode 不支持 map[int]...
	for uid, user := range users {
		result[strconv.Itoa(uid)] = user
	}

	buf, err := json.Marshal(result)

	if err != nil {
		logger.Errorln("[RecentCommentHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"ok": 0, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":`+string(buf)+`}`)
}

// 社区热门节点
// uri: /nodes/hot.json
func HotNodesHandler(rw http.ResponseWriter, req *http.Request) {
	nodes := service.FindHotNodes()
	hotNodes, err := json.Marshal(nodes)
	if err != nil {
		logger.Errorln("[HotNodesHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"ok": 0, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":`+string(hotNodes)+`}`)
}

// 活跃会员
// uri: /users/active.json
func ActiveUserHandler(rw http.ResponseWriter, req *http.Request) {
	activeUsers := service.FindActiveUsers(0, 9)
	buf, err := json.Marshal(activeUsers)
	if err != nil {
		logger.Errorln("[ActiveUserHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"ok": 0, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":`+string(buf)+`}`)
}

// 新加入会员
// uri: /users/newest.json
func NewestUserHandler(rw http.ResponseWriter, req *http.Request) {
	newestUsers := service.FindNewUsers(0, 9)
	buf, err := json.Marshal(newestUsers)
	if err != nil {
		logger.Errorln("[NewestUserHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"ok": 0, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"ok": 1, "data":`+string(buf)+`}`)
}

// 评论或回复 @ 某人 suggest
// uri: /at/users.json
func AtUsersHandler(rw http.ResponseWriter, req *http.Request) {
	term := req.FormValue("term")
	users := service.GetUserMentions(term, 10)
	buf, err := json.Marshal(users)
	if err != nil {
		logger.Errorln("[AtUsersHandler] json.marshal error:", err)
		fmt.Fprint(rw, `[]`)
		return
	}

	fmt.Fprint(rw, string(buf))
}

const maxImageSize = 5 << 20 // 5M

func UploadImageHandler(rw http.ResponseWriter, req *http.Request) {
	var (
		uri    string
		buf    []byte
		err    error
		reader io.Reader
	)

	origUrl := req.FormValue("url")
	if origUrl != "" {
		resp, err := http.Get(origUrl)
		if err != nil {
			fmt.Fprint(rw, `{"ok": 0, "error":"获取图片失败"}`)
			return
		}
		defer resp.Body.Close()

		buf, err := ioutil.ReadAll(resp.Body)

		ext := filepath.Ext(origUrl)
		if ext == "" {
			contentType := http.DetectContentType(buf)
			exts, _ := mime.ExtensionsByType(contentType)
			if len(exts) > 0 {
				ext = exts[0]
			}
		}

		uri = util.DateNow() + "/" + util.Md5Buf(buf) + ext

		reader = bytes.NewReader(buf)
	} else {

		file, fileHeader, err := req.FormFile("img")
		if err != nil {
			fmt.Fprint(rw, `{"ok": 0, "error":"非法文件上传！"}`)
			return
		}

		defer file.Close()

		// 如果是临时文件，存在硬盘中，则是 *os.File（大于32M），直接报错
		if _, ok := file.(*os.File); ok {
			fmt.Fprint(rw, `{"ok": 0, "error":"文件太大！"}`)
			return
		}

		reader = file

		buf, err := ioutil.ReadAll(file)
		imgDir := util.DateNow()
		if req.FormValue("avatar") != "" {
			imgDir = "avatar"
		}
		uri = imgDir + "/" + util.Md5Buf(buf) + filepath.Ext(fileHeader.Filename)
	}

	if err != nil {
		fmt.Fprint(rw, `{"ok": 0, "error":"文件读取失败！"}`)
		return
	}

	if len(buf) > maxImageSize {
		fmt.Fprint(rw, `{"ok": 0, "error":"文件太大！"}`)
		return
	}

	err = service.UploadMemoryFile(reader, uri)
	if err != nil {
		fmt.Fprint(rw, `{"ok": 0, "error":"文件上传失败！"}`)
		return
	}

	if origUrl != "" {
		uri = "http://studygolang.qiniudn.com/" + uri
	}

	fmt.Fprint(rw, `{"ok": 1, "uri":"`+uri+`"}`)
}
