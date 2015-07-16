package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

const CaddyExtVersion = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show caddyext's version",
	Long:  `Show caddyext's version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("CaddyExt's v%s\n", CaddyExtVersion)
	},
}
