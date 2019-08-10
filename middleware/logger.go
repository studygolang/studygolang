package middleware

import (
	"context"
	"fmt"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/logger"
	"github.com/twinj/uuid"
)

const HeaderKey = "X-Request-Id"

type LoggerConfig struct {
	// 是否输出 POST 参数，默认不输出
	OutputPost bool
	// 当 OutputPost 为 true 时，排除这些 path，避免包含敏感信息输出
	Excludes map[string]struct{}
}

var DefaultLoggerConfig = &LoggerConfig{}

func EchoLogger() echo.MiddlewareFunc {
	return EchoLoggerWitchConfig(DefaultLoggerConfig)
}

// EchoLoggerWitchConfig 用于 echo 框架的日志中间件
func EchoLoggerWitchConfig(loggerConfig *LoggerConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			start := time.Now()

			req := ctx.Request()
			resp := ctx.Response()

			objLogger := logger.GetLogger()
			ctx.Set("logger", objLogger)

			var params map[string][]string
			if loggerConfig.OutputPost {
				params, _ = ctx.FormParams()
				if len(loggerConfig.Excludes) > 0 {
					_, ok := loggerConfig.Excludes[req.URL.Path]
					if ok {
						params = ctx.QueryParams()
					}
				}
			} else {
				params = ctx.QueryParams()
			}
			objLogger.Infoln("request params:", params)

			remoteAddr := ctx.RealIP()

			id := func(ctx echo.Context) string {
				id := req.Header.Get(HeaderKey)
				if id == "" {
					id = ctx.FormValue("request_id")
					if id == "" {
						id = uuid.NewV4().String()
					}
				}

				ctx.Set("request_id", id)

				return id
			}(ctx)

			resp.Header().Set(HeaderKey, id)

			defer func() {
				method := req.Method
				path := req.URL.Path
				if path == "" {
					path = "/"
				}
				size := resp.Size
				code := resp.Status

				stop := time.Now()
				// [remoteAddr method path request_id "UA" code time size]
				uri := fmt.Sprintf(`[%s %s %s %s "%s" %d %s %d]`, remoteAddr, method, path, id, req.UserAgent(), code, stop.Sub(start), size)
				objLogger.SetContext(context.WithValue(context.Background(), "uri", uri))
				objLogger.Flush()
				logger.PutLogger(objLogger)
			}()

			if err := next(ctx); err != nil {
				return err
			}
			return nil
		}
	}
}
