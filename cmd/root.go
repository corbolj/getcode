package getcode

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "v0.0.1"

var rootCmd = &cobra.Command{
	Use:     "turndown",
	Version: version,
	Short:   "getgo - to help manage cloning and migrating git versioned code.",
	Long:    `getgo - to help manage cloning and migrating git versioned code.`,
	Run:     func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
