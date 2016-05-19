package middleware

import (
	"net/http"
	"util"

	. "http"

	"github.com/labstack/echo"
)

// EchoLogger 用于 echo 框架的日志中间件
func HTTPError() echo.MiddlewareFunc {
	return func(next echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(ctx echo.Context) error {
			if err := next.Handle(ctx); err != nil {

				if !ctx.Response().Committed() {
					if he, ok := err.(*echo.HTTPError); ok {
						switch he.Code {
						case http.StatusNotFound:
							if util.IsAjax(ctx) {
								return ctx.String(http.StatusOK, `{"ok":0,"error":"接口不存在"}`)
							}
							return Render(ctx, "404.html", nil)
						case http.StatusInternalServerError:
							if util.IsAjax(ctx) {
								return ctx.String(http.StatusOK, `{"ok":0,"error":"接口服务器错误"}`)
							}
							return Render(ctx, "500.html", nil)
						}
					}
				}
			}
			return nil
		})
	}
}
