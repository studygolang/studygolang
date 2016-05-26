package controller

import (
	"bytes"
	"db"
	"html/template"
	"net/http"

	"github.com/labstack/echo"
	"github.com/polaris1119/config"
	"github.com/polaris1119/goutils"
)

type InstallController struct{}

// 注册路由
func (self InstallController) RegisterRoute(g *echo.Group) {
	g.Get("/install", echo.HandlerFunc(self.Install))
	g.Match([]string{"GET", "POST"}, "/install/setup-config", echo.HandlerFunc(self.SetupConfig))
	g.Match([]string{"GET", "POST"}, "/install/do", echo.HandlerFunc(self.DoInstall))
}

func (InstallController) Install(ctx echo.Context) error {
	return renderInstall(ctx, "install/index.html", nil)
}

func (self InstallController) SetupConfig(ctx echo.Context) error {
	step := goutils.MustInt(ctx.QueryParam("step"))
	if step == 2 {
		err := self.genConfig(ctx)
		if err != nil {
			data := map[string]interface{}{"dbhost": ctx.FormValue("dbhost"), "dbport": ctx.FormValue("dbport")}
			return renderInstall(ctx, "isntall/setup-err.html", data)
		}
	}
	return renderInstall(ctx, "install/setup-config.html", map[string]interface{}{"step": step})
}

// DoInstall 执行安装，包括站点简单配置，安装数据库（创建数据库、表，填充基本数据）等
func (InstallController) DoInstall(ctx echo.Context) error {
	return nil
}

func (InstallController) genConfig(ctx echo.Context) error {
	env := ctx.FormValue("env")

	config.ConfigFile.SetSectionComments("global", "")
	config.ConfigFile.SetValue("global", "env", env)

	var (
		logLevel     = "DEBUG"
		domain       = "127.0.0.1"
		xormLogLevel = "0"
		xormShowSql  = "true"
	)
	if env == "pro" {
		logLevel = "INFO"
		domain = "studygolang.com"
		xormLogLevel = "1"
		xormShowSql = "false"
	}

	config.ConfigFile.SetValue("global", "log_level", logLevel)
	config.ConfigFile.SetValue("global", "domain", domain)
	config.ConfigFile.SetValue("global", "cookie_secret", goutils.RandString(10))
	config.ConfigFile.SetValue("global", "data_path", "data/max_online_num")

	config.ConfigFile.SetSectionComments("listen", "")
	config.ConfigFile.SetValue("listen", "host", "localhost")
	config.ConfigFile.SetValue("listen", "port", "8088")

	dbname := ctx.FormValue("dbname")
	uname := ctx.FormValue("uname")
	pwd := ctx.FormValue("pwd")
	dbhost := ctx.FormValue("dbhost")
	dbport := ctx.FormValue("dbport")

	config.ConfigFile.SetSectionComments("mysql", "")
	config.ConfigFile.SetValue("mysql", "host", dbhost)
	config.ConfigFile.SetValue("mysql", "port", dbport)
	config.ConfigFile.SetValue("mysql", "user", uname)
	config.ConfigFile.SetValue("mysql", "password", pwd)
	config.ConfigFile.SetValue("mysql", "dbname", dbname)
	config.ConfigFile.SetValue("mysql", "charset", "utf8")
	config.ConfigFile.SetKeyComments("mysql", "max_idle", "最大空闲连接数")
	config.ConfigFile.SetValue("mysql", "max_idle", "2")
	config.ConfigFile.SetKeyComments("mysql", "max_conn", "最大打开连接数")
	config.ConfigFile.SetValue("mysql", "max_conn", "27")

	config.ConfigFile.SetSectionComments("xorm", "")
	config.ConfigFile.SetValue("xorm", "show_sql", xormShowSql)
	config.ConfigFile.SetKeyComments("xorm", "log_level", "0-debug, 1-info, 2-warning, 3-error, 4-off, 5-unknow")
	config.ConfigFile.SetValue("xorm", "log_level", xormLogLevel)

	config.ConfigFile.SetSectionComments("security", "")
	config.ConfigFile.SetKeyComments("security", "unsubscribe_token_key", "退订邮件使用的 token key")
	config.ConfigFile.SetValue("security", "unsubscribe_token_key", goutils.RandString(18))
	config.ConfigFile.SetKeyComments("security", "activate_sign_salt", "注册激活邮件使用的 sign salt")
	config.ConfigFile.SetValue("security", "activate_sign_salt", goutils.RandString(18))

	config.ConfigFile.SetSectionComments("sensitive", "过滤广告")
	config.ConfigFile.SetKeyComments("sensitive", "title", "标题关键词")
	config.ConfigFile.SetValue("sensitive", "title", "")
	config.ConfigFile.SetKeyComments("sensitive", "content", "内容关键词")
	config.ConfigFile.SetValue("sensitive", "content", "")

	config.SaveConfigFile()

	return db.Init()
}

func renderInstall(ctx echo.Context, filename string, data map[string]interface{}) error {
	objLog := getLogger(ctx)

	if data == nil {
		data = make(map[string]interface{})
	}

	filename = config.TemplateDir + filename

	requestURI := ctx.Request().URI()
	tpl, err := template.ParseFiles(filename)
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
