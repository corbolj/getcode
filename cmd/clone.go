package getcode

import (
	getcode "github.com/corbolj/getcode/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	cobra.OnInitialize(init_viper_config)
	rootCmd.AddCommand(cloneCmd)
}

func init_viper_config() {

	viper.SetConfigName(".getcode")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	getcode.ErrorCheck(err, "Error reading config file.")
}
