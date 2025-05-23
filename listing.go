package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

// Entry holds the file information we need for printing
type Entry struct {
	Name       string
	FullPath   string
	Info       os.FileInfo
	LinkTarget string // if symlink, store the target here
}

// walk handles the logic of listing directories and optionally recursing
func walk(paths []string, flags lsFlags) error {
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("cannot access '%s': %v", path, err)
		}

		// Show the directory header only when needed
		if info.IsDir() && (len(paths) > 1 || flags.recursive || path != ".") {
			fmt.Printf("%s:\n", path)
		}

		if info.IsDir() {
			entries, err := listDirectory(path, flags)
			if err != nil {
				return err
			}
			sortEntries(entries, flags)
			printEntries(path, entries, flags)

			if flags.recursive {
				var subDirs []string
				for _, e := range entries {
					if e.Info.IsDir() && e.Name != "." && e.Name != ".." {
						subDirs = append(subDirs, e.FullPath)
					}
				}
				err := walk(subDirs, flags)
				if err != nil {
					return err
				}
			}
		} else {
			linkTarget := ""
			if info.Mode()&os.ModeSymlink != 0 {
				t, _ := os.Readlink(path)
				linkTarget = t
			}
			e := Entry{
				Name:       getBase(path),
				FullPath:   path,
				Info:       info,
				LinkTarget: linkTarget,
			}
			if flags.longListing {
				fmt.Printf("total 1\n")
				printLong(e)
			} else {
				fmt.Println(e.Name)
			}
		}
	}
	return nil
}

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

	if flags.showAll {
		dot, _ := os.Stat(dir)
		dotdot, _ := os.Stat(dir + "/..")

		entries = append(entries, Entry{
			Name:       ".",
			FullPath:   dir,
			Info:       dot,
			LinkTarget: "",
		})
		entries = append(entries, Entry{
			Name:       "..",
			FullPath:   dir + "/..",
			Info:       dotdot,
			LinkTarget: "",
		})
	}

	for _, fi := range files {
		if !flags.showAll && isHidden(fi.Name()) {
			continue
		}

		fullPath := dir + "/" + fi.Name()
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

	return entries, nil
}

func printEntries(_ string, entries []Entry, flags lsFlags) {
	if flags.longListing {
		total := int64(0)
		for _, e := range entries {
			st := e.Info.Sys()
			if st, ok := st.(*syscall.Stat_t); ok {
				total += int64(st.Blocks) * 512 / 1024
			}
		}
		fmt.Printf("total %d\n", total)
		for _, e := range entries {
			printLong(e)
		}
	} else {
		for _, e := range entries {
			fmt.Println(colorize(e.Name, e.Info.Mode()))
		}
		fmt.Println()
	}
}

func getBase(path string) string {
	if path == "" {
		return ""
	}
	slash := strings.Split(path, "/")
	return slash[len(slash)-1]
}
