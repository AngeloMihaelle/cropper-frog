package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ToSeconds(timeStr string) (float64, error) {
	parts := strings.Split(timeStr, ":")
	var seconds float64 = 0
	var err error

	switch len(parts) {
	case 2: // HH:MM
		h, errH := strconv.ParseFloat(parts[0], 64)
		m, errM := strconv.ParseFloat(parts[1], 64)
		if errH != nil || errM != nil {
			return 0, fmt.Errorf("invalid time format")
		}
		seconds = (h * 3600) + (m * 60)
	case 3: // HH:MM:SS
		h, errH := strconv.ParseFloat(parts[0], 64)
		m, errM := strconv.ParseFloat(parts[1], 64)
		s, errS := strconv.ParseFloat(parts[2], 64)
		if errH != nil || errM != nil || errS != nil {
			return 0, fmt.Errorf("invalid time format")
		}
		seconds = (h * 3600) + (m * 60) + s
	default:
		return 0, fmt.Errorf("invalid time format: must be HH:MM or HH:MM:SS")
	}

	return seconds, err
}
