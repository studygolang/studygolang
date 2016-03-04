package controller

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"
	"time"
	"util"

	"github.com/labstack/echo"
	"github.com/polaris1119/config"
	"github.com/polaris1119/logger"
)

const (
	TplFileKey = "tpl"

	DataKey = "data"
)

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
	"explode": func(s, sep string) []string {
		return strings.Split(s, sep)
	},
	"noescape": func(s string) template.HTML {
		return template.HTML(s)
	},
}

// render html 输出
func render(ctx *echo.Context) error {
	objLog := ctx.Get("logger").(*logger.Logger)

	contentFile := ctx.Get(TplFileKey).(string)
	// 为了使用自定义的模板函数，首先New一个以第一个模板文件名为模板名。
	// 这样，在ParseFiles时，新返回的*Template便还是原来的模板实例
	htmlFiles := []string{config.TemplateDir + "common/layout.html", config.TemplateDir + contentFile}
	tpl, err := template.New("layout.html").Funcs(funcMap).ParseFiles(htmlFiles...)
	if err != nil {
		objLog.Errorf("解析模板出错（ParseFiles）：[%q] %s\n", ctx.Request().RequestURI, err)
		return err
	}
	// 如果没有定义css和js模板，则定义之
	if jsTpl := tpl.Lookup("js"); jsTpl == nil {
		tpl.Parse(`{{define "js"}}{{end}}`)
	}
	if jsTpl := tpl.Lookup("css"); jsTpl == nil {
		tpl.Parse(`{{define "css"}}{{end}}`)
	}

	data := ctx.Get(DataKey).(map[string]interface{})

	// 当前用户信息
	// me, _ := CurrentUser(req)
	data["me"] = map[string]interface{}{}

	// websocket主机
	data["wshost"] = "127.0.0.1"
	data["build"] = map[string]string{
		"version": "1.0",        // version.Version,
		"date":    "2016-01-16", // version.Date,
	}

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, data)
	if err != nil {
		objLog.Errorln("excute template error:", err)
		return err
	}

	return ctx.HTML(http.StatusOK, buf.String())
}
