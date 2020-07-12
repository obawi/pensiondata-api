package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"gihtub.com/obawi/pensiondata-api"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestScraperAuthRequired(t *testing.T) {
	t.Run("return unauthorized error when SCRAPER-KEY header is not set", func(t *testing.T) {
		_ = os.Setenv("SCRAPER_KEY", "s3cr3t")
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.CreateQuoteFn = func(isin string, quote pensiondata.ScraperCreateQuote) (pensiondata.PublicQuote, error) {
			return pensiondata.PublicQuote{}, errors.New("internal error")
		}

		InitQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		scraperCreateQuote := pensiondata.ScraperCreateQuote{
			Date:  "2020-06-30",
			Price: decimal.NewFromFloat(7.99),
		}
		jsonScraperCreateQuote, _ := json.Marshal(scraperCreateQuote)

		req, _ := http.NewRequest("POST", "/funds/BE123/quotes", bytes.NewBuffer(jsonScraperCreateQuote))

		r.ServeHTTP(resp, req)

		if http.StatusUnauthorized != resp.Code {
			t.Errorf("want %d, got %d", http.StatusUnauthorized, resp.Code)
		}

		_ = os.Unsetenv("SCRAPER_KEY")
	})

	t.Run("return unauthorized error when SCRAPER_KEY env var is not set", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.CreateQuoteFn = func(isin string, quote pensiondata.ScraperCreateQuote) (pensiondata.PublicQuote, error) {
			return pensiondata.PublicQuote{}, errors.New("internal error")
		}

		InitQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		scraperCreateQuote := pensiondata.ScraperCreateQuote{
			Date:  "2020-06-30",
			Price: decimal.NewFromFloat(7.99),
		}
		jsonScraperCreateQuote, _ := json.Marshal(scraperCreateQuote)

		req, _ := http.NewRequest("POST", "/funds/BE123/quotes", bytes.NewBuffer(jsonScraperCreateQuote))
		req.Header.Add("SCRAPER-KEY", "s3cr3t")

		r.ServeHTTP(resp, req)

		if http.StatusUnauthorized != resp.Code {
			t.Errorf("want %d, got %d", http.StatusUnauthorized, resp.Code)
		}
	})
}
