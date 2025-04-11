package util

import "strings"

func JoinNames(names []string) string {
	return strings.Join(names, ", ")
}
