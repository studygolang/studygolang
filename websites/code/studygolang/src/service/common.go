package service

import (
	"fmt"
	"net/url"
	"regexp"

	"model"
	"util"
)

func updateSetClause(form url.Values, fields []string) (query string, args []interface{}) {
	stringBuilder := util.NewBuffer()

	for _, field := range fields {
		if _, ok := form[field]; !ok {
			continue
		}
		stringBuilder.Append(",").Append(field).Append("=?")
		args = append(args, form.Get(field))
	}

	if stringBuilder.Len() > 0 {
		query = stringBuilder.String()[1:]
	}

	return
}

// 构造update语句中的set部分子句
func GenSetClause(form url.Values, fields []string) string {
	stringBuilder := util.NewBuffer()
	for _, field := range fields {
		if form.Get(field) != "" {
			stringBuilder.Append(",").Append(field).Append("=").Append(form.Get(field))
		}
	}
	if stringBuilder.Len() > 0 {
		return stringBuilder.String()[1:]
	}
	return ""
}

// @某人
func parseAtUser(content string) string {
	user := model.NewUser()

	reg := regexp.MustCompile(`@([^\s@]{4,20})`)
	return reg.ReplaceAllStringFunc(content, func(matched string) string {
		username := matched[1:]

		// 校验 username 是否存在
		err := user.Where("username=?", username).Find()
		if err != nil {
			return matched
		}

		if user.Username != username {
			return matched
		}
		return fmt.Sprintf(`<a href="/user/%s" title="%s">%s</a>`, username, matched, matched)
	})
}
