package cmd

import (
	"bytes"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var rootCmd = &cobra.Command{
	Use:   "pochama",
	Short: "Pochama is a packaging tool for some tedious operations.",
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
	rootCmd.PersistentFlags().BoolVar(&initCfg, "initconfig", false, "create default config file only if the file does not exist")

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

func createCfgIFNotExists(cfg any) error {
	yamlCfg, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = viper.ReadConfig(bytes.NewBuffer(yamlCfg))
	if err != nil {
		return err
	}
	err = viper.SafeWriteConfig()
	if err != nil {
		return err
	}
	return nil
}
