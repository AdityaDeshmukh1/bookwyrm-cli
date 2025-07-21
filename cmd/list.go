package cmd

import (
    "fmt"
    "log"

    "github.com/spf13/cobra"
    "github.com/adityadeshmukh1/bookwyrm-cli/client"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
    Use:   "list [username]",
    Short: "List books from a user's public shelves",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        user := args[0]
        shelves := map[string]string{
            "To Read":         fmt.Sprintf("https://bookwyrm.social/user/%s/books/to-read?page=1", user),
            "Reading":         fmt.Sprintf("https://bookwyrm.social/user/%s/books/reading?page=1", user),
            "Read":            fmt.Sprintf("https://bookwyrm.social/user/%s/books/read?page=1", user),
            "Stopped Reading": fmt.Sprintf("https://bookwyrm.social/user/%s/books/stopped-reading?page=1", user),
        }

        for name, url := range shelves {
            fmt.Printf("\n%s Shelf:\n", name)
            page, err := client.FetchShelf(url)
            if err != nil {
                log.Printf("Error fetching %s shelf: %v\n", name, err)
                continue
            }

            for _, book := range page.OrderedItems {
                title, err := client.FetchBookTitle(book.ID)
                if err != nil {
                    log.Printf("Error fetching title for book %s: %v\n", book.ID, err)
                    continue
                }
                fmt.Printf("  - %s (ID: %s)\n", title, book.ID)
            }
        }
    },
}

func init() {
    rootCmd.AddCommand(listCmd)
}
