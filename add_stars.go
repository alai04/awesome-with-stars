package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	patternAwesomeRepo = `(?m)^\s*\*\s+\[[^\]]+\]\(https://github\.com/([\w-\.]+)/([\w-\.]+)(/.*)?\)`
	maxWorkers         = 10
)

func urlREADME(repo string) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/master/README.md", repo)
}

func filename(repo string) string {
	return fmt.Sprintf("README-%x.md", md5.Sum([]byte(repo)))
}

func getREADME(awesomeRepo string) (b []byte, err error) {
	res, err := http.Get(urlREADME(awesomeRepo))
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		return
	}

	b, err = ioutil.ReadAll(res.Body)
	return
}

func getRepoListFromREADME(awesomeRepo string) (rList []string, err error) {
	b, err := getREADME(awesomeRepo)
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
		log.Debug(repoFullName)
	}
	return
}

var nFinished uint32

func worker(id int, jobs <-chan string, results chan<- RepoInfo) {
	for repo := range jobs {
		rInfo, err := getInfoOfRepo(repo)
		if err == nil {
			results <- rInfo
		} else {
			log.Errorf("Get info of %s error: %v", repo, err)
			results <- RepoInfo{}
		}
	}
	atomic.AddUint32(&nFinished, 1)
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

	for atomic.LoadUint32(&nFinished) < maxWorkers {
		rInfo := <-chResult
		if rInfo.ID > 0 {
			log.Debugf("Get a repo info: %v", rInfo)
			reposInfo = append(reposInfo, rInfo)
		} else {
			log.Debug("Get a null repo info")
		}
	}
	return
}

func addStarsToREADME(awesomeRepo string) (err error) {
	b, err := getREADME(awesomeRepo)
	if err != nil {
		return
	}

	re := regexp.MustCompile(patternAwesomeRepo)
	matches := re.FindAllSubmatchIndex(b, -1)
	if len(matches) <= 0 {
		err = errors.New("Awesome repository not found")
		return
	}

	var cur int
	var rInfo RepoInfo
	store := NewMongoStore(false)
	out, err := os.Create(filename(awesomeRepo))
	if err != nil {
		return
	}
	defer out.Close()
	for _, match := range matches {
		rInfo.RepoName = string(b[match[2]:match[5]])
		// fmt.Printf("%s %v\n", rInfo.FullName, match)
		_, err = out.Write(b[cur:match[1]])
		if err != nil {
			return
		}
		err2 := store.Load(&rInfo)
		if err2 == nil {
			out.WriteString(fmt.Sprintf(" (★%d)", rInfo.Stars))
		}
		cur = match[1]
	}
	out.Write(b[cur:])
	return
}

func saveRepoInfos(awesomerepo string) (err error) {
	rList, err := getRepoListFromREADME(awesomerepo)
	if err != nil {
		log.Errorf("Get README of %s error: %v", awesomerepo, err)
		return
	}
	log.Debugf("There are %d repos in README file", len(rList))
	start := time.Now()
	rInfos := getReposInfo(rList)
	end := time.Now()
	log.Debugf("GetReposInfo use %v\n", end.Sub(start))

	store := NewMongoStore(false)
	for _, rInfo := range rInfos {
		store.Save(rInfo)
	}
	return
}

var mutexAddStars sync.Mutex

func prepareGetREADME(awesomeRepo string) {
	log.Debugf("Try to get README of %s", awesomeRepo)
	go func() {
		mutexAddStars.Lock()
		err := saveRepoInfos(awesomeRepo)
		if err == nil {
			err = addStarsToREADME(awesomeRepo)
		}
		if err != nil {
			log.Errorf("Error at addStarsToREADME: %v", err)
		}
		mutexAddStars.Unlock()
	}()
	return
}
