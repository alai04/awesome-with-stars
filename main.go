package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, awesome-with-stars")
	// getInfoOfRepo("shimberger/gohls")
	rList, _ := getRepoListFromREADME("avelino/awesome-go")
	start := time.Now()
	fmt.Printf("Start at: %v\n", start)
	rInfos := getReposInfo(rList[:100])
	end := time.Now()
	fmt.Printf("Duration: %v\n", end.Sub(start))

	store := NewMongoStore(false)
	for _, rInfo := range rInfos {
		store.Save(rInfo)
	}
}
