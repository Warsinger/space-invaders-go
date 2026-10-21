//go:build !js

package score

import (
	"os"
	"strconv"
)

const File = "score/highscore.txt"

func Load() int {
	data, err := os.ReadFile(File)
	if err != nil {
		return 0
	}
	return parse(string(data))
}

func Save(value int) error {
	if err := os.MkdirAll("score", 0755); err != nil {
		return err
	}
	return os.WriteFile(File, []byte(strconv.Itoa(max(0, value))), 0644)
}
