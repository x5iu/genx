//go:build windows

package cmd

import "path/filepath"

func isAbs(path string) bool {
	return filepath.IsAbs(path)
}
