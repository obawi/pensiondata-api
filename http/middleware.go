package http

import (
	"os"

	"github.com/gin-gonic/gin"
)

// ScraperAuthRequired is a middleware to check if the request has been send by the scraper app
func ScraperAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		scraperKey, isScraperKeySet := os.LookupEnv("SCRAPER_KEY")
		if !isScraperKeySet || c.GetHeader("SCRAPER-KEY") != scraperKey {
			c.AbortWithStatus(401)
			return
		}

		c.Next()
	}
}
