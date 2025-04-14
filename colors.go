package main

import (
	"fmt"
	"os"
)

func colorize(name string, mode os.FileMode) string {
	if mode&os.ModeSymlink != 0 {
		return fmt.Sprintf("\033[1;36m%s\033[0m", name) // Cyan (bold) for symlinks
	} else if mode.IsDir() {
		return fmt.Sprintf("\033[1;34m%s\033[0m", name) // Blue (bold) for directories
	} else if mode&0111 != 0 {
		return fmt.Sprintf("\033[1;32m%s\033[0m", name) // Green (bold) for executables
	} else if mode&os.ModeDevice != 0 {
		return fmt.Sprintf("\033[1;33m%s\033[0m", name) // Yellow (bold) for devices
	}
	return name // No color - Default color
}
