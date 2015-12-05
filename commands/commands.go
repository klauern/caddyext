package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var CaddyExt = &cobra.Command{
	Use:   "caddyext",
	Short: "Caddy's directive/extension manager",
	Long:  `Caddy's directive/extension manager`,
}

var caddyPath string = "github.com/mholt/caddy"
var caddyRegistry string = "http://raw.githubusercontent.com/caddyserver/buildsrv/master/features/registry.go"

var (
	ErrMissingArguments = errors.New("Missing arguments")
)

func Execute() {
	if len(os.Getenv("CADDYPATH")) > 0 {
		caddyPath = os.Getenv("CADDYPATH")
	}

	gopaths := strings.Split(os.Getenv("GOPATH"), string(filepath.ListSeparator))
	found := false
	for _, gopath := range gopaths {
		gopath = filepath.Join(gopath, "src")
		fpath := filepath.Join(gopath, caddyPath, "caddy/directives.go")
		if _, err := os.Stat(fpath); err == nil {
			found = true
			caddyPath = filepath.Join(gopath, caddyPath)
			break
		}
	}

	if found == false {
		fmt.Println("Caddy's source not found on any $GOPATH directories.")
		fmt.Println("Set CADDYPATH on your enviroment to a valid caddy source.")
		return
	}

	CaddyExt.AddCommand(buildCmd)
	CaddyExt.AddCommand(installCmd)
	CaddyExt.AddCommand(removeCmd)
	CaddyExt.AddCommand(stackCmd)
	CaddyExt.AddCommand(enableCmd)
	CaddyExt.AddCommand(disableCmd)
	CaddyExt.AddCommand(moveCmd)
	CaddyExt.AddCommand(resetCmd)
	CaddyExt.AddCommand(versionCmd)
	CaddyExt.Execute()
}

func resolveGoPath(path string) (foundPath string, found bool) {
	gopaths := strings.Split(os.Getenv("GOPATH"), string(filepath.ListSeparator))
	found = false
	foundPath = ""
	for _, gopath := range gopaths {
		gopath = filepath.Join(gopath, "src")
		fpath := filepath.Join(gopath, path)
		if _, err := os.Stat(fpath); err == nil {
			found = true
			foundPath = fpath
			break
		}
	}
	return
}
