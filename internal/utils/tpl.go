package utils

import (
	"html/template"
	"strings"
)

// FuncMap is a map of functions that can be used in templates.
var FuncMap = template.FuncMap{
	"lower":            strings.ToLower,
	"upper":            strings.ToUpper,
	"title":            strings.ToTitle,
	"lowerUnderline":   ToLowerUnderlinedNamer,
	"plural":           ToPlural,
	"lowerPlural":      ToLowerPlural,
	"lowerSpacePlural": ToLowerSpacePlural,
	"lowerCamel":       ToLowerCamel,
	"lowerSpace":       ToLowerSpacedNamer,
	"titleSpace":       ToTitleSpaceNamer,
	"convIfCond":       tplConvToIfCond,
	"convSwaggerType":  tplConvToSwaggerType,
	"raw":              func(s string) template.HTML { return template.HTML(s) },
}

func tplConvToIfCond(t string) template.HTML {
	cond := `v != nil`
	if strings.HasPrefix(t, "*") {
		cond = `v != nil`
	} else if t == "string" {
		cond = `len(v) > 0`
	} else if strings.Contains(t, "int") {
		cond = `v != 0`
	} else if strings.Contains(t, "float") {
		cond = `v != 0`
	} else if t == "time.Time" {
		cond = `!v.IsZero()`
	}
	return template.HTML(cond)
}

func tplConvToSwaggerType(t string) string {
	if strings.Contains(t, "int") || strings.Contains(t, "float") {
		return "number"
	}
	return "string"
}
