package commands

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/caddyserver/caddyext/directives"
	"github.com/spf13/cobra"
)

var moveCmd = &cobra.Command{
	Use:   "move <name> <stack-index>",
	Short: "Move target's index on Caddy's stack",
	Long:  `Move target's index on Caddy's stack`,
	Run:   MoveExtension,
}

func init() {
}

func MoveExtension(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmdError(cmd, ErrMissingArguments)
	}

	dir, err := directives.NewFrom(filepath.Join(caddyPath, "caddy/directives.go"))
	if err != nil {
		cmdError(cmd, err)
	}

	name := args[0]
	index, _ := strconv.Atoi(args[1])

	err = dir.MoveDirective(name, index)
	if err != nil {
		cmdError(cmd, err)
	}

	err = dir.Save()
	if err != nil {
		cmdError(cmd, err)
	}

	fmt.Println(name, "successfully moved.")
}
