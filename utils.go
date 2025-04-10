package main

func isHidden(name string) bool {
	return len(name) > 1 && name[0] == '.'
}
