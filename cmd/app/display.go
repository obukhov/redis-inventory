package app

import (
	"encoding/json"
	"os"

	"github.com/obukhov/redis-inventory/src/logger"
	"github.com/obukhov/redis-inventory/src/renderer"
	"github.com/obukhov/redis-inventory/src/trie"
	"github.com/spf13/cobra"
)

var displayCmd = &cobra.Command{
	Use:   "display",
	Short: "Reads cached usage data created with `index` command and displays according to output and output params",
	Long:  "It can be useful to play with different display params, for example: depth, padding, human readable formatting",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		consoleLogger := logger.NewConsoleLogger(logLevel)
		consoleLogger.Info().Msg("Loading index")

		f, err := os.Open(indexFileName)
		if err != nil {
			consoleLogger.Fatal().Err(err).Msg("Can't create renderer")
		}

		var root *trie.Node
		dec := json.NewDecoder(f)
		err = dec.Decode(&root)
		if err != nil {
			consoleLogger.Fatal().Err(err).Msg("Can't decode file")
		}

		r, err := renderer.NewRenderer(output, outputParams, consoleLogger)
		if err != nil {
			consoleLogger.Fatal().Err(err).Msg("Can't create renderer")
		}

		err = r.Render(root)
		if err != nil {
			consoleLogger.Fatal().Err(err).Msg("Can't render report")
		}

		consoleLogger.Info().Msg("Done")
	},
}

func init() {
	RootCmd.AddCommand(displayCmd)
	displayCmd.Flags().StringVarP(&indexFileName, "indexFile", "i", "./redis-inventory.json", "Name of the file where the data is saved")
	displayCmd.Flags().StringVarP(&output, "output", "o", "table", "One of possible outputs: json, jsonp, table")
	displayCmd.Flags().StringVarP(&outputParams, "output-params", "p", "", "Parameters specific for output type")
	displayCmd.Flags().StringVarP(&logLevel, "logLevel", "l", "info", "Level of logs to be displayed")
}
