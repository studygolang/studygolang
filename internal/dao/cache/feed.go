package cache

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/garyburd/redigo/redis"
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

	if s == "notop" {
		return []*model.Feed{}
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

	val := "notop"
	if len(feeds) > 0 {
		b, _ := json.Marshal(feeds)
		val = string(b)
	}

	redisClient.SET("feed:top", val, 300)
}

func feedListKey(p int, tab string) string {
	return "feed:list:" + tab + ":" + strconv.Itoa(p)
}

func (feedCache) GetList(ctx context.Context, p int, tab string) []*model.Feed {
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()

	s := redisClient.GET(feedListKey(p, tab))
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

func (feedCache) SetList(ctx context.Context, p int, tab string, feeds []*model.Feed) {
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()

	b, _ := json.Marshal(feeds)
	redisClient.SET(feedListKey(p, tab), string(b), 300)
}

// InvalidateAll 清空所有 tab、所有页码的 feed 列表与置顶缓存。
func (feedCache) InvalidateAll(ctx context.Context) {
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()

	redisClient.DEL("feed:top")
	if redisClient.Conn == nil {
		return
	}

	cursor := 0
	for {
		result, err := redis.Values(redisClient.Do("SCAN", cursor, "MATCH", nosql.KeyPrefix+"feed:list:*", "COUNT", 100))
		if err != nil {
			return
		}

		var keys []string
		if _, err = redis.Scan(result, &cursor, &keys); err != nil {
			return
		}
		if len(keys) > 0 {
			if _, err = redisClient.Do("DEL", redis.Args{}.AddFlat(keys)...); err != nil {
				return
			}
		}
		if cursor == 0 {
			return
		}
	}
}
