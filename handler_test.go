package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const CONTENT_TYPE_JSON = "application/json; charset=utf-8"

func TestListFunds(t *testing.T) {
	t.Run("return list of funds successfully", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		// Mock Storage and ListFunds()
		s := MockStore{}
		s.ListFundsFn = func() ([]Fund, error) {
			return testFunds(), nil
		}

		r.GET("/funds", ListFunds(s))

		// Create a response recorder
		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds", nil)

		r.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Errorf("want %d, got %d", http.StatusOK, resp.Code)
		}
		if CONTENT_TYPE_JSON != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", CONTENT_TYPE_JSON, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return internal error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := MockStore{}
		s.ListFundsFn = func() ([]Fund, error) {
			return []Fund{}, errors.New("internal error")
		}

		r.GET("/funds", ListFunds(s))

		// Create a response recorder
		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds", nil)

		r.ServeHTTP(resp, req)

		if http.StatusInternalServerError != resp.Code {
			t.Errorf("want %d, got %d", http.StatusInternalServerError, resp.Code)
		}
		if CONTENT_TYPE_JSON != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", CONTENT_TYPE_JSON, resp.Header().Get("Content-Type"))
		}
	})
}

func TestGetFundByIsin(t *testing.T) {
	t.Run("return fund successfully", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return testFund(), nil
		}

		r.GET("/funds/:isin", GetFundByIsin(s))

		// Create a response recorder
		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123", nil)

		r.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Errorf("want %d, got %d", http.StatusOK, resp.Code)
		}
		if CONTENT_TYPE_JSON != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", CONTENT_TYPE_JSON, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return not found error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return Fund{}, ErrNotFound
		}

		r.GET("/funds/:isin", GetFundByIsin(s))

		// Create a response recorder
		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123", nil)

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

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return Fund{}, errors.New("internal error")
		}

		r.GET("/funds/:isin", GetFundByIsin(s))

		// Create a response recorder
		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123", nil)

		r.ServeHTTP(resp, req)

		if http.StatusInternalServerError != resp.Code {
			t.Errorf("want %d, got %d", http.StatusInternalServerError, resp.Code)
		}
		if CONTENT_TYPE_JSON != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", CONTENT_TYPE_JSON, resp.Header().Get("Content-Type"))
		}
	})
}

