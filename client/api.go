package client
import (
    "encoding/json"
    "net/http"

    "github.com/adityadeshmukh1/bookwyrm-cli/models"
)

func FetchShelf(shelfURL string) (*models.ShelfPage, error) {
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

    var page models.ShelfPage
    if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
        return nil, err
    }
    return &page, nil
}


func FetchBookTitle(bookURL string) (string, error) {
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
	var book models.Book
	if err := json.NewDecoder(resp.Body).Decode(&book); err != nil {
		return "", err
	}

	return book.Title, nil
}
