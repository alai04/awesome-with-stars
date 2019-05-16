package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
	if _, err := strconv.ParseBool(os.Getenv("DEBUG")); err == nil {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	fmt.Println("Hello, awesome-with-stars")
	// rInfo, _ := getInfoOfRepo("jaraco/path.py")
	// store := NewMongoStore(false)
	// store.Save(rInfo)
	// getRepoListFromREADME("vinta/awesome-python")
	// saveRepoInfos("vinta/awesome-python")
	// addStarsToREADME("vinta/awesome-python")

	runServer()
}

func saveRepoInfos(awesomerepo string) {
	rList, _ := getRepoListFromREADME(awesomerepo)
	log.Debugf("There are %d repos in README file", len(rList))
	start := time.Now()
	fmt.Printf("Start at: %v\n", start)
	rInfos := getReposInfo(rList)
	end := time.Now()
	fmt.Printf("Duration: %v\n", end.Sub(start))

	store := NewMongoStore(false)
	for _, rInfo := range rInfos {
		store.Save(rInfo)
	}
}

func runServer() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "80"
	}
	server := newServer()
	server.Run(":" + port)
}
