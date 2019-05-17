package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// RepoInfo holds repository information retrieved from api.github.com
type RepoInfo struct {
	RepoName string
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
	req.Header.Add("Authorization", "token "+os.Getenv("GITHUB_TOKEN")) //85c75c2e5ca6e7337b1fe87b5d9c254385bd0e0a")

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

	rInfo.RepoName = repo
	log.Debugf("%v", rInfo)
	return
}

// RepoInfoStore is interface for RepoInfo Store
type RepoInfoStore interface {
	Save(RepoInfo) error
	Load(*RepoInfo) error
}
