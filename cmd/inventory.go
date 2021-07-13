package cmd

import (
	"context"
	"fmt"
	"github.com/mediocregopher/radix/v4"
	"github.com/obukhov/redis-inventory/src/logger"
	"github.com/obukhov/redis-inventory/src/scanner"
	"github.com/obukhov/redis-inventory/src/trie"
	"github.com/spf13/cobra"
	"log"
)

var scanCmd = &cobra.Command{
	Use:   "inventory [sourceHost:port]",
	Short: "",
	Long:  "",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Start scanning")

		clientSource, err := (radix.PoolConfig{}).New(context.Background(), "tcp", args[0])
		if err != nil {
			log.Fatal(err)
		}

		redisScanner := scanner.NewScanner(
			clientSource,
			logger.NewConsoleLogger(),
		)

		result := trie.NewTrie(trie.NewPunctuationSplitter(':'), 10)
		redisScanner.Scan(scanner.ScanOptions{ScanCount: 1000}, result)

		fmt.Println("Finish scanning")
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
