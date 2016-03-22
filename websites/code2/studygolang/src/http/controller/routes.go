package controller

import "github.com/labstack/echo"

func RegisterRoutes(e *echo.Echo) {
	new(TopicController).RegisterRoute(e)
	new(AccountController).RegisterRoute(e)
	new(SidebarController).RegisterRoute(e)
	new(CommentController).RegisterRoute(e)
}
