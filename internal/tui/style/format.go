package style

import (
	"fmt"
	"strconv"
	"strings"
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

func FormatRangeInUnit(raw, unit string) string {
	if len(raw) == 0 {
		return "-"
	}
	val, err := strconv.ParseFloat(raw, 64)
	if err != nil || val < 0 {
		return FormatRange(raw)
	}
	switch unit {
	case "au":
		return fmt.Sprintf("%.4f AU", val/149_597_870.7)
	case "lmin":
		return fmt.Sprintf("%.2f light-min", val/17_987_547.48)
	case "lhour":
		return fmt.Sprintf("%.4f light-hr", val/1_079_252_848.8)
	default:
		return FormatRange(raw)
	}
}

var sparklineBars = []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

func Sparkline(values []float64) string {
	if len(values) < 2 {
		return ""
	}
	min, max := values[0], values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	var b strings.Builder
	for _, v := range values {
		idx := 0
		if max > min {
			idx = int((v - min) / (max - min) * 7)
			if idx > 7 {
				idx = 7
			}
		}
		b.WriteString(sparklineBars[idx])
	}
	return b.String()
}
