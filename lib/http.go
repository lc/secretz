package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// GitHubAPI is used to unmarshal Github's API response. This holds usernames.
type GitHubAPI []struct {
	Login string `json:"login"`
}

var key = GetAPIKey()

// Secretz is the http client for the tool.
var Secretz = &http.Client{
	Timeout: time.Second * 10,
}

// OrgMembers queries GitHub's API and gets a list of members for a given org.
func OrgMembers(target string) (g *GitHubAPI) {
	memStruct := new(GitHubAPI)
	target = fmt.Sprintf("https://api.github.com/orgs/%s/public_members", target)
	resp, err := Secretz.Get(target)
	if err != nil {
		log.Fatalf("Error creating HTTP Request: %v", err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&memStruct)
	if err != nil {
		log.Fatalf("Error parsing response from GitHub: %v", err)
	}
	return memStruct
}

// QueryApi is used to query Travis-CI's API
func QueryApi(target string) (body []byte) {
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		log.Fatalf("Error creating HTTP Request: %v", err)
	}
	req.Header.Add("User-Agent", `API Explorer`)
	if key != "" {
		token := fmt.Sprintf("token %s", key)
		req.Header.Add("Authorization", token)
	}
	req.Header.Add("Travis-API-Version", `3`)
	resp, err := Secretz.Do(req)
	if err != nil {
		log.Fatalf("Could not request API: %v", err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("\nTravisCI responded with a non-200 statuscode. You're likely being rate-limited")
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error parsing response from TravisCI: %v", err)
	}
	return bytes
}
