package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	ips map[string]*rateLimiter
	mu  sync.Mutex
	rps int
}

type rateLimiter struct {
	count    int
	lastSeen time.Time
}

func NewRateLimiter(rps int) *RateLimiter {
	rl := &RateLimiter{
		ips: make(map[string]*rateLimiter),
		rps: rps,
	}

	// Очистка старых записей каждую минуту
	go func() {
		for {
			time.Sleep(time.Minute)
			rl.mu.Lock()
			for ip, limiter := range rl.ips {
				if time.Since(limiter.lastSeen) > time.Minute {
					delete(rl.ips, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()

	return rl
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		limiter, exists := rl.ips[ip]
		if !exists {
			limiter = &rateLimiter{}
			rl.ips[ip] = limiter
		}

		// Сброс счётчика, если прошла секунда
		if time.Since(limiter.lastSeen) > time.Second {
			limiter.count = 0
		}

		limiter.count++
		limiter.lastSeen = time.Now()

		if limiter.count > rl.rps {
			rl.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}

		rl.mu.Unlock()
		c.Next()
	}
}
