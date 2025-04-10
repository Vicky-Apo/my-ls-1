package main

import (
	"fmt"
	"os"
)

// A small helper to track errors
func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	flags, paths := parseFlags()
	err := walk(paths, flags)
	if err != nil {
		checkErr(err)
	}
}
