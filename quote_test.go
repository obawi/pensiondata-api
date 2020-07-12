package pensiondata

import (
	"errors"
	"github.com/shopspring/decimal"
	"reflect"
	"testing"
	"time"
)

func TestGetQuote(t *testing.T) {
	t.Run("return quote successfully", func(t *testing.T) {
		date, _ := time.Parse("2006-01-02", "2020-06-27")
		want := Quote{Date: date, Price: decimal.NewFromFloat(5.99)}

		fund := Fund{Isin: "BE123"}
		fundRepo := FundRepositoryMock{}
		fundRepo.FindByISINFn = func(isin string) (Fund, error) {
			return fund, nil
		}

		quoteRepo := QuoteRepositoryMock{}
		quoteRepo.FindByISINAndDateFn = func(isin, date string) (Quote, error) {
			return want, nil
		}

		s := NewQuoteService(fundRepo, quoteRepo)
		got, _ := s.GetQuote("BE123", "2020-06-27")

		if !reflect.DeepEqual(newPublicQuote(want), got) {
			t.Errorf("want %v, got %v", newPublicQuote(want), got)
		}
	})

	t.Run("return error for fund", func(t *testing.T) {
		fundRepo := FundRepositoryMock{}
		fundRepo.FindByISINFn = func(isin string) (Fund, error) {
			return Fund{}, errors.New("error")
		}

		quoteRepo := QuoteRepositoryMock{}
		quoteRepo.FindByISINAndDateFn = func(isin, date string) (Quote, error) {
			return Quote{}, nil
		}

		s := NewQuoteService(fundRepo, quoteRepo)
		_, err := s.GetQuote("BE123", "2020-06-27")

		if err == nil {
			t.Errorf("want error")
		}
	})

	t.Run("return error for quote", func(t *testing.T) {
		fundRepo := FundRepositoryMock{}
		fund := Fund{Isin: "BE123"}
		fundRepo.FindByISINFn = func(isin string) (Fund, error) {
			return fund, nil
		}

		quoteRepo := QuoteRepositoryMock{}
		quoteRepo.FindByISINAndDateFn = func(isin, date string) (Quote, error) {
			return Quote{}, errors.New("error")
		}

		s := NewQuoteService(fundRepo, quoteRepo)
		_, err := s.GetQuote("BE123", "2020-06-27")

		if err == nil {
			t.Errorf("want error")
		}
	})
}

func TestLatestQuote(t *testing.T) {
	t.Run("return quote successfully", func(t *testing.T) {
		date, _ := time.Parse("2006-01-02", "2020-06-27")
		want := Quote{Date: date, Price: decimal.NewFromFloat(5.99)}

		fundRepo := FundRepositoryMock{}
		quoteRepo := QuoteRepositoryMock{}
		quoteRepo.FindByDateDescFn = func(isin string) (Quote, error) {
			return want, nil
		}

		s := NewQuoteService(fundRepo, quoteRepo)
		got, _ := s.GetLatestQuote("BE123")

		if !reflect.DeepEqual(newPublicQuote(want), got) {
			t.Errorf("want %v, got %v", newPublicQuote(want), got)
		}
	})

	t.Run("return error", func(t *testing.T) {
		fundRepo := FundRepositoryMock{}
		quoteRepo := QuoteRepositoryMock{}
		quoteRepo.FindByDateDescFn = func(isin string) (Quote, error) {
			return Quote{}, errors.New("error")
		}

		s := NewQuoteService(fundRepo, quoteRepo)
		_, err := s.GetLatestQuote("BE123")

		if err == nil {
			t.Errorf("want error")
		}
	})
}

