package utils

import (
	"bytes"
	"strings"

	"github.com/jinzhu/inflection"
)

var commonInitialisms = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
var commonInitialismsReplacer *strings.Replacer

func init() {
	var commonInitialismsForReplacer []string
	for _, initialism := range commonInitialisms {
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, strings.Title(strings.ToLower(initialism)))
	}
	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
}

func ToLowerUnderlinedNamer(name string) string {
	const (
		lower = false
		upper = true
	)

	if name == "" {
		return ""
	}

	var (
		value                                    = commonInitialismsReplacer.Replace(name)
		buf                                      = bytes.NewBufferString("")
		lastCase, currCase, nextCase, nextNumber bool
	)

	for i, v := range value[:len(value)-1] {
		nextCase = bool(value[i+1] >= 'A' && value[i+1] <= 'Z')
		nextNumber = bool(value[i+1] >= '0' && value[i+1] <= '9')

		if i > 0 {
			if currCase == upper {
				if lastCase == upper && (nextCase == upper || nextNumber == upper) {
					buf.WriteRune(v)
				} else {
					if value[i-1] != '_' && value[i+1] != '_' {
						buf.WriteRune('_')
					}
					buf.WriteRune(v)
				}
			} else {
				buf.WriteRune(v)
				if i == len(value)-2 && (nextCase == upper && nextNumber == lower) {
					buf.WriteRune('_')
				}
			}
		} else {
			currCase = upper
			buf.WriteRune(v)
		}
		lastCase = currCase
		currCase = nextCase
	}

	buf.WriteByte(value[len(value)-1])

	s := strings.ToLower(buf.String())
	return s
}

func ToPlural(v string) string {
	return inflection.Plural(v)
}

func toLowerPlural(v, sep string) string {
	ss := strings.Split(ToLowerUnderlinedNamer(v), "_")
	if len(ss) > 0 {
		ss[len(ss)-1] = ToPlural(ss[len(ss)-1])
	}
	return strings.Join(ss, sep)
}

func ToLowerPlural(v string) string {
	return toLowerPlural(v, "")
}

func ToLowerSpacePlural(v string) string {
	return toLowerPlural(v, " ")
}

func ToLowerHyphensPlural(v string) string {
	return toLowerPlural(v, "-")
}

func ToLowerCamel(v string) string {
	if v == "" {
		return ""
	}
	return strings.ToLower(v[:1]) + v[1:]
}

func ToLowerSpacedNamer(v string) string {
	return strings.Replace(ToLowerUnderlinedNamer(v), "_", " ", -1)
}

func ToTitleSpaceNamer(v string) string {
	vv := strings.Split(ToLowerUnderlinedNamer(v), "_")
	if len(vv) > 0 && len(vv[0]) > 0 {
		vv[0] = strings.ToUpper(vv[0][:1]) + vv[0][1:]
	}
	return strings.Join(vv, " ")
}
