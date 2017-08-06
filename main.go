package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

const (
	packageJSONURL = "https://raw.githubusercontent.com/%s/%s/master/package.json"

	indexHTML = `<!DOCTYPE html>
<html lang="en">
	<head>
		<style>html{font:sans-serif}</style>
	</head>
	<body>
		<h1>npm-feed</h1>
		<p>This program provides a feed for dependencies of a package.</p>
		<p>You can use it by specifying a URL containing the following structure:</p>
		<div><pre>/organization/project/atom.xml</pre></div>
		<div><pre>/organization/project/rss.xml</pre></div>
		<h2>Examples:</h2>
		<ul>
		<li><a href="/facebook/react/rss.xml">React Dependencies RSS Feed</a></li>
		<li><a href="/angular/angular/atom.xml">Angular Dependencies Atom Feed</a></li>
		</ul>
	</body>
</html>
`
)

var (
	listenAddr    = ":8000"
	githubToken   = ""
	timeout       = 10 * time.Second
	cacheDuration = 5 * time.Minute

	cache = gocache.New(cacheDuration, 30*time.Second)
)

func init() {
	eStringVar(&listenAddr, "listen", "Address to listen on.")
	eStringVar(&githubToken, "github-token", "GitHub token.")
	eDurationVar(&timeout, "timeout", "HTTP timeout")
	eDurationVar(&cacheDuration, "cache", "Duration of cache.")

	flag.Parse()
}

func main() {
	server := http.Server{
		Addr:         listenAddr,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		Handler:      http.HandlerFunc(serve),
	}

	log.Fatalln(server.ListenAndServe())
}

func serve(rw http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		rw.Header().Add("Content-Type", "text/html; charset=utf8")
		rw.Write([]byte(indexHTML))
	default:
		path := r.URL.Path
		if len(path) > 1 && path[0] == '/' {
			path = path[1:]
		}

		p := strings.Split(path, "/")
		if len(p) != 3 {
			rw.WriteHeader(404)
			return
		}

		org, project, feed := p[0], p[1], p[2]
		rw.Write(getData(org, project, feed, getAbsURL(r)))
	}
}

func getAbsURL(r *http.Request) string {
	uri := url.URL{}

	uri.Scheme = "http"
	if r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
		uri.Scheme = "https"
	}

	uri.Host = r.Host
	uri.Path = r.URL.Path

	return uri.String()
}

func getData(org, project, feed, uri string) []byte {
	cacheKey := fmt.Sprintf("%s:%s:%s", org, project, feed)
	data, found := cache.Get(cacheKey)

	content, ok := data.([]byte)

	if found && ok {
		return content
	}

	if feed != "rss.xml" && feed != "atom.xml" {
		return nil
	}

	releases := Releases{}

	deps := getPackageJSONDeps(org, project)
	for _, dep := range deps {
		pkgrels := getNodePackageReleases(dep)
		for _, rel := range pkgrels {
			releases = append(releases, rel)
		}
	}

	sort.Sort(sort.Reverse(releases))

	if len(releases) > 20 {
		releases = releases[0:20]
	}

	switch feed {
	case "rss.xml":
		content = rss(releases, uri)
	case "atom.xml":
		content = atom(releases, uri)
	}

	cache.Set(cacheKey, content, cacheDuration)

	return content
}
