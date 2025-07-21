package models

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
