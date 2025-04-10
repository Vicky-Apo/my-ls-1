package main

import (
	"fmt"
	"os"
)

// Flags structure
type lsFlags struct {
	showAll     bool // -a
	longListing bool // -l
	recursive   bool // -R
	reverse     bool // -r
	sortByTime  bool // -t
}

// parseFlags parses our required flags (-l, -R, -a, -r, -t) from os.Args.
// This function also handles combined flags like '-la', '-art', etc.
func parseFlags() (lsFlags, []string) {
	flags := lsFlags{}
	// We will manually parse arguments so that we allow bundled flags (e.g. '-la')
	// Alternatively, you can use the standard library `flag` in a more advanced way.
	args := os.Args[1:]
	var paths []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if len(arg) > 1 && arg[0] == '-' {
			// It's a flag set, e.g., "-laR"
			for _, ch := range arg[1:] {
				switch ch {
				case 'l':
					flags.longListing = true
				case 'R':
					flags.recursive = true
				case 'a':
					flags.showAll = true
				case 'r':
					flags.reverse = true
				case 't':
					flags.sortByTime = true
				default:
					// If you have extra flags, handle them or ignore them
					fmt.Fprintf(os.Stderr, "Warning: ignoring unknown flag '-%c'\n", ch)
				}
			}
		} else {
			// It's presumably a path
			paths = append(paths, arg)
		}
	}

	// If no paths provided, default to "."
	if len(paths) == 0 {
		paths = []string{"."}
	}

	return flags, paths
}
