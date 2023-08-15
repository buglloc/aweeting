package main

import (
	"fmt"
	"os"

	_ "go.uber.org/automaxprocs"

	"github.com/buglloc/aweeting/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "aweeting: %v\n", err)
		os.Exit(1)
	}
}
