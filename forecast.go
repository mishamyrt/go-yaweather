package yaweather

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Forecast represents the forecast for a day.
type Forecast struct {
	// Date is the date of the forecast, in the local time zone.
	Date time.Time `json:"date"`
	// Week is the day of the week, 0 (Sunday) to 6 (Saturday).
	Week int `json:"week"`
	// MoonPhase is the phase of the moon, "убывающая", "старая", "растущая",
	MoonPhase MoonPhase `json:"moon_phase"`
	// Sunrise time in the local time zone.
	Sunrise time.Time `json:"sunrise"`
	// Sunset time in the local time zone.
	Sunset time.Time `json:"sunset"`
	// Forecasts for the morning.
	Morning ForecastPart `json:"morning"`
	// Forecasts for the day.
	Day ForecastPart `json:"day"`
	// Forecasts for the evening.
	Evening ForecastPart `json:"evening"`
	// Forecasts for the night.
	Night ForecastPart `json:"night"`
}

// ForecastPart is the forecast for a part of the day.
type ForecastPart struct {
	// Temperature is the fact temperature in the degrees Celsius.
	Temperature int16 `json:"temperature"`
	// FeelsLike is the feels-like temperature in the degrees Celsius.
	FeelsLike int16 `json:"feels_like"`
	// Wind speed and direction.
	Wind Wind `json:"wind"`
	// Humidity is the relative humidity in percent.
	Humidity float32 `json:"humidity"`
	// Pressure is the atmospheric pressure in mmHg.
	Pressure int16 `json:"pressure"`
	// Condition is the weather condition.
	Condition Condition `json:"condition"`
}

func (f Forecast) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Morning: %s\n", f.Morning)
	fmt.Fprintf(&builder, "Day: %s\n", f.Day)
	fmt.Fprintf(&builder, "Evening: %s\n", f.Evening)
	fmt.Fprintf(&builder, "Night: %s", f.Night)
	return builder.String()
}

func (f ForecastPart) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%d°C", f.Temperature)
	if f.FeelsLike != f.Temperature {
		fmt.Fprintf(&builder, " (feels like %d°C)", f.FeelsLike)
	}
	fmt.Fprintf(&builder, ", %s", f.Condition)
	fmt.Fprintf(&builder, ", %d m/s %s", f.Wind.Speed, f.Wind.Direction)
	fmt.Fprintf(&builder, ", humidity %d%%", int(f.Humidity*100))
	fmt.Fprintf(&builder, ", pressure %d mmHg", f.Pressure)
	return builder.String()
}

// Parse returns daily forecasts in page order, ignoring ads and duplicate cards.
func parseForecast(r io.Reader) ([]Forecast, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidHTML, err)
	}

	var days []Forecast
	var now time.Time
	seen := make(map[string]struct{})
	for n := range doc.Descendants() {
		if n.Type != html.ElementNode {
			continue
		}
		if n.Data == "time" && now.IsZero() {
			dt, err := time.Parse(time.RFC3339, attr(n, "datetime"))
			if err != nil {
				return nil, wrapErrInvalidDate(err)
			}
			now = dt
		}
		if n.Data != "li" || n.Parent.Data != "ul" {
			continue
		}
		if attr(n.Parent, "data-forecast") != "true" {
			continue
		}
		article := child(child(n, "a"), "article")
		if article == nil {
			continue
		}
		id := attr(article, "data-day")
		if _, isSeen := seen[id]; isSeen {
			continue
		}
		day, err := parseDay(article, now)
		if err != nil {
			return nil, wrapErrUnexpectedHTML(err)
		}
		days = append(days, day)
		seen[id] = struct{}{}
	}
	if len(days) == 0 {
		return nil, ErrForecastsNotFound
	}
	return days, nil
}

