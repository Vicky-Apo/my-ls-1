package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Entry holds the file information we need for printing
type Entry struct {
	Name       string
	FullPath   string
	Info       fs.FileInfo
	LinkTarget string // if symlink, store the target here
}

// walk handles the logic of listing directories and optionally recursing
func walk(paths []string, flags lsFlags) error {
	for idx, path := range paths {
		// If multiple paths, print a header like real ls does
		if len(paths) > 1 {
			if idx > 0 {
				fmt.Println()
			}
			fmt.Printf("%s:\n", path)
		}

		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("cannot access '%s': %v", path, err)
		}

		// If path is directory, list it. If path is a symlink to directory and user appended '/', we also list it.
		// If it's a file, just show the file.
		if info.IsDir() {
			// list directory
			entries, err := listDirectory(path, flags)
			if err != nil {
				return err
			}
			printEntries(path, entries, flags)

			// If recursive, we traverse subdirectories
			if flags.recursive {
				// For each subdir, do the same
				for _, e := range entries {
					if e.Info.IsDir() && e.Name != "." && e.Name != ".." {
						// Recurse
						subPaths := []string{e.FullPath}
						err := recurseDir(subPaths, flags)
						if err != nil {
							return err
						}
					}
				}
			}
		} else {
			// It's a file or symlink to file => just print it
			// Make a pseudo-Entry so we can reuse print logic
			linkTarget := ""
			if info.Mode()&os.ModeSymlink != 0 {
				t, _ := os.Readlink(path)
				linkTarget = t
			}
			e := Entry{
				Name:       filepath.Base(path),
				FullPath:   path,
				Info:       info,
				LinkTarget: linkTarget,
			}
			if flags.longListing {
				// We can mimic "ls -l file"
				fmt.Printf("total 1\n")
				printLong(e)
			} else {
				fmt.Println(e.Name)
			}
		}
	}
	return nil
}

// listDirectory reads the directory contents and returns a list of Entry objects
func listDirectory(dir string, flags lsFlags) ([]Entry, error) {
	entries := []Entry{}

	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	for _, fi := range files {
		// If not showAll, skip hidden files (files that start with '.')
		if !flags.showAll && strings.HasPrefix(fi.Name(), ".") {
			continue
		}

		fullPath := filepath.Join(dir, fi.Name())
		linkTarget := ""
		if fi.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(fullPath)
			if err == nil {
				linkTarget = target
			}
		}

		e := Entry{
			Name:       fi.Name(),
			FullPath:   fullPath,
			Info:       fi,
			LinkTarget: linkTarget,
		}
		entries = append(entries, e)
	}

	// Sort the entries
	sortEntries(entries, flags)

	return entries, nil
}

// recurseDir is used internally for printing sub-directories in -R mode
func recurseDir(paths []string, flags lsFlags) error {
	for _, path := range paths {
		fmt.Printf("\n%s:\n", path)

		entries, err := listDirectory(path, flags)
		if err != nil {
			return err
		}
		printEntries(path, entries, flags)

		// Then go deeper
		for _, e := range entries {
			if e.Info.IsDir() && e.Name != "." && e.Name != ".." {
				subPaths := []string{e.FullPath}
				err := recurseDir(subPaths, flags)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
