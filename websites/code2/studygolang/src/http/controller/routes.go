package controller

import "github.com/labstack/echo"

func RegisterRoutes(router *echo.Echo) {
	topicController := &TopicController{}
	router.Get("/topics/:view", topicController.Topics)
}
