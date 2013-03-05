package main

import (
	. "controller"
	"controller/admin"
	"github.com/gorilla/mux"
)

func initRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/topics/new", NewTopicHandler)
	// admin 子系统
	router.HandleFunc("/admin", admin.IndexHandler) // 支持/admin访问
	subrouter := router.PathPrefix("/admin").Subrouter()
	subrouter.HandleFunc("/", admin.IndexHandler)
	// 用户管理
	subrouter.HandleFunc("/users", admin.UsersHandler)
	subrouter.HandleFunc("/newuser", admin.NewUserHandler)
	subrouter.HandleFunc("/adduser", admin.AddUserHandler)
	subrouter.HandleFunc("/profiler", admin.ProfilerHandler)
	return router
}
