package utils

import "strings"

// CX is a helper function to generate a string of class names based
// on a map of class names and their conditions.
//
// CX is not guaranteed to be ordered, use CXX for that.
func CX(classes map[string]bool) string {
	b := strings.Builder{}
	for class, condition := range classes {
		if condition {
			b.WriteString(class)
			b.WriteString(" ")
		}
	}

	return strings.TrimSpace(b.String())
}

// CXX takes an alternating string and boolean arguments.
// Odd values must be string and even values must be boolean.
//
// Function panics if above conditions are not fulfilled.
//
// Example:
//
//	utils.CXX("my-0", true, "my-1", cond)
func CXX(classAndConds ...any) string {
	s := strings.Builder{}
	for i, j := 0, 1; j < len(classAndConds); i, j = i+2, j+2 {
		class := classAndConds[i].(string)
		b := classAndConds[j].(bool)
		if b {
			s.WriteString(class)
			s.WriteByte(' ')
		}
	}
	return s.String()
}
