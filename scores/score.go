// Package score stores the high score separately for native and browser runs.
package score

import (
	"strconv"
	"strings"
)

func parse(data string) int {
	value, err := strconv.Atoi(strings.TrimSpace(data))
	if err != nil || value < 0 {
		return 0
	}
	return value
}
