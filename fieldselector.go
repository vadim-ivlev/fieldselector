// A primitive GraphQL query parser.
// Function GetSelectedFields
// returns a list of fields inside curly braces
// selected to be returned from a GraphQL function.
package fieldselector

import (
	"regexp"
	"strings"
)

// consecutive whitespace characters
var regSpaces, _ = regexp.Compile(`\s+`)

// Returns positions of a parenthesis and of its paired parenthesis.
// - brackets - What brackets to look for.For example "()" or "{}"
// - s - string to search for brackets
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

// Returns the text inside the brackets
// - s - string to search for brackets
func insideBrackets(brackets, s string) string {
	var i0, i1 = findBrackets(brackets, s)
	if i0 == -1 || i1 == -1 {
		return ""
	}
	return s[i0+1 : i1]
}

// Returns the text outside the brackets (cuts out the inner part).
// - s - string to search for brackets
func outsideBrackets(brackets, s string) string {
	var i0, i1 = findBrackets(brackets, s)
	if i0 == -1 || i1 == -1 {
		return s
	}
	return s[:i0] + s[i1+1:]
}

// returns the text inside the curly braces located after a substring
// that matches the regular expression
func selectionAfter(regexpStr, s string) string {
	var re, _ = regexp.Compile(regexpStr)
	var loc = re.FindStringIndex(s)
	if loc == nil {
		return ""
	}
	var ss = s[loc[1]:]
	return insideBrackets("{}", ss)
}

// removes leading, trailing and extra spaces within the string
func removeSpaces(s string) string {
	var ss = regSpaces.ReplaceAllString(s, " ")
	return strings.Trim(ss, " ")
}

// GetSelectedFields returns a list of fields in the GraphQL query
// selected to return from a function inside curly braces.
// - selectionPath - path, including a function name, for which the list is returned.
// May be incomplete or missing.
// - query - GraphQL query text
func GetSelectedFields(selectionPath []string, query string) []string {
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
