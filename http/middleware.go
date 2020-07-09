package http

import (
	"os"

	"github.com/gin-gonic/gin"
)

// ScraperAuthRequired is a middleware to check if the request has been send by the scraper app
func ScraperAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("SCRAPER-KEY") != os.Getenv("SCRAPER_KEY") {
			c.AbortWithStatus(401)
			return
		}

		c.Next()
	}
}
