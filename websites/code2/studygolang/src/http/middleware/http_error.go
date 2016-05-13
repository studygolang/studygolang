package middleware

import (
	"net/http"

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
							return Render(ctx, "404.html", nil)
						}
					}
				}
			}
			return nil
		})
	}
}
