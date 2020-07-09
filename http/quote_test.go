package http

import (
	"bytes"
	"encoding/json"
	// "bytes"
	// "encoding/json"
	"errors"
	"gihtub.com/obawi/pensiondata-api"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetQuotes(t *testing.T) {
	t.Run("return list of quotes successfully", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.GetQuotesFn = func(isin string) ([]pensiondata.PublicQuote, error) {
			return testPublicQuotes(), nil
		}

		NewQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123/quotes", nil)

		r.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Errorf("want %d, got %d", http.StatusOK, resp.Code)
		}
		if CONTENT_TYPE_JSON != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", CONTENT_TYPE_JSON, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return not found error for fund", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.GetQuotesFn = func(isin string) ([]pensiondata.PublicQuote, error) {
			return []pensiondata.PublicQuote{}, pensiondata.ErrFundNotFound
		}

		NewQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123/quotes", nil)

		r.ServeHTTP(resp, req)

		if http.StatusNotFound != resp.Code {
			t.Errorf("want %d, got %d", http.StatusNotFound, resp.Code)
		}
		if CONTENT_TYPE_JSON != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", CONTENT_TYPE_JSON, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return internal error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.GetQuotesFn = func(isin string) ([]pensiondata.PublicQuote, error) {
			return []pensiondata.PublicQuote{}, errors.New("internal error")
		}

		NewQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123/quotes", nil)

		r.ServeHTTP(resp, req)

		if http.StatusInternalServerError != resp.Code {
			t.Errorf("want %d, got %d", http.StatusInternalServerError, resp.Code)
		}
		if CONTENT_TYPE_JSON != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", CONTENT_TYPE_JSON, resp.Header().Get("Content-Type"))
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

		NewQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123/quotes/2020-06-28", nil)

		r.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Errorf("want %d, got %d", http.StatusOK, resp.Code)
		}
		if CONTENT_TYPE_JSON != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", CONTENT_TYPE_JSON, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return not found error for found", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.GetQuoteFn = func(isin, date string) (pensiondata.PublicQuote, error) {
			return pensiondata.PublicQuote{}, pensiondata.ErrFundNotFound
		}

		NewQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123/quotes/2020-06-28", nil)

		r.ServeHTTP(resp, req)

		if http.StatusNotFound != resp.Code {
			t.Errorf("want %d, got %d", http.StatusNotFound, resp.Code)
		}
		if CONTENT_TYPE_JSON != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", CONTENT_TYPE_JSON, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return not found error for quote", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.GetQuoteFn = func(isin, date string) (pensiondata.PublicQuote, error) {
			return pensiondata.PublicQuote{}, pensiondata.ErrQuoteNotFound
		}

		NewQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123/quotes/2020-06-28", nil)

		r.ServeHTTP(resp, req)

		if http.StatusNotFound != resp.Code {
			t.Errorf("want %d, got %d", http.StatusNotFound, resp.Code)
		}
		if CONTENT_TYPE_JSON != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", CONTENT_TYPE_JSON, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return internal error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.GetQuoteFn = func(isin, date string) (pensiondata.PublicQuote, error) {
			return pensiondata.PublicQuote{}, errors.New("internal error")
		}

		NewQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123/quotes/2020-06-28", nil)

		r.ServeHTTP(resp, req)

		if http.StatusInternalServerError != resp.Code {
			t.Errorf("want %d, got %d", http.StatusInternalServerError, resp.Code)
		}
		if CONTENT_TYPE_JSON != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", CONTENT_TYPE_JSON, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return latest quote successfully", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.GetLatestQuoteFn = func(isin string) (pensiondata.PublicQuote, error) {
			return testPublicQuote(), nil
		}

		NewQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123/quotes/latest", nil)

		r.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Errorf("want %d, got %d", http.StatusOK, resp.Code)
		}
		if CONTENT_TYPE_JSON != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", CONTENT_TYPE_JSON, resp.Header().Get("Content-Type"))
		}
	})
}

func TestCreateQuote(t *testing.T) {
	t.Run("create quote successfully", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.CreateQuoteFn = func(isin string, quote pensiondata.ScraperCreateQuote) (pensiondata.PublicQuote, error) {
			return testPublicQuote(), nil
		}

		NewQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		scraperCreateQuote := pensiondata.ScraperCreateQuote{
			Date:  "2020-06-30",
			Price: decimal.NewFromFloat(7.99),
		}
		jsonScraperCreateQuote, _ := json.Marshal(scraperCreateQuote)

		req, _ := http.NewRequest("POST", "/funds/BE123/quotes", bytes.NewBuffer(jsonScraperCreateQuote))

		r.ServeHTTP(resp, req)

		if http.StatusCreated != resp.Code {
			t.Errorf("want %d, got %d", http.StatusCreated, resp.Code)
		}
		if CONTENT_TYPE_JSON != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", CONTENT_TYPE_JSON, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return bad request error for invalid quote JSON binding", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.CreateQuoteFn = func(isin string, quote pensiondata.ScraperCreateQuote) (pensiondata.PublicQuote, error) {
			return pensiondata.PublicQuote{}, nil
		}

		NewQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		jsonScraperCreateQuote, _ := json.Marshal(false)

		req, _ := http.NewRequest(http.MethodPost, "/funds/BE123/quotes", bytes.NewBuffer(jsonScraperCreateQuote))

		r.ServeHTTP(resp, req)

		if http.StatusBadRequest != resp.Code {
			t.Errorf("want %d, got %d", http.StatusBadRequest, resp.Code)
		}
		if CONTENT_TYPE_JSON != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", CONTENT_TYPE_JSON, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return not found error for fund", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.CreateQuoteFn = func(isin string, quote pensiondata.ScraperCreateQuote) (pensiondata.PublicQuote, error) {
			return pensiondata.PublicQuote{}, pensiondata.ErrFundNotFound
		}

		NewQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		scraperCreateQuote := pensiondata.ScraperCreateQuote{
			Date:  "2020-06-30",
			Price: decimal.NewFromFloat(7.99),
		}
		jsonScraperCreateQuote, _ := json.Marshal(scraperCreateQuote)

		req, _ := http.NewRequest("POST", "/funds/BE123/quotes", bytes.NewBuffer(jsonScraperCreateQuote))

		r.ServeHTTP(resp, req)

		if http.StatusNotFound != resp.Code {
			t.Errorf("want %d, got %d", http.StatusNotFound, resp.Code)
		}
		if CONTENT_TYPE_JSON != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", CONTENT_TYPE_JSON, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return internal error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		quoteService := pensiondata.QuoteServiceMock{}
		quoteService.CreateQuoteFn = func(isin string, quote pensiondata.ScraperCreateQuote) (pensiondata.PublicQuote, error) {
			return pensiondata.PublicQuote{}, errors.New("internal error")
		}

		NewQuoteHandler(r, quoteService)

		resp := httptest.NewRecorder()

		scraperCreateQuote := pensiondata.ScraperCreateQuote{
			Date:  "2020-06-30",
			Price: decimal.NewFromFloat(7.99),
		}
		jsonScraperCreateQuote, _ := json.Marshal(scraperCreateQuote)

		req, _ := http.NewRequest("POST", "/funds/BE123/quotes", bytes.NewBuffer(jsonScraperCreateQuote))

		r.ServeHTTP(resp, req)

		if http.StatusInternalServerError != resp.Code {
			t.Errorf("want %d, got %d", http.StatusInternalServerError, resp.Code)
		}
		if CONTENT_TYPE_JSON != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", CONTENT_TYPE_JSON, resp.Header().Get("Content-Type"))
		}
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
