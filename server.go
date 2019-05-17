package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shurcooL/github_flavored_markdown"
	"github.com/shurcooL/github_flavored_markdown/gfmstyle"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func newServer() *negroni.Negroni {
	n := negroni.Classic()
	mx := mux.NewRouter()
	mx.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(gfmstyle.Assets)))
	mx.HandleFunc("/readme/{user}/{repo}", readmeHandler).Methods("GET")
	n.UseHandler(mx)
	return n
}

func readmeHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	repo := fmt.Sprintf("%s/%s", vars["user"], vars["repo"])
	log.Infof("GET readme of %s", repo)
	b, err := ioutil.ReadFile(filename(repo))
	if err != nil {
		msg := fmt.Sprintf("%s\nRetry after minutes ...", err.Error())
		prepareGetREADME(repo)
		http.Error(w, msg, http.StatusNotFound)
		return
	}
	w.Write(github_flavored_markdown.Markdown(b))
}
