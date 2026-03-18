package domain

import "errors"

var (
	ErrRepositoryNotFound = errors.New("Repository's not found")
	ErrGitHubAPI          = errors.New("Github api error")
	ErrInvalidInput       = errors.New("Invalid input: owner and repo are required")
)
