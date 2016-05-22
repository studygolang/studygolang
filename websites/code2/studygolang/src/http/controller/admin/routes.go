package admin

import "github.com/labstack/echo"

func RegisterRoutes(g *echo.Group) {
	new(AuthorityController).RegisterRoute(g)
	new(ArticleController).RegisterRoute(g)
}
