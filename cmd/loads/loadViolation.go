package main

import "github.com/gl1n0m3c/IT_LAB_INIT/internal/exloads/exfuncs"

func main() {
	err := exfuncs.LoadViolation()
	if err != nil {
		panic(err)
	}
}
