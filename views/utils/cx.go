package utils

import "strings"

// CX is a helper function to generate a string of class names based
// on a map of class names and their conditions.
func CX(classes map[string]bool) string {
	b := strings.Builder{}
	for class, condition := range classes {
		if condition {
			b.WriteString(class)
			b.WriteString(" ")
		}
	}

	return b.String()
}
