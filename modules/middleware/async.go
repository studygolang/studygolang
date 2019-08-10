package middleware

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/polaris1119/goutils"
)

// EchoAsync 用于 echo 框架的异步处理中间件
func EchoAsync() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			req := ctx.Request()

			if req.Method != "GET" {
				// 是否异步执行
				async := goutils.MustBool(ctx.FormValue("async"), false)
				if async {
					go next(ctx)

					result := map[string]interface{}{
						"code": 0,
						"msg":  "ok",
						"data": nil,
					}
					return ctx.JSON(http.StatusOK, result)
				}
			}

			if err := next(ctx); err != nil {
				return err
			}

			return nil
		}
	}
}
