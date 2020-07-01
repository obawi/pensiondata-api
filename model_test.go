package main

import (
	"github.com/shopspring/decimal"
	"testing"
	"time"
)

func TestNewPublicFund(t *testing.T) {
	t.Run("return correctly formatted PublicFund", func(t *testing.T) {
		date, _ := time.Parse("2006-01-02", "2020-06-27")
		fund := Fund{
			Isin:       "BE123",
			Name:       "First Fund",
			Bank:       "Banka",
			LaunchDate: date,
			Currency:   "EUR",
		}

		given := NewPublicFund(fund)

		if fund.Isin != given.Isin {
			t.Errorf("want %s, got %s", fund.Isin, given.LaunchDate)
		}
		if fund.Name != given.Name {
			t.Errorf("want %s, got %s", fund.Name, given.Name)
		}
		if fund.Bank != given.Bank {
			t.Errorf("want %s, got %s", fund.Bank, given.Bank)
		}
		if "2020-06-27" != given.LaunchDate {
			t.Errorf("want %s, got %s", "2020-06-27", given.LaunchDate)
		}
		if fund.Currency != given.Currency {
			t.Errorf("want %s, got %s", fund.Currency, given.Currency)
		}
	})
}

func TestNewPublicQuote(t *testing.T) {
	t.Run("return correctly formatted PublicQuote", func(t *testing.T) {
		date, _ := time.Parse("2006-01-02", "2020-06-27")
		quote := Quote{
			Date:  date,
			Price: decimal.NewFromFloat(5.99),
		}

		given := NewPublicQuote(quote)

		if !quote.Price.Equal(decimal.NewFromFloat(given.Price)) {
			t.Errorf("want %s, got %.2f", quote.Price.String(), given.Price)
		}
		if "2020-06-27" != given.Date {
			t.Errorf("want %s, got %s", "2020-06-27", given.Date)
		}
	})
}
