package cmd

import (
	"os"
	"strings"
)

// Environ calls os.environ and converts the returned slice into a map.
func Environ() map[string]string {
	env := os.Environ()
	res := make(map[string]string, len(env))
	for _, s := range env {
		ss := strings.Split(s, "=")
		res[ss[0]] = ss[1]
	}
	return res
}
