package commands

import (
	"fmt"
	"path/filepath"

	"github.com/pedronasser/caddyext/directives"
	"github.com/spf13/cobra"
)

var disableCmd = &cobra.Command{
	Use:   "disable <name>",
	Short: "Disables a installed directive or extension",
	Long:  `Disables a installed directive or extension`,
	Run:   DisableExtension,
}

func init() {
}

func DisableExtension(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmdError(cmd, ErrMissingArguments)
	}

	dir, err := directives.NewFrom(filepath.Join(caddyPath, "config/directives.go"))
	if err != nil {
		cmdError(cmd, err)
	}

	name := args[0]

	err = dir.DisableDirective(name)
	if err != nil {
		cmdError(cmd, err)
	}

	err = dir.Save()
	if err != nil {
		cmdError(cmd, err)
	}

	fmt.Println(name, "successfully disabled.")
}
