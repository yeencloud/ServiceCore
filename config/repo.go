package config

type Repository struct {
	Repository string
	URL        string
}

func (config *Config) GetRepository() *Repository {
	if config.repository == nil {
		repository := config.GetEnvStringOrDefault("GITHUB_REPOSITORY", "")
		repositoryURL := config.GetEnvStringOrDefault("GITHUB_REPOSITORY_URL", "")

		repo := &Repository{
			Repository: repository,
			URL:        repositoryURL,
		}

		config.repository = repo

		return repo
	}

	return config.repository
}