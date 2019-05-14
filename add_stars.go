package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

const (
	patternAwesomeRepo = `(?m)^\*\s+\[.+\]\(https://github\.com/([\w-]+)/([\w-]+)(/.+)?\)`
	maxWorkers         = 10
)

func urlREADME(repo string) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/master/README.md", repo)
}

func getRepoListFromREADME(repo string) (rList []string, err error) {
	res, err := http.Get(urlREADME(repo))
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		return
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	re := regexp.MustCompile(patternAwesomeRepo)
	matches := re.FindAllSubmatch(b, -1)
	if len(matches) <= 0 {
		err = errors.New("Awesome repository not found")
		return
	}
	for _, match := range matches {
		repoFullName := fmt.Sprintf("%s/%s", match[1], match[2])
		rList = append(rList, repoFullName)
	}
	return
}

var nFinished int

func worker(id int, jobs <-chan string, results chan<- RepoInfo) {
	for repo := range jobs {
		rInfo, err := getInfoOfRepo(repo)
		if err == nil {
			results <- rInfo
		} else {
			log.Printf("Get info of %s error: %v", repo, err)
			results <- RepoInfo{}
		}
	}
	nFinished++
}

func getReposInfo(repos []string) (reposInfo []RepoInfo) {
	chRequest := make(chan string, len(repos))
	chResult := make(chan RepoInfo, len(repos))

	for w := 0; w < maxWorkers; w++ {
		go worker(w, chRequest, chResult)
	}

	for _, repo := range repos {
		chRequest <- repo
	}
	close(chRequest)

	for nFinished < maxWorkers {
		rInfo := <-chResult
		fmt.Println(rInfo)
		reposInfo = append(reposInfo, rInfo)
	}
	return
}
