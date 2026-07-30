package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Client struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}

var (
	clients = make(map[string]*Client)
	mu      sync.Mutex
)

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	client, exists := clients[ip]

	if !exists {
		limiter := rate.NewLimiter(1, 2)
		clients[ip] = &Client{
			Limiter:  limiter,
			LastSeen: time.Now(),
		}
		return limiter
	}
	client.LastSeen = time.Now()
	return client.Limiter
}

func cleanUpClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()

		for ip, client := range clients {
			if time.Since(client.LastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}
func init() {
	go cleanUpClients()
}

func RateLimiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		limiter := getLimiter(ip)
		allowed := limiter.Allow()
		fmt.Println("IP:", ip)
		fmt.Println("Allowed:", allowed)
		if !limiter.Allow() {
			ctx.JSON(429, gin.H{
				"error": "Too Many Requests",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
