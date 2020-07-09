package pensiondata

// FundRepositoryMock for tests
type FundRepositoryMock struct {
	FindByISINFn func(string) (Fund, error)
	FindAllFn    func() ([]Fund, error)
}

// FundServiceMock for tests
type FundServiceMock struct {
	GetFundByISINFn func(string) (PublicFund, error)
	GetFundsFn      func() ([]PublicFund, error)
}

// FindByISIN mock
func (r FundRepositoryMock) FindByISIN(isin string) (Fund, error) {
	return r.FindByISINFn(isin)
}

// FindAll mock
func (r FundRepositoryMock) FindAll() ([]Fund, error) {
	return r.FindAllFn()
}

// GetFundByISIN mock
func (s FundServiceMock) GetFundByISIN(isin string) (PublicFund, error) {
	return s.GetFundByISINFn(isin)
}

// GetFunds mock
func (s FundServiceMock) GetFunds() ([]PublicFund, error) {
	return s.GetFundsFn()
}
