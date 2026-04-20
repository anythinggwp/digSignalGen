package cmd

import (
	"github.com/anythinggwp/digSignalGen/internal"
	"github.com/spf13/cobra"
)

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

func init() {
	genCmd.Flags().String("alpha", "0.2|0.7", "setup alpha's by input alpha1|alpha2 float64 type")
	genCmd.Flags().Uint64("length", 100, "setup wave length")
	genCmd.Flags().Uint64("parts", 1, "setup wave parts")
	genCmd.Flags().String("init-cond", "0.6|-0.1", "setup started x's for generated wave x[n-1]|x[n-2]")
	genCmd.Flags().String("save-file", "", "path for saving resul")
	genCmd.Flags().Bool("disable-output", false, "disable gui output")
	genCmd.Flags().Bool("growing-graph", false, "generate growing wave")
	genCmd.Flags().Bool("decrease-graph", false, "generate decreasing wave")
}
