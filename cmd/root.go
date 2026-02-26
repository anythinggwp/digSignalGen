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

	genCmd.Flags().String("alpha", "0.2|0.7", "setup alpha's by input alpha1|alpha2 float64 type")
	genCmd.Flags().Uint64("length", 100, "setup wave length")
	genCmd.Flags().String("init-cond", "0.6|-0.1", "setup started x's for generated wave x[n-1]|x[n-2]")

	if err = rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
