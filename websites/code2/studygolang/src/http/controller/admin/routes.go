package admin

import "github.com/labstack/echo"

func RegisterRoutes(g *echo.Group) {
	new(AuthorityController).RegisterRoute(g)
	new(ArticleController).RegisterRoute(g)
	new(RuleController).RegisterRoute(g)
	new(ReadingController).RegisterRoute(g)
	new(ToolController).RegisterRoute(g)
}
