package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pedronasser/caddyext/directives"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install <name> <source>",
	Short: "Install and enables a extension",
	Long:  `Install and enables a extension`,
	Run:   InstallExtension,
}

// Flags
var (
	flagUpdate bool
)

// Errors
var (
	ErrInstallNoSource       = errors.New("Undefined extension source")
	ErrInstallSourceNotFound = errors.New("Extension source doesn't exist inside current GOPATH")
)

func init() {
	installCmd.PersistentFlags().BoolVar(&flagUpdate, "u", false, "update source")
}

func InstallExtension(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmdError(cmd, ErrMissingArguments)
	}

	dir, err := directives.NewFrom(filepath.Join(caddyPath, "caddy/directives.go"))
	if err != nil {
		cmdError(cmd, err)
	}

	name := args[0]
	source := args[1]

	gopaths := strings.Split(os.Getenv("GOPATH"), string(filepath.ListSeparator))
	found := false
	for _, gopath := range gopaths {
		gopath = filepath.Join(gopath, "src")
		fpath := filepath.Join(gopath, source)
		if _, err := os.Stat(fpath); err == nil {
			found = true
			caddyPath = filepath.Join(gopath, caddyPath)
			break
		}
	}

	if found == false {
		cmdError(cmd, ErrInstallSourceNotFound)
	}

	getExtension(source, flagUpdate)

	err = dir.AddDirective(name, source)
	if err != nil {
		cmdError(cmd, err)
	}

	err = dir.Save()
	if err != nil {
		cmdError(cmd, err)
	}

	fmt.Println(name, "successfully added to Caddy.")
}

func getExtension(source string, update bool) error {
	var updateFlag string
	if update {
		updateFlag = "-u"
	}
	_, err := exec.Command("go", "get", updateFlag, source).Output()
	return err
}
