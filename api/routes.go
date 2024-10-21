package api

func (s *ApiServer) CreateRoutes() {
	s.Router.HandleFunc("/ping", pongHandler)
	s.Router.HandleFunc("/publicGithubRepos", s.getGithubPublicRepositories)
}