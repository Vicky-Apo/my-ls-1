package main

import (
	"fmt"
	"os"
)

func colorize(name string, mode os.FileMode) string {
	if mode&os.ModeSymlink != 0 {
		return fmt.Sprintf("\033[36m%s\033[0m", name) // Cyan for symlinks
	} else if mode.IsDir() {
		return fmt.Sprintf("\033[34m%s\033[0m", name) // Blue for directories
	} else if mode&0111 != 0 {
		return fmt.Sprintf("\033[32m%s\033[0m", name) // Green for executables
	}
	return name // No color
}
