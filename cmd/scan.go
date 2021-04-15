package cmd

import (
	"context"
	"fmt"
	"github.com/mediocregopher/radix/v4"
	"github.com/obukhov/redis-inventory/src/scanner"
	"github.com/spf13/cobra"
	"log"
)

var scanCmd = &cobra.Command{
	Use:   "scan [sourceHost:port]",
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
		)

		redisScanner.Scan()

		fmt.Println("Finish scanning")
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
