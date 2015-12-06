package commands

import (
	"fmt"
	"path/filepath"

	"github.com/caddyserver/caddyext/directives"
	"github.com/spf13/cobra"
)

var stackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Show stack of directives/extensions",
	Long:  `Show stack of directives/extensions`,
	Run:   StackExtension,
}

func init() {
}

func StackExtension(cmd *cobra.Command, args []string) {
	dir, _ := directives.NewFrom(filepath.Join(caddyPath, "caddy/directives.go"))
	list := dir.List()
	fmt.Println("\nAvailable Caddy directives/extensions:")
	fmt.Println("   (✓) ENABLED | (-) DISABLED\n")
	for i, d := range list {
		active := "✓"
		isCore := ""
		if !d.Active {
			active = "-"
		}
		if d.Core {
			isCore = "(core)"
		}
		fmt.Printf("   %d. (%s) %s %s\n", i, active, d.Name, isCore)
	}
	fmt.Println("")
}
