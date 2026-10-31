package yaweather

import (
	"strconv"
	"strings"
	"time"
)

// parseMonthWord parses a month word in the format "января"
func parseMonthWord(w string) int {
	switch w {
	case "января":
		return 1
	case "февраля":
		return 2
	case "марта":
		return 3
	case "апреля":
		return 4
	case "мая":
		return 5
	case "июня":
		return 6
	case "июля":
		return 7
	case "августа":
		return 8
	case "сентября":
		return 9
	case "октября":
		return 10
	case "ноября":
		return 11
	case "декабря":
		return 12
	default:
		return 0
	}
}

// parseStringDate parses a date string in the format "Завтра, 1 ноября" or "Понедельник, 1 ноября"
func parseStringDate(d string, base time.Time) (time.Time, error) {
	_, date, hasComma := strings.Cut(strings.Trim(d, " "), ", ")
	if !hasComma {
		return time.Time{}, newErrInvalidDate("missing parts")
	}
	dayNumber, monthWord, hasSpace := strings.Cut(date, " ")
	if !hasSpace {
		return time.Time{}, newErrInvalidDate("missing nbsp")
	}
	month := parseMonthWord(monthWord)
	if month == 0 {
		return time.Time{}, newErrInvalidDate("invalid month")
	}
	day, err := strconv.Atoi(dayNumber)
	if err != nil {
		return time.Time{}, newErrInvalidDate("invalid day value")
	}

	result := time.Date(base.Year(), time.Month(month), day, 0, 0, 0, 0, base.Location())
	if base.Month() > time.Month(month) {
		result = result.AddDate(1, 0, 0)
	}

	return result, nil
}

// timeOnDate returns UTC time on the given date at the given clock time.
func timeOnDate(date time.Time, clock string, tzOffset int) (time.Time, error) {
	hours, minutes, hasColon := strings.Cut(clock, ":")
	if !hasColon {
		return time.Time{}, newErrInvalidDate("missing colon")
	}
	h, err := strconv.Atoi(hours)
	if err != nil || h > 23 {
		return time.Time{}, newErrInvalidDate("invalid hours value")
	}
	m, err := strconv.Atoi(minutes)
	if err != nil || m > 59 {
		return time.Time{}, newErrInvalidDate("invalid minutes value")
	}
	h -= tzOffset
	if h < 0 {
		h += 24
		date = date.AddDate(0, 0, -1)
	}
	return time.Date(date.Year(), date.Month(), date.Day(), h, m, 0, 0, time.UTC), nil
}
