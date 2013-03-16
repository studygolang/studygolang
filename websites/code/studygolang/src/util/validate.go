package util

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// 校验表单数据
func Validate(data url.Values, rules map[string]map[string]map[string]string) (errMsg string) {
	for field, rule := range rules {
		val := data.Get(field)
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
			validEmail := regexp.MustCompile(`^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+((\.[a-zA-Z0-9_-]{2,3}){1,2})$`)
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
		max = MustInt(parts[1])
		if src > max {
			errMsg = fmt.Sprintf(msg, max)
			return
		}
	}
	if parts[1] == "" {
		min = MustInt(parts[0])
		if src < min {
			errMsg = fmt.Sprintf(msg, min)
			return
		}
	}
	if min == 0 {
		min = MustInt(parts[0])
	}
	if max == 0 {
		max = MustInt(parts[1])
	}
	if src < min || src > max {
		errMsg = fmt.Sprintf(msg, min, max)
		return
	}
	return
}
