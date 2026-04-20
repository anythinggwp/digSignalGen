package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dsg",
	Short: "util for gen digital signal",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Init() {
	var err error
	rootCmd.AddCommand(genCmd)

	if err = rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
