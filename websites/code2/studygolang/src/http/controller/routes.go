package controller

import "github.com/labstack/echo"

func RegisterRoutes(e *echo.Group) {
	new(IndexController).RegisterRoute(e)
	new(AccountController).RegisterRoute(e)
	new(TopicController).RegisterRoute(e)
	new(ArticleController).RegisterRoute(e)
	new(ProjectController).RegisterRoute(e)
	new(ResourceController).RegisterRoute(e)
	new(ReadingController).RegisterRoute(e)
	new(WikiController).RegisterRoute(e)
	new(UserController).RegisterRoute(e)
	new(LikeController).RegisterRoute(e)
	new(FavoriteController).RegisterRoute(e)
	new(MessageController).RegisterRoute(e)
	new(SidebarController).RegisterRoute(e)
	new(CommentController).RegisterRoute(e)
	new(WideController).RegisterRoute(e)
	new(ImageController).RegisterRoute(e)
	new(CaptchaController).RegisterRoute(e)
	new(WebsocketController).RegisterRoute(e)
}
