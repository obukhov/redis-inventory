package app

import (
	"context"
	"github.com/mediocregopher/radix/v4"
	"github.com/obukhov/redis-inventory/src/adapter"
	"github.com/obukhov/redis-inventory/src/logger"
	"github.com/obukhov/redis-inventory/src/renderer"
	"github.com/obukhov/redis-inventory/src/scanner"
	"github.com/obukhov/redis-inventory/src/trie"
	"github.com/spf13/cobra"
	"os"
)

var indexCmd = &cobra.Command{
	Use:   "index [sourceHost:port]",
	Short: "Scan keys space and save the index as temporary file for further display with display command",
	Long:  "Keep in mind that there are scanning options and displaying options, if the instance was indexed with maxChildren=10 it cannot be changed in display unlike depth parameter",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		consoleLogger := logger.NewConsoleLogger(logLevel)
		consoleLogger.Info().Msg("Start indexing")

		clientSource, err := (radix.PoolConfig{}).New(context.Background(), "tcp", args[0])
		if err != nil {
			consoleLogger.Fatal().Err(err).Msg("Can't create redis client")
		}

		redisScanner := scanner.NewScanner(
			adapter.NewRedisService(clientSource),
			adapter.NewPrettyProgressWriter(os.Stdout),
			consoleLogger,
		)

		resultTrie := trie.NewTrie(trie.NewPunctuationSplitter([]rune(separators)...), maxChildren)
		redisScanner.Scan(
			adapter.ScanOptions{
				ScanCount: scanCount,
				Pattern:   pattern,
				Throttle:  throttleNs,
			},
			resultTrie,
		)

		indexFileName := os.TempDir() + "/redis-inventory.json"
		f, err := os.Create(indexFileName)
		if err != nil {
			consoleLogger.Fatal().Err(err).Msg("Can't create renderer")
		}

		r := renderer.NewJSONRenderer(f, renderer.JSONRendererParams{})

		err = r.Render(resultTrie.Root())
		if err != nil {
			consoleLogger.Fatal().Err(err).Msg("Can't write to file")
		}

		consoleLogger.Info().Msgf("Finish scanning and saved index as a file %s", indexFileName)
	},
}

func init() {
	RootCmd.AddCommand(indexCmd)
	indexCmd.Flags().StringVarP(&logLevel, "logLevel", "l", "info", "Level of logs to be displayed")
	indexCmd.Flags().StringVarP(&separators, "separators", "s", ":", "Symbols that logically separate levels of the key")
	indexCmd.Flags().IntVarP(&maxChildren, "maxChildren", "m", 10, "Maximum children node can have before start aggregating")
	indexCmd.Flags().StringVarP(&pattern, "pattern", "k", "*", "Glob pattern limiting the keys to be aggregated")
	indexCmd.Flags().IntVarP(&scanCount, "scanCount", "c", 1000, "Number of keys to be scanned in one iteration (argument of scan command)")
	indexCmd.Flags().IntVarP(&throttleNs, "throttle", "t", 0, "Throttle: number of nanoseconds to sleep between keys")
}
