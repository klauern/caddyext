package commands

import (
	"fmt"
	"path/filepath"

	"github.com/pedronasser/caddyext/directives"
	"github.com/spf13/cobra"
)

var enableCmd = &cobra.Command{
	Use:   "enable <name>",
	Short: "Enables a installed directive or extension",
	Long:  `Enables a installed directive or extension`,
	Run:   EnableExtension,
}

func init() {
}

func EnableExtension(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmdError(cmd, ErrMissingArguments)
	}

	dir, err := directives.NewFrom(filepath.Join(caddyPath, "config/directives.go"))
	if err != nil {
		cmdError(cmd, err)
	}

	name := args[0]

	err = dir.EnableDirective(name)
	if err != nil {
		cmdError(cmd, err)
	}

	err = dir.Save()
	if err != nil {
		cmdError(cmd, err)
	}

	fmt.Println(name, "successfully enabled.")
}
