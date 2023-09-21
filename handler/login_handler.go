package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/katerji/bank/input"
	"github.com/katerji/bank/model"
	"github.com/katerji/bank/service"
)

const LoginPath = "/login"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User         model.UserOutput `json:"user"`
	Token        string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
}

func LoginHandler(c *gin.Context) {
	request := &LoginRequest{}
	err := c.BindJSON(request)
	if err != nil {
		sendBadRequest(c)
		return
	}
	authService := service.AuthService{}
	authInput := input.AuthInput{
		Email:    request.Email,
		Password: request.Password,
	}
	user, err := authService.Login(authInput)
	if err != nil {
		sendErrorMessage(c, err.Error())
		return
	}
	jwtService := service.JWTService{}
	token, err := jwtService.CreateJwt(user)
	if err != nil {
		sendErrorMessage(c, "")
		return
	}
	refreshToken, err := jwtService.CreateRefreshJwt(user)
	if err != nil {
		sendErrorMessage(c, "")
		return
	}
	response := LoginResponse{
		User:         user.ToOutput(),
		Token:        token,
		RefreshToken: refreshToken,
	}
	sendJSONResponse(c, response)
	return
}
