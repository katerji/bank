package handler

import "github.com/gin-gonic/gin"

const LandingPath = "/"

func LandingController(c *gin.Context) {
	c.String(200, "Welcome to the bank!")
	return
}
