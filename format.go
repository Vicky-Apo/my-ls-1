package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// printEntries prints out the list of entries
// If flags.longListing is false, we just show names. If it is true, we show details.
func printEntries(dir string, entries []Entry, flags lsFlags) {
	if flags.longListing {
		// First, print "total" line that `ls -l` typically includes
		// We'll sum up st_blocks (512-byte blocks). On many systems, that's in syscall.Stat_t.
		total := int64(0)
		for _, e := range entries {
			st := e.Info.Sys()
			if st, ok := st.(*syscall.Stat_t); ok {
				total += int64(st.Blocks) * 512 / 1024
			}
		}
		// If you want the exact block size as `ls`, you may need to multiply or configure differently.
		fmt.Printf("total %d\n", total)

		// Now print each entry in "long" format
		for _, e := range entries {
			printLong(e)
		}
	} else {
		// Normal (multi-column in real ls) but we’ll just do single-line for each in this example
		for _, e := range entries {
			colored := colorize(e.Name, e.Info.Mode())
			fmt.Printf("%s  ", colored)
		}
		fmt.Println()
	}
}

// printLong prints one entry in "ls -l" style
func printLong(e Entry) {
	st, ok := e.Info.Sys().(*syscall.Stat_t)
	if !ok {
		// fallback
		fmt.Printf("%s\n", e.Name)
		return
	}

	// 1) File mode string (e.g. drwxr-xr-x)
	modeStr := modeToString(e.Info.Mode())

	// 2) Number of hard links
	nlink := st.Nlink

	// 3) Owner name
	usrName := getUserName(st.Uid)
	grpName := getGroupName(st.Gid)

	// 4) File size (for normal files) or major/minor for device
	sizeOrDevice := strconv.FormatInt(e.Info.Size(), 10)

	// 5) Time
	modTime := formatTime(e.Info.ModTime())

	// 6) Name
	nameStr := colorize(e.Name, e.Info.Mode())
	// If it's a symlink, append -> target
	if e.Info.Mode()&os.ModeSymlink != 0 {
		nameStr += " -> " + e.LinkTarget
	}

	// Typical `ls -l` output layout:
	//
	// drwxr-xr-x  2 user  group  4096 Jun  4 12:34 myfolder
	// or
	// -rw-r--r--  1 user  group  1042 Oct 19  2022 somefile.txt
	//
	// We'll do something like:
	fmt.Printf("%s %3d %-8s %-8s %8s %s %s\n",
		modeStr,
		nlink,
		usrName,
		grpName,
		sizeOrDevice,
		modTime,
		nameStr)
}

// formatTime tries to replicate the “ls -l” time format
// If the file’s mtime is older than ~6 months, standard `ls` shows the year instead of HH:MM.

func formatTime(t time.Time) string {
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)
	if t.Before(sixMonthsAgo) {
		// older than 6 months => "Mon _2  2006"
		return t.Format("Jan  2  2006")
	}
	// else => "Mon _2 15:04"
	return t.Format("Jan  2 15:04")
}

func modeToString(m os.FileMode) string {
	var str strings.Builder

	// 1. File type
	switch {
	case m.IsDir():
		str.WriteRune('d')
	case m&os.ModeSymlink != 0:
		str.WriteRune('l')
	case m&os.ModeNamedPipe != 0:
		str.WriteRune('p')
	case m&os.ModeSocket != 0:
		str.WriteRune('s')
	case m&os.ModeDevice != 0:
		if m&os.ModeCharDevice != 0 {
			str.WriteRune('c')
		} else {
			str.WriteRune('b')
		}
	default:
		str.WriteRune('-')
	}

	// 2. Owner bits
	str.WriteByte(rBit(m, 0400))
	str.WriteByte(wBit(m, 0200))
	str.WriteByte(xBit(m, 0100, m&os.ModeSetuid != 0))

	// 3. Group bits
	str.WriteByte(rBit(m, 0040))
	str.WriteByte(wBit(m, 0020))
	str.WriteByte(xBit(m, 0010, m&os.ModeSetgid != 0))

	// 4. Others bits
	str.WriteByte(rBit(m, 0004))
	str.WriteByte(wBit(m, 0002))
	str.WriteByte(xBit(m, 0001, m&os.ModeSticky != 0))

	return str.String()
}

func rBit(m os.FileMode, mask os.FileMode) byte {
	if m&mask != 0 {
		return 'r'
	}
	return '-'
}

func wBit(m os.FileMode, mask os.FileMode) byte {
	if m&mask != 0 {
		return 'w'
	}
	return '-'
}

func xBit(m os.FileMode, mask os.FileMode, special bool) byte {
	if m&mask != 0 {
		if special {
			// setuid, setgid, or sticky bit set
			// typical ls outputs 's' or 't'
			return 's'
		} else {
			return 'x'
		}
	} else {
		if special {
			return 'S'
		} else {
			return '-'
		}
	}
}

// A more thorough approach for rwx bits:
func init() {
	// We can rewrite the modeToString more robustly:
	// but for brevity, let's keep it short.
	// See explanation in the code comment below for a better approach.
}

// getUserName returns the username for a given UID
func getUserName(uid uint32) string {
	u, err := user.LookupId(strconv.FormatUint(uint64(uid), 10))
	if err != nil {
		// fallback to UID
		return strconv.FormatUint(uint64(uid), 10)
	}
	return u.Username
}

// getGroupName returns the group name for a given GID
func getGroupName(gid uint32) string {
	g, err := user.LookupGroupId(strconv.FormatUint(uint64(gid), 10))
	if err != nil {
		// fallback
		return strconv.FormatUint(uint64(gid), 10)
	}
	return g.Name
}
