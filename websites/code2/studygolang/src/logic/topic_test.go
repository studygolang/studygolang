package logic_test

import (
	"logic"
	"testing"
)

func TestFindAll(t *testing.T) {
	paginator := logic.NewPaginator(2)
	topicsMap := logic.DefaultTopic.FindAll(nil, paginator, "", "")
	t.Fatal(topicsMap)
}
