package main

import (
	"bytes"
	"strconv"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/bradfitz/iter"
	"github.com/gen0cide/laforge/ent"
	"github.com/iancoleman/strcase"
)

type TempleteContext struct {
	Build              *ent.Build
	Competition        *ent.Competition
	Environment        *ent.Environment
	Host               *ent.Host
	DNS                *ent.DNS
	Network            *ent.Network
	Script             *ent.Script
	Team               *ent.Team
	ProvisionedNetwork *ent.ProvisionedNetwork
	ProvisionedHost    *ent.ProvisionedHost
	ProvisioningStep   *ent.ProvisioningStep
}

// TemplateFuncLib is a standard template library of functions
var TemplateFuncLib = template.FuncMap{
	"hclstring":            QuotedHCLString,
	"N":                    iter.N,
	"UnsafeAtoi":           UnsafeStringAsInt,
	"Decr":                 Decr,
	"ToUpper":              strings.ToUpper,
	"Contains":             strings.Contains,
	"HasPrefix":            strings.HasPrefix,
	"HasSuffix":            strings.HasSuffix,
	"Join":                 strings.Join,
	"Replace":              strings.Replace,
	"Repeat":               strings.Repeat,
	"Split":                strings.Split,
	"Title":                strings.Title,
	"ToLower":              strings.ToLower,
	"ToSnake":              strcase.ToSnake,
	"ToScreamingSnake":     strcase.ToScreamingSnake,
	"ToKebab":              strcase.ToKebab,
	"ToScreamingKebab":     strcase.ToScreamingKebab,
	"ToDelimited":          strcase.ToDelimited,
	"ToScreamingDelimited": strcase.ToScreamingDelimited,
	"ToCamel":              strcase.ToCamel,
	"ToLowerCamel":         strcase.ToLowerCamel,
	"Incr":                 Incr,
}

// UnsafeStringAsInt is a template helper function that will return -1 if it cannot convert the string to an integer.
func UnsafeStringAsInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return i
}

// Decr is a template helper function to non-destructively decrement an integer
func Decr(i int) int {
	return i - 1
}

// Incr is a template helper function to non-destructively increment an integer
func Incr(i int) int {
	return i + 1
}

// QuotedHCLString is a template function to render safe HCLv2 strings
func QuotedHCLString(s string) string {
	e := new(bytes.Buffer)

	//nolint:gosec,errcheck
	e.WriteByte('"')

	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] {
				i++
				continue
			}
			if start < i {

				//nolint:gosec,errcheck
				e.WriteString(s[start:i])

			}
			switch b {
			case '\\', '"':

				//nolint:gosec,errcheck
				e.WriteByte('\\')

				//nolint:gosec,errcheck
				e.WriteByte(b)

			case '\n':

				//nolint:gosec,errcheck
				e.WriteByte('\\')

				//nolint:gosec,errcheck
				e.WriteByte('n')

			case '\r':

				//nolint:gosec,errcheck
				e.WriteByte('\\')

				//nolint:gosec,errcheck
				e.WriteByte('r')

			case '\t':

				//nolint:gosec,errcheck
				e.WriteByte('\\')

				//nolint:gosec,errcheck
				e.WriteByte('t')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				//nolint:gosec,errcheck
				e.WriteString(`\u00`)

				//nolint:gosec,errcheck
				e.WriteByte(hex[b>>4])

				//nolint:gosec,errcheck
				e.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				//nolint:gosec,errcheck
				e.WriteString(s[start:i])
			}
			//nolint:gosec,errcheck
			e.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}

		if c == '\u2028' || c == '\u2029' {
			if start < i {
				//nolint:gosec,errcheck
				e.WriteString(s[start:i])
			}
			//nolint:gosec,errcheck
			e.WriteString(`\u202`)
			//nolint:gosec,errcheck
			e.WriteByte(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		//nolint:gosec,errcheck
		e.WriteString(s[start:])
	}
	//nolint:gosec,errcheck
	e.WriteByte('"')
	return e.String()
}

var hex = "0123456789abcdef"

var htmlSafeSet = [utf8.RuneSelf]bool{
	' ':  true,
	'!':  true,
	'"':  false,
	'#':  true,
	'$':  true,
	'%':  true,
	'&':  false,
	'\'': true,
	'(':  true,
	')':  true,
	'*':  true,
	'+':  true,
	',':  true,
	'-':  true,
	'.':  true,
	'/':  true,
	'0':  true,
	'1':  true,
	'2':  true,
	'3':  true,
	'4':  true,
	'5':  true,
	'6':  true,
	'7':  true,
	'8':  true,
	'9':  true,
	':':  true,
	';':  true,
	'<':  false,
	'=':  true,
	'>':  false,
	'?':  true,
	'@':  true,
	'A':  true,
	'B':  true,
	'C':  true,
	'D':  true,
	'E':  true,
	'F':  true,
	'G':  true,
	'H':  true,
	'I':  true,
	'J':  true,
	'K':  true,
	'L':  true,
	'M':  true,
	'N':  true,
	'O':  true,
	'P':  true,
	'Q':  true,
	'R':  true,
	'S':  true,
	'T':  true,
	'U':  true,
	'V':  true,
	'W':  true,
	'X':  true,
	'Y':  true,
	'Z':  true,
	'[':  true,
	'\\': false,
	']':  true,
	'^':  true,
	'_':  true,
	'`':  true,
	'a':  true,
	'b':  true,
	'c':  true,
	'd':  true,
	'e':  true,
	'f':  true,
	'g':  true,
	'h':  true,
	'i':  true,
	'j':  true,
	'k':  true,
	'l':  true,
	'm':  true,
	'n':  true,
	'o':  true,
	'p':  true,
	'q':  true,
	'r':  true,
	's':  true,
	't':  true,
	'u':  true,
	'v':  true,
	'w':  true,
	'x':  true,
	'y':  true,
	'z':  true,
	'{':  true,
	'|':  true,
	'}':  true,
	'~':  true,
}
