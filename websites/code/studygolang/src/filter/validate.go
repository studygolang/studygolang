// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package filter

import (
	"fmt"
	"github.com/studygolang/mux"
	"logger"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"util"
)

type FormValidateFilter struct {
	*mux.EmptyFilter
}

func (this *FormValidateFilter) PreFilter(rw http.ResponseWriter, req *http.Request) bool {
	req.ParseForm()
	uri := req.RequestURI
	// 后台管理中添加用户和前台注册用同样的验证规则
	if uri == "/admin/adduser.json" {
		uri = "/account/register.json"
	}
	// 验证表单
	errMsg := Validate(req.Form, Rule(uri))
	logger.Debugln("validate:", errMsg)
	// 验证失败
	// TODO:暂时规定，要验证的表单提交都采用异步方式，期望返回json数据
	if errMsg != "" {
		header := rw.Header()
		header.Set("Content-Type", "application/json; charset=utf-8")
		// 必须post方式请求
		if req.Method != "POST" {
			errMsg = "非法请求！"
		}
		fmt.Fprint(rw, `{"errno": 1,"error":"`, errMsg, `"}`)
		return false
	}
	return true
}

// 校验表单数据
func Validate(data url.Values, rules map[string]map[string]map[string]string) (errMsg string) {
	for field, rule := range rules {
		val := strings.TrimSpace(data.Get(field))
		// 检查【必填】
		if requireInfo, ok := rule["require"]; ok {
			if val == "" {
				errMsg = requireInfo["error"]
				return
			}
		}
		// 检查【长度】
		if lengthInfo, ok := rule["length"]; ok {
			valLen := len(val)
			if lenRange, ok := lengthInfo["range"]; ok {
				errMsg = checkRange(valLen, lenRange, lengthInfo["error"])
				if errMsg != "" {
					return
				}
			}
		}
		// 检查【int类型】以及可能的范围
		if intInfo, ok := rule["int"]; ok {
			valInt, err := strconv.Atoi(val)
			if err != nil {
				errMsg = field + "类型错误！"
				return
			}
			if intRange, ok := intInfo["range"]; ok {
				errMsg = checkRange(valInt, intRange, intInfo["error"])
				if errMsg != "" {
					return
				}
			}
		}
		// 检查【邮箱】
		if emailInfo, ok := rule["email"]; ok {
			validEmail := regexp.MustCompile(`^(\w)+([\w\.-])*@([\w-])+((\.[\w-]{2,3}){1,2})$`)
			if !validEmail.MatchString(val) {
				errMsg = emailInfo["error"]
				return
			}
		}
		// 检查【两值比较】
		if compareInfo, ok := rule["compare"]; ok {
			compared := compareInfo["field"] // 被比较的字段
			// 比较规则
			switch compareInfo["rule"] {
			case "=":
				if val != data.Get(compared) {
					errMsg = compareInfo["error"]
					return
				}
			case "<":
			case ">":
			default:

			}
		}
		// 检查【正则表达式】
		if regexInfo, ok := rule["regex"]; ok {
			validVal := regexp.MustCompile(regexInfo["pattern"])
			if !validVal.MatchString(val) {
				errMsg = regexInfo["error"]
				return
			}
		}
	}
	return
}

// checkRange 检查范围值是否合法。
// src为要检查的值；destRange为目标范围；msg出错时信息参数
func checkRange(src int, destRange string, msg string) (errMsg string) {
	parts := strings.SplitN(destRange, ",", 2)
	parts[0] = strings.TrimSpace(parts[0])
	parts[1] = strings.TrimSpace(parts[1])
	min, max := 0, 0
	if parts[0] == "" {
		max = util.MustInt(parts[1])
		if src > max {
			errMsg = fmt.Sprintf(msg, max)
		}
		return
	}
	if parts[1] == "" {
		min = util.MustInt(parts[0])
		if src < min {
			errMsg = fmt.Sprintf(msg, min)
		}
		return
	}
	if min == 0 {
		min = util.MustInt(parts[0])
	}
	if max == 0 {
		max = util.MustInt(parts[1])
	}
	if src < min || src > max {
		errMsg = fmt.Sprintf(msg, min, max)
		return
	}
	return
}
