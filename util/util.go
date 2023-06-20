package util

import (
	"math"
	"strings"
)

func MathRoundPrecision(value float64, precision int) float64 {
	return math.Round(value*(math.Pow10(precision))) / math.Pow10(precision)
}

func FormatTitle(value string) string {
	return strings.Title(strings.Join(strings.Fields(value), " "))
}
