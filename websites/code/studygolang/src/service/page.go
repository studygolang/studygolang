package service

import (
	"util"
)

// 每页显示多少条
const PAGE_NUM = 15

// 构造分页html
// curPage 当前页码；total总记录数
func GetPageHtml(curPage, total int) string {
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
	stringBuilder.Append(`<li class="prev previous_page disabled"><a href="#">← 上一页</a></li>`)
	before := 5
	showPages := 8
	for i := 0; i < pageCount; i++ {
		if i >= showPages {
			break
		}
		if curPage == i+1 {
			stringBuilder.Append(`<li class="active"><a href="/topics?p=`).AppendInt(i + 1).Append(`">`).AppendInt(i + 1).Append("</a></li>")
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
		stringBuilder.Append(`<li><a href="/topics?p=`).AppendInt(i + 1).Append(`">`).AppendInt(i + 1).Append("</a></li>")
	}
	stringBuilder.Append(`<li class="next next_page "><a rel="next" href="/topics?page=2">下一页 →</a></li>`)
	return stringBuilder.String()
}
