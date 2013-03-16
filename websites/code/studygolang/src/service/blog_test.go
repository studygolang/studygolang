package service_test

import (
	. "service"
	"testing"
)

func TestFindNewBlogs(t *testing.T) {
	articleList := FindNewBlogs()
	if len(articleList) == 0 {
		t.Fatal("xxxx")
	}
	t.Log(len(articleList))
	for k, article := range articleList {
		t.Log(k, article)
		t.Log("===")
	}
}
