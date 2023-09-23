package cache

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/polaris1119/nosql"
	"github.com/studygolang/studygolang/internal/model"
)

type feedCache struct{}

var Feed feedCache

func (feedCache) GetTop(ctx context.Context) []*model.Feed {
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()

	s := redisClient.GET("feed:top")
	if s == "" {
		return nil
	}

	feeds := make([]*model.Feed, 0)
	err := json.Unmarshal([]byte(s), &feeds)
	if err != nil {
		return nil
	}

	return feeds
}

func (feedCache) SetTop(ctx context.Context, feeds []*model.Feed) {
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()

	b, _ := json.Marshal(feeds)
	redisClient.SET("feed:top", string(b), 300)
}

func (feedCache) GetList(ctx context.Context, p int) []*model.Feed {
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()

	s := redisClient.GET("feed:list:" + strconv.Itoa(p))
	if s == "" {
		return nil
	}

	feeds := make([]*model.Feed, 0)
	err := json.Unmarshal([]byte(s), &feeds)
	if err != nil {
		return nil
	}

	return feeds
}

func (feedCache) SetList(ctx context.Context, p int, feeds []*model.Feed) {
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()

	b, _ := json.Marshal(feeds)
	redisClient.SET("feed:list:"+strconv.Itoa(p), string(b), 300)
}
