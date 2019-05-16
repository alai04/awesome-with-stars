package main

import (
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
	mx.HandleFunc("/readme/{repo}", readmeHandler).Methods("GET")
	n.UseHandler(mx)
	return n
}

func readmeHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	repo := vars["repo"]
	log.Printf("GET readme of %s", repo)
	b, err := ioutil.ReadFile("output/README.md")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.Write(github_flavored_markdown.Markdown(b))
}
