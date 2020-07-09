package http

import (
	"fmt"
	"gihtub.com/obawi/pensiondata-api"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

var internalErrorMessage = "An internal error occurred, please try again later. " +
	"If the problem persists drop us a line at hello@pensiondata.eu"

type FundHandler struct {
	s pensiondata.FundService
}

// NewFundHandler return a new FundHandler and register routes
func NewFundHandler(router *gin.Engine, service pensiondata.FundService) *FundHandler {
	h := &FundHandler{s: service}

	// setup routes
	router.GET("/funds", h.GetFunds())
	router.GET("/funds/:isin", h.GetFundByISIN())

	return h
}

// GetFunds return all funds
func (h FundHandler) GetFunds() gin.HandlerFunc {
	return func(context *gin.Context) {
		publicFunds, err := h.s.GetFunds()
		if err != nil {
			log.Printf("Error while listing funds: %s", err)
			context.JSON(http.StatusInternalServerError, gin.H{
				"error":   http.StatusInternalServerError,
				"message": internalErrorMessage,
			})
			return
		}
		context.JSON(http.StatusOK, publicFunds)
	}
}

// GetFundByISIN return a fund for the given isin
func (h FundHandler) GetFundByISIN() gin.HandlerFunc {
	return func(context *gin.Context) {
		isin := strings.ToUpper(context.Params.ByName("isin"))
		publicFund, err := h.s.GetFundByISIN(isin)
		if err != nil {
			if err == pensiondata.ErrFundNotFound {
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
		context.JSON(http.StatusOK, publicFund)
	}
}
