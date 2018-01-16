// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:polaris	polaris@studygolang.com

package logic

import (
	"strings"

	"github.com/polaris1119/goutils"
)

// 每页显示多少条
const PerPage = 50

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

const maxShow = 5

// GetPageHtml 构造分页html, uri 当前uri, total 记录总数
func (this *Paginator) GetPageHtml(uri string, total ...int) string {
	if len(total) > 0 {
		this.SetTotal(int64(total[0]))
	}

	if this.totalPage < 2 {
		return ""
	}

	if !strings.Contains(uri, "?") {
		uri += "?"
	}

	// 显示5页，然后显示...，接着显示最后两页
	stringBuilder := goutils.NewBuffer()
	stringBuilder.Append(`<li class="prev previous_page">`)
	// 当前是第一页
	if this.curPage != 1 {
		stringBuilder.Append(`<a href="`).Append(uri).Append("p=").Append(this.curPage - 1).Append(`">&laquo;</a>`)
	}
	stringBuilder.Append(`</li>`)

	// 当前页前后至少保持两页可以翻看；
	// 离最开始或最后5页内，不显示 ...，否则显示 ...
	if this.curPage <= maxShow {
		// 当前页之前的页
		for i := 0; i < this.curPage-1; i++ {
			this.appendPage(stringBuilder, uri, i+1)
		}
		// 当前页
		this.appendCurPage(stringBuilder, uri)

		// 当前页之后的页

		// 全部显示
		if this.totalPage-this.curPage <= maxShow {
			for i := this.curPage; i < this.totalPage; i++ {
				this.appendPage(stringBuilder, uri, i+1)
			}
		} else {
			// 显示当前页后两页、...和最后两页
			for i := this.curPage; i < this.curPage+2; i++ {
				this.appendPage(stringBuilder, uri, i+1)
			}
			stringBuilder.Append(`<li class="disabled"><a href="#"><span class="gap">…</span></a></li>`)
			for i := this.totalPage - 2; i < this.totalPage; i++ {
				this.appendPage(stringBuilder, uri, i+1)
			}
		}
	} else {
		// 显示最开始两页，然后是...和挨着当前页的两页
		for i := 0; i < 2; i++ {
			this.appendPage(stringBuilder, uri, i+1)
		}
		stringBuilder.Append(`<li class="disabled"><a href="#"><span class="gap">…</span></a></li>`)
		for i := this.curPage - 2; i < this.curPage; i++ {
			this.appendPage(stringBuilder, uri, i)
		}
		// 当前页
		this.appendCurPage(stringBuilder, uri)

		// 当前页之后的页

		// 全部显示
		if this.totalPage-this.curPage <= maxShow {
			for i := this.curPage; i < this.totalPage; i++ {
				this.appendPage(stringBuilder, uri, i+1)
			}
		} else {
			// 显示当前页后两页、...和最后两页
			for i := this.curPage; i < this.curPage+2; i++ {
				this.appendPage(stringBuilder, uri, i+1)
			}
			stringBuilder.Append(`<li class="disabled"><a href="#"><span class="gap">…</span></a></li>`)
			for i := this.totalPage - 2; i < this.totalPage; i++ {
				this.appendPage(stringBuilder, uri, i+1)
			}
		}
	}

	// 处理next
	stringBuilder.Append(`<li class="next next_page ">`)
	if this.curPage < this.totalPage {
		stringBuilder.Append(`<a href="`).Append(uri).Append("p=").Append(this.curPage + 1).Append(`">&raquo;</a>`)
	}
	stringBuilder.Append(`</li>`)
	return stringBuilder.String()
}

func (this *Paginator) appendPage(stringBuilder *goutils.Buffer, uri string, page int) {
	stringBuilder.Append(`<li><a href="`).Append(uri).Append("p=").Append(page).Append(`">`).Append(page).Append("</a></li>")
}

// appendCurPage 当前页
func (this *Paginator) appendCurPage(stringBuilder *goutils.Buffer, uri string) {
	stringBuilder.Append(`<li class="active"><a href="`).Append(uri).Append("p=").Append(this.curPage).Append(`">`).Append(this.curPage).Append("</a></li>")
}

func (this *Paginator) SetTotal(total int64) *Paginator {
	this.total = int(total)
	this.totalPage = this.total / this.perPage
	if this.total%this.perPage != 0 {
		this.totalPage++
	}

	return this
}

func (this *Paginator) GetTotal() int {
	return this.total
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

func (this *Paginator) HasMorePage() bool {
	return this.totalPage > this.curPage
}
