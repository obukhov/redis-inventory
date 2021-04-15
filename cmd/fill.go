package cmd

import (
	"context"
	"fmt"
	"github.com/mediocregopher/radix/v4"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var cycles int

var fillCmd = &cobra.Command{
	Use:   "fill [host:port]",
	Short: "Create random keys in redis instance",
	Args:  cobra.MinimumNArgs(1),
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Filling redis with random data")

		keys := [][]string{
			{
				"dev",
				"prod",
			},
			{
				"users",
				"counters",
				"foobar",
			},
			{
				"foo",
				"bar",
				"hey",
			},
		}

		rand.Seed(time.Now().UTC().UnixNano())

		ctx := context.Background()
		redisClient, err := (radix.PoolConfig{}).New(ctx, "tcp", args[0])
		if err != nil {
			log.Fatal(err)
		}
		defer redisClient.Close()

		for i := 0; i < cycles; i++ {
			key := ""
			for n := 0; n < len(keys); n++ {
				key += keys[n][rand.Intn(len(keys[n]))] + ":"
			}
			key += strconv.Itoa(int(rand.Int31()))
			value := strings.Repeat("N", rand.Intn(1000))

			setErr := redisClient.Do(ctx, radix.Cmd(nil, "SET", key, value))
			if setErr != nil {
				log.Println(setErr)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(fillCmd)
	fillCmd.Flags().IntVar(&cycles, "cycles", 1000, "Cycles count to perform")
}
