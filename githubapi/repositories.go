package githubapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func FetchRepos() (result Repositories, err error) {
	url := "https://api.github.com/repositories"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error when creating a new request: %v", err)
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Could not read response body %v with error: %v", resp.Body, err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("Failed to parse JSON response: %v", err)
	}

	return
}

// func FetchLanguages(owner string, repo string) (result )