package middleware

import (
	"CloudDisk/internal/rds"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 登陆校验
func VailToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie("token")
		username, err := rds.QueryToken(token)

		if err != nil || username == "" {

			c.JSONP(http.StatusBadGateway, gin.H{
				"msg": token,
				// "err":   err.Error(),
				// "token": token,
			})
			c.Abort()

		} else {
			c.Next()
		}
	}
}
