package postgres

import (
	"database/sql"
	"gihtub.com/obawi/pensiondata-api"
)

// QuoteRepository is the struct used to implement the pensiondata.QuoteRepository interface for Postgres
type QuoteRepository struct {
	DB *sql.DB
}

// NewQuoteRepository return a new QuoteRepository for Postgres
func NewQuoteRepository(db *sql.DB) pensiondata.QuoteRepository {
	return &QuoteRepository{DB: db}
}

// FindByISINAndDate return a quote for the given fund isin and date
func (r QuoteRepository) FindByISINAndDate(isin, date string) (pensiondata.Quote, error) {
	row := r.DB.QueryRow("SELECT date, price FROM quotes WHERE fund_isin = $1 AND DATE(date) = $2;", isin, date)

	var quote pensiondata.Quote
	if err := row.Scan(&quote.Date, &quote.Price); err != nil {
		if err == sql.ErrNoRows {
			return pensiondata.Quote{}, pensiondata.ErrQuoteNotFound
		}
		return pensiondata.Quote{}, err
	}

	return quote, nil
}

// FindByDateDesc return a quote for the given isin order by date desc
func (r QuoteRepository) FindByDateDesc(isin string) (pensiondata.Quote, error) {
	row := r.DB.QueryRow("SELECT date, price FROM quotes WHERE fund_isin = $1 ORDER BY date DESC LIMIT 1;", isin)

	var quote pensiondata.Quote
	if err := row.Scan(&quote.Date, &quote.Price); err != nil {
		if err == sql.ErrNoRows {
			return pensiondata.Quote{}, pensiondata.ErrQuoteNotFound
		}
		return pensiondata.Quote{}, err
	}

	return quote, nil
}

// FindAll return all the quotes for the given isin
func (r QuoteRepository) FindAll(isin string) ([]pensiondata.Quote, error) {
	rows, err := r.DB.Query("SELECT date, price FROM quotes WHERE fund_isin = $1 ORDER BY date DESC;", isin)
	if err != nil {
		return []pensiondata.Quote{}, err
	}
	defer rows.Close()

	var quotes []pensiondata.Quote
	for rows.Next() {
		var quote pensiondata.Quote
		if err := rows.Scan(&quote.Date, &quote.Price); err != nil {
			return []pensiondata.Quote{}, err
		}
		quotes = append(quotes, quote)
	}

	if err = rows.Err(); err != nil {
		return []pensiondata.Quote{}, err
	}

	return quotes, nil
}

// Create return the newly created quote
func (r QuoteRepository) Create(isin string, quote pensiondata.Quote) (pensiondata.Quote, error) {
	statement, err := r.DB.Prepare("INSERT INTO quotes (price, date, fund_isin) VALUES ($1, $2, $3);")
	if err != nil {
		return pensiondata.Quote{}, err
	}

	if _, err := statement.Exec(quote.Price, quote.Date, isin); err != nil {
		return pensiondata.Quote{}, err
	}

	createdQuote, err := r.FindByISINAndDate(isin, quote.Date.Format("2006-01-02"))
	if err != nil {
		return pensiondata.Quote{}, err
	}

	return createdQuote, nil
}
