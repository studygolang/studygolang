package logic_test

import (
	"logic"
	"testing"
)

func TestFindAll(t *testing.T) {
	topicsMap := new(logic.TopicLogic).FindAll(nil)
	t.Fatal(topicsMap)
}
