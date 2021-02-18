// Примитивный парсер GraphQL запроса.
// Экспортируемая функция GetSelectedFields
// возвращает список полей в тексте запроса GraphQL
// выбранных для возврата из функции внутри фигурных скобок.

package fieldselector

import (
	"fmt"
	"regexp"
	"strings"
)

// идущие подряд пробельные символы
var regSpaces, _ = regexp.Compile(`\s+`)

// Находит скобку и возвращает ее позицию и позицию парной ей скобки.
// - brackets - Какие скобки искать. Например "()" или "{}"
// - s - строка где нужно найти скобки
func findBrackets(brackets, s string) (ind0, ind1 int) {
	ind0, ind1 = -1, -1
	if len(brackets) < 2 {
		return ind0, ind1
	}
	ind0 = strings.Index(s, brackets[:1])
	if ind0 == -1 {
		return ind0, ind1
	}

	n := 0
	for i := ind0; i < len(s); i++ {
		if s[i] == brackets[0] {
			n++
		}
		if s[i] == brackets[1] {
			n--
		}
		if n == 0 {
			ind1 = i
			return
		}
	}
	return
}

// Возвращает текст внутри скобок
// - s - строка где нужно найти скобки
func insideBrackets(brackets, s string) string {
	var i0, i1 = findBrackets(brackets, s)
	if i0 == -1 || i1 == -1 {
		return ""
	}
	return s[i0+1 : i1]
}

// Возвращает текст снаружи скобок (вырезает то что внутри скобок).
// - s - строка где нужно найти скобки
func outsideBrackets(brackets, s string) string {
	var i0, i1 = findBrackets(brackets, s)
	if i0 == -1 || i1 == -1 {
		return s
	}
	return s[:i0] + s[i1+1:]
}

// возвращает текст внутри фигурных скобок найденных после подстроки
// удовлетворяющей регулярному выражению
func selectionAfter(regexpStr, s string) string {
	var re, _ = regexp.Compile(regexpStr)
	var loc = re.FindStringIndex(s)
	if loc == nil {
		return ""
	}
	var ss = s[loc[1]:]
	return insideBrackets("{}", ss)
}

// удаляет ведущие, лидирующие и лишние пробелы внутри строки
func removeSpaces(s string) string {
	var ss = regSpaces.ReplaceAllString(s, " ")
	return strings.Trim(ss, " ")
}

// GetSelectedFields возвращает список полей в тексте запроса GraphQL
// выбранных для возврата из функции внутри фигурных скобок.
// - selectionPath - путь, включая имя функции, для которого возвращается список.
// может быть неполным или с пропусками.
// - query - текст GraphQL запроса
func GetSelectedFields(selectionPath []string, query string) []string {
	fmt.Println("v001 change 4")
	// (?s) - multiline search, .+? not greedy search
	var regexpStr = `(?s)` + strings.Join(selectionPath, `.+?`)

	var ss = selectionAfter(regexpStr, query)
	if ss == "" {
		return []string{}
	}

	// remove all ()
	for o := outsideBrackets("()", ss); o != ss; o = outsideBrackets("()", ss) {
		ss = o
	}
	// remove all {}
	for o := outsideBrackets("{}", ss); o != ss; o = outsideBrackets("{}", ss) {
		ss = o
	}
	ss = removeSpaces(ss)
	return strings.Split(ss, " ")
}
