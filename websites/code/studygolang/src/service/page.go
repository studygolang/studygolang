// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"strings"
	"util"
)

// 每页显示多少条
const PAGE_NUM = 15

// 构造分页html
// curPage 当前页码；total总记录数；uri 当前uri
func GetPageHtml(curPage, total int, uri string) string {
	// 总页数
	pageCount := total / PAGE_NUM
	if total%PAGE_NUM != 0 {
		pageCount++
	}
	if pageCount < 2 {
		return ""
	}
	// 显示5页，然后显示...，接着显示最后两页
	stringBuilder := util.NewBuffer()
	stringBuilder.Append(`<li class="prev previous_page">`)
	// 当前是第一页
	if curPage != 1 {
		stringBuilder.Append(`<a href="`).Append(uri).Append("?p=").AppendInt(curPage - 1).Append(`">&laquo;</a>`)
	}
	stringBuilder.Append(`</li>`)
	before := 5
	showPages := 8
	for i := 0; i < pageCount; i++ {
		if i >= showPages {
			break
		}
		if curPage == i+1 {
			stringBuilder.Append(`<li class="active"><a href="`).Append(uri).Append("?p=").AppendInt(i + 1).Append(`">`).AppendInt(i + 1).Append("</a></li>")
			continue
		}
		// 分界点
		if curPage >= before {
			if curPage >= 7 {
				before = 2
			} else {
				before = curPage + 2
			}
			showPages += 2
		}
		if i == before {
			stringBuilder.Append(`<li class="disabled"><a href="#"><span class="gap">…</span></a></li>`)
			continue
		}
		stringBuilder.Append(`<li><a href="`).Append(uri).Append("?p=").AppendInt(i + 1).Append(`">`).AppendInt(i + 1).Append("</a></li>")
	}
	stringBuilder.Append(`<li class="next next_page ">`)
	// 最后一页
	if curPage < pageCount {
		stringBuilder.Append(`<a href="`).Append(uri).Append("?p=").AppendInt(curPage + 1).Append(`">&raquo;</a>`)
	}
	stringBuilder.Append(`</li>`)
	return stringBuilder.String()
}

// 构造分页html(new)
// curPage 当前页码；pageNum 每页记录数；total总记录数；uri 当前uri
func GenPageHtml(curPage, pageNum, total int, uri string) string {
	// 总页数
	pageCount := total / pageNum
	if total%pageNum != 0 {
		pageCount++
	}
	if pageCount < 2 {
		return ""
	}

	needQues := true
	if strings.Contains(uri, "?") {
		needQues = false
	}

	// 显示5页，然后显示...，接着显示最后两页
	stringBuilder := util.NewBuffer()
	// 当前是第一页
	if curPage != 1 {
		stringBuilder.Append(`<li><a href="`).Append(uri)
		if needQues {
			stringBuilder.Append("?")
		} else {
			stringBuilder.Append("&")
		}
		stringBuilder.Append("p=").AppendInt(curPage - 1).Append(`">&laquo;</a>`)
	} else {
		stringBuilder.Append(`<li class="disabled"><a href="#">&laquo;</a>`)
	}
	stringBuilder.Append(`</li>`)

	before := 5
	showPages := 8
	for i := 0; i < pageCount; i++ {
		if i >= showPages {
			break
		}
		if curPage == i+1 {
			stringBuilder.Append(`<li class="active"><a href="`).Append(uri)
			if needQues {
				stringBuilder.Append("?")
			} else {
				stringBuilder.Append("&")
			}

			stringBuilder.Append("p=").AppendInt(i + 1).Append(`">`).AppendInt(i + 1).Append("</a></li>")
			continue
		}
		// 分界点
		if curPage >= before {
			if curPage >= 7 {
				before = 2
			} else {
				before = curPage + 2
			}
			showPages += 2
		}
		if i == before {
			stringBuilder.Append(`<li class="disabled"><a href="#"><span class="gap">…</span></a></li>`)
			continue
		}
		stringBuilder.Append(`<li><a href="`).Append(uri)
		if needQues {
			stringBuilder.Append("?")
		} else {
			stringBuilder.Append("&")
		}
		stringBuilder.Append("p=").AppendInt(i + 1).Append(`">`).AppendInt(i + 1).Append("</a></li>")
	}

	// 最后一页
	if curPage < pageCount {
		stringBuilder.Append(`<li><a href="`).Append(uri)
		if needQues {
			stringBuilder.Append("?")
		} else {
			stringBuilder.Append("&")
		}
		stringBuilder.Append("p=").AppendInt(curPage + 1).Append(`">&raquo;</a>`)
	} else {
		stringBuilder.Append(`<li class="disabled"><a href="#">&raquo;</a>`)
	}
	stringBuilder.Append(`</li>`)

	return stringBuilder.String()
}
