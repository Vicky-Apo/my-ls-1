package main

import "strings"

func sortEntries(entries []Entry, flags lsFlags) {
	// Sort by name by default (alphabetical, case-insensitive, ignoring dot prefix)
	for i := 0; i < len(entries); i++ {
		for j := 0; j < len(entries)-1-i; j++ {
			ni := strings.ToLower(strings.TrimPrefix(entries[j].Name, "."))
			nj := strings.ToLower(strings.TrimPrefix(entries[j+1].Name, "."))
			if ni > nj {
				entries[j], entries[j+1] = entries[j+1], entries[j]
			}
		}
	}

	// If sorting by time
	if flags.sortByTime {
		for i := 0; i < len(entries); i++ {
			for j := 0; j < len(entries)-1-i; j++ {
				t1 := entries[j].Info.ModTime()
				t2 := entries[j+1].Info.ModTime()

				if t1.Before(t2) {
					entries[j], entries[j+1] = entries[j+1], entries[j]
				} else if t1.Equal(t2) {
					n1 := strings.ToLower(entries[j].Name)
					n2 := strings.ToLower(entries[j+1].Name)
					if n1 > n2 {
						entries[j], entries[j+1] = entries[j+1], entries[j]
					}
				}
			}
		}
	}

	// Reverse if needed
	if flags.reverse {
		for i := 0; i < len(entries)/2; i++ {
			j := len(entries) - 1 - i
			entries[i], entries[j] = entries[j], entries[i]
		}
	}
}
