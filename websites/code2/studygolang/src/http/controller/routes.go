package controller

import "github.com/labstack/echo"

func RegisterRoutes(e *echo.Echo) {
	// e = e.Group("", middleware.NeedLogin)

	new(TopicController).RegisterRoute(e)
	new(AccountController).RegisterRoute(e)
}
