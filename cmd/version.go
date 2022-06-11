package cmd

import (
	"github.com/spf13/cobra"
)

const Version = "unknown version"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of pochama",
	Long:  `All software has versions. This is pochama's`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("pochama version is %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
