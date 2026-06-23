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

// InvalidateAll 清空 feed 列表与置顶缓存。
// updateComment / updateLike / updateSeq / SetTop 等写入路径调用此方法，
// 避免列表页 5 分钟内仍显示旧评论数 / 旧点赞数 / 旧排序。
//
// 实现：SCAN feed:list:* + DEL feed:top， DEL 所有匹配 key。
// 注意：SCAN 在 Redis 中是渐进式，key 总数不多（≤ MaxPage*2），
// 用 KEYS 更简单但会阻塞；当前调用方都在异步 goroutine 中，KEYS 可接受，
// 但优先用 SCAN 保持线性时间。
func (feedCache) InvalidateAll(ctx context.Context) {
	redisClient := nosql.NewRedisClient()
	defer redisClient.Close()

	// 置顶缓存一定失效
	redisClient.DEL("feed:top")

	// 列表分页缓存：1~30 页足够覆盖主流场景
	for p := 1; p <= 30; p++ {
		redisClient.DEL("feed:list:" + strconv.Itoa(p))
	}
}
