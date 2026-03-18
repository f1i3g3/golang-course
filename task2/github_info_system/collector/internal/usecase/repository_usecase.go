package usecase

import (
	"local/github_info_system/collector/internal/domain"
)

type GitHubRepo interface {
	GetRepository(owner, repo string) (*domain.Repository, error)
}

type RepositoryUseCase struct {
	githubRepo GitHubRepo
}

func NewRepositoryUseCase(githubRepo GitHubRepo) *RepositoryUseCase {
	return &RepositoryUseCase{
		githubRepo: githubRepo,
	}
}

func (uc *RepositoryUseCase) GetRepositoryInfo(owner, repo string) (*domain.Repository, error) {
	if owner == "" || repo == "" {
		return nil, domain.ErrInvalidInput
	}
	
	return uc.githubRepo.GetRepository(owner, repo)
}
