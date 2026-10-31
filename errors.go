package yaweather

import (
	"errors"
	"fmt"
)

// ErrHTTP is returned when the HTTP request fails.
var ErrHTTP = errors.New("http error")

// ErrUnknownCondition is returned when the weather condition for icon is unknown.
var ErrUnknownCondition = errors.New("unknown weather condition")

// ErrInvalidDate is returned when the date is invalid.
var ErrInvalidDate = errors.New("invalid date format")

// ErrInvalidHTML is returned when the HTML is invalid.
var ErrInvalidHTML = errors.New("invalid html")

// ErrUnexpectedHTML is returned when the HTML is unexpected.
var ErrUnexpectedHTML = errors.New("unexpected html")

// ErrForecastsNotFound is returned when the forecasts are not found in the HTML.
var ErrForecastsNotFound = errors.New("forecasts not found")

// ErrUnknownMoonPhase is returned when the moon phase is unknown.
var ErrUnknownMoonPhase = errors.New("unknown moon phase")

// ErrInvalidJSON is returned when the JSON is invalid.
var ErrInvalidJSON = errors.New("invalid json")

func newErrHTTP(code int) error {
	return fmt.Errorf("%w: got code %d", ErrHTTP, code)
}

func newErrUnknownCondition(condition string) error {
	return fmt.Errorf("%w: %s", ErrUnknownCondition, condition)
}

func newErrUnknownMoonPhase(phase string) error {
	return fmt.Errorf("%w: %s", ErrUnknownMoonPhase, phase)
}

func newErrInvalidDate(date string) error {
	return fmt.Errorf("%w: %s", ErrInvalidDate, date)
}

func wrapErrInvalidDate(err error) error {
	return fmt.Errorf("%w: %w", ErrInvalidDate, err)
}

func wrapErrUnexpectedHTML(err error) error {
	return fmt.Errorf("%w: %w", ErrUnexpectedHTML, err)
}

func wrapErrInvalidJSON(err error) error {
	return fmt.Errorf("%w: %w", ErrInvalidHTML, err)
}
