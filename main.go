package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_NAME"), os.Getenv("DATABASE_SSLMODE"))

	db, err := sql.Open(os.Getenv("DATABASE_DRIVER"), connStr)
	if err != nil {
		log.Fatal(err)
	}

	store := PostgresStore{db: db}

	router := gin.Default()

	fundGroup := router.Group("/funds")
	{
		fundGroup.GET("/", ListFunds(store))
		fundGroup.GET("/:isin", GetFundByIsin(store))
		fundGroup.GET("/:isin/quotes", ListQuotes(store))
		fundGroup.GET("/:isin/quotes/:date", GetQuoteByDate(store))

		createQuoteGroup := fundGroup.Group("/")
		createQuoteGroup.Use(ScraperAuthRequired())
		{
			createQuoteGroup.POST("/:isin/quotes", CreateQuote(store))
		}
	}

	if len(os.Getenv("ALWAYSDATA_HTTPD_IP")) != 0 && len(os.Getenv("ALWAYSDATA_HTTPD_PORT")) != 0 {
		router.Run(os.Getenv("ALWAYSDATA_HTTPD_IP") + ":" + os.Getenv("ALWAYSDATA_HTTPD_PORT"))
	}

	router.Run()
}
