package commands

import (
	"fmt"
	"path/filepath"

	"github.com/pedronasser/caddyext/directives"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove an extension from caddy's directives source (only 3rd-party)",
	Long:  `Remove an extension from caddy's directives source (only 3rd-party)`,
	Run:   RemoveExtension,
}

func init() {
}

func RemoveExtension(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmdError(cmd, ErrMissingArguments)
	}

	dir, err := directives.NewFrom(filepath.Join(caddyPath, "config/directives.go"))
	if err != nil {
		cmdError(cmd, err)
	}

	name := args[0]

	err = dir.RemoveDirective(name)
	if err != nil {
		cmdError(cmd, err)
	}

	err = dir.Save()
	if err != nil {
		cmdError(cmd, err)
	}

	fmt.Println(name, "successfully removed from Caddy.")
}
