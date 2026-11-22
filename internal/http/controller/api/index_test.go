package api

import (
	"testing"

	"github.com/studygolang/studygolang/internal/model"
)

func TestSanitizeHomeDataForPublic(t *testing.T) {
	privateUsers := []*model.User{
		{Email: "feed@example.com"},
		{Email: "topic@example.com"},
		{Email: "article@example.com"},
		{Email: "doc@example.com"},
		{Email: "subject@example.com"},
	}
	publicUser := &model.User{Email: "public@example.com", Open: 1}
	data := map[string]interface{}{
		"feeds":    []*model.Feed{{User: privateUsers[0]}},
		"topics":   []map[string]interface{}{{"user": privateUsers[1]}},
		"articles": []*model.Article{{LastReplyUser: privateUsers[2]}},
		"users":    map[int]*model.User{1: privateUsers[3], 2: publicUser},
		"subjects": []*model.Subject{{User: privateUsers[4]}},
	}

	sanitizeHomeDataForPublic(data)

	for i, user := range privateUsers {
		if user.Email != "" {
			t.Errorf("private user %d email = %q, want empty", i, user.Email)
		}
	}
	if publicUser.Email != "public@example.com" {
		t.Errorf("public user email = %q, want preserved", publicUser.Email)
	}
}
