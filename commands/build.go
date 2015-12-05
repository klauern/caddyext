package commands

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [output path]",
	Short: "Build caddy from the current state",
	Long:  `Build caddy from the current state`,
	Run:   BuildCaddy,
}

// Errors
var (
	ErrBuildSourceNotFound          = errors.New("Caddy source doesn't exist inside current GOPATH")
	ErrBuildOutputDirectoryNotFound = errors.New("Output directory doesn't exist")
)

func init() {
}

func BuildCaddy(cmd *cobra.Command, args []string) {
	outputPath, _ := os.Getwd()

	if len(args) > 0 {
		if filepath.IsAbs(args[0]) {
			outputPath = args[0]
		} else {
			outputPath = filepath.Join(outputPath, args[0])
		}
		args = args[1:]
	}

	outputPath = filepath.Join(outputPath, "customCaddy")

	if err := caddyBuild(caddyPath, outputPath, args...); err != nil {
		cmdError(cmd, err)
	}
}

func caddyBuild(caddy, output string, args ...string) error {
	cmd := exec.Command("go", append([]string{"build", "-o", output}, args...)...)
	cmd.Dir = caddy
	errBuf := new(bytes.Buffer)
	cmd.Stderr = errBuf
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Println("new caddy build:", output)

	return nil
}
