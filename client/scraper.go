package client

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "os"
    "path/filepath"
    "strings"

    "github.com/PuerkitoBio/goquery"
)

const (
    cookieFilePath = ".bookwyrm/cookies.json"
    bookwyrmDomain = "bookwyrm.social"
)

func cookieStoragePath() (string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    dir := filepath.Join(home, ".bookwyrm")
    os.MkdirAll(dir, 0700)
    return filepath.Join(dir, "cookies.json"), nil
}

// LoginAndGetClient logs into Bookwyrm and saves cookies for future use
func LoginAndGetClient(username, password string) (*http.Client, error) {
    baseURL := "https://bookwyrm.social"
    loginURL := baseURL + "/login"

    jar, err := cookiejar.New(nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create cookie jar: %w", err)
    }

    client := &http.Client{Jar: jar}

    // Step 1: Fetch CSRF token
    resp, err := client.Get(loginURL)
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

    // Step 2: Prepare and send login POST
    formData := url.Values{}
    formData.Set("localname", username)
    formData.Set("password", password)
    formData.Set("csrfmiddlewaretoken", csrfToken)

    req, err := http.NewRequest("POST", loginURL, strings.NewReader(formData.Encode()))
    if err != nil {
        return nil, fmt.Errorf("failed to create POST request: %w", err)
    }

    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Referer", loginURL)

    resp, err = client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("login POST failed: %w", err)
    }
    defer resp.Body.Close()

    // Detect failed login
    if resp.Request.URL.Path == "/login" {
        return nil, errors.New("login failed, possibly incorrect credentials")
    }

    // Save cookies
    err = saveCookies(jar)
    if err != nil {
        return nil, fmt.Errorf("failed to save cookies: %w", err)
    }

    return client, nil
}

// GetAuthenticatedClient loads cookies from disk to reuse the session
func GetAuthenticatedClient() (*http.Client, error) {
    jar, err := cookiejar.New(nil)
    if err != nil {
        return nil, err
    }

    client := &http.Client{Jar: jar}

    err = loadCookies(jar)
    if err != nil {
        return nil, err
    }

    return client, nil
}

func saveCookies(jar *cookiejar.Jar) error {
    u, _ := url.Parse("https://" + bookwyrmDomain)
    cookies := jar.Cookies(u)

    data, err := json.MarshalIndent(cookies, "", "  ")
    if err != nil {
        return err
    }

    path, err := cookieStoragePath()
    if err != nil {
        return err
    }

    return os.WriteFile(path, data, 0600)
}

func loadCookies(jar *cookiejar.Jar) error {
    path, err := cookieStoragePath()
    if err != nil {
        return err
    }

    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }

    var cookies []*http.Cookie
    if err := json.Unmarshal(data, &cookies); err != nil {
        return err
    }

    u, _ := url.Parse("https://" + bookwyrmDomain)
    jar.SetCookies(u, cookies)

    return nil
}

