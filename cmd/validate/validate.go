package validate

import (
	"github.com/sikalabs/vali/cmd/root"
	"github.com/sikalabs/vali/pkg/vali"
	"github.com/spf13/cobra"
)

var FlagConfig string

var Cmd = &cobra.Command{
	Use:     "validate",
	Aliases: []string{"vali"},
	Args:    cobra.NoArgs,
	Run: func(c *cobra.Command, args []string) {
		vali.ValidateCommand(FlagConfig)
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
}
