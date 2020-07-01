package main

import (
	"github.com/shopspring/decimal"
	"time"
)

// Fund is Fund's representation in the database
type Fund struct {
	Isin       string
	Name       string
	Bank       string
	LaunchDate time.Time
	Currency   string
}

// Quote is Quote's representation in the database
type Quote struct {
	Date  time.Time
	Price decimal.Decimal
}

// PublicFund is Fund's representation to be returned by the API
type PublicFund struct {
	Isin       string `json:"isin"`
	Name       string `json:"name"`
	Bank       string `json:"bank"`
	LaunchDate string `json:"launch_date"`
	Currency   string `json:"currency"`
}

// NewPublicFund return a PublicFund based on a Fund
func NewPublicFund(fund Fund) PublicFund {
	return PublicFund{
		Isin:       fund.Isin,
		Name:       fund.Name,
		Bank:       fund.Bank,
		LaunchDate: fund.LaunchDate.Format("2006-01-02"),
		Currency:   fund.Currency,
	}
}

// PublicQuote is Quote's representation to be returned by the API
type PublicQuote struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
}

// NewPublicQuote return a PublicQuote based on a Quote
func NewPublicQuote(quote Quote) PublicQuote {
	price, _ := quote.Price.Float64()
	return PublicQuote{
		Date:  quote.Date.Format("2006-01-02"),
		Price: price,
	}
}

// ScraperCreateQuote is Quote's representation send by the scraper to be created
type ScraperCreateQuote struct {
	Date  string          `json:"date"`
	Price decimal.Decimal `json:"price"`
}
