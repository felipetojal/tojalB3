package main

import (
	"log"

	"github.com/felipetojal/tojalB3/cmd"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	return cmd.Execute()
}
