package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/katerji/bank/input"
	"github.com/katerji/bank/service"
)

const RegisterPath = "/register"

const (
	errorMessageEmailAlreadyExists = "Email already exists."
	userRegisteredSuccessfully     = "User registered successfully."
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(c *gin.Context) {
	request := &RegisterRequest{}
	err := c.BindJSON(request)
	if err != nil {
		sendBadRequest(c)
		return
	}
	registerUserInput := input.AuthInput{
		Email:    request.Email,
		Password: request.Password,
	}
	userService := service.AuthService{}
	_, err = userService.Register(registerUserInput)
	if err != nil {
		sendErrorMessage(c, errorMessageEmailAlreadyExists)
		return
	}
	sendResponseMessage(c, userRegisteredSuccessfully)
	return
}
