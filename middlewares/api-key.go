package middlewares

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ApiKeyMiddleware() gin.HandlerFunc {

	apiKeyServer := os.Getenv("API_KEY")
	if apiKeyServer == "" {
		panic("API_KEY environment variable is not set")
	}

	return func(ctx *gin.Context) {
		apiKeyHeader := ctx.GetHeader("X-API-Key")

		if apiKeyHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":  http.StatusBadRequest,
				"error": "API Key is required",
				"path":  ctx.FullPath(),
			})

			return
		}

		if apiKeyHeader != apiKeyServer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  http.StatusUnauthorized,
				"error": "invalid API Key",
				"path":  ctx.FullPath(),
			})

			return
		}

		ctx.Next()
	}
}
