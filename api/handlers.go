package api

import (
	"encoding/json"
	"net/http"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/sclng-backend-test-v1/dto"
	"github.com/Scalingo/sclng-backend-test-v1/githubapi"
	log "github.com/sirupsen/logrus"
)

func (s *ApiServer) getRepos(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	repos, err := githubapi.FetchRepos()
	if err != nil {
		log.Fatalf("Error while fetching for repositories: %v", err)
	}

	var reposDto dto.RepositoriesDto

	count := 0

	for _, repo := range repos {
		repoDto := dto.RepositoryDto{
			FullName: repo.FullName,
			Owner: repo.Owner.Login,
			Repository: repo.Name,
			Languages: make(map[string]dto.LanguageDto),
		}

		// languages, err := FetchLanguages(repo.LanguagesURL, s.GithubToken)
		// if err == nil {
		// 	for lang, bytes := range languages {
		// 		repoDto.Languages[lang] = dto.LanguageDto{Bytes: bytes}
		// 	}
		// }

		reposDto.Repositories = append(reposDto.Repositories, repoDto)
		count++
	}

	log.Infof("Repos: %d", count)

	response, _ := json.Marshal(reposDto)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

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