package name_normalizers

import (
	"strings"
)

// NormalizeRealSpaces убирает \t \n \r - как пробельные символы
func NormalizeRealSpaces(name string) (string, bool) {
	s := strings.ReplaceAll(name, "\t", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")

	// нормализуем пробелы
	s = strings.Join(strings.Fields(s), " ")

	return s, s != name && s != ""
}

// NormalizeFakeSpaces убирает \t \n \r - как символы, а не пробелы
func NormalizeFakeSpaces(name string) (string, bool) {
	s := strings.ReplaceAll(name, "\\t", " ")
	s = strings.ReplaceAll(s, "\\n", " ")
	s = strings.ReplaceAll(s, "\\r", " ")

	// нормализуем пробелы
	s = strings.Join(strings.Fields(s), " ")

	return s, s != name && s != ""
}

// NormalizePunctuation заменяет лишние знаки препинания и кавычки
func NormalizePunctuation(name string) (string, bool) {
	s := strings.Map(func(r rune) rune {
		switch r {
		case ',', ';':
			return ' '
		case '«', '»', '\'':
			return '"'
		}
		return r
	}, name)

	// нормализуем пробелы
	s = strings.Join(strings.Fields(s), " ")

	return s, s != name && s != ""
}
