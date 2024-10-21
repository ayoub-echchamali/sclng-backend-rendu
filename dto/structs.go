package dto

type RepositoryDto struct {
	FullName     string `json:"full_name"`
	URL          string `json:"url"`
	Owner        string `json:"owner"`
	Repository   string `json:"repository"`
	Languages    map[string]LanguageDto `json:"languages"`
	License      string `json:"license"`
}

type LanguageDto struct {
	Bytes int `json:"bytes"`
}

type RepositoriesDto struct {
	Repositories []RepositoryDto `json:"repositories"`
	TotalItems int `json:"total_items"`
}