package ui

import (
	"testing"
)

func TestFormatFloat(t *testing.T) {
	var tests = []struct {
		f    float64
		want string
	}{
		{1.23456789, "1.2"},
		{1234567892222, "1234.6B"},
		{123456789, "123.5M"},
		{12345, "12.3K"},
	}

	for _, test := range tests {
		if got := formatFloat(test.f); got != test.want {
			t.Errorf("formatFloat(%f) = %s; want %s", test.f, got, test.want)
		}
	}
}
