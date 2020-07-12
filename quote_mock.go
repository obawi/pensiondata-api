package pensiondata

// QuoteRepositoryMock used for tests
type QuoteRepositoryMock struct {
	FindByISINAndDateFn func(string, string) (Quote, error)
	FindByDateDescFn    func(string) (Quote, error)
	FindAllFn           func(string) ([]Quote, error)
	CreateFn            func(string, Quote) (Quote, error)
}

// QuoteServiceMock used for tests
type QuoteServiceMock struct {
	GetQuoteFn       func(string, string) (PublicQuote, error)
	GetLatestQuoteFn func(string) (PublicQuote, error)
	GetQuotesFn      func(string) ([]PublicQuote, error)
	CreateQuoteFn    func(string, ScraperCreateQuote) (PublicQuote, error)
}

// FindByISINAndDate mock
func (q QuoteRepositoryMock) FindByISINAndDate(isin, date string) (Quote, error) {
	return q.FindByISINAndDateFn(isin, date)
}

// FindByDateDesc mock
func (q QuoteRepositoryMock) FindByDateDesc(isin string) (Quote, error) {
	return q.FindByDateDescFn(isin)
}

// FindAll mock
func (q QuoteRepositoryMock) FindAll(isin string) ([]Quote, error) {
	return q.FindAllFn(isin)
}

// Create mock
func (q QuoteRepositoryMock) Create(isin string, quote Quote) (Quote, error) {
	return q.CreateFn(isin, quote)
}

// GetQuote mock
func (s QuoteServiceMock) GetQuote(isin, date string) (PublicQuote, error) {
	return s.GetQuoteFn(isin, date)
}

// GetLatestQuote mock
func (s QuoteServiceMock) GetLatestQuote(isin string) (PublicQuote, error) {
	return s.GetLatestQuoteFn(isin)
}

// GetQuotes mock
func (s QuoteServiceMock) GetQuotes(isin string) ([]PublicQuote, error) {
	return s.GetQuotesFn(isin)
}

// GetQuote mock
func (s QuoteServiceMock) CreateQuote(isin string, scraperCreateQuote ScraperCreateQuote) (PublicQuote, error) {
	return s.CreateQuoteFn(isin, scraperCreateQuote)
}
