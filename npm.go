package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type (
	// Package represents an NPM package.
	Package struct {
		ID   string               `json:"_id"`
		Time map[string]time.Time `json:"time"`
	}

	// Release represents a release on NPM.
	Release struct {
		Name    string
		Version string
		Date    time.Time
	}

	// Releases represents a set of releases.
	Releases []Release
)

func (r Releases) Len() int {
	return len(r)
}

func (r Releases) Less(i, j int) bool {
	return r[j].Date.After(r[i].Date)
}

func (r Releases) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func getNodePackageURL(name string) string {
	pkg := url.URL{
		Scheme: "https",
		Host:   "www.npmjs.com",
		Path:   "/package/" + name,
	}

	return pkg.String()
}

func getNodePackageRegistryURL(name string) string {
	pkg := url.URL{
		Scheme: "https",
		Host:   "registry.npmjs.com",
		Path:   "/" + name,
	}

	return pkg.String()
}

func getNodePackageReleases(name string) []Release {
	cl := http.Client{
		Timeout: timeout,
	}

	url := getNodePackageRegistryURL(name)

	// Create request.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	// Execute request.
	resp, err := cl.Do(req)
	if err != nil {
		panic(err)
	}

	// Parse response body.
	pkg := Package{}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&pkg)
	if err != nil {
		panic(err)
	}

	releases := make([]Release, 0, len(pkg.Time))
	for ver, time := range pkg.Time {
		if ver == "modified" || ver == "created" {
			continue
		}

		releases = append(releases, Release{
			Name:    pkg.ID,
			Version: ver,
			Date:    time,
		})
	}

	return releases
}
