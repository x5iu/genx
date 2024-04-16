package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	list bool
	run  string
)

func init() {
	Generate.SetHelpCommand(&cobra.Command{Hidden: true})
	Generate.Flags().BoolVarP(&list, "list", "l", false, "list commands without running \"go generate\"")
	Generate.Flags().StringVarP(&run, "run", "r", "", "specifies a regular expression to select directives whose full original source text matches the expression")
}

var Generate = &cobra.Command{
	Use:           "genx",
	Version:       "v0.6.0",
	SilenceUsage:  true,
	SilenceErrors: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   true,
		DisableDescriptions: true,
		HiddenDefaultCmd:    true,
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}

		queue := NewQueue(args)
		if queue.Len() == 0 {
			queue.Push(".")
		}

		items := make([]*generateItem, 0, 10)
		for queue.Len() > 0 {
			dir := queue.Pop()
			if !isAbs(dir) {
				dir = filepath.Join(pwd, dir)
			}

			assert(isAbs(dir), "not an abs path: "+strconv.Quote(dir))
			stat, err := os.Stat(dir)
			if err != nil {
				return err
			}

			var files []fs.DirEntry
			if !stat.IsDir() {
				files = []fs.DirEntry{
					&fileEntry{Filename: dir},
				}
			} else {
				files, err = os.ReadDir(dir)
				if err != nil {
					return err
				}
			}

			for _, file := range files {
				name := file.Name()
				if !isAbs(name) {
					name = filepath.Join(dir, name)
				}

				assert(isAbs(name), "not an abs path: "+strconv.Quote(name))
				if file.IsDir() {
					queue.Push(name)
					continue
				}

				if ext := filepath.Ext(name); ext != goExt {
					continue
				}

				r, err := os.Open(name)
				if err != nil {
					return err
				}

				commands := getGoGenerateCommands(r)
				for _, command := range commands {
					items = append(items, &generateItem{
						File:    name,
						Command: command,
					})
				}
			}
		}

		if run != "" {
			var re *regexp.Regexp
			if re, err = regexp.Compile(run); err != nil {
				return err
			}
			copied := make([]*generateItem, len(items))
			copy(copied, items)
			items = items[:0]
			for _, item := range copied {
				if re != nil && !re.MatchString(item.Command.Cmd) {
					continue
				}
				items = append(items, item)
			}
		}

		align(pwd, items)
		generated := make(map[string]struct{}, len(items))
		for _, s := range items {
			fmt.Println(s.Repr)
			if !list {
				dir := filepath.Dir(s.File)
				if _, exists := generated[dir]; !exists {
					if err = generate(dir); err != nil {
						return err
					}
					generated[dir] = struct{}{}
				}
			}
		}

		return nil
	},
}

type fileEntry struct {
	Filename string
}

func (f *fileEntry) Name() string               { return f.Filename }
func (f *fileEntry) IsDir() bool                { return false }
func (f *fileEntry) Type() fs.FileMode          { panic("unreachable") }
func (f *fileEntry) Info() (fs.FileInfo, error) { panic("unreachable") }

const (
	goExt            = ".go"
	goGeneratePrefix = "//go:generate "
)

type command struct {
	Pos int
	Cmd string
}

func getGoGenerateCommands(r io.Reader) (commands []*command) {
	commands = make([]*command, 0, 10)
	scanner := bufio.NewScanner(r)
	for pos := 1; scanner.Scan(); pos++ {
		if text := scanner.Text(); strings.HasPrefix(text, goGeneratePrefix) {
			commands = append(commands, &command{
				Pos: pos,
				Cmd: text,
			})
		}
	}
	return commands
}

func generate(dir string) (err error) {
	if err = os.Chdir(dir); err != nil {
		return err
	}
	args := make([]string, 0, 4)
	args = append(args, "generate")
	if run != "" {
		args = append(args, "-run", run)
	}
	args = append(args, ".")
	cmd := exec.Command("go", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type generateItem struct {
	File    string
	Command *command
	Repr    string
}
