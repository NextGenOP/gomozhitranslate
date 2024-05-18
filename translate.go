package main

import (
    "encoding/json"
    "errors"
    "math/rand"
    "net/http"
    "net/url"
    "time"
)

// Translator represents a translation service.
type Translator struct {
    URLs   []string
    Engine string
}

// NewTranslator creates a new Translator instance.
func NewTranslator(urls []string, engine string) *Translator {
    return &Translator{
        URLs:   urls,
        Engine: engine,
    }
}

// buildURL constructs the API endpoint URL with query parameters.
func (t *Translator) buildURL(endpoint string, queryParams map[string]string) string {
    // Randomly select a URL from the list
    rand.Seed(time.Now().UnixNano())
    currentURL := t.URLs[rand.Intn(len(t.URLs))]

    // Create a URL structure
    u, err := url.Parse(currentURL)
    if err != nil {
        return ""
    }

    // Set the engine query parameter
    q := u.Query()
    q.Set("engine", t.Engine)

    // Set other query parameters
    for key, value := range queryParams {
        q.Set(key, value)
    }

    // Set the path to the API endpoint
    u.Path = "/api/" + endpoint

    // Encode the query parameters and return the full URL
    u.RawQuery = q.Encode()
    return u.String()
}

// makeRequest performs an HTTP GET request and handles errors.
func (t *Translator) makeRequest(url string) (*http.Response, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    if resp.StatusCode != http.StatusOK {
        return nil, errors.New("HTTP request failed with status code: " + resp.Status)
    }
    return resp, nil
}

// Languages retrieves a map of supported languages.
func (t *Translator) Languages() (map[string]string, error) {
    url := t.buildURL("source_languages", nil)
    resp, err := t.makeRequest(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var allLanguages []struct {
        Name string `json:"Name"`
        Id   string `json:"Id"`
    }
    err = json.NewDecoder(resp.Body).Decode(&allLanguages)
    if err != nil {
        return nil, err
    }

    lang := make(map[string]string)
    for _, language := range allLanguages {
        lang[language.Name] = language.Id
    }
    return lang, nil
}

// Translate performs text translation from source to target language.
func (t *Translator) Translate(source, target, text string) (string, error) {
    queryParams := map[string]string{
        "from": source,
        "to":   target,
        "text": text,
    }

    url := t.buildURL("translate", queryParams)
    resp, err := t.makeRequest(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var result struct {
        TranslatedText string `json:"translated-text"`
    }
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        return "", err
    }

    return result.TranslatedText, nil
}

