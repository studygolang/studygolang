package logic

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/polaris1119/config"
)

func TestFindAtomFeedsUsesCurrentEngineURL(t *testing.T) {
	oldEngineURL := config.ConfigFile.MustValue("search", "engine_url")
	t.Cleanup(func() { config.ConfigFile.SetValue("search", "engine_url", oldEngineURL) })

	newSolr := func(title string) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/select" {
				t.Errorf("Solr 请求路径 = %q, want /select", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprintf(w, `{"response":{"numFound":1,"start":0,"docs":[{"title":%q}]}}`, title)
		}))
	}

	first := newSolr("first")
	defer first.Close()
	second := newSolr("second")
	defer second.Close()

	config.ConfigFile.SetValue("search", "engine_url", first.URL)
	respBody, err := DefaultSearcher.FindAtomFeeds(1)
	if err != nil {
		t.Fatalf("首次查询失败：%v", err)
	}
	if len(respBody.Docs) != 1 || respBody.Docs[0].Title != "first" {
		t.Fatalf("首次查询返回 %+v, want first", respBody.Docs)
	}

	config.ConfigFile.SetValue("search", "engine_url", second.URL)
	respBody, err = DefaultSearcher.FindAtomFeeds(1)
	if err != nil {
		t.Fatalf("运行时修改配置后查询失败：%v", err)
	}
	if len(respBody.Docs) != 1 || respBody.Docs[0].Title != "second" {
		t.Fatalf("修改配置后返回 %+v, want second", respBody.Docs)
	}
}

func TestDoSearchRejectsInvalidSolrResponses(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
	}{
		{name: "http error", status: http.StatusServiceUnavailable, body: `<html>unavailable</html>`},
		{name: "missing response with highlighting", status: http.StatusOK, body: `{"highlighting":{"01":{"title":["title"]}}}`},
		{name: "null response", status: http.StatusOK, body: `{"response":null}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			solr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer solr.Close()

			oldEngineURL := config.ConfigFile.MustValue("search", "engine_url")
			config.ConfigFile.SetValue("search", "engine_url", solr.URL)
			t.Cleanup(func() { config.ConfigFile.SetValue("search", "engine_url", oldEngineURL) })

			if _, err := DefaultSearcher.DoSearch("", "", 0, 10); err == nil {
				t.Fatal("异常 Solr 响应必须返回 error")
			}
		})
	}
}

func TestDoSearchSkipsNullDocsAndHighlighting(t *testing.T) {
	solr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response":{"numFound":2,"start":0,"docs":[null,{"id":"01","objid":1,"objtype":0,"title":"topic","content":"content"}]},
			"highlighting":{"01":null}
		}`))
	}))
	defer solr.Close()

	oldEngineURL := config.ConfigFile.MustValue("search", "engine_url")
	config.ConfigFile.SetValue("search", "engine_url", solr.URL)
	t.Cleanup(func() { config.ConfigFile.SetValue("search", "engine_url", oldEngineURL) })

	respBody, err := DefaultSearcher.DoSearch("", "", 0, 10)
	if err != nil {
		t.Fatalf("包含空文档的响应应被安全处理：%v", err)
	}
	if len(respBody.Docs) != 1 {
		t.Fatalf("文档数量 = %d, want 1", len(respBody.Docs))
	}
	if respBody.Docs[0].Title != "topic" || respBody.Docs[0].HlTitle != "topic" {
		t.Errorf("文档或高亮标题未正确保留：%+v", respBody.Docs[0])
	}
}
