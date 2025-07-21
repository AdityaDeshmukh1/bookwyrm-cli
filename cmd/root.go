package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "bookwyrm-cli",
    Short: "A CLI for interacting with BookWyrm instances",
    Long:  `bookwyrm-cli is a tool to view and manage BookWyrm shelves and books from the terminal.`,
}

func Execute() error {
    return rootCmd.Execute()
}
