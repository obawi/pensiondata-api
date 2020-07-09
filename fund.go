package pensiondata

import (
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

type FundRepository interface {
	FindByISIN(string) (Fund, error)
	FindAll() ([]Fund, error)
}

// FundService is the use cases for Fund
type FundService interface {
	GetFundByISIN(string) (PublicFund, error)
	GetFunds() ([]PublicFund, error)
}

// FundService is the implementation of FundService
type FundServiceImpl struct {
	repo FundRepository
}

func NewFundService(repo FundRepository) FundService {
	return &FundServiceImpl{repo: repo}
}

func (s FundServiceImpl) GetFundByISIN(isin string) (PublicFund, error) {
	fund, err := s.repo.FindByISIN(isin)
	if err != nil {
		return PublicFund{}, err
	}

	return newPublicFund(fund), nil
}

func (s FundServiceImpl) GetFunds() ([]PublicFund, error) {
	funds, err := s.repo.FindAll()
	if err != nil {
		return []PublicFund{}, err
	}

	var publicFunds []PublicFund
	for _, fund := range funds {
		publicFunds = append(publicFunds, newPublicFund(fund))
	}

	return publicFunds, nil
}

// PublicFund is Fund's representation to be returned by the API
type PublicFund struct {
	Isin       string `json:"isin"`
	Name       string `json:"name"`
	Bank       string `json:"bank"`
	LaunchDate string `json:"launch_date"`
	Currency   string `json:"currency"`
}

// newPublicFund return a PublicFund based on a Fund
func newPublicFund(fund Fund) PublicFund {
	return PublicFund{
		Isin:       fund.Isin,
		Name:       fund.Name,
		Bank:       fund.Bank,
		LaunchDate: fund.LaunchDate.Format("2006-01-02"),
		Currency:   fund.Currency,
	}
}
