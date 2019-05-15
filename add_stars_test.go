package main

import "testing"

const (
	awesome1      = "avelino/awesome-go"
	awesomeWrong1 = "avelino2/awesome-go"
)

func TestGetReposInfo(t *testing.T) {
	rList, err := getRepoListFromREADME(awesome1)
	if err != nil {
		t.Errorf("Get repo list from %s error: %v", awesome1, err)
		return
	}
	if len(rList) < 100 {
		t.Errorf("Repos in %s is too less, only %d repos found", awesome1, len(rList))
		return
	}

	const nTestRepos = 10
	rInfos := getReposInfo(rList[:nTestRepos])
	if len(rInfos) != nTestRepos {
		t.Errorf("Request %d repos info, but get %d", nTestRepos, len(rInfos))
	}
}

func TestGetRepoListFromREADMEWrong(t *testing.T) {
	_, err := getRepoListFromREADME(awesomeWrong1)
	if err == nil {
		t.Errorf("Get repo list from %s should error, but successful", awesome1)
	}
}
