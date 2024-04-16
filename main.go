package main

import (
	"github.com/spf13/cobra"
	"github.com/x5iu/genx/cmd"
)

func main() {
	cobra.CheckErr(cmd.Generate.Execute())
}
