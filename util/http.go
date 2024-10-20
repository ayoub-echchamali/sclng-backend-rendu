package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func SendRequest(method, url string, headers map[string]string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return &http.Response{}, fmt.Errorf("error when creating a new request: %w", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return response, fmt.Errorf("error making request: %w", err)
	}

	// Check for HTTP error responses
	if response.StatusCode != http.StatusOK {
		return response, fmt.Errorf("received non-200 response: %s", response.Status)
	}

	return response, nil
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error":message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}