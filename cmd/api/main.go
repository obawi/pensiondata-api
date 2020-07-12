package main

import (
	"github.com/obawi/pensiondata-api/http"
	"github.com/obawi/pensiondata-api/postgres"
	"log"
	"os"

	"github.com/obawi/pensiondata-api"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := postgres.NewConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	fundRepo := postgres.NewFundRepository(db)
	fundService := pensiondata.NewFundService(fundRepo)
	http.InitFundHandler(router, fundService)

	quoteRepo := postgres.NewQuoteRepository(db)
	quoteService := pensiondata.NewQuoteService(fundRepo, quoteRepo)
	http.InitQuoteHandler(router, quoteService)

	if len(os.Getenv("ALWAYSDATA_HTTPD_IP")) != 0 && len(os.Getenv("ALWAYSDATA_HTTPD_PORT")) != 0 {
		router.Run(os.Getenv("ALWAYSDATA_HTTPD_IP") + ":" + os.Getenv("ALWAYSDATA_HTTPD_PORT"))
	}

	router.Run()
}
