package yaweather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Location represents a pair of latitude and longitude
type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

// IsZero reports whether the coordinates are zero.
func (l Location) IsZero() bool {
	return l.Latitude == 0 && l.Longitude == 0
}

// URL returns the URL of the weather page for the given coordinates
func (l Location) url() string {
	lon := strconv.FormatFloat(l.Longitude, 'f', 6, 64)
	lat := strconv.FormatFloat(l.Latitude, 'f', 6, 64)
	return fmt.Sprintf("https://yandex.ru/pogoda/ru?lat=%s&lon=%s", lat, lon)
}

// GetForecast fetches forecasts
func (l Location) GetForecast(ctx context.Context) ([]Forecast, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, l.url(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint: errcheck
	if resp.StatusCode != http.StatusOK {
		return nil, newErrHTTP(resp.StatusCode)
	}

	return parseForecast(resp.Body)
}

// SearchResult represents single search result
type SearchResult struct {
	Title    string `json:"title"`
	Location `json:"coords"`
}

// Find searches for locations
func Find(ctx context.Context, q string) ([]SearchResult, error) {
	part := url.QueryEscape(q)
	url := fmt.Sprintf("https://yandex.ru/weather/api/suggest?part=%s&lang=en&type=weather", part)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint: errcheck
	if resp.StatusCode != http.StatusOK {
		return nil, newErrHTTP(resp.StatusCode)
	}
	var result []SearchResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, wrapErrInvalidJSON(err)
	}
	return result, nil
}
