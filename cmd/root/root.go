package root

import (
	"github.com/sikalabs/vali/version"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "vali",
	Short: "vali, " + version.Version,
}
