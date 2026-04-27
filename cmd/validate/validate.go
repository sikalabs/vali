package validate

import (
	"fmt"
	"log"

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
		handleError(err)

		out := vali.Validate(config)
		if out.OK == false {
			log.Fatalln("Validation failed")
		}

		fmt.Println("OK")
	},
}

func init() {
	root.Cmd.AddCommand(Cmd)
}

func handleError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
