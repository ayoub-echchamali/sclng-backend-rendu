package api

import (
	"encoding/json"
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

func (s *ApiServer) getRepos(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	languageParam := r.URL.Query().Get("languages")

	filterByLanguage := false

	if languageParam != "" {
		filterByLanguage = true
	}

	repos, err := githubapi.FetchRepos(s.Config.GithubToken)
	if err != nil {
		message := "Error while fetching repositories"
		log.WithError(err).Error(message)
		util.RespondWithError(w, http.StatusInternalServerError, message)
		return err
	}

	var reposDto dto.RepositoriesDto

	workerCount := 100
	resultChan := make(chan dto.RepositoryDto, workerCount)
	defer close(resultChan)

	var mu sync.Mutex

	var remainingRequests int
	var resetTime time.Time
	var rateLimitHit bool

	var wg sync.WaitGroup

	for _, repo := range repos {
		wg.Add(1)

		go func(repo githubapi.Repository) {
			match := false

			defer wg.Done()

			repoDto := dto.RepositoryDto{
				FullName:   repo.FullName,
				Owner:      repo.Owner.Login,
				Repository: repo.Name,
				Languages:  make(map[string]dto.LanguageDto),
			}

			languages, err := githubapi.FetchLanguagesWithRateLimit(repo.Owner.Login, repo.Name, s.Config.GithubToken, &remainingRequests, &resetTime, &mu, &rateLimitHit)

			if err == nil {
				for lang, bytes := range languages {
					if filterByLanguage && !match && strings.EqualFold(lang,languageParam) {
						match = true
					}
					repoDto.Languages[lang] = dto.LanguageDto{Bytes: bytes}
				}
			} else {
				log.Errorf("Failed to fetch language for %s/%s with error: %v", repo.Owner.Login, repo.FullName, err)
			}

			if filterByLanguage {
				if match {
					resultChan <- repoDto
				}
			} else {
				resultChan <- repoDto
			}

		}(repo)
	}

	go func() {
		for repoDto := range resultChan {
			reposDto.Repositories = append(reposDto.Repositories, repoDto)
		}
	}()

	wg.Wait()

	util.RespondWithJSON(w, http.StatusOK, reposDto)

	return nil
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
