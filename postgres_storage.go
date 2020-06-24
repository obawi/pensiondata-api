package main

import (
	"database/sql"
	"strings"
)

// PostgresStore is a storage specific struct to handle operations with Postgresql
type PostgresStore struct {
	db *sql.DB
}

// ListFunds return all funds available
func (p PostgresStore) ListFunds() ([]Fund, error) {
	var funds []Fund
	rows, err := p.db.Query("SELECT isin, name, bank, launch_date, currency FROM funds ORDER BY name ASC;")
	if err != nil {
		return []Fund{}, err
	}

	defer rows.Close()

	for rows.Next() {
		fund := Fund{}
		if err := rows.Scan(&fund.Isin, &fund.Name, &fund.Bank, &fund.LaunchDate, &fund.Currency); err != nil {
			return []Fund{}, err
		}

		funds = append(funds, fund)
	}

	if err = rows.Err(); err != nil {
		return []Fund{}, err
	}

	return funds, nil
}

// GetFundByIsin return the fund for the given isin
func (p PostgresStore) GetFundByIsin(isin string) (Fund, error) {
	isin = strings.ToUpper(isin)
	row := p.db.QueryRow("SELECT isin, name, bank, launch_date, currency FROM funds WHERE isin = $1;", isin)

	var fund Fund
	if err := row.Scan(&fund.Isin, &fund.Name, &fund.Bank, &fund.LaunchDate, &fund.Currency); err != nil {
		if err == sql.ErrNoRows {
			return Fund{}, ErrNotFound
		}

		return Fund{}, err
	}

	return fund, nil
}

// ListQuotes return all the quotes available for the given isin
func (p PostgresStore) ListQuotes(isin string) ([]Quote, error) {
	rows, err := p.db.Query("SELECT date, price FROM quotes WHERE fund_isin = $1 ORDER BY date DESC;", isin)
	if err != nil {
		return []Quote{}, err
	}

	defer rows.Close()

	var quotes []Quote
	for rows.Next() {
		var quote Quote
		if err := rows.Scan(&quote.Date, &quote.Price); err != nil {
			return []Quote{}, err
		}

		quotes = append(quotes, quote)
	}

	if err = rows.Err(); err != nil {
		return []Quote{}, err
	}

	return quotes, nil
}

// GetQuoteByDate return the quote for the given date and isin
func (p PostgresStore) GetQuoteByDate(isin, date string) (Quote, error) {
	isin = strings.ToUpper(isin)
	row := p.db.QueryRow("SELECT date, price FROM quotes WHERE fund_isin = $1 AND DATE(date) = $2;", isin, date)

	var quote Quote
	if err := row.Scan(&quote.Date, &quote.Price); err != nil {
		if err == sql.ErrNoRows {
			return Quote{}, ErrNotFound
		}

		return Quote{}, err
	}

	return quote, nil
}

// GetLatestQuote return the latest quote available for the given isin
func (p PostgresStore) GetLatestQuote(isin string) (Quote, error) {
	isin = strings.ToUpper(isin)
	row := p.db.QueryRow("SELECT date, price FROM quotes WHERE fund_isin = $1 ORDER BY date DESC LIMIT 1;", isin)

	var quote Quote
	if err := row.Scan(&quote.Date, &quote.Price); err != nil {
		if err == sql.ErrNoRows {
			return Quote{}, ErrNotFound
		}

		return Quote{}, err
	}

	return quote, nil
}

// CreateQuote store a new Quote and return it
func (p PostgresStore) CreateQuote(isin string, scraperCreateQuote ScraperCreateQuote) (Quote, error) {
	statement, err := p.db.Prepare("INSERT INTO quotes (price, date, fund_isin) VALUES ($1, $2, $3);")
	if err != nil {
		return Quote{}, err
	}

	if _, err := statement.Exec(scraperCreateQuote.Price, scraperCreateQuote.Date, isin); err != nil {
		return Quote{}, err
	}

	quote, err := p.GetQuoteByDate(isin, scraperCreateQuote.Date)
	if err != nil {
		return Quote{}, err
	}

	return quote, nil
}
