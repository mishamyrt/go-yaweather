package yaweather

// MoonPhase represents moon phase.
type MoonPhase string

//revive:disable
const (
	MoonPhaseUnknown          MoonPhase = ""
	MoonPhaseNewMoon          MoonPhase = "new"
	MoonPhaseWaxingCrescent   MoonPhase = "waxing-crescent"
	MoonPhaseFirstQuarterMoon MoonPhase = "first-quarter"
	MoonPhaseWaxingGibbous    MoonPhase = "waxing-gibbous"
	MoonPhaseFullMoon         MoonPhase = "full"
	MoonPhaseWaningGibbous    MoonPhase = "waning-gibbous"
	MoonPhaseLastQuarterMoon  MoonPhase = "last-quarter"
	MoonPhaseWaningCrescent   MoonPhase = "waning-crescent"
)

//revive:enable

func parseMoonPhase(phase string) (MoonPhase, error) {
	switch phase {
	case "новолуние":
		return MoonPhaseNewMoon, nil
	case "молодая":
		return MoonPhaseWaxingCrescent, nil
	case "первая 1/4":
		return MoonPhaseFirstQuarterMoon, nil
	case "растущая":
		return MoonPhaseWaxingGibbous, nil
	case "полнолуние":
		return MoonPhaseFullMoon, nil
	case "убывающая":
		return MoonPhaseWaningGibbous, nil
	case "последняя 1/4":
		return MoonPhaseLastQuarterMoon, nil
	case "cтарая": // First letter is latin "c"
		return MoonPhaseWaningCrescent, nil
	}
	return MoonPhaseUnknown, newErrUnknownMoonPhase(phase)
}
