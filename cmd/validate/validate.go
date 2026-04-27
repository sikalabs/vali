package validate

import (
	"fmt"
	"os"

	"github.com/sikalabs/vali/cmd/root"
	"github.com/sikalabs/vali/pkg/vali"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "validate",
	Aliases: []string{"vali"},
	Args:    cobra.NoArgs,
	Run: func(c *cobra.Command, args []string) {
		config, err := vali.ReadConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		out := vali.Validate(config)
		vali.PrintValidateOutput(out)
		if !out.OK {
			os.Exit(1)
		}
	},
}

func init() {
	root.Cmd.AddCommand(Cmd)
}
