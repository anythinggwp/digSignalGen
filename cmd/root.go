package cmd

import (
	"os"

	"github.com/anythinggwp/digSignalGen/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dsg",
	Short: "util for gen digital signal",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "util for gen digital signal",
	RunE: func(cmd *cobra.Command, args []string) error {
		builder, err := internal.NewWaveBuilder(cmd)
		if err != nil {
			return err
		}
		return builder.BuildWave()
	},
}

func Init() {
	var err error
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().String("alpha", "", "")

	if err = rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
