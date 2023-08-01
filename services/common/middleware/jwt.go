package middleware

import (
	"net/http"
	appctx "smartkid/services/common/context"
	"smartkid/services/common/infra/logger"
	jwt "smartkid/services/common/util"

	"github.com/gin-gonic/gin"
)

func JWTAuth(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			log.Debug("Token is required")
			Response(c, http.StatusUnauthorized, false, "Invalid request", gin.H{
				"reload": true,
			}, nil)
			c.Abort()
			return
		}
		j := jwt.NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == jwt.TokenExpired {
				log.Debug("Token is expired")
				Response(c, http.StatusUnauthorized, false, "Invalid request", gin.H{
					"reload": true,
				}, nil)
				c.Abort()
				return
			}
			c.Abort()
			return
		}

		ctx := appctx.FromContext(c)
		ctx = ctx.WithEntity(claims.UserId, nil)
		c.Next()
	}
}

func Response(c *gin.Context, httpCode int, success bool, message string, data interface{}, err error) {
	c.JSON(httpCode, gin.H{
		"success": success,
		"message": message,
		"data":    data,
		"error":   err,
	})
}
