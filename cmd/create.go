package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Main command to create a template or spesific module",
	Long:  `Main command to create a template or spesific module`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
