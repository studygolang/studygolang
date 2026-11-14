package api

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/studygolang/studygolang/internal/model"

	studyhttp "github.com/studygolang/studygolang/internal/http"

	"github.com/fatih/structs"
	echo "github.com/labstack/echo/v4"
)

func TestIsValidObjType(t *testing.T) {
	tests := []struct {
		name     string
		objtype  int
		expected bool
	}{
		// 合法枚举值
		{"topic=0 is valid (regression for C3)", 0, true},
		{"article", 1, true},
		{"resource", 2, true},
		{"wiki", 3, true},
		{"project", 4, true},
		{"book", 5, true},
		{"interview", 6, true},
		// 非法值
		{"-1 rejected", -1, false},
		{"7 rejected (out of enum)", 7, false},
		{"100 rejected (out of enum)", 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidObjType(tt.objtype); got != tt.expected {
				t.Errorf("isValidObjType(%d) = %v, want %v", tt.objtype, got, tt.expected)
			}
		})
	}
}

// TestNormalizeReplies 锁定评论 map 的 key 大小写。
// 生产环境 cmd/studygolang/main.go 的 init() 把 structs.DefaultTagName 设为 "json"，
// 若按 Go 字段名（大写）取值，id/uid/floor/ctime 会全部为 nil，
// 前端点赞、回复、编辑随之失效——这是静默故障，必须由测试守住。
func TestNormalizeReplies(t *testing.T) {
	origTagName := structs.DefaultTagName
	structs.DefaultTagName = "json"
	defer func() { structs.DefaultTagName = origTagName }()

	comment := &model.Comment{
		Cid:     12,
		Objid:   34,
		Objtype: model.TypeTopic,
		Content: "原始文本 @polaris",
		Uid:     56,
		Floor:   7,
		Likenum: 3,
		Flag:    model.FlagNormal,
	}

	// 复刻 CommentLogic.FindObjComments 的组装顺序
	raw := structs.Map(comment)
	raw["raw_content"] = comment.Content
	raw["content"] = template.HTML(`<a href="/user/polaris">@polaris</a>`)
	raw["user"] = &model.User{Username: "polaris", Avatar: "avatar.png"}

	ctx := echo.New().NewContext(
		httptest.NewRequest(http.MethodGet, "/", nil),
		httptest.NewRecorder(),
	)

	got := normalizeReplies(ctx, []map[string]interface{}{raw}, 0)
	if len(got) != 1 {
		t.Fatalf("normalizeReplies 返回 %d 条，期望 1 条", len(got))
	}

	item := got[0]
	wantInts := map[string]int{
		"id":      12,
		"objid":   34,
		"objtype": model.TypeTopic,
		"uid":     56,
		"floor":   7,
		"likenum": 3,
		"flag":    model.FlagNormal,
	}
	for key, want := range wantInts {
		if item[key] != want {
			t.Errorf("%s = %v, want %d", key, item[key], want)
		}
	}

	wantStrings := map[string]string{
		"raw_content": "原始文本 @polaris",
		"name":        "polaris",
		"username":    "polaris",
		"avatar":      "avatar.png",
	}
	for key, want := range wantStrings {
		if item[key] != want {
			t.Errorf("%s = %v, want %q", key, item[key], want)
		}
	}

	if _, ok := item["ctime"]; !ok {
		t.Error("缺少 ctime 字段")
	}
	if item["liked"] != nil {
		t.Errorf("未登录时不应返回 liked，got %v", item["liked"])
	}
}

func TestNormalizeRepliesEmpty(t *testing.T) {
	ctx := echo.New().NewContext(
		httptest.NewRequest(http.MethodGet, "/", nil),
		httptest.NewRecorder(),
	)

	got := normalizeReplies(ctx, nil, 0)
	if got == nil {
		t.Fatal("空输入应返回空切片而非 nil，否则前端 JSON 拿到 null")
	}
	if len(got) != 0 {
		t.Fatalf("空输入返回 %d 条，期望 0 条", len(got))
	}
}

// TestFailReturnsNonNilError 锁定 fail() 的错误语义。
// 曾经 fail() 直接返回 ctx.JSON(...)，而 Echo 写入成功时返回 nil，
// 导致 `return 0, fail(...)` 退化成 (0, nil)，所有 `if err != nil` 鉴权守卫失效，
// 未登录请求会继续跑业务逻辑并二次写响应。
func TestFailReturnsNonNilError(t *testing.T) {
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)

	if err := fail(ctx, "请先登录", studyhttp.NeedReLoginCode); err == nil {
		t.Fatal("fail() 必须返回非 nil error，否则调用方的 err != nil 守卫会穿透")
	}

	// 响应体只能有一份
	body := rec.Body.String()
	if strings.Count(body, "\"code\"") != 1 {
		t.Errorf("响应体应只含一个 code 字段，got %q", body)
	}
}

// TestParseAuthUIDRejectsAnonymous 未登录时 parseAuthUID 必须返回错误，
// 且 uid 为 0，调用方据此中断。
func TestParseAuthUIDRejectsAnonymous(t *testing.T) {
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)

	uid, err := parseAuthUID(ctx)
	if err == nil {
		t.Fatal("未登录时 parseAuthUID 必须返回错误")
	}
	if uid != 0 {
		t.Errorf("uid = %d, want 0", uid)
	}
	if !strings.Contains(rec.Body.String(), "请先登录") {
		t.Errorf("响应体缺少登录提示，got %q", rec.Body.String())
	}
}

// TestRequireAuthRejectsAnonymous requireAuth 失败时必须返回 (nil, error)，
// 否则调用方会 nil deref me。
func TestRequireAuthRejectsAnonymous(t *testing.T) {
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/", nil), rec)

	me, err := requireAuth(ctx)
	if err == nil {
		t.Fatal("未登录时 requireAuth 必须返回错误")
	}
	if me != nil {
		t.Errorf("me = %+v, want nil", me)
	}
}
