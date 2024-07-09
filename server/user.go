package server

import (
	database "mazi-bet/database"
	"mazi-bet/handlers"
	"mazi-bet/middlewares"
	"mazi-bet/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func WithDBConnection(handlerFunc func(c echo.Context, db *gorm.DB) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		dbConn, err := database.GetConnection()
		if err != nil {
			return c.JSON(http.StatusBadGateway, models.Response{ResponseCode: 502, Message: "Can't Connect To Database"})
		}
		return handlerFunc(c, dbConn)
	}
}

func userRoutes(e *echo.Echo) {
	e.POST("/users/login", WithDBConnection(handlers.LoginHandler))
	e.POST("/users/register", WithDBConnection(handlers.RegisterHandler))
	e.POST("/users/deposit", WithDBConnection(handlers.DepositHandler), middlewares.IsLoggedIn)
	e.POST("/users/withdraw", WithDBConnection(handlers.WithdrawHandler), middlewares.IsLoggedIn)
}
