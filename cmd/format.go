package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"strconv"
	"strings"
)

func shorten(pwd string, name string) string {
	if strings.HasPrefix(name, pwd) {
		assert(isAbs(name), "not an abs path: "+strconv.Quote(name))
		return strings.TrimPrefix(name, pwd+"/")
	}
	return name
}

var (
	boldGreen      = color.New(color.FgGreen, color.Bold).SprintfFunc()
	green          = color.New(color.FgGreen).SprintfFunc()
	italicBoldBlue = color.New(color.FgBlue, color.Italic, color.Bold).SprintfFunc()
	boldHiWhite    = color.New(color.FgHiWhite, color.Bold).SprintfFunc()
)

func align(pwd string, items []*generateItem) {
	var width int
	for _, item := range items {
		if l := len(shorten(pwd, item.File)) + len(strconv.Itoa(item.Command.Pos)); l > width {
			width = l
		}
	}
	for _, item := range items {
		file := "[" + shorten(pwd, item.File) + ":" + strconv.Itoa(item.Command.Pos) + "]"
		if l := len(file) - 3; l < width {
			file += strings.Repeat(" ", width-l)
		}
		item.Repr = fmt.Sprintf("%s %s", boldGreen(file), format(item.Command.Cmd))
	}
}

func format(command string) string {
	command = strings.TrimPrefix(command, goGeneratePrefix)
	args := split(command)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if i == 0 {
			args[i] = italicBoldBlue(arg)
		} else if strings.HasPrefix(arg, "-") {
			if eq := strings.Index(arg, "="); eq >= 0 {
				args[i] = boldHiWhite(arg[:eq]) + arg[eq:]
			} else {
				args[i] = boldHiWhite(arg)
			}
		} else if strings.HasPrefix(arg, "\"") || strings.HasPrefix(arg, "'") {
			args[i] = green(arg)
		}
	}
	return strings.Join(args, " ")
}

func split(line string) (args []string) {
	var (
		singleQuoted bool
		doubleQuoted bool
		arg          []byte
	)

	for i := 0; i < len(line); i++ {
		switch ch := line[i]; ch {
		case ' ':
			if doubleQuoted || singleQuoted {
				arg = append(arg, ch)
			} else if len(arg) > 0 {
				args = append(args, string(arg))
				arg = arg[:0]
			}
		case '"':
			if !(i > 0 && line[i-1] == '\\' || singleQuoted) {
				doubleQuoted = !doubleQuoted
			}
			arg = append(arg, ch)
		case '\'':
			if !(i > 0 && line[i-1] == '\\' || doubleQuoted) {
				singleQuoted = !singleQuoted
			}
			arg = append(arg, ch)
		default:
			arg = append(arg, ch)
		}
	}

	if len(arg) > 0 {
		args = append(args, string(arg))
	}

	return args
}
