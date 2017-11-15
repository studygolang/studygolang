// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package http

import (
	"bytes"
	"encoding/json"
	"global"
	"html/template"
	"logic"
	"math/rand"
	"model"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"util"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/polaris1119/config"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
	"github.com/polaris1119/times"
)

var Store = sessions.NewCookieStore([]byte(config.ConfigFile.MustValue("global", "cookie_secret")))

func SetLoginCookie(ctx echo.Context, username string) {
	Store.Options.HttpOnly = true

	session := GetCookieSession(ctx)
	if ctx.FormValue("remember_me") != "1" {
		// 浏览器关闭，cookie删除，否则保存30天(github.com/gorilla/sessions 包的默认值)
		session.Options = &sessions.Options{
			Path:     "/",
			HttpOnly: true,
		}
	}
	session.Values["username"] = username
	req := Request(ctx)
	resp := ResponseWriter(ctx)
	session.Save(req, resp)
}

func SetCookie(ctx echo.Context, key, value string) {
	Store.Options.HttpOnly = true

	session := GetCookieSession(ctx)
	session.Values[key] = value
	req := Request(ctx)
	resp := ResponseWriter(ctx)
	session.Save(req, resp)
}

func GetFromCookie(ctx echo.Context, key string) string {
	session := GetCookieSession(ctx)
	val, ok := session.Values[key]
	if ok {
		return val.(string)
	}
	return ""
}

// 必须是 http.Request
func GetCookieSession(ctx echo.Context) *sessions.Session {
	session, _ := Store.Get(Request(ctx), "user")
	return session
}

func Request(ctx echo.Context) *http.Request {
	return ctx.Request().(*standard.Request).Request
}

func ResponseWriter(ctx echo.Context) http.ResponseWriter {
	return ctx.Response().(*standard.Response).ResponseWriter
}

// 自定义模板函数
var funcMap = template.FuncMap{
	// 获取gravatar头像
	"gravatar": util.Gravatar,
	// 转为前端显示需要的时间格式
	"formatTime": func(i interface{}) string {
		ctime, ok := i.(string)
		if !ok {
			return ""
		}
		t, _ := time.Parse("2006-01-02 15:04:05", ctime)
		return t.Format(time.RFC3339) + "+08:00"
	},
	"hasPrefix": func(s, prefix string) bool {
		if strings.HasPrefix(s, prefix) {
			return true
		}
		return false
	},
	"substring": util.Substring,
	"add": func(nums ...interface{}) int {
		total := 0
		for _, num := range nums {
			if n, ok := num.(int); ok {
				total += n
			}
		}
		return total
	},
	"mod": func(num1, num2 int) int {
		if num1 == 0 {
			num1 = rand.Intn(500)
		}

		return num1 % num2
	},
	"explode": func(s, sep string) []string {
		return strings.Split(s, sep)
	},
	"noescape": func(s string) template.HTML {
		return template.HTML(s)
	},
	"timestamp": func(ts ...time.Time) int64 {
		if len(ts) > 0 {
			return ts[0].Unix()
		}
		return time.Now().Unix()
	},
	"distanceDay": func(i interface{}) int {
		var (
			t   time.Time
			err error
		)
		switch val := i.(type) {
		case string:
			t, err = time.ParseInLocation("2006-01-02 15:04:05", val, time.Local)
			if err != nil {
				return 0
			}
		case time.Time:
			t = val
		case model.OftenTime:
			t = time.Time(val)
		}

		return int(time.Now().Sub(t).Hours() / 24)
	},
	"canEdit":    logic.CanEdit,
	"canPublish": logic.CanPublish,
	"parseJSON": func(str string) map[string]interface{} {
		result := make(map[string]interface{})
		json.Unmarshal([]byte(str), &result)
		return result
	},
	"safeHtml": util.SafeHtml,
}

func tplInclude(file string, dot map[string]interface{}) template.HTML {
	var buffer = &bytes.Buffer{}
	tpl, err := template.New(filepath.Base(file)).Funcs(funcMap).ParseFiles(config.TemplateDir + file)
	// tpl, err := template.ParseFiles(config.TemplateDir + file)
	if err != nil {
		logger.Errorf("parse template file(%s) error:%v\n", file, err)
		return ""
	}
	err = tpl.Execute(buffer, dot)
	if err != nil {
		logger.Errorf("template file(%s) syntax error:%v", file, err)
		return ""
	}
	return template.HTML(buffer.String())
}

