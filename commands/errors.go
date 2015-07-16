package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func cmdError(cmd *cobra.Command, err error) {
	cmd.Usage()
	fmt.Printf("`%s` error: %s\n", cmd.Name(), err)
	os.Exit(0)
}
