package main

import "testing"

const (
	testRepo1 = "shimberger/gohls"
	testRepo2 = "shimberger/gohls2"
	testRepo3 = "shimberger2/gohls"
)

func TestGetInfoOfRepo(t *testing.T) {
	rInfo, err := getInfoOfRepo(testRepo1)
	if err != nil {
		t.Errorf("Get information of %s error: %v", testRepo1, err)
	}
	if rInfo.ID != 39337344 || rInfo.FullName != testRepo1 {
		t.Errorf("Get information of %s return: %v", testRepo1, rInfo)
	}
}

func TestGetInfoOfWrongRepo(t *testing.T) {
	_, err := getInfoOfRepo(testRepo2)
	if err == nil {
		t.Errorf("Get information of %s should error, but successful", testRepo2)
	}

	_, err = getInfoOfRepo(testRepo3)
	if err == nil {
		t.Errorf("Get information of %s should error, but successful", testRepo3)
	}
}
