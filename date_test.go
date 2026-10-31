package yaweather

import (
	"testing"
	"time"
)

func TestParseStringDate(t *testing.T) {
	base := time.Date(2026, 12, 30, 0, 0, 0, 0, time.UTC)
	cases := map[string]time.Time{
		"Сегодня, 30 декабря":   base,
		"Завтра, 31 декабря":    base.AddDate(0, 0, 1),
		"Понедельник, 1 января": base.AddDate(0, 0, 2),
		"Вторник, 2 января":     base.AddDate(0, 0, 3),
	}
	for value, expected := range cases {
		t.Run(value, func(t *testing.T) {
			parsed, err := parseStringDate(value, base)
			if err != nil {
				t.Fatal(err)
			}
			if !parsed.Equal(expected) {
				t.Errorf("Expected %v, got %v", expected, parsed)
			}
		})
	}
}
