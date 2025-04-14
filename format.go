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

	// 4) File size (or device info)
	sizeOrDevice := strconv.FormatInt(e.Info.Size(), 10)

	// 5) Time
	modTime := formatTime(e.Info.ModTime())

	// 6) Name
	nameStr := colorize(e.Name, e.Info.Mode())
	if e.Info.Mode()&os.ModeSymlink != 0 && e.LinkTarget != "" {
		targetInfo, err := os.Stat(e.LinkTarget)
		if err != nil {
			// fallback: just print uncolored target
			nameStr += " -> " + e.LinkTarget
		} else {
			nameStr += " -> " + colorize(e.LinkTarget, targetInfo.Mode())
		}
	}

	fmt.Printf("%s %3d %-8s %-8s %8s %s %s\n",
		modeStr,
		nlink,
		usrName,
		grpName,
		sizeOrDevice,
		modTime,
		nameStr)
}

func getUserName(uid uint32) string {
	usr, err := user.LookupId(strconv.Itoa(int(uid)))
	if err != nil {
		return strconv.Itoa(int(uid))
	}
	return usr.Username
}

func getGroupName(gid uint32) string {
	grp, err := user.LookupGroupId(strconv.Itoa(int(gid)))
	if err != nil {
		return strconv.Itoa(int(gid))
	}
	return grp.Name
}

func formatTime(t time.Time) string {
	now := time.Now()
	if now.Sub(t) > (time.Hour*24*30*6) || t.After(now) {
		return t.Format("Jan 02  2006")
	}
	return t.Format("Jan 02 15:04")
}

func modeToString(mode os.FileMode) string {
	var b strings.Builder

	switch {
	case mode.IsDir():
		b.WriteByte('d')
	case mode&os.ModeSymlink != 0:
		b.WriteByte('l')
	default:
		b.WriteByte('-')
	}

	perms := []struct {
		bit  os.FileMode
		char byte
	}{
		{mode & 0400, 'r'}, {mode & 0200, 'w'}, {mode & 0100, 'x'},
		{mode & 0040, 'r'}, {mode & 0020, 'w'}, {mode & 0010, 'x'},
		{mode & 0004, 'r'}, {mode & 0002, 'w'}, {mode & 0001, 'x'},
	}

	for _, p := range perms {
		if p.bit != 0 {
			b.WriteByte(p.char)
		} else {
			b.WriteByte('-')
		}
	}

	return b.String()
}
