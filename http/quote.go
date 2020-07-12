package http

import (
	"fmt"
	"github.com/obawi/pensiondata-api"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

// QuoteHandler handle all the HTTP requests for Quote
type QuoteHandler struct {
	s pensiondata.QuoteService
}

// InitQuoteHandler initialize a new QuoteHandler and register routes
func InitQuoteHandler(router *gin.Engine, service pensiondata.QuoteService) *QuoteHandler {
	h := &QuoteHandler{s: service}

	router.GET("/funds/:isin/quotes", h.GetQuotes())
	router.GET("/funds/:isin/quotes/:date", h.GetQuoteByDate())
	router.POST("/funds/:isin/quotes", ScraperAuthRequired(), h.CreateQuote())

	return h
}

// GetQuotes return all quotes for the given fund
func (h QuoteHandler) GetQuotes() gin.HandlerFunc {
	return func(context *gin.Context) {
		isin := strings.ToUpper(context.Params.ByName("isin"))
		publicQuotes, err := h.s.GetQuotes(isin)

		if err != nil {
			if err == pensiondata.ErrFundNotFound {
				context.JSON(http.StatusNotFound, gin.H{
					"error":   http.StatusNotFound,
					"message": fmt.Sprintf("The fund %s was not found", isin),
				})
				return
			}

			log.Printf("Error while listing quotes: %s", err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMessage})
			return
		}

		context.JSON(http.StatusOK, publicQuotes)
	}
}

// GetQuoteByDate return the quote for the given date
func (h QuoteHandler) GetQuoteByDate() gin.HandlerFunc {
	return func(context *gin.Context) {
		isin := strings.ToUpper(context.Params.ByName("isin"))
		date := context.Params.ByName("date")

		var publicQuote pensiondata.PublicQuote
		var err error

		// Hard code the "latest" route to avoid a wildcard route conflict in Gin
		if date == "latest" {
			publicQuote, err = h.s.GetLatestQuote(isin)
		} else {
			publicQuote, err = h.s.GetQuote(isin, date)
		}

		if err != nil {
			if err == pensiondata.ErrFundNotFound {
				context.JSON(http.StatusNotFound, gin.H{
					"error":   http.StatusNotFound,
					"message": fmt.Sprintf("The fund %s was not found", isin),
				})
				return
			} else if err == pensiondata.ErrQuoteNotFound {
				context.JSON(http.StatusNotFound, gin.H{
					"error":   http.StatusNotFound,
					"message": fmt.Sprintf("The quote for fund %s on %s was not found", isin, date),
				})
				return
			}

			log.Printf("Error while getting quote by date: %s", err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMessage})
			return
		}

		context.JSON(http.StatusOK, publicQuote)
	}
}

// CreateQuote create a new quote
func (h QuoteHandler) CreateQuote() gin.HandlerFunc {
	return func(context *gin.Context) {
		isin := strings.ToUpper(context.Params.ByName("isin"))

		var createQuote pensiondata.ScraperCreateQuote
		if err := context.BindJSON(&createQuote); err != nil {
			log.Printf("Error while binding request body to ScraperCreateQuote struct %v: %s", createQuote, err)
			context.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}

		publicQuote, err := h.s.CreateQuote(isin, createQuote)
		if err != nil {
			if err == pensiondata.ErrFundNotFound {
				context.JSON(http.StatusNotFound, gin.H{
					"error":   http.StatusNotFound,
					"message": fmt.Sprintf("The fund %s was not found", isin),
				})
				return
			}

			log.Printf("Error while creating quote: %s", err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMessage})
			return
		}

		context.JSON(http.StatusCreated, publicQuote)
	}
}
