package cmd

import (
	"context"
	"github.com/mediocregopher/radix/v4"
	"github.com/obukhov/redis-inventory/src/logger"
	"github.com/obukhov/redis-inventory/src/renderer"
	"github.com/obukhov/redis-inventory/src/scanner"
	"github.com/obukhov/redis-inventory/src/trie"
	"github.com/spf13/cobra"
)

var output, outputParams string

var scanCmd = &cobra.Command{
	Use:   "inventory [sourceHost:port]",
	Short: "",
	Long:  "",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		consoleLogger := logger.NewConsoleLogger()
		consoleLogger.Info().Msg("Start scanning")

		clientSource, err := (radix.PoolConfig{}).New(context.Background(), "tcp", args[0])
		if err != nil {
			consoleLogger.Fatal().Err(err).Msg("Can't create redis client")
		}

		redisScanner := scanner.NewScanner(clientSource, consoleLogger)

		resultTrie := trie.NewTrie(trie.NewPunctuationSplitter(':'), 10)
		redisScanner.Scan(scanner.ScanOptions{ScanCount: 1000}, resultTrie)

		r, err := renderer.NewRenderer(output, outputParams)
		if err != nil {
			consoleLogger.Fatal().Err(err).Msg("Can't create renderer")
		}

		err = r.Render(resultTrie)
		if err != nil {
			consoleLogger.Fatal().Err(err).Msg("Can't render report")
		}

		consoleLogger.Info().Msg("Finish scanning")
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&output, "output", "o", "table", "One of possible outputs: json, jsonp, table")
	scanCmd.Flags().StringVarP(&outputParams, "output-params", "p", "", "Parameters specific for output type")
}
