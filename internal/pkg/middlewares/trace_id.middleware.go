package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/logger"
)

const HEADER_TRACE_ID string = "X-Trace-Id"

func TraceIdMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceId := ctx.GetHeader(HEADER_TRACE_ID)

		if traceId == "" {
			traceId = uuid.New().String()
		}

		contextValue := context.WithValue(ctx.Request.Context(), logger.TraceIdKey, traceId)
		ctx.Request = ctx.Request.WithContext(contextValue)

		ctx.Writer.Header().Set(HEADER_TRACE_ID, traceId)

		ctx.Set(string(logger.TraceIdKey), traceId)

		ctx.Next()
	}
}
