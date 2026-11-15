package logic

import (
	"os"
	"strings"
	"sync"
	"testing"
	"text/template"
	"time"

	"github.com/studygolang/studygolang/internal/model"
)

func TestSitemapOutputIncludesWiki(t *testing.T) {
	oldSitemapPath := sitemapPath
	oldSitemapTpl := sitemapTpl
	sitemapPath = t.TempDir() + "/"
	sitemapTpl = template.Must(template.New("sitemap.xml").Funcs(funcMap).ParseFiles("../../template/sitemap.xml"))
	sitemapTplOnce = sync.Once{}
	sitemapTplOnce.Do(func() {})
	defer func() {
		sitemapPath = oldSitemapPath
		sitemapTpl = oldSitemapTpl
		sitemapTplOnce = sync.Once{}
		if oldSitemapTpl != nil {
			sitemapTplOnce.Do(func() {})
		}
	}()

	wiki := &model.Wiki{Uri: "go-concurrency", Mtime: time.Date(2026, 11, 14, 0, 0, 0, 0, time.UTC)}
	data := map[string]interface{}{
		"home": map[string]string{"loc": "https://studygolang.com"},
		"wiki": []*model.Wiki{wiki},
	}
	if err := output("wiki.xml", data); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(sitemapPath + "wiki.xml")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "https://studygolang.com/wiki/go-concurrency") {
		t.Fatalf("wiki URL missing from sitemap: %s", content)
	}
}
