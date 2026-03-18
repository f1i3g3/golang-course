package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"local/github_info_system/collector/internal/domain"
)

type Client struct {
	httpClient *http.Client
}

type GithubRepoResponse struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Stars       int       `json:"stargazers_count"`
	Forks       int       `json:"forks_count"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (c *Client) GetRepository(owner, repo string) (*domain.Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Network error: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var githubRepo GithubRepoResponse
		if err := json.NewDecoder(resp.Body).Decode(&githubRepo); err != nil {
			return nil, fmt.Errorf("Json parsing error: %w", err)
		}
		
		return &domain.Repository{
			Name:        githubRepo.Name,
			Description: githubRepo.Description,
			Stars:       githubRepo.Stars,
			Forks:       githubRepo.Forks,
			CreatedAt:   githubRepo.CreatedAt,
		}, nil
		
	case http.StatusNotFound:
		return nil, domain.ErrRepositoryNotFound
		
	default:
		return nil, fmt.Errorf("%w: status %d", domain.ErrGitHubAPI, resp.StatusCode)
	}
}
