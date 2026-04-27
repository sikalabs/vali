package cmd

import (
	"github.com/sikalabs/vali/cmd/root"
	_ "github.com/sikalabs/vali/cmd/validate"
	_ "github.com/sikalabs/vali/cmd/version"
	"github.com/spf13/cobra"
)

func Execute() {
	cobra.CheckErr(root.Cmd.Execute())
}
