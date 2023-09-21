package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/katerji/bank/model"
)

const UserInfoPath = "/user"

type UserInfoResponse struct {
	User model.UserOutput `json:"user"`
}

func UserInfoHandler(c *gin.Context) {
	user := c.MustGet("user").(model.User)
	response := UserInfoResponse{
		User: user.ToOutput(),
	}
	sendJSONResponse(c, response)
	return
}
