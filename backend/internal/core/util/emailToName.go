package util

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func EmailToName(email string) string {
	parts := strings.Split(email, "@")
	name := cases.Title(language.English, cases.Compact).String(parts[0])
	return name
}
