package domain

type ServiceError struct {
	Code  int
	Error error

	EmbeddedError error `json:",omitempty"`
}

func (se ServiceError) Embed(err error) *ServiceError {
	serr := se
	serr.EmbeddedError = err
	return &serr
}