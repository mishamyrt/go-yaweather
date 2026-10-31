package yaweather

import (
	"strconv"
	"strings"
)

// Wind represents wind data.
type Wind struct {
	// Speed in m/s.
	Speed uint16 `json:"speed"`
	// Direction of wind.
	Direction WindDirection `json:"direction"`
}

// WindDirection represents wind direction.
type WindDirection string

//revive:disable
const (
	WindDirectionNone      WindDirection = ""
	WindDirectionNorth     WindDirection = "n"
	WindDirectionNorthEast WindDirection = "ne"
	WindDirectionEast      WindDirection = "e"
	WindDirectionSouthEast WindDirection = "se"
	WindDirectionSouth     WindDirection = "s"
	WindDirectionSouthWest WindDirection = "sw"
	WindDirectionWest      WindDirection = "w"
	WindDirectionNorthWest WindDirection = "nw"
	WindDirectionCalm      WindDirection = "c"
)

//revive:enable

// ParseWind parses wind data from value and direction strings.
func parseWind(v string, dir string) (Wind, error) {
	speed, err := strconv.ParseUint(strings.Trim(v, " "), 10, 16)
	if err != nil {
		return Wind{}, err
	}
	var direction WindDirection
	if speed != 0 {
		direction = parseWindDirection(strings.Trim(dir, " "))
	} else {
		direction = WindDirectionCalm
	}

	return Wind{
		Speed:     uint16(speed),
		Direction: direction,
	}, nil
}

// parseWindDirection parses wind direction from string.
func parseWindDirection(s string) WindDirection {
	switch s {
	case "Ю":
		return WindDirectionSouth
	case "ЮВ":
		return WindDirectionSouthEast
	case "В":
		return WindDirectionEast
	case "СВ":
		return WindDirectionNorthEast
	case "С":
		return WindDirectionNorth
	case "СЗ":
		return WindDirectionNorthWest
	case "З":
		return WindDirectionWest
	case "ЮЗ":
		return WindDirectionSouthWest
	default:
		return WindDirectionCalm
	}
}
