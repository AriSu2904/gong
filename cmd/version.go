package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const appVersion = "1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gong version: %s\n\n", appVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
