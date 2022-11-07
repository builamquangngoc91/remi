package database

import (
	"strconv"
	"strings"
)

func GeneratePlaceholders(n int) string {
	if n <= 0 {
		return ""
	}

	var builder strings.Builder
	sep := ", "
	for i := 1; i <= n; i++ {
		if i == n {
			sep = ""
		}
		builder.WriteString("$" + strconv.Itoa(i) + sep)
	}

	return builder.String()
}
