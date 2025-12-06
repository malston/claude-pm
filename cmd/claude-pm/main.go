// ABOUTME: Entry point for the claude-pm CLI tool
// ABOUTME: Initializes and executes the root command
package main

import (
	"fmt"
	"os"

	"github.com/malston/claude-pm/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
