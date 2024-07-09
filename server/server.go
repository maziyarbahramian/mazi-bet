package server

import (
	"log"
	database "mazi-bet/database"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var e *echo.Echo

func init() {
	e = echo.New()
}

func StartServer() {
	_, err := database.GetConnection()
	if err != nil {
		log.Fatal(err)
	}

	e.GET("/swagger/*", echoSwagger.EchoWrapHandler())

	userRoutes(e)
	betEventRoutes(e)
	log.Fatal(e.Start(":8080"))
}
