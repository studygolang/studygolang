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
