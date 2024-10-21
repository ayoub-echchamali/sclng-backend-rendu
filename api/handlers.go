package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/sclng-backend-test-v1/dto"
	"github.com/Scalingo/sclng-backend-test-v1/githubapi"
	"github.com/Scalingo/sclng-backend-test-v1/util"
	log "github.com/sirupsen/logrus"
)

func (s *ApiServer) getGithubPublicRepositories(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	// this returns an array of strings
	// typical url is /repos?language=ruby&language=javascript
	languagesParam := r.URL.Query()["language"]
	licensesParam := r.URL.Query()["license"]

	filterByLanguage := false
	filterByLicense := false

	if len(languagesParam) > 0 {
		filterByLanguage = true
	}
	if len(licensesParam) > 0 {
		filterByLicense = true
	}

	// Fetching the 100 repos
	repos, err := githubapi.FetchPublicRepos(s.Config.GithubToken)
	if err != nil {
		message := "Error while fetching repositories"
		log.Error(fmt.Sprintf("%v: %v", message, err))
		// If error, respond with a 500
		util.RespondWithError(w, http.StatusInternalServerError, message)
		return err
	}

	// response dto
	var reposDto dto.RepositoriesDto

	// Github concurrent rate limiting is 100, we will specify then 100 workers
	workerCount := 60
	if s.Config.GithubToken != "" {
		workerCount = 100
	}

	// each worker will push the result of his work into a channel 
	resultChan := make(chan dto.RepositoryDto, workerCount)
	defer close(resultChan)

	// These variables help keep track of the api calls rate 
	var mu sync.Mutex
	var remainingRequests int
	var resetTime time.Time
	var rateLimitHit bool
	count := 0

	var wg sync.WaitGroup

	for _, repo := range repos {
		wg.Add(1)

		go func(repo githubapi.Repository) {
			// This flag helps check if at least one language matches the query params
			matchLanguage := false
			matchLicense := false

			defer wg.Done()

			repoDto := dto.RepositoryDto{
				FullName:   repo.FullName,
				Owner:      repo.Owner.Login,
				Repository: repo.Name,
				URL:        repo.URL,
				License:    repo.License,
				Languages:  make(map[string]dto.LanguageDto),
			}

			languages, err := githubapi.FetchRepoLanguages(repo.Owner.Login, repo.Name, s.Config.GithubToken, &remainingRequests, &resetTime, &mu, &rateLimitHit)

			if err == nil {
				/* 
					Here we do two things:
					- if language filtering is enabled, we loop through the languages specified to see if there is a match
					- build the dto for the response
				*/
				for lang, bytes := range languages {
					if filterByLanguage {
						for _, langParam := range languagesParam{
							if !matchLanguage && strings.EqualFold(lang,langParam) {
								matchLanguage = true
								break
							}
						}
					}
					repoDto.Languages[lang] = dto.LanguageDto{Bytes: bytes}
				}
				if filterByLicense {
					for _, license := range licensesParam{
						if !matchLicense && strings.EqualFold(license, repo.License) {
							matchLicense = true
							break
						}
					}
				}
			} else {
				log.Errorf("Failed to fetch language for %s/%s with error: %v", repo.Owner.Login, repo.FullName, err)
			}

			/* 
				If filtering is enabled, and we have a match, the dto is pushed into the channel
				If no match, the repo is ignored.
				If filtering is not enabled, the repo is directly pushed.
			*/
			if shouldIncludeRepo(filterByLanguage, filterByLicense, matchLanguage, matchLicense) {
				resultChan <- repoDto
			}
		}(repo)
	}

	// This routine reads from the channel to build the final dto
	go func() {
		for repoDto := range resultChan {
			reposDto.Repositories = append(reposDto.Repositories, repoDto)
			count++
		}
	}()

	wg.Wait()

	reposDto.TotalItems = count

	util.RespondWithJSON(w, http.StatusOK, reposDto)

	// defer close channel when function exits
	return nil
}

func shouldIncludeRepo(filterByLanguage, filterByLicense, matchLanguage, matchLicense bool) bool {
	// if any filter is enabled
	if filterByLanguage || filterByLicense {
		// return matchLanguage if filterByLanguage is enabled and filterByLicense is not
		if filterByLanguage && !filterByLicense {
			return matchLanguage
		}
		if filterByLanguage && filterByLicense {
			return matchLanguage && matchLicense
		}
		// return matchLicense if filterByLicense is enabled and filterByLanguage is not
		if !filterByLanguage && filterByLicense {
			return matchLicense
		}
	}
	// if no filtering is specified, include repo
	return true
}

func pongHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(map[string]string{"status": "pong"})
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
	}
	return nil
}
