package yaweather

import (
	"testing"
)

func TestParseWind(t *testing.T) {
	for _, tc := range []struct {
		name, speed, direction string
		want                   Wind
	}{
		{"north", "5", "С", Wind{Speed: 5, Direction: WindDirectionNorth}},
		{"north east", "5", "СВ", Wind{Speed: 5, Direction: WindDirectionNorthEast}},
		{"east", "5", "В", Wind{Speed: 5, Direction: WindDirectionEast}},
		{"south east", "5", "ЮВ", Wind{Speed: 5, Direction: WindDirectionSouthEast}},
		{"south", "5", "Ю", Wind{Speed: 5, Direction: WindDirectionSouth}},
		{"south west", "5", "ЮЗ", Wind{Speed: 5, Direction: WindDirectionSouthWest}},
		{"west", "5", "З", Wind{Speed: 5, Direction: WindDirectionWest}},
		{"north west", "5", "СЗ", Wind{Speed: 5, Direction: WindDirectionNorthWest}},
		{"calm", "0", "Штиль", Wind{Direction: WindDirectionCalm}},
		{"zero overrides direction", "0", "С", Wind{Direction: WindDirectionCalm}},
		{"zero without direction", "0", "", Wind{Direction: WindDirectionCalm}},
		{"spaces", "  12  ", "  ЮЗ  ", Wind{Speed: 12, Direction: WindDirectionSouthWest}},
		{"leading zeros", "005", "С", Wind{Speed: 5, Direction: WindDirectionNorth}},
		{"maximum speed", "65535", "ЮВ", Wind{Speed: 65535, Direction: WindDirectionSouthEast}},
		{"unknown direction", "5", "unknown", Wind{Speed: 5, Direction: WindDirectionCalm}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseWind(tc.speed, tc.direction)
			if err != nil || got != tc.want {
				t.Fatalf("parseWind(%q, %q) = %+v, %v; want %+v", tc.speed, tc.direction, got, err, tc.want)
			}
		})
	}
}
