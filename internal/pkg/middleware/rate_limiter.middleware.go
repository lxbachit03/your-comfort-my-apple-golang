package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lxbachit03/ygz-microservices-golang/internal/pkg/logger"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*client)
)

func getRateLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	c, exists := clients[ip]
	if !exists {
		// Limit to 5 requests per second with a burst size of 10
		limiter := rate.NewLimiter(rate.Limit(5), 10)
		c = &client{limiter: limiter}
		clients[ip] = c
	}

	c.lastSeen = time.Now()
	return c.limiter
}

func RateLimiterMiddleware(rateLimterLogger *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientIp := ctx.ClientIP()
		if clientIp == "" {
			// If client use proxy to fake ip
			clientIp = ctx.Request.RemoteAddr
		}

		limiter := getRateLimiter(clientIp)

		if !limiter.Allow() {
			if shoudLogRateLimit(clientIp) {
				traceId := logger.GetTraceID(ctx.Request.Context())

				rateLimterLogger.Warn().
					Str("trace_id", traceId).
					Str("method", ctx.Request.Method).
					Str("path", ctx.Request.URL.Path).
					Str("query", ctx.Request.URL.RawQuery).
					Str("client_ip", clientIp).
					Str("user_agent", ctx.Request.UserAgent()).
					Str("referer", ctx.Request.Referer()).
					Str("protocol", ctx.Request.Proto).
					Str("host", ctx.Request.Host).
					Str("remote_addr", ctx.Request.RemoteAddr).
					Str("request_uri", ctx.Request.RequestURI).
					Interface("headers", ctx.Request.Header).
					Msg("rate limiter exceeded")
			}

			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many request",
				"message": "Please try again later after few minutes",
			})

			return
		}

		ctx.Next()
	}
}

var rateLimitLogCache = sync.Map{}

const rateLimitLogTTL = 20 * time.Second

// To avoid logging too many requests, we use a cache to store the last time a request was logged
func shoudLogRateLimit(ip string) bool {
	now := time.Now()

	if val, ok := rateLimitLogCache.Load(ip); ok {
		if t, ok := val.(time.Time); ok && now.Sub(t) < rateLimitLogTTL {
			return false
		}
	}

	rateLimitLogCache.Store(ip, now)

	return true
}
