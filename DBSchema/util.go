package DBSchema

import (
	"fmt"
	generatorModels "github.com/lambda-platform/lambda/generator/models"
	"go/format"
	"strconv"
	"strings"
	"unicode"
)

var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
}

var intToWordMap = []string{
	"zero",
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
}

var Debug = false

func Generate(columnTypes []generatorModels.ColumnData, tableName string, structName string, pkgName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool, extraColumns string, extraStucts string) ([]byte, error) {
	var dbTypes string

	dbTypes, _, _ = generateStructTypes(columnTypes, 0, jsonAnnotation, gormAnnotation, gureguTypes)

	importTime := "import (\n\"time\"\n\"github.com/lambda-platform/lambda/DB\") \n var _ = time.Time{}  \n var _ = DB.Date{}  \n"

	src := fmt.Sprintf("package %s\n %s\n \ntype %s %s %s} %s",
		pkgName,
		importTime,
		structName,
		dbTypes,
		extraColumns, extraStucts)
	if gormAnnotation == true {
		tableNameFunc := "" +
			"func (" + strings.ToLower(string(structName[0])) + " *" + structName + ") TableName() string {\n" +
			"	return \"" + tableName + "\"" +
			"}"
		src = fmt.Sprintf("%s\n%s", src, tableNameFunc)
	}
	formatted, err := format.Source([]byte(src))
	if err != nil {
		err = fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
	}
	return formatted, err
}
func FieldName(s string) string {
	return FmtFieldName(strings.ToLower(StringifyFirstChar(s)))
}
func FmtFieldName(s string) string {

	return ToGo(s)

}
func ToGo(name string) string {
	if name == "_" {
		return "_"
	}
	runes := make([]rune, 0, len(name))

	wordWalker(name, wordWalkerFunc(false, &runes))

	return string(runes)
}

func UcFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func LcFirst(s string) string {
	if s == "" {
		return ""
	}

	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
func wordWalkerFunc(private bool, nameRunes *[]rune) func(*wordInfo) {
	return func(info *wordInfo) {
		word := info.Word

		switch {
		case private && info.WordOffset == 0:
			if strings.ToUpper(word) == word || strings.ToLower(word) == word {
				// ID → id, CAMEL → camel
				word = strings.ToLower(info.Word)
			} else {
				// ITicket → iTicket
				word = LcFirst(info.Word)
			}

		case info.MatchCommonInitial:
			word = strings.ToUpper(word)

		case !info.HasCommonInitial && (strings.ToUpper(word) == word || strings.ToLower(word) == word):
			// FOO or foo → Foo
			// FOo → FOo
			word = UcFirst(strings.ToLower(word))
		}

		*nameRunes = append(*nameRunes, []rune(word)...)
	}
}
func isDelimiter(c rune) bool {
	return c == '-' || c == '_' || unicode.IsSpace(c)
}

// This function is based on the following code.
// https://github.com/golang/lint/blob/06c8688daad7faa9da5a0c2f163a3d14aac986ca/lint.go#L679
func wordWalker(str string, f func(*wordInfo)) {
	runes := []rune(strings.TrimFunc(str, isDelimiter))
	w, i, wo := 0, 0, 0 // index of start of word, scan, word offset
	hasCommonInitial := false
	for i+1 <= len(runes) {
		eow := false // whether we hit the end of a word
		switch {
		case i+1 == len(runes):
			eow = true
		case isDelimiter(runes[i+1]):
			// underscore; shift the remainder forward over any run of underscores
			eow = true
			n := 1
			for i+n+1 < len(runes) && isDelimiter(runes[i+n+1]) {
				n++
			}

			// Leave at most one underscore if the underscore is between two digits
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}

			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		case unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]):
			// lower->non-lower
			eow = true
		}
		i++

		// [w,i) is a word.
		word := string(runes[w:i])
		if !eow && commonInitialisms[word] && !unicode.IsLower(runes[i]) {
			// through
			// split IDFoo → ID, Foo
			// but URLs → URLs
		} else if !eow {
			if commonInitialisms[word] {
				hasCommonInitial = true
			}
			continue
		}

		matchCommonInitial := false
		if commonInitialisms[strings.ToUpper(word)] {
			hasCommonInitial = true
			matchCommonInitial = true
		}

		f(&wordInfo{
			WordOffset:         wo,
			Word:               word,
			MatchCommonInitial: matchCommonInitial,
			HasCommonInitial:   hasCommonInitial,
		})
		hasCommonInitial = false
		w = i
		wo++
	}
}

type wordInfo struct {
	WordOffset         int
	Word               string
	MatchCommonInitial bool
	HasCommonInitial   bool
}

func lintFieldName(name string) string {
	// Fast path for simple cases: "_" and all lowercase.
	if name == "_" {
		return name
	}

	for len(name) > 0 && name[0] == '_' {
		name = name[1:]
	}

	allLower := true
	for _, r := range name {
		if !unicode.IsLower(r) {
			allLower = false
			break
		}
	}
	if allLower {
		runes := []rune(name)
		if u := strings.ToUpper(name); commonInitialisms[u] {
			copy(runes[0:], []rune(u))
		} else {
			runes[0] = unicode.ToUpper(runes[0])
		}
		return string(runes)
	}

	// Split camelCase at any lower->upper transition, and split on underscores.
	// Check each word for common initialisms.
	runes := []rune(name)
	w, i := 0, 0 // index of start of word, scan
	for i+1 <= len(runes) {
		eow := false // whether we hit the end of a word

		if i+1 == len(runes) {
			eow = true
		} else if runes[i+1] == '_' {
			// underscore; shift the remainder forward over any run of underscores
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}

			// Leave at most one underscore if the underscore is between two digits
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}

			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		} else if unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]) {
			// lower->non-lower
			eow = true
		}
		i++
		if !eow {
			continue
		}

		// [w,i) is a word.
		word := string(runes[w:i])
		if u := strings.ToUpper(word); commonInitialisms[u] {
			// All the common initialisms are ASCII,
			// so we can replace the bytes exactly.
			copy(runes[w:], []rune(u))

		} else if strings.ToLower(word) == word {
			// already all lowercase, and not the first word, so uppercase the first character.
			runes[w] = unicode.ToUpper(runes[w])
		}
		w = i
	}
	return string(runes)
}

func StringifyFirstChar(str string) string {
	if str != "" {
		first := str[:1]

		i, err := strconv.ParseInt(first, 10, 8)

		if err != nil {
			return str
		}

		return intToWordMap[i] + "_" + str[1:]
	} else {
		return ""
	}

}
