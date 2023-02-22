package ui

import (
	"strconv"
)

// format float64 to string and display value in billions, millions or thousands
func formatFloat(f float64) string {
	if f > 1000000000 {
		return strconv.FormatFloat(f/1000000000, 'f', 1, 64) + "B"
	} else if f > 1000000 {
		return strconv.FormatFloat(f/1000000, 'f', 1, 64) + "M"
	} else if f > 1000 {
		return strconv.FormatFloat(f/1000, 'f', 1, 64) + "K"
	} else {
		return strconv.FormatFloat(f, 'f', 1, 64)
	}
}

// fomat percentage to string
func formatPercent(f float64) string {
	if f > 0 {
		return strconv.FormatFloat(f, 'f', 2, 64) + "%"
	} else {
		return strconv.FormatFloat(f, 'f', 2, 64) + "%"
	}
}
