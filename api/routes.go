package api

func (s *ApiServer) createRoutes() {
	s.Router.HandleFunc("/ping", pongHandler)
	s.Router.HandleFunc("/repos", s.getRepos)
}