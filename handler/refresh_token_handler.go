package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/katerji/bank/service"
	"strings"
)

const RefreshTokenPath = "/refresh"

type RefreshTokenRequest struct{}

type RefreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func RefreshTokenHandler(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	refreshToken := strings.ReplaceAll(authorization, "Bearer ", "")
	if refreshToken == "" {
		sendErrorMessage(c, "Missing refresh token.")
		return
	}
	jwtService := service.JWTService{}
	user, err := jwtService.VerifyRefreshToken(refreshToken)
	if err != nil {
		sendErrorMessage(c, "Invalid refresh token.")
		return
	}
	token, err := jwtService.CreateJwt(user)
	if err != nil {
		sendErrorMessage(c, "Error creating token.")
	}
	refreshToken, err = jwtService.CreateRefreshJwt(user)
	if err != nil {
		sendErrorMessage(c, "Error creating token.")
	}
	response := RefreshTokenResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}
	sendJSONResponse(c, response)
}