// parseDay parses a day forecast from the article node.
func parseDay(article *html.Node, baseDate time.Time) (Forecast, error) {
	fields := make(map[string]string)
	var day Forecast
	for n := range article.Descendants() {
		if area := style(n, "grid-area"); area != "" {
			if _, exists := fields[area]; !exists {
				fields[area] = text(n)
				if strings.HasSuffix(area, "-icon") {
					fields[area] = style(n, "--icon")
				}
			}
		}
		if !classPrefix(n, "AppForecastDayDuration_item__") {
			continue
		}
		if day.Date.IsZero() {
			dt, err := parseStringDate(fields["day-title"], baseDate)
			if err != nil {
				return Forecast{}, err
			}
			day.Date = dt
		}
		caption := text(findClass(n, "AppForecastDayDuration_caption__"))
		value := text(findClass(n, "AppForecastDayDuration_value__"))
		switch caption {
		case "Фаза луны":
			mp, err := parseMoonPhase(value)
			if err != nil {
				return day, err
			}
			day.MoonPhase = mp
		case "":
			var times []string
			for item := range n.Descendants() {
				if classPrefix(item, "AppForecastDayDuration_text__") {
					times = append(times, text(item))
				}
			}
			if len(times) == 3 {
				var err error
				day.Sunrise, err = timeOnDate(day.Date, times[0], 4)
				if err != nil {
					return day, fmt.Errorf("sunrise: %w", err)
				}
				day.Sunset, err = timeOnDate(day.Date, times[2], 4)
				if err != nil {
					return day, fmt.Errorf("sunset: %w", err)
				}
			}
		}
	}
	for _, part := range []struct {
		prefix string
		value  *ForecastPart
	}{{"m", &day.Morning}, {"d", &day.Day}, {"e", &day.Evening}, {"n", &day.Night}} {
		if err := parsePart(fields, part.prefix, part.value); err != nil {
			return day, err
		}
	}
	if day.MoonPhase == "" || day.Sunrise.IsZero() || day.Sunset.IsZero() {
		return day, fmt.Errorf("missing moon phase or daylight data")
	}
	return day, nil
}

// parsePart parses a part of the forecast.
func parsePart(fields map[string]string, prefix string, part *ForecastPart) error {
	for _, field := range []struct {
		name  string
		value *int16
	}{{"temp", &part.Temperature}, {"feels", &part.FeelsLike}} {
		key := prefix + "-" + field.name
		value, err := parseTemperature(fields[key])
		if err != nil {
			return fmt.Errorf("%s: %w", key, err)
		}
		*field.value = value
	}
	humidityRaw := strings.TrimSuffix(fields[prefix+"-hum"], "%")
	humidity, err := strconv.ParseFloat(humidityRaw, 32)
	part.Humidity = float32(humidity / 100)
	if err != nil {
		return fmt.Errorf("%s-hum: %w", prefix, err)
	}
	pressure, err := strconv.Atoi(fields[prefix+"-press"])
	if err != nil {
		return fmt.Errorf("%s-press: %w", prefix, err)
	}
	part.Pressure = int16(pressure)
	part.Condition, _ = parseCondition(fields[prefix+"-icon"])
	wind, err := parseWind(strings.TrimSuffix(fields[prefix+"-wind"], " м/с"), fields[prefix+"-dir"])
	if err != nil {
		return fmt.Errorf("%s-wind: %w", prefix, err)
	}
	part.Wind = wind
	if fields[prefix+"-dir"] == "" || part.Condition == "" ||
		part.Humidity < 0 || part.Humidity > 100 || part.Pressure <= 0 {
		return fmt.Errorf("%s: missing or invalid weather fields", prefix)
	}
	return nil
}

func parseTemperature(raw string) (int16, error) {
	value := strings.Fields(strings.NewReplacer("−", "-", ",", ".").Replace(raw))
	if len(value) == 0 {
		return 0, fmt.Errorf("missing number")
	}
	n, err := strconv.ParseInt(strings.TrimSuffix(value[0], "°"), 10, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid number %q", raw)
	}
	return int16(n), nil
}

func attr(n *html.Node, name string) string {
	for _, a := range n.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

func style(n *html.Node, name string) string {
	for declaration := range strings.SplitSeq(attr(n, "style"), ";") {
		key, value, ok := strings.Cut(declaration, ":")
		if ok && strings.TrimSpace(key) == name {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func child(n *html.Node, tag string) *html.Node {
	if n != nil {
		for c := range n.ChildNodes() {
			if c.Type == html.ElementNode && c.Data == tag {
				return c
			}
		}
	}
	return nil
}

func classPrefix(n *html.Node, prefix string) bool {
	for _, class := range strings.Fields(attr(n, "class")) {
		if strings.HasPrefix(class, prefix) {
			return true
		}
	}
	return false
}

func findClass(n *html.Node, prefix string) *html.Node {
	for c := range n.Descendants() {
		if classPrefix(c, prefix) {
			return c
		}
	}
	return nil
}

func text(n *html.Node) string {
	return strings.Join(strings.Fields(rawText(n)), " ")
}

func rawText(n *html.Node) string {
	if n == nil {
		return ""
	}
	var out strings.Builder
	for c := range n.Descendants() {
		if c.Type == html.TextNode {
			out.WriteString(c.Data)
		} else if c.Type == html.ElementNode && c.Data == "br" {
			out.WriteByte('\n')
		}
	}
	return out.String()
}
