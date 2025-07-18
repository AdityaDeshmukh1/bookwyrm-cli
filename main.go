
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
)

// Define minimal structs to parse the JSON
type ShelfPage struct {
    OrderedItems []Book `json:"orderedItems"`
}

type Book struct {
    ID   string `json:"id"`
    Type string `json:"type"`
    Title string `json:"title"`
    // Add more fields if you want, e.g. author, cover, etc.
}

func fetchShelf(shelfURL string) (*ShelfPage, error) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", shelfURL, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Accept", "application/activity+json")

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var page ShelfPage
    if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
        return nil, err
    }
    return &page, nil
}


func fetchBookTitle(bookURL string) (string, error) {
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", bookURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/activity+json")

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode JSON response into BookDetails struct
	var book Book
	if err := json.NewDecoder(resp.Body).Decode(&book); err != nil {
		return "", err
	}

	return book.Title, nil
}

func main() {
    if len(os.Args) < 2 {
    	fmt.Println("Usage: go run main.go <username>")
    	return
    }
    
    user := os.Args[1]
    shelves := map[string]string{
        "To Read":        fmt.Sprintf("https://bookwyrm.social/user/%s/books/to-read?page=1", user),
        "Reading":        fmt.Sprintf("https://bookwyrm.social/user/%s/books/reading?page=1", user),
        "Read":           fmt.Sprintf("https://bookwyrm.social/user/%s/books/read?page=1", user),
        "Stopped Reading": fmt.Sprintf("https://bookwyrm.social/user/%s/books/stopped-reading?page=1", user),
    }

    for name, url := range shelves {
        fmt.Printf("\n%s Shelf:\n", name)
        page, err := fetchShelf(url)
        if err != nil {
            log.Printf("Error fetching %s shelf: %v\n", name, err)
            continue
        }
        for _, book := range page.OrderedItems {
        	title, err := fetchBookTitle(book.ID)
        	if err != nil {
        		log.Printf("Error fetching the book title for %s: %v\n", book.ID, err)
        		continue
        	}
            fmt.Printf("  - %s (ID: %s)\n", title, book.ID)
        }
    }
}

