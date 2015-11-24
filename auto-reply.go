package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/parkr/auto-reply/Godeps/_workspace/src/github.com/google/go-github/github"
	"github.com/parkr/auto-reply/Godeps/_workspace/src/golang.org/x/oauth2"
	"github.com/parkr/auto-reply/deprecate"
)

var (
	client *github.Client

	deprecatedRepos = []deprecate.RepoDeprecation{
		deprecate.RepoDeprecation{
			Nwo:     "jekyll/jekyll-help",
			Message: `This repository is no longer maintained. If you're still experiencing this problem, please search for your issue on [Jekyll Talk](https://talk.jekyllrb.com/), our new community forum. If it isn't there, feel free to post to the Help category and someone will assist you. Thanks!`,
		},
	}
)

func newClient() *github.Client {
	if token := os.Getenv("GITHUB_ACCESS_TOKEN"); token != "" {
		return github.NewClient(oauth2.NewClient(
			oauth2.NoContext,
			oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
			),
		))
	} else {
		log.Fatal("GITHUB_ACCESS_TOKEN required")
		return nil
	}
}

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "The port to serve to")
	flag.Parse()
	client = newClient()

	deprecationHandler := deprecate.NewHandler(client, deprecatedRepos)
	http.Handle("/_github/repos/deprecated", deprecationHandler)

	log.Printf("Listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
