package handlers

import (
	"math/rand"
	"mazi-bet/models"
	"mazi-bet/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type BetRequest struct {
	Amount float64 `json:"amount"`
}

type BetResponse struct {
	Message string  `json:"message"`
	Amount  float64 `json:"amount"`
}

// CreateBetEventHandler handles bet event creation
// @Summary create BetHandler event
// @Description Create BetHandler Event
// @Tags bet-event
// @Accept json
// @Produce json
// @Success 200 {object} BetResponse
// @Param body body BetRequest true "BetHandler Event Create request body"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router  /bet [post]
func BetHandler(c echo.Context, dbConn *gorm.DB) error {
	account := c.Get("account").(models.Account)
	var user models.User
	dbConn.First(&models.User{}, account.UserID).Scan(&user)

	body := BetRequest{}
	if err := c.Bind(&body); err != nil {
		errResponse := models.ErrorResponse{
			Message: "Invalid Request Payload",
		}
		return c.JSON(http.StatusBadRequest, errResponse)
	}
	if body.Amount <= 0 {
		errResponse := models.ErrorResponse{
			Message: "Amount should be greater than 0",
		}
		return c.JSON(http.StatusBadRequest, errResponse)
	}
	if user.Balance < body.Amount {
		errResponse := models.ErrorResponse{
			Message: "Insufficient Balance",
		}
		return c.JSON(http.StatusPaymentRequired, errResponse)
	}

	min := -body.Amount
	max := body.Amount * 2
	randomAmount := min + rand.Float64()*(max-min)

	if randomAmount < 0 {
		_, err := utils.Withdraw(dbConn, account.UserID, -randomAmount, "bet-lose")
		if err != nil {
			errResponse := models.ErrorResponse{
				Message: err.Error(),
			}
			return c.JSON(http.StatusConflict, errResponse)
		}
		errResponse := BetResponse{
			Message: "you lose",
			Amount:  randomAmount,
		}
		return c.JSON(http.StatusOK, errResponse)
	} else {
		_, err := utils.Deposit(dbConn, account.UserID, randomAmount, "bet-won")
		if err != nil {
			errResponse := models.ErrorResponse{
				Message: err.Error(),
			}
			return c.JSON(http.StatusConflict, errResponse)
		}
		errResponse := BetResponse{
			Message: "you won",
			Amount:  randomAmount,
		}
		return c.JSON(http.StatusOK, errResponse)
	}
}
