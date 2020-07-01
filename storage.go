package main

import "errors"

// ErrNotFound is returned when a resource was not found
var ErrNotFound = errors.New("resource not found")

// Storage is an interface providing the available operations on a storage
type Storage interface {
	ListFunds() ([]Fund, error)
	GetFundByIsin(isin string) (Fund, error)
	ListQuotes(string) ([]Quote, error)
	GetQuoteByDate(string, string) (Quote, error)
	GetLatestQuote(string) (Quote, error)

	CreateQuote(string, ScraperCreateQuote) (Quote, error)
}
