//go:build windows

package main

import "path/filepath"

func isAbs(path string) bool {
	return filepath.IsAbs(path)
}
