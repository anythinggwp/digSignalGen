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

var genCmd = &cobra.Command{
	Use:   "dsg",
	Short: "util for gen digital signal",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func Init() {
	var err error
	rootCmd.AddCommand()

	if err = rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
