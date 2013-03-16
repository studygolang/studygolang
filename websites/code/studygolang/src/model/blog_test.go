package model_test

import (
	. "model"
	"testing"
)

func TestNewBlog(t *testing.T) {
	articleList, err := NewArticle().FindAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(articleList) == 0 {
		t.Fatal("xxxx")
	}
}
