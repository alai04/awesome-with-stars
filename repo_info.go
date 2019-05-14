package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

// RepoInfo holds repository information retrieved from api.github.com
type RepoInfo struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Fork     bool   `json:"fork"`
	Stars    int    `json:"stargazers_count"`
	Forks    int    `json:"forks_count"`
}

// GET https://api.github.com/repos/:owner/:repo
func getInfoOfRepo(repo string) (rInfo RepoInfo, err error) {
	req, err := http.NewRequest("GET", "https://api.github.com/repos/"+repo, nil)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", "token 2f5668bcfea11f480b2581cc01488c9c07f758ed")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		return
	}

	err = json.NewDecoder(res.Body).Decode(&rInfo)
	if err != nil {
		return
	}

	// log.Printf("%v", rInfo)
	return
}

// RepoInfoStore is interface for RepoInfo Store
type RepoInfoStore interface {
	Save(RepoInfo) error
	Load(*RepoInfo) error
}
