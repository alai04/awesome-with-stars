package main

import (
	"encoding/json"
	"errors"
	"log"
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
	res, err := http.Get("https://api.github.com/repos/" + repo)
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

	log.Printf("%v", rInfo)
	return
}