const (
	LayoutTpl      = "common/layout.html"
	AdminLayoutTpl = "common.html"
)

// 是否访问这些页面
var filterPathes = map[string]struct{}{
	"/account/login":     {},
	"/account/register":  {},
	"/account/forgetpwd": {},
	"/account/resetpwd":  {},
	"/topics/new":        {},
	"/topics/modify":     {},
	"/resources/new":     {},
	"/resources/modify":  {},
	"/articles/new":      {},
	"/articles/modify":   {},
	"/project/new":       {},
	"/project/modify":    {},
	"/book/new":          {},
	"/wiki/new":          {},
	"/wiki/modify":       {},
}

// Render html 输出
func Render(ctx echo.Context, contentTpl string, data map[string]interface{}) error {
	if data == nil {
		data = map[string]interface{}{}
	}

	objLog := logic.GetLogger(ctx)

	contentTpl = LayoutTpl + "," + contentTpl
	// 为了使用自定义的模板函数，首先New一个以第一个模板文件名为模板名。
	// 这样，在ParseFiles时，新返回的*Template便还是原来的模板实例
	htmlFiles := strings.Split(contentTpl, ",")
	for i, contentTpl := range htmlFiles {
		htmlFiles[i] = config.TemplateDir + contentTpl
	}
	tpl, err := template.New("layout.html").Funcs(funcMap).
		Funcs(template.FuncMap{"include": tplInclude}).ParseFiles(htmlFiles...)
	if err != nil {
		objLog.Errorf("解析模板出错（ParseFiles）：[%q] %s\n", Request(ctx).RequestURI, err)
		return err
	}

	data["pos_ad"] = logic.DefaultAd.FindAll(ctx, ctx.Path())
	data["cur_time"] = times.Format("Y-m-d H:i:s")
	data["path"] = ctx.Path()
	data["filter"] = false
	if _, ok := filterPathes[ctx.Path()]; ok {
		data["filter"] = true
	}

	// TODO：每次查询有点影响性能
	hasLoginMisson := false
	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		// 每日登录奖励
		hasLoginMisson = logic.DefaultMission.HasLoginMission(ctx, me)
	}
	data["has_login_misson"] = hasLoginMisson

	return executeTpl(ctx, tpl, data)
}

// RenderAdmin html 输出
func RenderAdmin(ctx echo.Context, contentTpl string, data map[string]interface{}) error {
	if data == nil {
		data = map[string]interface{}{}
	}

	objLog := logic.GetLogger(ctx)

	contentTpl = AdminLayoutTpl + "," + contentTpl
	// 为了使用自定义的模板函数，首先New一个以第一个模板文件名为模板名。
	// 这样，在ParseFiles时，新返回的*Template便还是原来的模板实例
	htmlFiles := strings.Split(contentTpl, ",")
	for i, contentTpl := range htmlFiles {
		htmlFiles[i] = config.TemplateDir + "admin/" + contentTpl
	}

	requestURI := Request(ctx).RequestURI
	tpl, err := template.New("common.html").Funcs(funcMap).ParseFiles(htmlFiles...)
	if err != nil {
		objLog.Errorf("解析模板出错（ParseFiles）：[%q] %s\n", requestURI, err)
		return err
	}

	// 当前用户信息
	curUser := ctx.Get("user").(*model.Me)

	if menu1, menu2, curMenu1 := logic.DefaultAuthority.GetUserMenu(ctx, curUser, requestURI); menu2 != nil {
		data["menu1"] = menu1
		data["menu2"] = menu2
		data["uri"] = requestURI
		data["cur_menu1"] = curMenu1
	}

	return executeTpl(ctx, tpl, data)
}

// 后台 query 查询返回结果
func RenderQuery(ctx echo.Context, contentTpl string, data map[string]interface{}) error {
	objLog := logic.GetLogger(ctx)

	contentTpl = "common_query.html," + contentTpl
	contentTpls := strings.Split(contentTpl, ",")
	for i, contentTpl := range contentTpls {
		contentTpls[i] = config.TemplateDir + "admin/" + strings.TrimSpace(contentTpl)
	}

	requestURI := Request(ctx).RequestURI
	tpl, err := template.New("common_query.html").Funcs(funcMap).ParseFiles(contentTpls...)
	if err != nil {
		objLog.Errorf("解析模板出错（ParseFiles）：[%q] %s\n", requestURI, err)
		return err
	}

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, data)
	if err != nil {
		objLog.Errorf("执行模板出错（Execute）：[%q] %s\n", requestURI, err)
		return err
	}

	return ctx.HTML(http.StatusOK, buf.String())
}

