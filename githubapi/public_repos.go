package githubapi

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/Scalingo/sclng-backend-test-v1/util"
	log "github.com/sirupsen/logrus"
)

func FetchPublicRepos(token string) (Repositories, error) {
	// creating http request with necessary headers
	url := "https://api.github.com/repositories"
	headers := map[string]string{
		"Accept":                   "application/vnd.github+json",
		"X-GitHub-Api-Version":    "2022-11-28",
	}
	if token != "" {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", token)
	}
	response, err := util.SendRequest("GET", url, headers, nil)
	if err != nil {
		return Repositories{}, err
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return Repositories{}, fmt.Errorf("could not read response body: %w", err)
	}
	defer response.Body.Close()

	// First, check if the response might be an error
	var errorResponse ErrorResponse
	if err := json.Unmarshal(bodyBytes, &errorResponse); err == nil && errorResponse.Message != "" {
		return Repositories{}, fmt.Errorf("API error: %s (see: %s)", errorResponse.Message, errorResponse.DocumentationURL)
	}

	var result Repositories
	// If no error, unmarshal into the expected result
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return result, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return result, nil
}


func FetchRepoLanguages(owner, repo, token string, remainingRequests *int, resetTime *time.Time, mu *sync.Mutex, rateLimitHit *bool) (map[string]int, error) {
	// Check if rate limit has already been hit
	mu.Lock()
	if *rateLimitHit {
		waitDuration := time.Until(*resetTime)
		mu.Unlock()
		return nil, fmt.Errorf("rate limit already hit, please retry after %v at %s", waitDuration, resetTime.Format(time.RFC3339))
	}
	mu.Unlock()

	// If rate limit hasn't been hit yet, fetch languages from github api
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/languages", owner, repo)
	headers := map[string]string{
		"Accept":                   "application/vnd.github+json",
		"X-GitHub-Api-Version":    "2022-11-28",
	}

	if token != "" {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", token)
	}

	response, err := util.SendRequest("GET", url, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("error while sending github request for repository languages: %v", err)
	}

	// Update rate limit info
	mu.Lock()
	defer mu.Unlock()

	remaining, _ := strconv.Atoi(response.Header.Get("X-RateLimit-Remaining"))
	resetUnix, _ := strconv.ParseInt(response.Header.Get("X-RateLimit-Reset"), 10, 64)
	newResetTime := time.Unix(resetUnix, 0).UTC()

	*remainingRequests = remaining
	*resetTime = newResetTime

	// If rate limit has been hit, log warning and continue
	if *remainingRequests == 0 {
		*rateLimitHit = true
		waitDuration := time.Until(*resetTime)
		log.Warnf("rate limit exceeded, retry after %v (at %s)", waitDuration, resetTime.Format(time.RFC3339))
	}

	// parse result
	var result map[string]int
	var errorResponse ErrorResponse

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return map[string]int{}, fmt.Errorf("could not read response body: %w", err)
	}
	defer response.Body.Close()

	if err := json.Unmarshal(body, &errorResponse); err == nil && errorResponse.Message != "" {
		return result, fmt.Errorf("API error: %s (see: %s)", errorResponse.Message, errorResponse.DocumentationURL)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return result, nil
}
