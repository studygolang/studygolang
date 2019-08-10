package echoutils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	mycontext "github.com/studygolang/studygolang/context"

	"github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
	"github.com/polaris1119/nosql"
)

const logKey = "logger"

// GetLogger 由调用者确保 ctx 中存在 logger.Logger 对象
func GetLogger(ctx context.Context) *logger.Logger {
	return ctx.Value(logKey).(*logger.Logger)
}

// 是否异步处理
func IsAsync(ctx echo.Context) bool {
	return goutils.MustBool(ctx.FormValue("async"), false)
}

// WrapContext 返回一个 context.Context 实例
func WrapEchoContext(ctx echo.Context) context.Context {
	return mycontext.EchoContext(ctx)
}

// WrapContext 返回一个 context.Context 实例。如果 ctx == nil，需要确保 调用 logger.PutLogger()
func WrapContext(ctx context.Context) context.Context {
	var objLogger *logger.Logger
	if ctx == nil {
		ctx = context.Background()
		objLogger = logger.GetLogger()
	} else {
		objLogger = GetLogger(ctx)
	}
	return context.WithValue(ctx, logKey, objLogger)
}

func LogFlush(ctx context.Context) {
	objLogger := GetLogger(ctx)
	objLogger.Flush()
	logger.PutLogger(objLogger)
}

func Success(ctx echo.Context, data interface{}) error {
	result := map[string]interface{}{
		"code": 0,
		"msg":  "ok",
		"data": data,
	}

	b, err := json.Marshal(result)
	if err != nil {
		return err
	}

	go func(b []byte) {
		if cacheKey := ctx.Get(nosql.CacheKey); cacheKey != nil {
			logger.Debugln("cache save:", cacheKey, "now:", time.Now())
			nosql.DefaultLRUCache.CompressAndAdd(cacheKey, b, nosql.NewCacheData())
		}
	}(b)

	if ctx.Response().Committed {
		LogFlush(WrapEchoContext(ctx))
		return nil
	}

	return ctx.JSONBlob(http.StatusOK, b)
}

func Fail(ctx echo.Context, code int, msg string) error {
	if ctx.Response().Committed {
		LogFlush(WrapEchoContext(ctx))
		return nil
	}

	result := map[string]interface{}{
		"code": code,
		"msg":  msg,
	}

	GetLogger(WrapEchoContext(ctx)).Errorln("operate fail:", result)

	return ctx.JSON(http.StatusOK, result)
}

func AsyncResponse(ctx echo.Context, logicInstance interface{}, methodName string, args ...interface{}) error {
	wrapCtx := mycontext.EchoContext(ctx)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("async response panic:", err)
			}
		}()
		defer LogFlush(wrapCtx)

		instance := reflect.ValueOf(logicInstance)

		in := make([]reflect.Value, len(args)+1)
		in[0] = reflect.ValueOf(wrapCtx)
		for i, arg := range args {
			in[i+1] = reflect.ValueOf(arg)
		}

		instance.MethodByName(methodName).Call(in)
	}()

	return Success(ctx, nil)
}
