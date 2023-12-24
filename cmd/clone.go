package getcode

import (
	getcode "github.com/corbolj/getcode/pkg"
	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "clones down all repos in an organization or user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		i := args[0]
		getcode.Clone(i)
	},
}

func init() {
	// Pull configs into object
	rootCmd.AddCommand(cloneCmd)
}
