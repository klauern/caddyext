package commands

import (
	"fmt"
	"path/filepath"

	"github.com/caddyserver/caddyext/directives"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset caddy state",
	Long:  `Reset caddy state`,
	Run:   ResetCaddy,
}

// Errors
var ()

func init() {
}

func ResetCaddy(cmd *cobra.Command, args []string) {
	dir, err := directives.NewFrom(filepath.Join(caddyPath, "caddy/directives.go"))
	if err != nil {
		cmdError(cmd, err)
	}

	dir.Reset()

	err = dir.Save()
	if err != nil {
		cmdError(cmd, err)
	}

	fmt.Println("Active caddy state has been reseted.")
}
