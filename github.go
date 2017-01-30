package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	// Manifest specifies the package.json format.
	Manifest struct {
		DevDependencies map[string]string `json:"devDependencies"`
		Dependencies    map[string]string `json:"dependencies"`
	}
)

func getPackageJSONDeps(org, project string) []string {
	cl := http.Client{
		Timeout: timeout,
	}

	// Create request.
	req, err := http.NewRequest("GET", fmt.Sprintf(packageJSONURL, org, project), nil)
	if err != nil {
		panic(err)
	}

	// Add headers.
	req.Header.Add("Accept", "application/vnd.github.v3.raw")
	if githubToken != "" {
		req.Header.Add("Authorization", "token "+githubToken)
	}

	// Execute request.
	resp, err := cl.Do(req)
	if err != nil {
		panic(err)
	}

	// Parse response body.
	manifest := Manifest{}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&manifest)
	if err != nil {
		panic(err)
	}

	dependencies := make([]string, 0, len(manifest.Dependencies)+len(manifest.DevDependencies))
	for dep := range manifest.Dependencies {
		dependencies = append(dependencies, dep)
	}
	for dep := range manifest.DevDependencies {
		dependencies = append(dependencies, dep)
	}

	return dependencies
}
