package main

import (
	"fmt"
	"os"
	"strconv"

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

	fmt.Println("open http://localhost/readme/avelino/awesome-go or http://localhost/readme/vinta/awesome-python")
	runServer()
}

func runServer() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "80"
	}
	server := newServer()
	server.Run(":" + port)
}
