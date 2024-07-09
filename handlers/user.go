package handlers

import (
	"encoding/json"
	"fmt"
	"mazi-bet/models"
	"mazi-bet/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AccountResponse struct {
	ID       uint    `json:"ID"`
	Username string  `json:"Username"`
	Token    string  `json:"Token"`
	Balance  float64 `json:"Balance"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserCreateRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type DepositWithdrawRequest struct {
	Amount float64 `json:"amount"`
}

// @Summary Register a new user
// @Description Register a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param body body UserCreateRequest true "User registration details"
// @Success 200 {object} AccountResponse
// @Failure 400 {object} ErrorResponseRegisterLogin
// @Failure 422 {object} ErrorResponseRegisterLogin
// @Failure 500 {object} ErrorResponseRegisterLogin
// @Router /users/register [post]
func RegisterHandler(c echo.Context, dbConn *gorm.DB) error {
	// Read Request Body
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.ErrorResponseRegisterLogin{ResponseCode: 422, Message: "Invalid JSON"})
	}

	// check json format
	jsonFormatValidationMsg, jsonFormatErr := utils.ValidateJsonFormat(jsonBody, "firstname", "lastname", "email", "phone", "username", "password")
	if jsonFormatErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.ErrorResponseRegisterLogin{ResponseCode: 422, Message: jsonFormatValidationMsg})
	}

	// check user validation
	userFormatValidationMsg, user, userFormatErr := utils.ValidateUser(jsonBody)
	if userFormatErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.ErrorResponseRegisterLogin{ResponseCode: 422, Message: userFormatValidationMsg})
	}

	// check unique
	userUniqueMsg, userUniqueErr := utils.CheckUnique(user, jsonBody["username"].(string), dbConn)
	if userUniqueErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.ErrorResponseRegisterLogin{ResponseCode: 422, Message: userUniqueMsg})
	}

	// Insert User Object Into Database
	createdUser := dbConn.Create(&user)
	if createdUser.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponseRegisterLogin{ResponseCode: 500, Message: "User Creation Failed"})
	}

	// create account
	accountCreationMsg, account, accountCreationErr := utils.CreateAccount(int(user.ID), jsonBody["username"].(string), false, jsonBody["password"].(string), dbConn)
	if accountCreationErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.ErrorResponseRegisterLogin{ResponseCode: 422, Message: accountCreationMsg})
	}

	return c.JSON(http.StatusCreated, account)
}

// LoginHandler handles user login
// @Summary User login
// @Description Login with username and password
// @Tags users
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login request body"
// @Success 200 {object} AccountResponse
// @Failure 400 {object} ErrorResponseRegisterLogin
// @Failure 422 {object} ErrorResponseRegisterLogin
// @Router  /users/login [post]
func LoginHandler(c echo.Context, db *gorm.DB) error {
	// Read Request Body
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.ErrorResponseRegisterLogin{ResponseCode: 422, Message: "Invalid JSON"})
	}

	// find account based on username and check password correction
	findAccountMsg, account, findAccountErr := utils.Login(jsonBody["username"].(string), jsonBody["password"].(string), db)
	if findAccountErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.ErrorResponseRegisterLogin{ResponseCode: 422, Message: findAccountMsg})
	}

	return c.JSON(http.StatusOK, account)
}

// DepositHandler handles deposit to user balance
// @Summary Deposit Handler
// @Description Deposit an amount to user balance
// @Tags balance
// @Accept json
// @Produce json
// @Param body body DepositWithdrawRequest true "Deposit request body"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} AccountResponse
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router  /users/deposit [post]
func DepositHandler(c echo.Context, db *gorm.DB) error {
	account := c.Get("account").(models.Account)
	body := DepositWithdrawRequest{}

	if err := c.Bind(&body); err != nil {
		errResponse := models.ErrorResponse{
			Message: "Invalid Request Payload",
		}
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	user, err := utils.Deposit(db, account.UserID, body.Amount, "deposit")
	if err != nil {
		errResponse := models.ErrorResponse{
			Message: err.Error(),
		}
		return c.JSON(http.StatusConflict, errResponse)
	}
	fmt.Printf("User balance after deposit: %.2f", user.Balance)

	return c.JSON(http.StatusOK, user)
}

// WithdrawHandler handles withdraw from user balance
// @Summary Withdraw Handler
// @Description withdraw from user balance
// @Tags balance
// @Accept json
// @Produce json
// @Param body body DepositWithdrawRequest true "withdraw request body"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} AccountResponse
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router  /users/withdraw [post]
func WithdrawHandler(c echo.Context, db *gorm.DB) error {
	account := c.Get("account").(models.Account)
	body := DepositWithdrawRequest{}

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
	user, err := utils.Withdraw(db, account.UserID, body.Amount, "withdraw")
	if err != nil {
		errResponse := models.ErrorResponse{
			Message: err.Error(),
		}
		return c.JSON(http.StatusConflict, errResponse)
	}
	fmt.Printf("User balance after withdraw: %.2f", user.Balance)

	return c.JSON(http.StatusOK, user)
}
