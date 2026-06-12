package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate" // package for ratelimiting
)

type visitor struct {
	limiter *rate.Limiter // rate limiter for this visitor
}

var (
	visitors = make(map[string]*visitor) // make function to initialize the map
	mu       sync.Mutex                  // protects the visitors map
)

// RateLimit allows (rps) requests per second with a burst of (burst) per client IP.
func RateLimit(rps float64, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()                 // lock the mutex before accessing the visitors map
		v, exists := visitors[ip] // check if we already have a visitor for this IP
		if !exists {
			v = &visitor{limiter: rate.NewLimiter(rate.Limit(rps), burst)}
			visitors[ip] = v
		}
		mu.Unlock() // unlock the mutex after we're done accessing the visitors map

		// Check if the request is allowed by the rate limiter. If not, return a 429 Too Many Requests response.
		if !v.limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort() // stop processing the request further
			return
		}

		c.Next() // move to the next visitor in the chain if the request is allowed
	}
}