func executeTpl(ctx echo.Context, tpl *template.Template, data map[string]interface{}) error {
	objLog := logic.GetLogger(ctx)

	// 如果没有定义css和js模板，则定义之
	if jsTpl := tpl.Lookup("js"); jsTpl == nil {
		tpl.Parse(`{{define "js"}}{{end}}`)
	}
	if cssTpl := tpl.Lookup("css"); cssTpl == nil {
		tpl.Parse(`{{define "css"}}{{end}}`)
	}
	// 如果没有 seo 模板，则定义之
	if seoTpl := tpl.Lookup("seo"); seoTpl == nil {
		tpl.Parse(`{{define "seo"}}
			<meta name="keywords" content="` + logic.WebsiteSetting.SeoKeywords + `">
			<meta name="description" content="` + logic.WebsiteSetting.SeoDescription + `">
		{{end}}`)
	}

	// 当前用户信息
	curUser, ok := ctx.Get("user").(*model.Me)
	if ok {
		data["me"] = curUser
	} else {
		data["me"] = &model.Me{}
	}

	// websocket主机
	if global.OnlineEnv() {
		data["wshost"] = global.App.Domain
	} else {
		data["wshost"] = global.App.Host + ":" + global.App.Port
	}
	global.App.SetUptime()
	global.App.SetCopyright()

	isHttps := CheckIsHttps(ctx)
	cdnDomain := global.App.CDNHttp
	if isHttps {
		global.App.BaseURL = "https://" + global.App.Domain + "/"
		cdnDomain = global.App.CDNHttps
	} else {
		global.App.BaseURL = "http://" + global.App.Domain + "/"
	}

	data["app"] = global.App
	data["is_https"] = isHttps
	data["cdn_domain"] = cdnDomain
	data["use_cdn"] = config.ConfigFile.MustBool("global", "use_cdn", false)
	data["is_pro"] = global.OnlineEnv()

	data["online_users"] = map[string]int{"online": logic.Book.Len(), "maxonline": logic.MaxOnlineNum()}

	data["setting"] = logic.WebsiteSetting

	// 记录处理时间
	data["resp_time"] = time.Since(ctx.Get("req_start_time").(time.Time))

	buf := new(bytes.Buffer)
	err := tpl.Execute(buf, data)
	if err != nil {
		objLog.Errorln("excute template error:", err)
		return err
	}

	return ctx.HTML(http.StatusOK, buf.String())
}

func CheckIsHttps(ctx echo.Context) bool {
	isHttps := goutils.MustBool(ctx.Request().Header().Get("X-Https"))
	if logic.WebsiteSetting.OnlyHttps {
		isHttps = true
	}

	return isHttps
}

///////////////////////////////// APP 相关 //////////////////////////////

const (
	TokenSalt       = "b3%JFOykZx_golang_polaris"
	NeedReLoginCode = 600
)

func ParseToken(token string) (int, bool) {
	if len(token) < 32 {
		return 0, false
	}

	pos := strings.LastIndex(token, "uid")
	if pos == -1 {
		return 0, false
	}
	return goutils.MustInt(token[pos+3:]), true
}

func ValidateToken(token string) bool {
	_, ok := ParseToken(token)
	if !ok {
		return false
	}

	expireTime := time.Unix(goutils.MustInt64(token[:10]), 0)
	if time.Now().Before(expireTime) {
		return true
	}
	return false
}

func GenToken(uid int) string {
	expireTime := time.Now().Add(30 * 24 * time.Hour).Unix()

	buffer := goutils.NewBuffer().Append(expireTime).Append(uid).Append(TokenSalt)

	md5 := goutils.Md5(buffer.String())

	buffer = goutils.NewBuffer().Append(expireTime).Append(md5).Append("uid").Append(uid)
	return buffer.String()
}

func AccessControl(ctx echo.Context) {
	ctx.Response().Header().Add("Access-Control-Allow-Origin", "*")
}
