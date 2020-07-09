package pensiondata

import "errors"

// ErrFundNotFound is returned when a fund was not found
var ErrFundNotFound = errors.New("found not found")

// ErrQuoteNotFound is returned when a fund was not found
var ErrQuoteNotFound = errors.New("quote not found")
