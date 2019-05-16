package main

import "testing"

func TestSaveAndLoadRepoInfo(t *testing.T) {
	store := NewMongoStore(true)
	rInfo := RepoInfo{"mattn/go-colorable", 22405117, "mattn/go-colorable", false, 356, 47}
	err := store.Save(rInfo)
	if err != nil {
		t.Errorf("Save RepoInfo to mongo error: %v", err)
	}

	rInfo2 := RepoInfo{RepoName: "mattn/go-colorable"}
	err = store.Load(&rInfo2)
	if err != nil {
		t.Errorf("Load RepoInfo from mongo error: %v", err)
	}
	if rInfo2.ID != rInfo.ID {
		t.Errorf("Load RepoInfo expected %v, return %v", rInfo, rInfo2)
	}

	rInfo3 := RepoInfo{RepoName: "mattn2/go-colorable"}
	err = store.Load(&rInfo3)
	if err == nil {
		t.Errorf("Load RepoInfo from mongo should error, but successful")
	}
}

func TestSaveAndLoadRepoNameWithDot(t *testing.T) {
	store := NewMongoStore(true)
	rInfo := RepoInfo{"jaraco/path.py", 4093296, "jaraco/path.py", false, 862, 116}
	err := store.Save(rInfo)
	if err != nil {
		t.Errorf("Save RepoInfo to mongo error: %v", err)
	}

	rInfo2 := RepoInfo{RepoName: "jaraco/path.py"}
	err = store.Load(&rInfo2)
	if err != nil {
		t.Errorf("Load RepoInfo from mongo error: %v", err)
	}
	if rInfo2.ID != rInfo.ID {
		t.Errorf("Load RepoInfo expected %v, return %v", rInfo, rInfo2)
	}
}
