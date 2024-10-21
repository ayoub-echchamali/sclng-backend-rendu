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
	Name         string `json:"name"`
	FullName     string `json:"full_name"`
	URL          string `json:"url"`
	Owner        Owner  `json:"owner"`
	LanguagesURL string `json:"languages_url"`
	License      string `json:"license"`
}

type Repositories []Repository

type ErrorResponse struct {
	Message         string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

// LANGUAGES
