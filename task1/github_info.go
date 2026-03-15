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

func (r Repository) String() string {
	desc := r.Description
	if desc == "" {
		desc = "<no description>"
	}

	return fmt.Sprintf("\nRepository information:\n"+
		"\n"+
		"Name:\t\t%s\n"+
		"Description:\t%s\n"+
		"Stars:\t\t%d\n"+
		"Forks:\t\t%d\n"+
		"Creation date:\t%s",
		r.Name, desc, r.Stars, r.Forks, r.CreatedAt.Format("02.01.2006 15:04"))
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

	fmt.Println(repository)
}

func printUsage() {
	fmt.Println("Usage (from task1): github_info <owner> <repo>")
	fmt.Println("Sample: github_info golang go")
}
