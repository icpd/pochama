package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "pochama",
	Short: "pochama is a tool for developing",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cfgDir string
var initCfg bool

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgDir, "dir", "d", "", "config file directory (default is $HOME/.config/pochama)")
	rootCmd.PersistentFlags().BoolVar(&initCfg, "initcfg", false, "create default config file only if the file does not exist")

	initCfgDir()
}

func initCfgDir() {
	if cfgDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			cobra.CheckErr(err)
		}
		cfgDir = homeDir + "/.config/pochama"

		if _, err = os.Stat(cfgDir); os.IsNotExist(err) {
			err = os.MkdirAll(cfgDir, 0755)
			if err != nil {
				cobra.CheckErr(err)
			}
		}
	}

	viper.AddConfigPath(cfgDir)
	viper.SetConfigType("yaml")
}
