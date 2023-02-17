//go:build darwin || linux || unix

package main

import (
	gopath "path"
)

func isAbs(path string) bool {
	return gopath.IsAbs(path)
}
