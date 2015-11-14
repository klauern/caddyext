package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func cmdError(cmd *cobra.Command, err error) {
	fmt.Printf("\n`caddyext %s` error: %s\n", cmd.Name(), err)
	cmd.Usage()
	os.Exit(0)
}
