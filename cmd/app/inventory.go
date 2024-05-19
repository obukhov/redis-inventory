package app

import (
	"context"
	"os"

	"github.com/Perseus/redis-inventory/src/adapter"
	"github.com/Perseus/redis-inventory/src/logger"
	"github.com/Perseus/redis-inventory/src/renderer"
	"github.com/Perseus/redis-inventory/src/scanner"
	"github.com/Perseus/redis-inventory/src/trie"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "inventory redis://[:<password>@]<host>:<port>[/<dbIndex>]",
	Short: "Scan keys and display summary right away with selected output and output params",
	Long:  "Scan command builds prefix tree in memory and then displays the usage summary. To avoid scanning redis instance when trying different output formats use `index` and `display` commands",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		consoleLogger := logger.NewConsoleLogger(logLevel)
		consoleLogger.Info().Msg("Start scanning")

		rsvc := adapter.NewRedisService(consoleLogger)
		if isCluster {
			client := redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:    []string{args[0]},
				Password: password,
			})

			_, err := client.Ping(context.Background()).Result()
			if err != nil {
				consoleLogger.Fatal().Err(err).Msg("Can't create redis cluster client")
			}

			rsvc.SetClusterClient(*client)
		} else {
			client := redis.NewClient(&redis.Options{
				Addr:     args[0],
				Password: password,
			})

			_, err := client.Ping(context.Background()).Result()
			if err != nil {
				consoleLogger.Fatal().Err(err).Msg("Can't create redis client")
			}

			rsvc.SetClient(*client)
		}

		redisScanner := scanner.NewScanner(
			rsvc,
			adapter.NewPrettyProgressWriter(os.Stdout),
			consoleLogger,
		)

		resultTrie := trie.NewTrie(trie.NewPunctuationSplitter([]rune(separators)...), maxChildren)
		redisScanner.Scan(
			adapter.ScanOptions{
				ScanCount: scanCount,
				Pattern:   pattern,
				Throttle:  throttleNs,
				IsCluster: isCluster,
			},
			resultTrie,
		)

		r, err := renderer.NewRenderer(output, outputParams, consoleLogger)
		if err != nil {
			consoleLogger.Fatal().Err(err).Msg("Can't create renderer")
		}

		err = r.Render(resultTrie.Root())
		if err != nil {
			consoleLogger.Fatal().Err(err).Msg("Can't render report")
		}

		consoleLogger.Info().Msg("Finish scanning")
	},
}

func init() {
	RootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&output, "output", "o", "table", "One of possible outputs: json, jsonp, table")
	scanCmd.Flags().StringVarP(&password, "password", "a", "", "Password for redis instance")
	scanCmd.Flags().StringVarP(&outputParams, "output-params", "p", "", "Parameters specific for output type")
	scanCmd.Flags().StringVarP(&logLevel, "logLevel", "l", "info", "Level of logs to be displayed")
	scanCmd.Flags().StringVarP(&separators, "separators", "s", ":", "Symbols that logically separate levels of the key")
	scanCmd.Flags().IntVarP(&maxChildren, "maxChildren", "m", 10, "Maximum children node can have before start aggregating")
	scanCmd.Flags().StringVarP(&pattern, "pattern", "k", "*", "Glob pattern limiting the keys to be aggregated")
	scanCmd.Flags().IntVarP(&scanCount, "scanCount", "c", 1000, "Number of keys to be scanned in one iteration (argument of scan command)")
	scanCmd.Flags().IntVarP(&throttleNs, "throttle", "t", 0, "Throttle: number of nanoseconds to sleep between keys")
	scanCmd.Flags().BoolVarP(&isCluster, "cluster", "C", false, "Enable cluster mode")
}
