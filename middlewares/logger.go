package middlewares

import (
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		logsPath = "logs/http.log"

		ctx.Next()
	}
}
