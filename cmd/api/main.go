package main

import (
	"gihtub.com/obawi/pensiondata-api/http"
	"gihtub.com/obawi/pensiondata-api/postgres"
	"log"
	"os"

	"gihtub.com/obawi/pensiondata-api"
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
	http.NewFundHandler(router, fundService)

	quoteRepo := postgres.NewQuoteRepository(db)
	quoteService := pensiondata.NewQuoteService(fundRepo, quoteRepo)
	http.NewQuoteHandler(router, quoteService)

	if len(os.Getenv("ALWAYSDATA_HTTPD_IP")) != 0 && len(os.Getenv("ALWAYSDATA_HTTPD_PORT")) != 0 {
		router.Run(os.Getenv("ALWAYSDATA_HTTPD_IP") + ":" + os.Getenv("ALWAYSDATA_HTTPD_PORT"))
	}

	router.Run()
}
