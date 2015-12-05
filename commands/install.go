package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"go/ast"
	"go/parser"
	"go/token"

	"github.com/pedronasser/caddyext/directives"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:     "install <name:[repository]...>",
	Short:   "Install and enable extension(s)",
	Long:    `Install and enable extension(s)`,
	Example: `  caddyext install git search:github.com/pedronasser/caddy-search`,
	Run:     InstallExtension,
}

// Flags
var (
	flagUpdate bool
	flagAfter  string
	flagBefore string
)

// Errors
var (
	ErrInstallNoSource         = errors.New("Couldn't find extension's repository")
	ErrInstallSourceNotFound   = errors.New("Extension source doesn't exist inside current GOPATH")
	ErrInstallExtensionResolve = errors.New("Coundn't resolve that extension from Caddy's registry. Please provide a repository for extension.")
)

func init() {
	installCmd.PersistentFlags().BoolVar(&flagUpdate, "u", false, "update source")
	installCmd.PersistentFlags().StringVarP(&flagAfter, "after", "a", "", "directive that new directives would be installed AFTER")
	installCmd.PersistentFlags().StringVarP(&flagBefore, "before", "b", "", "directive that new directives would be installed BEFORE")
}

func InstallExtension(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmdError(cmd, ErrMissingArguments)
	}

	dir, err := directives.NewFrom(filepath.Join(caddyPath, "caddy/directives.go"))
	if err != nil {
		cmdError(cmd, err)
	}

	directives := dir.List()

	for _, directive := range args {
		dirparts := strings.Split(directive, ":")
		name := dirparts[0]

		var source string
		if len(dirparts) < 2 {
			fmt.Printf("trying to resolve `%s` from Caddy's registry\n", name)
			source = resolveExtension(name)
			if len(source) == 0 {
				cmdError(cmd, ErrInstallExtensionResolve)
			}
		} else {
			source = dirparts[1]
		}

		getExtension(source, flagUpdate)

		gopaths := strings.Split(os.Getenv("GOPATH"), string(filepath.ListSeparator))
		found := false
		for _, gopath := range gopaths {
			gopath = filepath.Join(gopath, "src")
			fpath := filepath.Join(gopath, source)
			if _, err := os.Stat(fpath); err == nil {
				found = true
				break
			}
		}

		if found == false {
			cmdError(cmd, ErrInstallSourceNotFound)
		}

		err = dir.AddDirective(name, source)
		if err != nil {
			cmdError(cmd, err)
		}

		if len(flagAfter) > 0 || len(flagBefore) > 0 {
			for i, d := range directives {
				if d.Name == flagBefore {
					dir.MoveDirective(name, i)
				} else if d.Name == flagAfter {
					dir.MoveDirective(name, i+1)
				}
			}
		}

		fmt.Printf("`%s` added to Caddy.\n", name)
	}

	err = dir.Save()
	if err != nil {
		cmdError(cmd, err)
	}
}

func resolveExtension(extension string) (resolved string) {
	resp, err := http.Get(caddyRegistry)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	f, err := parser.ParseFile(token.NewFileSet(), "", data, 0)
	node, ok := f.Scope.Lookup("Registry").Decl.(ast.Node)
	if !ok {
		return
	}

	c := node.(*ast.ValueSpec).Values[0].(*ast.CompositeLit)
	for _, m := range c.Elts {
		var directive *ast.BasicLit
		var directiveRepo *ast.BasicLit
		token := m.(*ast.CompositeLit).Elts

		if v, ok := token[0].(*ast.BasicLit); ok {
			directive = v
		} else {
			return
		}

		name, err := strconv.Unquote(directive.Value)
		if err != nil {
			return
		}

		if v, ok := token[1].(*ast.BasicLit); ok {
			directiveRepo = v
		} else {
			return
		}

		repo, err := strconv.Unquote(directiveRepo.Value)
		if err != nil {
			return
		}

		if name == extension {
			if len(repo) > 0 {
				resolved = repo
			}
			return
		}
	}

	return
}

func getExtension(source string, update bool) error {
	var updateFlag string
	if update {
		updateFlag = "-u"
	}
	_, err := exec.Command("go", "get", updateFlag, source).Output()
	return err
}
