package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

func main() {
	cobra.CheckErr(Root.Execute())
}

var (
	list bool
	run  string
)

func init() {
	Root.SetHelpCommand(&cobra.Command{Hidden: true})
	Root.Flags().BoolVarP(&list, "list", "l", false, "list commands without running \"go generate\"")
	Root.Flags().StringVarP(&run, "run", "r", "", "specifies a regular expression to select directives whose full original source text matches the expression")
}

var Root = &cobra.Command{
	Use:           "genx",
	Version:       "v0.3.1",
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

		items := make([]*GenerateItem, 0, 10)
		for queue.Len() > 0 {
			dir := queue.Pop()
			if !path.IsAbs(dir) {
				dir = path.Join(pwd, dir)
			}

			assert(path.IsAbs(dir), "not an abs path")
			stat, err := os.Stat(dir)
			if err != nil {
				return err
			}

			if !stat.IsDir() {
				return fmt.Errorf("%q is not an directory", dir)
			}

			files, err := os.ReadDir(dir)
			if err != nil {
				return err
			}

			for _, file := range files {
				name := file.Name()
				if !path.IsAbs(name) {
					name = path.Join(dir, name)
				}

				assert(path.IsAbs(name), "not an abs path")
				if file.IsDir() {
					queue.Push(name)
					continue
				}

				if ext := path.Ext(name); ext != GoExt {
					continue
				}

				r, err := os.Open(name)
				if err != nil {
					return err
				}

				commands := GoGenerateCommands(r)
				for _, command := range commands {
					items = append(items, &GenerateItem{
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
			copied := make([]*GenerateItem, len(items))
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
				dir := path.Dir(s.File)
				if _, exists := generated[dir]; !exists {
					if err = Generate(dir); err != nil {
						return err
					}
					generated[dir] = struct{}{}
				}
			}
		}

		return nil
	},
}

const (
	GoExt      = ".go"
	GoGenerate = "//go:generate "
)

type Command struct {
	Pos int
	Cmd string
}

func GoGenerateCommands(r io.Reader) (commands []*Command) {
	commands = make([]*Command, 0, 10)
	scanner := bufio.NewScanner(r)
	for pos := 1; scanner.Scan(); pos++ {
		if text := scanner.Text(); strings.HasPrefix(text, GoGenerate) {
			commands = append(commands, &Command{
				Pos: pos,
				Cmd: text,
			})
		}
	}
	return commands
}

func Generate(dir string) (err error) {
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

type GenerateItem struct {
	File    string
	Command *Command
	Repr    string
}
