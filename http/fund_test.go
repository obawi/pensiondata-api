package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/obawi/pensiondata-api"
)

const contentTypeJson = "application/json; charset=utf-8"

func TestGetFunds(t *testing.T) {
	t.Run("return list of funds successfully", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := pensiondata.FundServiceMock{}
		s.GetFundsFn = func() ([]pensiondata.PublicFund, error) {
			return testPublicFunds(), nil
		}

		InitFundHandler(r, s)

		// Create a response recorder
		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds", nil)

		r.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Errorf("want %d, got %d", http.StatusOK, resp.Code)
		}
		if contentTypeJson != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", contentTypeJson, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return internal error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := pensiondata.FundServiceMock{}
		s.GetFundsFn = func() ([]pensiondata.PublicFund, error) {
			return []pensiondata.PublicFund{}, errors.New("internal error")
		}

		InitFundHandler(r, s)

		// Create a response recorder
		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds", nil)

		r.ServeHTTP(resp, req)

		if http.StatusInternalServerError != resp.Code {
			t.Errorf("want %d, got %d", http.StatusInternalServerError, resp.Code)
		}
		if contentTypeJson != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", contentTypeJson, resp.Header().Get("Content-Type"))
		}
	})
}

func TestGetFundByISIN(t *testing.T) {
	t.Run("return fund successfully", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := pensiondata.FundServiceMock{}
		s.GetFundByISINFn = func(isin string) (pensiondata.PublicFund, error) {
			return testPublicFund(), nil
		}

		InitFundHandler(r, s)

		// Create a response recorder
		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123", nil)

		r.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Errorf("want %d, got %d", http.StatusOK, resp.Code)
		}
		if contentTypeJson != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", contentTypeJson, resp.Header().Get("Content-Type"))
		}
	})

	t.Run("return not found error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		s := pensiondata.FundServiceMock{}
		s.GetFundByISINFn = func(isin string) (pensiondata.PublicFund, error) {
			return pensiondata.PublicFund{}, pensiondata.ErrFundNotFound
		}

		InitFundHandler(r, s)

		// Create a response recorder
		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123", nil)

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

		s := pensiondata.FundServiceMock{}
		s.GetFundByISINFn = func(isin string) (pensiondata.PublicFund, error) {
			return pensiondata.PublicFund{}, errors.New("internal error")
		}

		InitFundHandler(r, s)

		// Create a response recorder
		resp := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/funds/BE123", nil)

		r.ServeHTTP(resp, req)

		if http.StatusInternalServerError != resp.Code {
			t.Errorf("want %d, got %d", http.StatusInternalServerError, resp.Code)
		}
		if contentTypeJson != resp.Header().Get("Content-Type") {
			t.Errorf("want %s, got %s", contentTypeJson, resp.Header().Get("Content-Type"))
		}
	})
}

func testPublicFund() pensiondata.PublicFund {
	return pensiondata.PublicFund{
		Isin:       "BE123",
		Name:       "First Fund",
		Bank:       "Banka",
		LaunchDate: "2020-06-27",
		Currency:   "EUR",
	}
}

func testPublicFunds() []pensiondata.PublicFund {
	return []pensiondata.PublicFund{
		{
			Isin:       "BE123",
			Name:       "First Fund",
			Bank:       "Banka",
			LaunchDate: "2020-06-27",
			Currency:   "EUR",
		},
		{
			Isin:       "LU123",
			Name:       "Second Fund",
			Bank:       "Banko",
			LaunchDate: "2020-06-27",
			Currency:   "EUR",
		},
	}
}
