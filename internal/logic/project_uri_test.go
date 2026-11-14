package logic

import "testing"

func TestProjectUriSlug(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"空格转连字符", "Go Web Framework", "go-web-framework"},
		{"斜杠不能进路径段", "gin-gonic/gin", "gin-gonic-gin"},
		{"井号会截断 fragment", "foo#bar", "foo-bar"},
		{"问号会截断 query", "foo?bar", "foo-bar"},
		{"百分号会破坏解码", "100%go", "100-go"},
		{"保留中文", "中文项目", "中文项目"},
		{"中英混排", "Go 语言中文网", "go-语言中文网"},
		{"保留点和下划线", "go.v1_beta", "go.v1_beta"},
		{"连续特殊字符折叠成一个连字符", "a  --  b", "a-b"},
		{"去掉首尾连字符", "  hello  ", "hello"},
		{"全特殊字符返回空", "???///", ""},
		{"空串", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := projectUriSlug(tt.in); got != tt.want {
				t.Errorf("projectUriSlug(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestProjectUriSlugIdempotent 已规范化的 uri 再跑一次必须不变，
// 否则编辑项目会不断改写 URL。
func TestProjectUriSlugIdempotent(t *testing.T) {
	for _, in := range []string{"Go Web Framework", "中文项目", "go.v1_beta", "gin-gonic/gin"} {
		once := projectUriSlug(in)
		if twice := projectUriSlug(once); twice != once {
			t.Errorf("projectUriSlug 不幂等：%q -> %q -> %q", in, once, twice)
		}
	}
}
