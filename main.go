package main

import (
	"os"
)

func read(fname string) (string, error) {
	contents, err := os.ReadFile(fname)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}
