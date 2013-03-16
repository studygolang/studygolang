package main

import (
	"config"
	. "controller"
	"controller/admin"
	"filter"
	"github.com/studygolang/mux"
)

func initRouter() *mux.Router {
	// 登录校验过滤器
	loginFilter := new(filter.LoginFilter)
	loginFilterChain := mux.NewFilterChain(loginFilter)

	router := mux.NewRouter()
	// 所有的页面都需要先检查用户cookie是否存在，以便在没登录时自动登录
	cookieFilter := new(filter.CookieFilter)
	// 大部分handler都需要页面展示
	frontViewFilter := filter.NewViewFilter()
	router.FilterChain(mux.NewFilterChain([]mux.Filter{cookieFilter, frontViewFilter}...))

	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/topics{view:(|/popular|/no_reply|/last)}", TopicsHandler)
	router.HandleFunc("/topics/{tid:[0-9]+}", TopicDetailHandler)
	router.HandleFunc("/topics/new{json:(|.json)}", NewTopicHandler).AppendFilterChain(loginFilterChain)

	// 某个节点下的话题
	router.HandleFunc("/topics/node{nid:[0-9]+}", NodesHandler)

	formValidateFilter := new(filter.FormValidateFilter)
	// 注册
	router.HandleFunc("/account/register{json:(|.json)}", RegisterHandler).AppendFilterChain(mux.NewFilterChain(formValidateFilter))
	// 登录
	router.HandleFunc("/account/login", LoginHandler)
	router.HandleFunc("/account/logout", LogoutHandler)

	router.HandleFunc("/account/edit{json:(|.json)}", AccountEditHandler).AppendFilterChain(loginFilterChain)
	router.HandleFunc("/account/changepwd.json", ChangePwdHandler).AppendFilterChain(loginFilterChain)

	router.HandleFunc("/account/forgetpwd", ForgetPasswdHandler)
	router.HandleFunc("/account/resetpwd", ResetPasswdHandler)

	// 用户相关
	router.HandleFunc("/users", UsersHandler)
	router.HandleFunc("/user/{username:\\w+}", UserHomeHandler)

	// 酷站
	router.HandleFunc("/sites", SitesHandler)
	// 资源
	router.HandleFunc("/resources", ResourcesHandler)

	// 评论
	router.HandleFunc("/comment/{objid:[0-9]+}.json", CommentHandler).AppendFilterChain(loginFilterChain)

	/////////////////// 异步请求 开始///////////////////////
	// 某节点下其他帖子
	router.HandleFunc("/topics/others/{nid:[0-9]+}_{tid:[0-9]+}.json", OtherTopicsHandler)
	// 统计信息
	router.HandleFunc("/topics/stat.json", StatHandler)
	/////////////////// 异步请求 结束 ///////////////////////

	// 管理后台权限检查过滤器
	adminFilter := new(filter.AdminFilter)
	backViewFilter := filter.NewViewFilter(config.ROOT + "/template/admin/common.html")
	adminFilterChain := mux.NewFilterChain([]mux.Filter{loginFilter, adminFilter, backViewFilter}...)
	// admin 子系统
	// router.HandleFunc("/admin", admin.IndexHandler).AppendFilterChain(loginFilterChain) // 支持"/admin访问"
	subrouter := router.PathPrefix("/admin").Subrouter()
	// 所有后台需要的过滤器链
	subrouter.FilterChain(adminFilterChain)
	subrouter.HandleFunc("/", admin.IndexHandler)

	// 帖子管理
	subrouter.HandleFunc("/topics", admin.TopicsHandler)
	subrouter.HandleFunc("/nodes", admin.NodesHandler)

	// 用户管理
	subrouter.HandleFunc("/users", admin.UsersHandler)
	subrouter.HandleFunc("/newuser", admin.NewUserHandler)
	subrouter.HandleFunc("/adduser", admin.AddUserHandler)
	subrouter.HandleFunc("/profiler", admin.ProfilerHandler)

	// 错误处理handler
	router.HandleFunc("/noauthorize", NoAuthorizeHandler) // 无权限handler
	// 404页面
	router.HandleFunc("/{*}", NotFoundHandler)

	return router
}
