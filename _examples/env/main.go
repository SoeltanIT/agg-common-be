package main

import (
	"os"

	common "github.com/SoeltanIT/agg-common-be"
)

func main() {
	_ = os.Setenv("FOO", "bar")

	common.GetEnv("FOO", "no bar") // bar
	common.GetEnv("BAR", "no bar") // no bar
}