func TestListQuotes(t *testing.T) {
	t.Run("return list of quotes successfully", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return testFund(), nil
		}
		s.ListQuotesFn = func(isin string) ([]Quote, error) {
			return testQuotes(), nil
		}

		r.GET("/funds/:isin/quotes", ListQuotes(s))

		// Create a response recorder
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

	t.Run("return not found error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return testFund(), ErrNotFound
		}

		r.GET("/funds/:isin/quotes", ListQuotes(s))

		// Create a response recorder
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

	t.Run("return internal error for fund", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return Fund{}, errors.New("internal error")
		}

		r.GET("/funds/:isin/quotes", ListQuotes(s))

		// Create a response recorder
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

	t.Run("return internal error for quotes", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return testFund(), nil
		}
		s.ListQuotesFn = func(isin string) ([]Quote, error) {
			return []Quote{}, errors.New("internal error")
		}

		r.GET("/funds/:isin/quotes", ListQuotes(s))

		// Create a response recorder
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

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return testFund(), nil
		}
		s.GetQuoteByDateFn = func(isin string, date string) (Quote, error) {
			return testQuote(), nil
		}

		r.GET("/funds/:isin/quotes/:date", GetQuoteByDate(s))

		// Create a response recorder
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

	t.Run("return not found error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return testFund(), ErrNotFound
		}

		r.GET("/funds/:isin/quotes/:date", GetQuoteByDate(s))

		// Create a response recorder
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

	t.Run("return internal error for fund", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return Fund{}, errors.New("internal error")
		}

		r.GET("/funds/:isin/quotes/:date", GetQuoteByDate(s))

		// Create a response recorder
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

	t.Run("return not found error for quote", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return testFund(), nil
		}
		s.GetQuoteByDateFn = func(isin string, date string) (Quote, error) {
			return Quote{}, ErrNotFound
		}

		r.GET("/funds/:isin/quotes/:date", GetQuoteByDate(s))

		// Create a response recorder
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

	t.Run("return internal error for quote", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return testFund(), nil
		}
		s.GetQuoteByDateFn = func(isin string, date string) (Quote, error) {
			return Quote{}, errors.New("internal error")
		}

		r.GET("/funds/:isin/quotes/:date", GetQuoteByDate(s))

		// Create a response recorder
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

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return testFund(), nil
		}
		s.GetLatestQuoteFn = func(isin string) (Quote, error) {
			return testQuote(), nil
		}

		r.GET("/funds/:isin/quotes/:date", GetQuoteByDate(s))

		// Create a response recorder
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

	t.Run("return not found error for latest quote", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return testFund(), nil
		}
		s.GetLatestQuoteFn = func(isin string) (Quote, error) {
			return Quote{}, ErrNotFound
		}

		r.GET("/funds/:isin/quotes/:date", GetQuoteByDate(s))

		// Create a response recorder
		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123/quotes/latest", nil)

		r.ServeHTTP(resp, req)

		if http.StatusNotFound != resp.Code {
			t.Errorf("want %d, got %d", http.StatusNotFound, resp.Code)
		}
		if CONTENT_TYPE_JSON != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", CONTENT_TYPE_JSON, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return internal error for latest quote", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return testFund(), nil
		}
		s.GetLatestQuoteFn = func(isin string) (Quote, error) {
			return Quote{}, errors.New("internal error")
		}

		r.GET("/funds/:isin/quotes/:date", GetQuoteByDate(s))

		// Create a response recorder
		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123/quotes/latest", nil)

		r.ServeHTTP(resp, req)

		if http.StatusInternalServerError != resp.Code {
			t.Errorf("want %d, got %d", http.StatusInternalServerError, resp.Code)
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

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return testFund(), nil
		}
		s.CreateQuoteFn = func(isin string, scraperCreateQuote ScraperCreateQuote) (Quote, error) {
			return testQuote(), nil
		}

		r.POST("/funds/:isin/quotes", CreateQuote(s))

		// Create a response recorder
		resp := httptest.NewRecorder()

		scraperCreateQuote := ScraperCreateQuote{
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

		s := MockStore{}
		r.POST("/funds/:isin/quotes", CreateQuote(s))

		// Create a response recorder
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

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return Fund{}, ErrNotFound
		}

		r.POST("/funds/:isin/quotes", CreateQuote(s))

		// Create a response recorder
		resp := httptest.NewRecorder()

		scraperCreateQuote := ScraperCreateQuote{
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

	t.Run("return internal error for fund", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return Fund{}, errors.New("internal error")
		}

		r.POST("/funds/:isin/quotes", CreateQuote(s))

		// Create a response recorder
		resp := httptest.NewRecorder()

		scraperCreateQuote := ScraperCreateQuote{
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

	t.Run("return not found error for quote", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := MockStore{}
		s.GetFundByIsinFn = func(isin string) (Fund, error) {
			return testFund(), nil
		}
		s.CreateQuoteFn = func(isin string, quote ScraperCreateQuote) (Quote, error) {
			return Quote{}, errors.New("internal error")
		}

		r.POST("/funds/:isin/quotes", CreateQuote(s))

		// Create a response recorder
		resp := httptest.NewRecorder()

		scraperCreateQuote := ScraperCreateQuote{
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

func testFund() Fund {
	date, _ := time.Parse("2006-01-02", "2020-06-27")
	return Fund{
		Isin:       "BE123",
		Name:       "First Fund",
		Bank:       "Banka",
		LaunchDate: date,
		Currency:   "EUR",
	}
}

func testFunds() []Fund {
	return []Fund{
		{
			Isin:       "BE123",
			Name:       "First Fund",
			Bank:       "Banka",
			LaunchDate: time.Now(),
			Currency:   "EUR",
		},
		{
			Isin:       "LU123",
			Name:       "Second Fund",
			Bank:       "Banko",
			LaunchDate: time.Now(),
			Currency:   "EUR",
		},
	}
}

func testQuote() Quote {
	date, _ := time.Parse("2006-01-02", "2020-06-28")
	return Quote{
		Price: decimal.NewFromFloat(5.99),
		Date:  date,
	}
}

func testQuotes() []Quote {
	return []Quote{
		{
			Date:  time.Now(),
			Price: decimal.NewFromFloat(5.99),
		},
		{
			Date:  time.Now(),
			Price: decimal.NewFromFloat(5.99),
		},
	}
}
