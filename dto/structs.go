package dto

type RepositoryDto struct {
	FullName     string `json:"full_name"`
	Owner        string `json:"owner"`
	Repository   string `json:"repository"`
	Languages    map[string]LanguageDto `json:"languages"`
}

type LanguageDto struct {
	Bytes int `json:"bytes"`
}

type RepositoriesDto struct {
	Repositories []RepositoryDto `json:"repositories"`
}