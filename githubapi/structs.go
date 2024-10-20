package githubapi

// REPOS

type Owner struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
	URL   string `json:"url"`
	Type  string `json:"type"`
}

type Repository struct {
	ID           int    `json:"id"`
	NodeID       string `json:"node_id"`
	Name         string `json:"name"`
	FullName     string `json:"full_name"`
	Private      bool   `json:"private"`
	Owner        Owner  `json:"owner"`
	HTMLURL      string `json:"html_url"`
	Description  string `json:"description"`
	Fork         bool   `json:"fork"`
	URL          string `json:"url"`
	LanguagesURL string `json:"languages_url"`
}

type Repositories []Repository

type ErrorResponse struct {
	Message         string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

// LANGUAGES
