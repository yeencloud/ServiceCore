package config

type Repository struct {
	Repository string
	URL        string
}

func (config *Config) getRepository() {
	if config.Repository != nil {
		return
	}

	repository := config.GetEnvStringOrDefault("GITHUB_REPOSITORY", "")
	repositoryURL := config.GetEnvStringOrDefault("GITHUB_REPOSITORY_URL", "")

	repo := &Repository{
		Repository: repository,
		URL:        repositoryURL,
	}

	config.Repository = repo
}
