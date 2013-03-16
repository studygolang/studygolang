package model_test

import (
	. "model"
	"testing"
)

func TestNewTopic(t *testing.T) {
	// err := NewTopic().Set("lastreplyuid=1").Where("uid=").Update()
	err := NewTopicEx().Where("tid=1").Increment("reply", 1)
	if err != nil {
		t.Fatal(err)
	}
}
