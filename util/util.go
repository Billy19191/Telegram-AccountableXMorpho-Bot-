package util

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func FormatNumberWithSeparator(number float64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%.2f", number)
}
