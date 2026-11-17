package api

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/studygolang/studygolang/internal/logic"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/config"
)

func TestRSSContract(t *testing.T) {
	solr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {
				"numFound": 2,
				"start": 0,
				"docs": [
					{"id":"011","objid":11,"objtype":0,"title":"话题标题","author":"polaris","content":"话题内容","created_at":"2024-01-02 03:04:05"},
					{"id":"612","objid":12,"objtype":6,"title":"面试题标题","author":"interviewer","content":"面试题内容","created_at":"2024-01-02 03:04:05"}
				]
			}
		}`))
	}))
	defer solr.Close()

	oldEngineURL := config.ConfigFile.MustValue("search", "engine_url")
	config.ConfigFile.SetValue("search", "engine_url", solr.URL)
	t.Cleanup(func() { config.ConfigFile.SetValue("search", "engine_url", oldEngineURL) })

	oldName := logic.WebsiteSetting.Name
	oldDomain := logic.WebsiteSetting.Domain
	oldSlogan := logic.WebsiteSetting.Slogan
	oldOnlyHTTPS := logic.WebsiteSetting.OnlyHttps
	logic.WebsiteSetting.Name = "Go语言中文网"
	logic.WebsiteSetting.Domain = "example.test"
	logic.WebsiteSetting.Slogan = "Golang 中文社区"
	logic.WebsiteSetting.OnlyHttps = false
	t.Cleanup(func() {
		logic.WebsiteSetting.Name = oldName
		logic.WebsiteSetting.Domain = oldDomain
		logic.WebsiteSetting.Slogan = oldSlogan
		logic.WebsiteSetting.OnlyHttps = oldOnlyHTTPS
	})

	recorder := httptest.NewRecorder()
	ctx := echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/api/v1/feed", nil), recorder)
	if err := (FeedController{}).RSS(ctx); err != nil {
		t.Fatalf("RSS() error = %v", err)
	}

	if !strings.HasPrefix(recorder.Body.String(), xml.Header) {
		t.Fatalf("RSS 缺少 XML 声明：%q", recorder.Body.String())
	}
	if got := recorder.Header().Get("Content-Type"); got != echo.MIMEApplicationXMLCharsetUTF8 {
		t.Errorf("Content-Type = %q, want %q", got, echo.MIMEApplicationXMLCharsetUTF8)
	}

	var feed rssFeed
	if err := xml.Unmarshal(recorder.Body.Bytes(), &feed); err != nil {
		t.Fatalf("RSS XML 无法解析：%v", err)
	}
	if feed.Version != "2.0" {
		t.Errorf("RSS version = %q, want 2.0", feed.Version)
	}
	if feed.Channel.Link != "http://example.test/" {
		t.Errorf("OnlyHttps=false 时 channel link = %q", feed.Channel.Link)
	}
	if len(feed.Channel.Items) != 2 {
		t.Fatalf("RSS item 数量 = %d, want 2", len(feed.Channel.Items))
	}

	topic := feed.Channel.Items[0]
	if topic.Link != "http://example.test/topics/11" {
		t.Errorf("话题链接 = %q", topic.Link)
	}
	if topic.GUID != topic.Link {
		t.Errorf("GUID = %q, want %q", topic.GUID, topic.Link)
	}
	if topic.Author != "polaris" {
		t.Errorf("作者 = %q, want polaris", topic.Author)
	}
	if _, err := time.Parse(time.RFC1123Z, topic.PubDate); err != nil {
		t.Errorf("pubDate %q 不符合 RFC822 数字时区格式：%v", topic.PubDate, err)
	}
	if !strings.HasSuffix(topic.PubDate, "+0000") {
		t.Errorf("pubDate 应统一为 UTC，got %q", topic.PubDate)
	}

	interview := feed.Channel.Items[1]
	if interview.Link != "http://example.test/interview/12" {
		t.Errorf("面试题链接 = %q, want http://example.test/interview/12", interview.Link)
	}
}
