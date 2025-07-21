
package main

import (
    "log"

    "github.com/adityadeshmukh1/bookwyrm-cli/cmd"
)

func main() {
    if err := cmd.Execute(); err != nil {
        log.Fatal(err)
    }
}

