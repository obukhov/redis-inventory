package main

import (
	"github.com/obukhov/redis-inventory/cmd/app"
	"github.com/spf13/cobra/doc"
	"log"
)

func main() {
	err := doc.GenMarkdownTree(app.RootCmd, "./docs/cobra")
	if err != nil {
		log.Fatal(err)
	}
}
