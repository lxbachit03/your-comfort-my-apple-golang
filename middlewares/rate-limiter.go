package middlewares

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type clientRateLimit struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clientIps = make(map[string]*clientRateLimit)
	mu        sync.Mutex
)

func getRateLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	clientIpRateLimit, exist := clientIps[ip]

	if !exist {
		limiter := rate.NewLimiter(5, 15) // 5 req every 1 second, 15 burst
		newClientRatelimit := &clientRateLimit{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		clientIps[ip] = newClientRatelimit
		return limiter
	}

	clientIpRateLimit.lastSeen = time.Now()

	return clientIpRateLimit.limiter
}

func getClientIp(ctx *gin.Context) string {
	clientIp := ctx.ClientIP()

	if clientIp == "" {
		clientIp = ctx.RemoteIP()
	}

	return clientIp
}

func CleanUpRateLimiter() {
	for {
		time.Sleep(time.Minute)

		mu.Lock()

		log.Printf("clientIps: %+v", clientIps)
		log.Print("=-=-=-=-=-=-=-=-=")
		for ip, clientIpRateLimit := range clientIps {
			log.Print("lastSeen: ", clientIpRateLimit.lastSeen)
			log.Print("since: ", time.Since(clientIpRateLimit.lastSeen))
			if time.Since(clientIpRateLimit.lastSeen) > 1*time.Minute {
				delete(clientIps, ip)
			}
		}

		mu.Unlock()
	}
}

// ab -n 20 -c 1 -H "X-API-Key:b3f6590f-20a9-4843-84df-9156d5f37306" http://localhost:8080/api/v1/users/
func RateLimiterMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientIp := getClientIp(ctx)

		limiter := getRateLimiter(clientIp)

		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": "Too many requests, please try again after few seconds",
			})

			return
		}

		ctx.Next()
	}
}
