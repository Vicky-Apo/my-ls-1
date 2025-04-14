package main

import (
	"sort"
	"strings"
)

func sortEntries(entries []Entry, flags lsFlags) {
	// By default, sort by name
	sort.Slice(entries, func(i, j int) bool {
		cleanName := func(name string) string {
			if strings.HasPrefix(name, ".") {
				return strings.ToLower(name[1:])
			}
			return strings.ToLower(name)
		}

		return cleanName(entries[i].Name) < cleanName(entries[j].Name)
	})

	// If sorting by time, override the sort
	if flags.sortByTime {
		sort.SliceStable(entries, func(i, j int) bool {

			t1 := entries[i].Info.ModTime()
			t2 := entries[j].Info.ModTime()

			if t1.Equal(t2) {
				return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
			}
			return t1.After(t2)
		})
	}

	// Reverse if needed
	if flags.reverse {
		for i := 0; i < len(entries)/2; i++ {
			j := len(entries) - 1 - i
			entries[i], entries[j] = entries[j], entries[i]
		}
	}
}
