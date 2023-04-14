package utils

import (
	"html/template"
	"strings"
)

// FuncMap is a map of functions that can be used in templates.
var FuncMap = template.FuncMap{
	"lower":           strings.ToLower,
	"upper":           strings.ToUpper,
	"title":           strings.ToTitle,
	"lowerUnderline":  ToLowerUnderlinedNamer,
	"plural":          ToPlural,
	"lowerPlural":     ToLowerPlural,
	"lowerCamel":      ToLowerCamel,
	"lowerSpace":      ToLowerSpacedNamer,
	"titleSpace":      ToTitleSpaceNamer,
	"convIfCond":      tplConvToIfCond,
	"convSwaggerType": tplConvToSwaggerType,
}

func tplConvToIfCond(t string) string {
	if strings.HasPrefix(t, "*") {
		return `v != nil`
	} else if t == "string" {
		return `v != ""`
	} else if strings.Contains(t, "int") {
		return `v != 0`
	} else if strings.Contains(t, "float") {
		return `v != 0`
	} else if t == "time.Time" {
		return `!v.IsZero()`
	} else {
		return `v != nil`
	}
}

func tplConvToSwaggerType(t string) string {
	if strings.Contains(t, "int") || strings.Contains(t, "float") {
		return `number`
	}
	return `string`
}
