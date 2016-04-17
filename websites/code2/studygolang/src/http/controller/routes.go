package controller

import "github.com/labstack/echo"

func RegisterRoutes(e *echo.Echo) {
	new(IndexController).RegisterRoute(e)
	new(AccountController).RegisterRoute(e)
	new(TopicController).RegisterRoute(e)
	new(ArticleController).RegisterRoute(e)
	new(ProjectController).RegisterRoute(e)
	new(ResourceController).RegisterRoute(e)
	new(ReadingController).RegisterRoute(e)
	new(SidebarController).RegisterRoute(e)
	new(CommentController).RegisterRoute(e)
}
