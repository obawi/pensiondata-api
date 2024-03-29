package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/obawi/pensiondata-api"
	"github.com/shopspring/decimal"
)

func TestGetQuotes(t *testing.T) {
	t.Run("return list of quotes successfully", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.GetQuotesFn = func(isin string) ([]pensiondata.PublicQuote, error) {
			return testPublicQuotes(), nil
		}

		InitQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123/quotes", nil)

		r.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Errorf("want %d, got %d", http.StatusOK, resp.Code)
		}
		if contentTypeJson != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", contentTypeJson, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return not found error for fund", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.GetQuotesFn = func(isin string) ([]pensiondata.PublicQuote, error) {
			return []pensiondata.PublicQuote{}, pensiondata.ErrFundNotFound
		}

		InitQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123/quotes", nil)

		r.ServeHTTP(resp, req)

		if http.StatusNotFound != resp.Code {
			t.Errorf("want %d, got %d", http.StatusNotFound, resp.Code)
		}
		if contentTypeJson != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", contentTypeJson, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return internal error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.GetQuotesFn = func(isin string) ([]pensiondata.PublicQuote, error) {
			return []pensiondata.PublicQuote{}, errors.New("internal error")
		}

		InitQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123/quotes", nil)

		r.ServeHTTP(resp, req)

		if http.StatusInternalServerError != resp.Code {
			t.Errorf("want %d, got %d", http.StatusInternalServerError, resp.Code)
		}
		if contentTypeJson != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", contentTypeJson, resp.Header().Get("Content-Type"))
		}
	})
}

func TestGetQuoteByDate(t *testing.T) {
	t.Run("return quote successfully", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.GetQuoteFn = func(isin, date string) (pensiondata.PublicQuote, error) {
			return testPublicQuote(), nil
		}

		InitQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123/quotes/2020-06-28", nil)

		r.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Errorf("want %d, got %d", http.StatusOK, resp.Code)
		}
		if contentTypeJson != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", contentTypeJson, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return not found error for found", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.GetQuoteFn = func(isin, date string) (pensiondata.PublicQuote, error) {
			return pensiondata.PublicQuote{}, pensiondata.ErrFundNotFound
		}

		InitQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123/quotes/2020-06-28", nil)

		r.ServeHTTP(resp, req)

		if http.StatusNotFound != resp.Code {
			t.Errorf("want %d, got %d", http.StatusNotFound, resp.Code)
		}
		if contentTypeJson != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", contentTypeJson, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return not found error for quote", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.GetQuoteFn = func(isin, date string) (pensiondata.PublicQuote, error) {
			return pensiondata.PublicQuote{}, pensiondata.ErrQuoteNotFound
		}

		InitQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123/quotes/2020-06-28", nil)

		r.ServeHTTP(resp, req)

		if http.StatusNotFound != resp.Code {
			t.Errorf("want %d, got %d", http.StatusNotFound, resp.Code)
		}
		if contentTypeJson != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", contentTypeJson, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return internal error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.GetQuoteFn = func(isin, date string) (pensiondata.PublicQuote, error) {
			return pensiondata.PublicQuote{}, errors.New("internal error")
		}

		InitQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123/quotes/2020-06-28", nil)

		r.ServeHTTP(resp, req)

		if http.StatusInternalServerError != resp.Code {
			t.Errorf("want %d, got %d", http.StatusInternalServerError, resp.Code)
		}
		if contentTypeJson != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", contentTypeJson, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return latest quote successfully", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.GetLatestQuoteFn = func(isin string) (pensiondata.PublicQuote, error) {
			return testPublicQuote(), nil
		}

		InitQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123/quotes/latest", nil)

		r.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Errorf("want %d, got %d", http.StatusOK, resp.Code)
		}
		if contentTypeJson != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", contentTypeJson, resp.Header().Get("Content-Type"))
		}
	})
}

func TestCreateQuote(t *testing.T) {
	t.Run("create quote successfully", func(t *testing.T) {
		_ = os.Setenv("SCRAPER_KEY", "s3cr3t")
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.CreateQuoteFn = func(isin string, quote pensiondata.ScraperCreateQuote) (pensiondata.PublicQuote, error) {
			return testPublicQuote(), nil
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

		if http.StatusCreated != resp.Code {
			t.Errorf("want %d, got %d", http.StatusCreated, resp.Code)
		}
		if contentTypeJson != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", contentTypeJson, resp.Header().Get("Content-Type"))
		}
		_ = os.Unsetenv("SCRAPER_KEY")
	})

	t.Run("return bad request error for invalid quote JSON binding", func(t *testing.T) {
		_ = os.Setenv("SCRAPER_KEY", "s3cr3t")
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.CreateQuoteFn = func(isin string, quote pensiondata.ScraperCreateQuote) (pensiondata.PublicQuote, error) {
			return pensiondata.PublicQuote{}, nil
		}

		InitQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		jsonScraperCreateQuote, _ := json.Marshal(false)

		req, _ := http.NewRequest(http.MethodPost, "/funds/BE123/quotes", bytes.NewBuffer(jsonScraperCreateQuote))
		req.Header.Add("SCRAPER-KEY", "s3cr3t")

		r.ServeHTTP(resp, req)

		if http.StatusBadRequest != resp.Code {
			t.Errorf("want %d, got %d", http.StatusBadRequest, resp.Code)
		}
		if contentTypeJson != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", contentTypeJson, resp.Header().Get("Content-Type"))
		}
		_ = os.Unsetenv("SCRAPER_KEY")
	})

	t.Run("return not found error for fund", func(t *testing.T) {
		_ = os.Setenv("SCRAPER_KEY", "s3cr3t")
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.CreateQuoteFn = func(isin string, quote pensiondata.ScraperCreateQuote) (pensiondata.PublicQuote, error) {
			return pensiondata.PublicQuote{}, pensiondata.ErrFundNotFound
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

		if http.StatusNotFound != resp.Code {
			t.Errorf("want %d, got %d", http.StatusNotFound, resp.Code)
		}
		if contentTypeJson != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", contentTypeJson, resp.Header().Get("Content-Type"))
		}
		_ = os.Unsetenv("SCRAPER_KEY")
	})

	t.Run("return internal error", func(t *testing.T) {
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
		req.Header.Add("SCRAPER-KEY", "s3cr3t")

		r.ServeHTTP(resp, req)

		if http.StatusInternalServerError != resp.Code {
			t.Errorf("want %d, got %d", http.StatusInternalServerError, resp.Code)
		}
		if contentTypeJson != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", contentTypeJson, resp.Header().Get("Content-Type"))
		}
		_ = os.Unsetenv("SCRAPER_KEY")
	})
}

func testPublicQuote() pensiondata.PublicQuote {
	return pensiondata.PublicQuote{Price: 5.99, Date: "2020-06-28"}
}

func testPublicQuotes() []pensiondata.PublicQuote {
	return []pensiondata.PublicQuote{
		{Date: "2020-06-28", Price: 5.99},
		{Date: "2020-06-27", Price: 5.99},
	}
}
