package admin

import (
	"html/template"
	"fmt"
	"net/http"
	"net/url"
)

func UsersHandler(rw http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles(ROOT + "/template/admin/common.html", ROOT + "/template/admin/users.html")
	if err != nil {
		fmt.Fprintln(rw, err)
		return
	}
	tpl.Execute(rw, nil)
}

// 添加新用户表单页面
func NewUserHandler(rw http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles(ROOT + "/template/admin/common.html", ROOT + "/template/admin/newuser.html")
	if err != nil {
		fmt.Fprintln(rw, err)
		return
	}
	tpl.Execute(rw, nil)
}

func Validate(data url.Values, rules map[string]map[string]string) (errMsg string) {
	for field, rule := range rules {
	    val := data.Get(field)
	    if _, ok := rule["require"]; ok {
	        if val == "" {
	            errMsg = field + "不能为空"
		    return
	        }
	    }

	    if valRange, ok := rule["range"]; ok {
	        
	    }
	}
}

// 执行添加新用户（异步请求，返回json）
func AddUserHandler(rw http.ResponseWriter, req *http.Request) {
	errMsg := ""

	username := req.FormValue("username")
	email := req.FormValue("email")
	// name := req.FormValue("name")
	passwd := req.FormValue("passwd")
	pass2 := req.FormValue("pass2")
	// send_password := req.FormValue("send_password")
	
	if username == "" {
	    errMsg = "用户名不能空!"
	}
	
	if email == "" {
	    errMsg = "邮箱不能空!"
	}
	
	if passwd == "" {
	    errMsg = "密码不能空!"
	}
	
	if passwd != pass2 {
	    errMsg = "两次密码不一致!"
	}
	// TODO: 其他判断
	header := rw.Header()
	header.Set("Content-Type", "application/json; charset=utf-8")
	if errMsg != "" {
	    fmt.Fprint(rw, `{"error":"`, errMsg, `"}`)
	    return
	}
	
	// 入库
}

func ProfilerHandler(rw http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles(ROOT + "/template/admin/common.html", ROOT + "/template/admin/profiler.html")
	if err != nil {
		fmt.Fprintln(rw, err)
		return
	}
	tpl.Execute(rw, nil)
}