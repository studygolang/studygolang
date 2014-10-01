package service

import (
	"net/url"
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
