package controller

import "github.com/labstack/echo"

func RegisterRoutes(router *echo.Echo) {
	new(TopicController).RegisterRoute(router)
}
