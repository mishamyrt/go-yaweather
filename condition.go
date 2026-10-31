package yaweather

// Condition represents weather condition.
type Condition string

//revive:disable
const (
	UnknownCondition     Condition = ""
	Clear                Condition = "clear"
	Cloudy               Condition = "cloudy"
	RainLight            Condition = "rain-light"
	RainMedium           Condition = "rain-medium"
	RainHeavy            Condition = "rain-heavy"
	RainWithThunderstorm Condition = "rain-with-thunderstorm"
	SnowRain             Condition = "snow-rain"
	SnowLight            Condition = "snow-light"
	SnowMedium           Condition = "snow-medium"
	SnowHeavy            Condition = "snow-heavy"
	SnowWithThunderstorm Condition = "snow-with-thunderstorm"
	Hail                 Condition = "hail"
	HailWithThunderstorm Condition = "hail-with-thunderstorm"
	Thunderstorm         Condition = "thunderstorm"
)

//revive:enable

// parseCondition parses an icon index to weather condition.
func parseCondition(code string) (Condition, error) {
	switch code {
	case "0", "6", "12":
		return RainLight, nil
	case "1", "7", "13":
		return RainMedium, nil
	case "2", "8", "14":
		return RainHeavy, nil
	case "3", "9", "15":
		return SnowLight, nil
	case "4", "10", "16":
		return SnowMedium, nil
	case "5", "11", "17":
		return SnowHeavy, nil
	case "18":
		return SnowRain, nil
	case "19":
		return Hail, nil
	case "20":
		return Thunderstorm, nil
	case "21":
		return SnowWithThunderstorm, nil
	case "22":
		return HailWithThunderstorm, nil
	case "23", "24", "25":
		return Cloudy, nil
	case "26", "27":
		return Clear, nil
	}
	return UnknownCondition, newErrUnknownCondition(code)
}
