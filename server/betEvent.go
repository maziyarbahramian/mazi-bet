package server

import (
	"mazi-bet/handlers"
	"mazi-bet/middlewares"

	"github.com/labstack/echo/v4"
)

func betEventRoutes(e *echo.Echo) {
	e.POST("/bet", WithDBConnection(handlers.BetHandler), middlewares.IsLoggedIn)
}
