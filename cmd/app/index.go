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

var indexCmd = &cobra.Command{
	Use:   "index redis://[:<password>@]<host>:<port>[/<dbIndex>]",
	Short: "Scan keys and save prefix tree in a temporary file for further rendering with display command",
	Long:  "Keep in mind that some options are scanning (index) options that cannot be redefined later. For example, `maxChildren` changes the way index data is built, unlike `depth` parameter only influencing rendering",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		consoleLogger := logger.NewConsoleLogger(logLevel)
		consoleLogger.Info().Msg("Start indexing")

		rsvc := adapter.NewRedisService(consoleLogger)
		if isCluster {
			client := redis.NewClusterClient(&redis.ClusterOptions{
				Addrs: []string{args[0]},
			})

			err := client.Ping(context.Background())
			if err != nil {
				consoleLogger.Fatal().Err(err.Err()).Msg("Can't create redis cluster client")
			}

			rsvc.SetClusterClient(*client)
		} else {
			client := redis.NewClient(&redis.Options{
				Addr: args[0],
			})

			err := client.Ping(context.Background())
			if err != nil {
				consoleLogger.Fatal().Err(err.Err()).Msg("Can't create redis client")
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
	indexCmd.Flags().BoolVarP(&isCluster, "cluster", "C", false, "Enable cluster mode")
}
