package main

import "path/filepath"

func check(filename string) bool {
	extStart := -1
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			extStart = i
			break
		}
	}
	if extStart == -1 && checkExtensions > 0 {
		return false
	} else if extStart != -1 {
		ext := filename[extStart:]
		if extensions[ext]&checkExtensions != checkExtensions {
			return false
		}
	}
	for _, p := range exclude {
		if ok, _ := filepath.Match(p, filename); ok {
			return false
		}
	}
	checkRestrict := len(restrict) == 0
	for _, p := range restrict {
		if ok, _ := filepath.Match(p, filename); ok {
			checkRestrict = true
		}
	}
	return checkRestrict
}
