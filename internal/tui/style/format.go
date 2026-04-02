package style

import (
	"fmt"
	"strconv"
)

func DefaultIfEmpty(value, defaultValue string) string {
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func DashIfEmpty(value string) string {
	if len(value) == 0 {
		return "-"
	}
	return value
}

func FormatDataRate(raw string) string {
	if len(raw) == 0 {
		return "-"
	}
	val, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return raw
	}
	if val < 0 {
		return raw
	}
	if val > 1000000 {
		return fmt.Sprintf("%.2f Mb/s", val/1000000)
	}
	if val > 1000 {
		return fmt.Sprintf("%.2f kb/s", val/1000)
	}
	return fmt.Sprintf("%.2f b/s", val)
}

func FormatRange(raw string) string {
	if len(raw) == 0 {
		return "-"
	}
	val, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return raw
	}
	if val < 0 {
		return raw
	}
	if val > 1000000000 {
		return fmt.Sprintf("%.2f billion km", val/1000000000)
	}
	if val > 1000000 {
		return fmt.Sprintf("%.2f million km", val/1000000)
	}
	if val > 1000 {
		return fmt.Sprintf("%.2f thousand km", val/1000)
	}
	return fmt.Sprintf("%.2f km", val)
}

func FormatRTLT(raw string) string {
	if len(raw) == 0 {
		return "-"
	}
	val, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return raw
	}
	if val < 0 {
		return raw
	}
	if val > 3600 {
		return fmt.Sprintf("%.2f hours", val/3600)
	}
	if val > 60 {
		return fmt.Sprintf("%.2f minutes", val/60)
	}
	return fmt.Sprintf("%.2f seconds", val)
}

func FormatAngle(raw string) string {
	if len(raw) == 0 {
		return "-"
	}
	return raw + "°"
}

func FormatWind(raw string) string {
	if len(raw) == 0 {
		return "-"
	}
	return raw + " km/h"
}

func FormatPowerTx(raw string) string {
	if len(raw) == 0 {
		return "-"
	}
	return raw + " kW"
}

func FormatPowerRx(raw string) string {
	if len(raw) == 0 {
		return "-"
	}
	return raw + " dBm"
}
