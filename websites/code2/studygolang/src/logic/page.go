// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package logic

import (
	"strings"

	"github.com/polaris1119/goutils"
)

// 每页显示多少条
const PerPage = 15

// Paginator 分页器
type Paginator struct {
	curPage int
	perPage int

	total     int
	totalPage int
}

func NewPaginator(curPage int) *Paginator {
	return NewPaginatorWithPerPage(curPage, PerPage)
}

func NewPaginatorWithPerPage(curPage, perPage int) *Paginator {
	return &Paginator{curPage: curPage, perPage: perPage}
}

// GetPageHtml 构造分页html, uri 当前uri, total 记录总数
func (this *Paginator) GetPageHtml(uri string, total ...int) string {
	if len(total) > 0 {
		this.SetTotal(int64(total[0]))
	}

	if this.totalPage < 2 {
		return ""
	}

	// 显示5页，然后显示...，接着显示最后两页
	stringBuilder := goutils.NewBuffer()
	stringBuilder.Append(`<li class="prev previous_page">`)
	// 当前是第一页
	if this.curPage != 1 {
		stringBuilder.Append(`<a href="`).Append(uri).Append("?p=").Append(this.curPage - 1).Append(`">&laquo;</a>`)
	}
	stringBuilder.Append(`</li>`)
	before := 5
	showPages := 8
	for i := 0; i < this.totalPage; i++ {
		if i >= showPages {
			break
		}
		if this.curPage == i+1 {
			stringBuilder.Append(`<li class="active"><a href="`).Append(uri).Append("?p=").Append(i + 1).Append(`">`).Append(i + 1).Append("</a></li>")
			continue
		}
		// 分界点
		if this.curPage >= before {
			if this.curPage >= 7 {
				before = 2
			} else {
				before = this.curPage + 2
			}
			showPages += 2
		}
		if i == before {
			stringBuilder.Append(`<li class="disabled"><a href="#"><span class="gap">…</span></a></li>`)
			continue
		}
		stringBuilder.Append(`<li><a href="`).Append(uri).Append("?p=").Append(i + 1).Append(`">`).Append(i + 1).Append("</a></li>")
	}
	stringBuilder.Append(`<li class="next next_page ">`)
	// 最后一页
	if this.curPage < this.totalPage {
		stringBuilder.Append(`<a href="`).Append(uri).Append("?p=").Append(this.curPage + 1).Append(`">&raquo;</a>`)
	}
	stringBuilder.Append(`</li>`)
	return stringBuilder.String()
}

func (this *Paginator) SetTotal(total int64) *Paginator {
	this.total = int(total)
	this.totalPage = this.total / this.perPage
	if this.total%this.perPage != 0 {
		this.totalPage++
	}

	return this
}

func (this *Paginator) Offset() (offset int) {
	if this.curPage > 1 {
		offset = (this.curPage - 1) * this.perPage
	}
	return
}

func (this *Paginator) SetPerPage(perPage int) {
	this.perPage = perPage
}

func (this *Paginator) PerPage() int {
	return this.perPage
}

// 构造分页html(new)
// curPage 当前页码；pageNum 每页记录数；total总记录数；uri 当前uri
func GenPageHtml(curPage, pageNum, total int, uri string) string {
	// 总页数
	pageCount := int(total / pageNum)
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
	stringBuilder := goutils.NewBuffer()
	// 当前是第一页
	if curPage != 1 {
		stringBuilder.Append(`<li><a href="`).Append(uri)
		if needQues {
			stringBuilder.Append("?")
		} else {
			stringBuilder.Append("&")
		}
		stringBuilder.Append("p=").Append(curPage - 1).Append(`">&laquo;</a>`)
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

			stringBuilder.Append("p=").Append(i + 1).Append(`">`).Append(i + 1).Append("</a></li>")
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
		stringBuilder.Append("p=").Append(i + 1).Append(`">`).Append(i + 1).Append("</a></li>")
	}

	// 最后一页
	if curPage < pageCount {
		stringBuilder.Append(`<li><a href="`).Append(uri)
		if needQues {
			stringBuilder.Append("?")
		} else {
			stringBuilder.Append("&")
		}
		stringBuilder.Append("p=").Append(curPage + 1).Append(`">&raquo;</a>`)
	} else {
		stringBuilder.Append(`<li class="disabled"><a href="#">&raquo;</a>`)
	}
	stringBuilder.Append(`</li>`)

	return stringBuilder.String()
}
