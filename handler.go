package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var internalErrorMessage = "An internal error occurred, please try again later. " +
	"If the problem persists drop us a line at hello@pensiondata.be"

// ListFunds return all funds available
func ListFunds(store Storage) gin.HandlerFunc {
	return func(context *gin.Context) {
		funds, err := store.ListFunds()

		if err != nil {
			log.Printf("Error while listing funds: %s", err)
			context.JSON(http.StatusInternalServerError, gin.H{
				"error":   http.StatusInternalServerError,
				"message": internalErrorMessage,
			})

			return
		}

		var publicFunds []PublicFund
		for _, fund := range funds {
			publicFunds = append(publicFunds, NewPublicFund(fund))
		}

		context.JSON(http.StatusOK, publicFunds)
	}
}

// GetFundByIsin return the fund for the given isin
func GetFundByIsin(store Storage) gin.HandlerFunc {
	return func(context *gin.Context) {
		isin := strings.ToUpper(context.Params.ByName("isin"))
		fund, err := store.GetFundByIsin(isin)

		if err != nil {
			if err == ErrNotFound {
				context.JSON(http.StatusNotFound, gin.H{
					"error":   http.StatusNotFound,
					"message": fmt.Sprintf("The fund %s was not found", isin),
				})

				return
			}

			log.Printf("Error while getting fund with isin %s: %s", isin, err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMessage})
			return
		}

		context.JSON(http.StatusOK, NewPublicFund(fund))
	}
}

// ListQuotes return all available quotes for the given fund
func ListQuotes(store Storage) gin.HandlerFunc {
	return func(context *gin.Context) {
		isin := strings.ToUpper(context.Params.ByName("isin"))
		_, err := store.GetFundByIsin(isin)

		if err != nil {
			if err == ErrNotFound {
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

		quotes, err := store.ListQuotes(isin)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error":   http.StatusInternalServerError,
				"message": internalErrorMessage,
			})

			return
		}

		var publicQuotes []PublicQuote
		for _, quote := range quotes {
			publicQuotes = append(publicQuotes, NewPublicQuote(quote))
		}

		context.JSON(http.StatusOK, publicQuotes)
	}
}

// GetQuoteByDate return the quote for the given date
func GetQuoteByDate(store Storage) gin.HandlerFunc {
	return func(context *gin.Context) {
		isin := strings.ToUpper(context.Params.ByName("isin"))
		date := context.Params.ByName("date")
		_, err := store.GetFundByIsin(isin)

		if err != nil {
			if err == ErrNotFound {
				context.JSON(http.StatusNotFound, gin.H{
					"error":   http.StatusNotFound,
					"message": fmt.Sprintf("The fund %s was not found", isin),
				})

				return
			}

			log.Printf("Error while getting quote by date: %s", err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMessage})

			return
		}

		// Hard code the "latest" route to avoid a wildcard route conflict in Gin
		if date == "latest" {
			getLatestQuote(context, store, isin)
			return
		}

		quote, err := store.GetQuoteByDate(isin, date)
		if err != nil {
			if err == ErrNotFound {
				context.JSON(http.StatusNotFound, gin.H{
					"error":   http.StatusNotFound,
					"message": fmt.Sprintf("The quote for fund %s on %s was not found", isin, date),
				})

				return
			}

			context.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMessage})
			return
		}

		context.JSON(http.StatusOK, NewPublicQuote(quote))
	}
}

func getLatestQuote(context *gin.Context, store Storage, isin string) {
	quote, err := store.GetLatestQuote(isin)

	if err != nil {
		if err == ErrNotFound {
			context.JSON(http.StatusNotFound, gin.H{
				"error":   http.StatusNotFound,
				"message": fmt.Sprintf("The latest quote for fund %s was not found", isin),
			})

			return
		}

		log.Printf("Error while getting latest quote: %s", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMessage})

		return
	}

	context.JSON(http.StatusOK, NewPublicQuote(quote))
}

// CreateQuote create a new quote in the storage
func CreateQuote(store Storage) gin.HandlerFunc {
	return func(context *gin.Context) {
		isin := strings.ToUpper(context.Params.ByName("isin"))

		var createQuote ScraperCreateQuote
		if err := context.BindJSON(&createQuote); err != nil {
			log.Printf("Error while binding request body to ScraperCreateQuote struct %v: %s", createQuote, err)
			context.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}

		fund, err := store.GetFundByIsin(isin)
		if err != nil {
			if err == ErrNotFound {
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

		quote, err := store.CreateQuote(fund.Isin, createQuote)
		if err != nil {
			log.Printf("Error while creating quote: %s", err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMessage})
			return
		}

		context.JSON(http.StatusCreated, NewPublicQuote(quote))
	}
}
