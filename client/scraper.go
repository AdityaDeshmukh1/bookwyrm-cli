package client

import (
    "errors"
    "fmt"
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "strings"

    "github.com/PuerkitoBio/goquery"
)

// LoginAndGetClient logs into Bookwyrm using username and password,
// returns an authenticated HTTP client with session cookies.
func LoginAndGetClient(username, password string) (*http.Client, error) {
    baseURL := "https://bookwyrm.social"
    loginURL := baseURL + "/login"

    jar, err := cookiejar.New(nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create cookie jar: %w", err)
    }

    client := &http.Client{Jar: jar}

    // Step 1: GET login page to fetch CSRF token and cookies
    resp, err := client.Get(loginURL)
    oldURL := loginURL

    if err != nil {
        return nil, fmt.Errorf("GET login page error: %w", err)
    }
    defer resp.Body.Close()

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to parse login page HTML: %w", err)
    }

    csrfToken, exists := doc.Find("input[name='csrfmiddlewaretoken']").Attr("value")
    if !exists {
        return nil, errors.New("csrfmiddlewaretoken not found on login page")
    }

    // Step 2: Prepare login form data
    formData := url.Values{}
    formData.Set("localname", username)
    formData.Set("password", password)
    formData.Set("csrfmiddlewaretoken", csrfToken)

    // Step 3: Create POST request for login
    req, err := http.NewRequest("POST", loginURL, strings.NewReader(formData.Encode()))
    if err != nil {
        return nil, fmt.Errorf("failed to create POST request: %w", err)
    }

    // Required headers for CSRF validation
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Referer", loginURL)

    // Step 4: Perform login POST
    resp, err = client.Do(req)
    currURL := resp.Request.URL.String()
    if err != nil {
        return nil, fmt.Errorf("login POST request failed: %w", err)
    }
    defer resp.Body.Close()

    // Clever hack: if login is successfull, redirects to /
    // else stays on the same page
    if oldURL == currURL {
        return nil, fmt.Errorf("login failed, status code: %d", resp.StatusCode)
    }

    // Client now has the session cookie set
    return client, nil
}

