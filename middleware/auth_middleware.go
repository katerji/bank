package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/katerji/bank/service"
	"strings"
)

func GetAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		token := strings.ReplaceAll(authorization, "Bearer ", "")
		jwtService := service.JWTService{}
		userObject, err := jwtService.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"message": "Unauthorized.",
			})
			return
		}
		c.Set("user", userObject)
		c.Next()
	}
}
