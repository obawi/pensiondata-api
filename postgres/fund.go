package postgres

import (
	"database/sql"
	"gihtub.com/obawi/pensiondata-api"
)

// FundRepository is the struct used to implement the pensiondata.FundRepository interface for Postgres
type FundRepository struct {
	DB *sql.DB
}

// NewFundRepository return a new FundRepository for Postgres
func NewFundRepository(db *sql.DB) *FundRepository {
	return &FundRepository{DB: db}
}

// FindByISIN return the fund for the given isin
func (r FundRepository) FindByISIN(isin string) (pensiondata.Fund, error) {
	row := r.DB.QueryRow("SELECT isin, name, bank, launch_date, currency FROM funds WHERE isin = $1;", isin)

	var fund pensiondata.Fund
	if err := row.Scan(&fund.Isin, &fund.Name, &fund.Bank, &fund.LaunchDate, &fund.Currency); err != nil {
		if err == sql.ErrNoRows {
			return pensiondata.Fund{}, pensiondata.ErrFundNotFound
		}
		return pensiondata.Fund{}, err
	}

	return fund, nil
}

// FindAll return all funds
func (r FundRepository) FindAll() ([]pensiondata.Fund, error) {
	var funds []pensiondata.Fund
	rows, err := r.DB.Query("SELECT isin, name, bank, launch_date, currency FROM funds ORDER BY name ASC;")
	if err != nil {
		return []pensiondata.Fund{}, err
	}

	defer rows.Close()

	for rows.Next() {
		fund := pensiondata.Fund{}
		if err := rows.Scan(&fund.Isin, &fund.Name, &fund.Bank, &fund.LaunchDate, &fund.Currency); err != nil {
			return []pensiondata.Fund{}, err
		}
		funds = append(funds, fund)
	}

	if err = rows.Err(); err != nil {
		return []pensiondata.Fund{}, err
	}

	return funds, nil
}
