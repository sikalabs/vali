package validate_and_run

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/sikalabs/vali/cmd/root"
	"github.com/sikalabs/vali/pkg/vali"
	"github.com/spf13/cobra"
)

var FlagConfig string

var Cmd = &cobra.Command{
	Use:     "validate-and-run",
	Aliases: []string{"vali-and-run"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(c *cobra.Command, args []string) {
		vali.ValidateCommand(FlagConfig)

		bin, err := exec.LookPath(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		if err := syscall.Exec(bin, args, os.Environ()); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	root.Cmd.AddCommand(Cmd)
	Cmd.Flags().StringVarP(
		&FlagConfig,
		"config",
		"c",
		"",
		"path to config file",
	)
	Cmd.Flags().SetInterspersed(false)
}
