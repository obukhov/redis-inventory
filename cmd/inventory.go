package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mediocregopher/radix/v4"
	"github.com/obukhov/redis-inventory/src/logger"
	"github.com/obukhov/redis-inventory/src/scanner"
	"github.com/obukhov/redis-inventory/src/trie"
	"github.com/obukhov/redis-inventory/src/trieoutput"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var output string

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

		resultTrie := trie.NewTrie(trie.NewPunctuationSplitter(':'), 10)
		redisScanner.Scan(scanner.ScanOptions{ScanCount: 1000}, resultTrie)

		switch output {

		case "table":
			trieoutput.NewTableTrieOutput(os.Stdout, 10, "").Render(resultTrie)
		case "json":
			j, _ := json.Marshal(resultTrie.Root())
			fmt.Println(string(j))

		case "jsonp":
			j, _ := json.MarshalIndent(resultTrie.Root(), "", " ")
			fmt.Println(string(j))
		default:
			panic("Unknown output format: " + output)
		}

		fmt.Println("Finish scanning")
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&output, "output", "o", "table", "One of possible outputs: json, jsonp, table")
}
