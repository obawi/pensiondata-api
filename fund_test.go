package pensiondata

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestGetFundByISIN(t *testing.T) {
	t.Run("return fund successfully", func(t *testing.T) {
		date, _ := time.Parse("2006-02-01", "2020-07-07")
		want := Fund{
			Isin:       "BE123",
			Name:       "First Fund",
			Bank:       "Banka",
			LaunchDate: date,
			Currency:   "EUR",
		}

		r := FundRepositoryMock{}
		r.FindByISINFn = func(isin string) (Fund, error) {
			return want, nil
		}

		fundService := NewFundService(r)
		got, _ := fundService.GetFundByISIN("BE123")

		if !reflect.DeepEqual(newPublicFund(want), got) {
			t.Errorf("want %v, got %v", newPublicFund(want), got)
		}
	})

	t.Run("return error", func(t *testing.T) {
		r := FundRepositoryMock{}
		r.FindByISINFn = func(isin string) (Fund, error) {
			return Fund{}, errors.New("error")
		}

		fundService := NewFundService(r)
		_, err := fundService.GetFundByISIN("BE123")

		if err == nil {
			t.Errorf("want error")
		}
	})
}

func TestGetFunds(t *testing.T) {
	t.Run("return funds successfully", func(t *testing.T) {
		date, _ := time.Parse("2006-02-01", "2020-07-07")
		wants := []Fund{
			{Isin: "BE123", Name: "First Fund", Bank: "Banka", LaunchDate: date, Currency: "EUR"},
			{Isin: "LU123", Name: "Second Fund", Bank: "Banko", LaunchDate: date, Currency: "EUR"},
		}

		r := FundRepositoryMock{}
		r.FindAllFn = func() ([]Fund, error) {
			return wants, nil
		}

		fundService := NewFundService(r)
		got, _ := fundService.GetFunds()

		if len(wants) != len(got) {
			t.Errorf("want %d, got %d", len(wants), len(got))
		}

		var wantPublicFunds []PublicFund
		for _, want := range wants {
			wantPublicFunds = append(wantPublicFunds, newPublicFund(want))
		}

		if !reflect.DeepEqual(wantPublicFunds, got) {
			t.Errorf("want %v, got %v", wantPublicFunds, got)
		}
	})

	t.Run("return error", func(t *testing.T) {
		r := FundRepositoryMock{}
		r.FindAllFn = func() ([]Fund, error) {
			return []Fund{}, errors.New("error")
		}

		fundService := NewFundService(r)
		_, err := fundService.GetFunds()

		if err == nil {
			t.Errorf("want error")
		}
	})
}

func TestNewPublicFund(t *testing.T) {
	t.Run("return correctly formatted PublicFund", func(t *testing.T) {
		date, _ := time.Parse("2006-01-02", "2020-06-27")
		want := Fund{
			Isin:       "BE123",
			Name:       "First Fund",
			Bank:       "Banka",
			LaunchDate: date,
			Currency:   "EUR",
		}

		got := newPublicFund(want)

		if want.Isin != got.Isin {
			t.Errorf("want %s, got %s", want.Isin, got.Isin)
		}
		if want.Name != got.Name {
			t.Errorf("want %s, got %s", want.Name, got.Name)
		}
		if want.Bank != got.Bank {
			t.Errorf("want %s, got %s", want.Bank, got.Bank)
		}
		if want.LaunchDate.Format("2006-01-02") != got.LaunchDate {
			t.Errorf("want %s, got -%s", want.LaunchDate.Format("2006-01-02"), got.LaunchDate)
		}
		if want.Currency != got.Currency {
			t.Errorf("want %s, got %s", want.Currency, got.Currency)
		}
	})
}