func TestGetQuotes(t *testing.T) {
	t.Run("return quotes successfully", func(t *testing.T) {
		date, _ := time.Parse("2006-01-02", "2020-06-27")
		wants := []Quote{{Date: date, Price: decimal.NewFromFloat(5.99)},
			{Date: date, Price: decimal.NewFromFloat(7.99)}}

		fundRepo := FundRepositoryMock{}
		fundRepo.FindByISINFn = func(isin string) (Fund, error) {
			return Fund{}, nil
		}

		quoteRepo := QuoteRepositoryMock{}
		quoteRepo.FindAllFn = func(isin string) ([]Quote, error) {
			return wants, nil
		}

		s := NewQuoteService(fundRepo, quoteRepo)
		got, _ := s.GetQuotes("BE123")

		var wantPublicQuote []PublicQuote
		for _, want := range wants {
			wantPublicQuote = append(wantPublicQuote, newPublicQuote(want))
		}

		if !reflect.DeepEqual(wantPublicQuote, got) {
			t.Errorf("want %v, got %v", wantPublicQuote, got)
		}
	})

	t.Run("return error for fund", func(t *testing.T) {
		fundRepo := FundRepositoryMock{}
		fundRepo.FindByISINFn = func(isin string) (Fund, error) {
			return Fund{}, errors.New("error")
		}

		quoteRepo := QuoteRepositoryMock{}
		quoteRepo.FindAllFn = func(isin string) ([]Quote, error) {
			return []Quote{}, nil
		}

		s := NewQuoteService(fundRepo, quoteRepo)
		_, err := s.GetQuotes("BE123")

		if err == nil {
			t.Errorf("want error")
		}
	})

	t.Run("return error for quote", func(t *testing.T) {
		fundRepo := FundRepositoryMock{}
		fundRepo.FindByISINFn = func(isin string) (Fund, error) {
			return Fund{}, nil
		}

		quoteRepo := QuoteRepositoryMock{}
		quoteRepo.FindAllFn = func(isin string) ([]Quote, error) {
			return []Quote{}, errors.New("error")
		}

		s := NewQuoteService(fundRepo, quoteRepo)
		_, err := s.GetQuotes("BE123")

		if err == nil {
			t.Errorf("want error")
		}
	})
}

func TestCreateQuote(t *testing.T) {
	t.Run("create quote successfully", func(t *testing.T) {
		want := ScraperCreateQuote{Date: "2020-06-27", Price: decimal.NewFromFloat(5.99)}

		fundRepo := FundRepositoryMock{}
		fundRepo.FindByISINFn = func(isin string) (Fund, error) {
			return Fund{}, nil
		}

		date, _ := time.Parse("2006-01-02", "2020-06-27")
		quote := Quote{Date: date, Price: decimal.NewFromFloat(5.99)}
		quoteRepo := QuoteRepositoryMock{}
		quoteRepo.CreateFn = func(isin string, q Quote) (Quote, error) {
			return quote, nil
		}

		s := NewQuoteService(fundRepo, quoteRepo)
		got, _ := s.CreateQuote("BE123", want)

		if !reflect.DeepEqual(newPublicQuote(quote), got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	t.Run("return error for time parsing", func(t *testing.T) {
		want := ScraperCreateQuote{Date: "2020-065-274", Price: decimal.NewFromFloat(5.99)}

		fundRepo := FundRepositoryMock{}
		fundRepo.FindByISINFn = func(isin string) (Fund, error) {
			return Fund{}, nil
		}

		quoteRepo := QuoteRepositoryMock{}
		quoteRepo.CreateFn = func(isin string, q Quote) (Quote, error) {
			return Quote{}, nil
		}

		s := NewQuoteService(fundRepo, quoteRepo)
		_, err := s.CreateQuote("BE123", want)

		if err == nil {
			t.Errorf("want error")
		}
	})

	t.Run("return error for fund", func(t *testing.T) {
		fundRepo := FundRepositoryMock{}
		fundRepo.FindByISINFn = func(isin string) (Fund, error) {
			return Fund{}, errors.New("error")
		}

		quoteRepo := QuoteRepositoryMock{}
		quoteRepo.CreateFn = func(isin string, q Quote) (Quote, error) {
			return Quote{}, nil
		}

		s := NewQuoteService(fundRepo, quoteRepo)
		_, err := s.CreateQuote("BE123", ScraperCreateQuote{})

		if err == nil {
			t.Errorf("want error")
		}
	})

	t.Run("return error for quote", func(t *testing.T) {
		want := ScraperCreateQuote{Date: "2020-06-27", Price: decimal.NewFromFloat(5.99)}

		fundRepo := FundRepositoryMock{}
		fundRepo.FindByISINFn = func(isin string) (Fund, error) {
			return Fund{}, nil
		}

		quoteRepo := QuoteRepositoryMock{}
		quoteRepo.CreateFn = func(isin string, q Quote) (Quote, error) {
			return Quote{}, errors.New("error")
		}

		s := NewQuoteService(fundRepo, quoteRepo)
		_, err := s.CreateQuote("BE123", want)

		if err == nil {
			t.Errorf("want error")
		}
	})
}

func TestNewPublicQuote(t *testing.T) {
	t.Run("return correctly formatted PublicQuote", func(t *testing.T) {
		date, _ := time.Parse("2006-01-02", "2020-06-27")
		want := Quote{
			Date:  date,
			Price: decimal.NewFromFloat(5.99),
		}

		got := newPublicQuote(want)

		if want.Date.Format("2006-01-02") != got.Date {
			t.Errorf("want %s, got %s", want.Date.Format("2006-01-02"), got.Date)
		}
		if !want.Price.Equal(decimal.NewFromFloat(got.Price)) {
			t.Errorf("want %s, got %.2f", want.Price, got.Price)
		}
	})
}
