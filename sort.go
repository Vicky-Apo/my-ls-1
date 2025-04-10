package main

import (
	"sort"
	"strings"
)

// sortEntries sorts the slice of entries in-place according to the flags
func sortEntries(entries []Entry, flags lsFlags) {
	// By default, sort by name (lexicographically)
	sort.Slice(entries, func(i, j int) bool {
		// For convenience, use strings.Compare
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})

	// If sorting by time, we re-sort by ModTime
	if flags.sortByTime {
		sort.Slice(entries, func(i, j int) bool {
			t1 := entries[i].Info.ModTime()
			t2 := entries[j].Info.ModTime()
			// Newest first (like 'ls -t' does by default)
			if t1.Equal(t2) {
				// Tie-break by name to match stable sorting
				return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
			}
			return t1.After(t2)
		})
	}

	// If reverse is on, we simply reverse the slice
	if flags.reverse {
		for i := 0; i < len(entries)/2; i++ {
			opp := len(entries) - 1 - i
			entries[i], entries[opp] = entries[opp], entries[i]
		}
	}
}
