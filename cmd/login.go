
package cmd

import (
    "fmt"
    "io"
    "log"
    "strings"

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

        // Fetch home page to confirm login
        resp, err := httpClient.Get("https://bookwyrm.social/home")
        if err != nil {
            log.Fatalf("Failed to fetch home page after login: %v", err)
        }
        defer resp.Body.Close()

        bodyBytes, err := io.ReadAll(resp.Body)
        if err != nil {
            log.Fatalf("Failed to read home page content: %v", err)
        }

        bodyStr := string(bodyBytes)

        if strings.Contains(strings.ToLower(bodyStr), "logout") {
            fmt.Printf("✅ Login successful! Welcome, %s.\n", loginUsername)
        } else {
            fmt.Printf("❌ Login may have failed — 'logout' not found in page. Status code: %d\n", resp.StatusCode)
        }
    },
}

func init() {
    rootCmd.AddCommand(loginCmd)

    loginCmd.Flags().StringVarP(&loginUsername, "username", "u", "", "Bookwyrm username")
    loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "Bookwyrm password")
}

