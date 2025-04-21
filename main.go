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
	sortStrings(paths)
	err := walk(paths, flags)
	if err != nil {
		checkErr(err)
	}
}

func sortStrings(slice []string) {
	n := len(slice)
	for i := 0; i < n; i++ {
		for j := 0; j < n-1-i; j++ {
			if slice[j] > slice[j+1] {
				slice[j], slice[j+1] = slice[j+1], slice[j]
			}
		}
	}
}
