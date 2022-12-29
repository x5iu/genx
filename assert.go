package main

import "fmt"

func assert(expr bool, format string, args ...any) {
	if !expr {
		panic(fmt.Sprintf(format, args...))
	}
}
