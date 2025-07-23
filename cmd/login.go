package cmd

import (
    "fmt"
    "log"

    "github.com/spf13/cobra"
    "github.com/adityadeshmukh1/bookwyrm-cli/client"
)

var (
    loginUsername string
    loginPassword string
)

var loginCmd = &cobra.Command{
    Use:   "login",
    Short: "Login to Bookwyrm with username and password",
    Run: func(cmd *cobra.Command, args []string) {
        if loginUsername == "" || loginPassword == "" {
            log.Fatal("Username and password must be provided with --username and --password flags")
        }

        httpClient, err := client.LoginAndGetClient(loginUsername, loginPassword)
        if err != nil {
            log.Fatalf("Login failed: %v", err)
        }

        // Simple test: fetch home page to confirm session
        resp, err := httpClient.Get("https://bookwyrm.social/home")
        if err != nil {
            log.Fatalf("Failed to fetch home page after login: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode == 200 {
            fmt.Printf("Login successful! Welcome, %s.\n", loginUsername)
        } else {
            fmt.Printf("Login might have failed, status code: %d\n", resp.StatusCode)
        }
    },
}

func init() {
    rootCmd.AddCommand(loginCmd)

    loginCmd.Flags().StringVarP(&loginUsername, "username", "u", "", "Bookwyrm username")
    loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "Bookwyrm password")
}
