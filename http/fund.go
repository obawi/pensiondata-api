package http

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/obawi/pensiondata-api"
)

// FundHandler handle all the HTTP requests for Fund
type FundHandler struct {
	s pensiondata.FundService
}

// InitFundHandler initialize a new FundHandler and register routes
func InitFundHandler(router *gin.Engine, service pensiondata.FundService) {
	h := &FundHandler{s: service}

	// setup routes
	router.GET("/funds", h.GetFunds())
	router.GET("/funds/:isin", h.GetFundByISIN())
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

// GetFundByISIN return the fund for the given isin
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
