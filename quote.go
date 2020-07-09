package pensiondata

import (
	"github.com/shopspring/decimal"
	"time"
)

// Quote is Quote's representation in the database
type Quote struct {
	Date  time.Time
	Price decimal.Decimal
}

type QuoteRepository interface {
	FindByISINAndDate(string, string) (Quote, error)
	FindByDateDesc(string) (Quote, error)
	FindAll(string) ([]Quote, error)
	Create(string, Quote) (Quote, error)
}

type QuoteService interface {
	GetQuote(string, string) (PublicQuote, error)
	GetLatestQuote(string) (PublicQuote, error)
	GetQuotes(string) ([]PublicQuote, error)
	CreateQuote(string, ScraperCreateQuote) (PublicQuote, error)
}

type QuoteServiceImpl struct {
	fundRepo  FundRepository
	quoteRepo QuoteRepository
}

func NewQuoteService(fundRepo FundRepository, quoteRepo QuoteRepository) QuoteService {
	return &QuoteServiceImpl{fundRepo: fundRepo, quoteRepo: quoteRepo}
}

func (s QuoteServiceImpl) GetQuote(isin string, date string) (PublicQuote, error) {
	if _, err := s.fundRepo.FindByISIN(isin); err != nil {
		return PublicQuote{}, err
	}

	quote, err := s.quoteRepo.FindByISINAndDate(isin, date)
	if err != nil {
		return PublicQuote{}, err
	}

	return newPublicQuote(quote), nil
}

func (s QuoteServiceImpl) GetLatestQuote(isin string) (PublicQuote, error) {
	quote, err := s.quoteRepo.FindByDateDesc(isin)

	if err != nil {
		return PublicQuote{}, err
	}

	return newPublicQuote(quote), nil
}

func (s QuoteServiceImpl) GetQuotes(isin string) ([]PublicQuote, error) {
	if _, err := s.fundRepo.FindByISIN(isin); err != nil {
		return []PublicQuote{}, err
	}

	quotes, err := s.quoteRepo.FindAll(isin)
	if err != nil {
		return []PublicQuote{}, err
	}

	var publicQuotes []PublicQuote
	for _, quote := range quotes {
		publicQuotes = append(publicQuotes, newPublicQuote(quote))
	}

	return publicQuotes, nil
}

func (s QuoteServiceImpl) CreateQuote(isin string, scraperQuote ScraperCreateQuote) (PublicQuote, error) {
	if _, err := s.fundRepo.FindByISIN(isin); err != nil {
		return PublicQuote{}, err
	}

	date, err := time.Parse("2006-01-02", scraperQuote.Date)
	if err != nil {
		return PublicQuote{}, err
	}

	quote := Quote{Date: date, Price: scraperQuote.Price}
	createdQuote, err := s.quoteRepo.Create(isin, quote)
	if err != nil {
		return PublicQuote{}, err
	}

	return newPublicQuote(createdQuote), nil
}

// PublicQuote is Quote's representation to be returned by the API
type PublicQuote struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
}

// ScraperCreateQuote is Quote's representation send by the scraper to be created
type ScraperCreateQuote struct {
	Date  string          `json:"date"`
	Price decimal.Decimal `json:"price"`
}

// newPublicQuote return a PublicQuote based on a Quote
func newPublicQuote(quote Quote) PublicQuote {
	price, _ := quote.Price.Float64()
	return PublicQuote{
		Date:  quote.Date.Format("2006-01-02"),
		Price: price,
	}
}
