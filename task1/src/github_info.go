package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Repository struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Stars       int       `json:"stargazers_count"`
	Forks       int       `json:"forks_count"`
	CreatedAt   time.Time `json:"created_at"`
}

func main() {
	if len(os.Args) != 3 {
		printUsage()
		os.Exit(1)
	}

	owner := os.Args[1]
	repo := os.Args[2]

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Network error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		fmt.Printf("Repository %s/%s not found\n", owner, repo)
		os.Exit(1)
	default:
		fmt.Printf("HTTP error %d: %s\n", resp.StatusCode, resp.Status)
		os.Exit(1)
	}

	var repository Repository
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		fmt.Printf("JSON parsing error: %v\n", err)
		os.Exit(1)
	}

	printRepositoryInfo(repository)
}

func printUsage() {
	fmt.Println("Usage (from task1/src): go run github_info.go <owner> <repo>")
	fmt.Println("Sample: github_info.go golang go")
}

func printRepositoryInfo(repo Repository) {
	fmt.Println("\nRepository information:")
	fmt.Println("")
	fmt.Printf("Name:\t\t%s\n", repo.Name)
	fmt.Printf("Description:\t%s\n", getDescription(repo.Description))
	fmt.Printf("Stars:\t\t%d\n", repo.Stars)
	fmt.Printf("Forks:\t\t%d\n", repo.Forks)
	fmt.Printf("Creation date:\t%s\n", repo.CreatedAt.Format("02.01.2006 15:04"))
}
func getDescription(desc string) string {
	if desc == "" {
		return "<no description>"
	}
	return desc
}