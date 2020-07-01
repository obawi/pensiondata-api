package main

// MockStore used for tests
type MockStore struct {
	ListFundsFn      func() ([]Fund, error)
	GetFundByIsinFn  func(string) (Fund, error)
	ListQuotesFn     func(string) ([]Quote, error)
	GetQuoteByDateFn func(string, string) (Quote, error)
	GetLatestQuoteFn func(string) (Quote, error)
	CreateQuoteFn    func(string, ScraperCreateQuote) (Quote, error)
}

// ListFunds mock
func (s MockStore) ListFunds() ([]Fund, error) {
	return s.ListFundsFn()
}

// GetFundByIsin mock
func (s MockStore) GetFundByIsin(isin string) (Fund, error) {
	return s.GetFundByIsinFn(isin)
}

// ListQuotes mock
func (s MockStore) ListQuotes(isin string) ([]Quote, error) {
	return s.ListQuotesFn(isin)
}

// GetQuoteByDate mock
func (s MockStore) GetQuoteByDate(isin, date string) (Quote, error) {
	return s.GetQuoteByDateFn(isin, date)
}

// GetLatestQuote mock
func (s MockStore) GetLatestQuote(isin string) (Quote, error) {
	return s.GetLatestQuoteFn(isin)
}

// CreateQuote mock
func (s MockStore) CreateQuote(isin string, scraperCreateQuote ScraperCreateQuote) (Quote, error) {
	return s.CreateQuoteFn(isin, scraperCreateQuote)
}
