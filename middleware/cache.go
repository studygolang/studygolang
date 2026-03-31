package middleware

import (
	"net/http"
	"sort"
	"sync"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
	"github.com/polaris1119/nosql"
	"golang.org/x/sync/singleflight"
)

type CacheKeyAlgorithm interface {
	GenCacheKey(echo.Context) string
}

type CacheKeyFunc func(echo.Context) string

func (self CacheKeyFunc) GenCacheKey(ctx echo.Context) string {
	return self(ctx)
}

var CacheKeyAlgorithmMap = make(map[string]CacheKeyAlgorithm)

var LruCache = nosql.DefaultLRUCache

// singleflight 防止缓存雪崩
var sfGroup singleflight.Group
var sfMutex sync.RWMutex

// EchoCache 用于 echo 框架的缓存中间件。支持自定义 cache 数量
func EchoCache(cacheMaxEntryNum ...int) echo.MiddlewareFunc {

	if len(cacheMaxEntryNum) > 0 {
		LruCache = nosql.NewLRUCache(cacheMaxEntryNum[0])
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			req := ctx.Request()

			if req.Method == "GET" {
				cacheKey := getCacheKey(ctx)

				if cacheKey != "" {
					ctx.Set(nosql.CacheKey, cacheKey)

					value, compressor, ok := LruCache.GetAndUnCompress(cacheKey)
					if ok {
						cacheData, ok := compressor.(*nosql.CacheData)
						if ok {

							// 1分钟更新一次
							if time.Now().Sub(cacheData.StoreTime) >= time.Minute {
								// 使用 singleflight 防止缓存雪崩
								// 只允许一个请求重建缓存，其他请求使用旧数据
								_, _, shared := sfGroup.Do(cacheKey, func() (interface{}, error) {
									// 重建缓存完成后立即移除，允许下次重建
									defer sfGroup.Forget(cacheKey)

									// 重新执行 handler 生成新缓存
									logger.Debugln("rebuilding cache for:", cacheKey)
									err := next(ctx)
									return nil, err
								})

								if shared {
									// 如果是共享的结果（即其他请求已处理），使用旧缓存
									logger.Debugln("cache hit (stale, another request rebuilding):", cacheData.StoreTime, "now:", time.Now())
									return ctx.JSONBlob(http.StatusOK, value)
								}
								// 如果是当前请求执行的重建，next(ctx) 已经返回了响应
								return nil
							}

							logger.Debugln("cache hit:", cacheData.StoreTime, "now:", time.Now())
							return ctx.JSONBlob(http.StatusOK, value)
						}
					}
				}
			}

			if err := next(ctx); err != nil {
				return err
			}

			return nil
		}
	}
}

func getCacheKey(ctx echo.Context) string {
	cacheKey := ""
	if cacheKeyAlgorithm, ok := CacheKeyAlgorithmMap[ctx.Path()]; ok {
		// nil 表示不缓存
		if cacheKeyAlgorithm != nil {
			cacheKey = cacheKeyAlgorithm.GenCacheKey(ctx)
		}
	} else {
		cacheKey = defaultCacheKeyAlgorithm(ctx)
	}

	return cacheKey
}

func defaultCacheKeyAlgorithm(ctx echo.Context) string {
	filter := map[string]bool{
		"from":      true,
		"sign":      true,
		"nonce":     true,
		"timestamp": true,
	}
	form, err := ctx.FormParams()
	if err != nil {
		return ""
	}

	var keys = make([]string, 0, len(form))
	for key := range form {
		if _, ok := filter[key]; !ok {
			keys = append(keys, key)
		}
	}

	sort.Sort(sort.StringSlice(keys))

	buffer := goutils.NewBuffer()
	for _, k := range keys {
		buffer.Append(k).Append("=").Append(ctx.FormValue(k))
	}

	req := ctx.Request()
	return goutils.Md5(req.Method + req.URL.Path + buffer.String())
}
