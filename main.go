package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"os"
	"sync"
	"time"

	"secretz/lib"
)

import jsoniter "github.com/json-iterator/go"

var repos []string

// TravisCI is the struct that holds repo slugs.
type TravisCI struct {
	Repositories []struct {
		Slug   string `json:"slug"`
		Active bool   `json:"active"`
	} `json:"repositories"`
}

// Builds is the struct that build information from TravisCI
type Builds struct {
	Type       string `json:"@type"`
	Pagination struct {
		IsFirst bool `json:"is_first"`
		IsLast  bool `json:"is_last"`
	} `json:"@pagination"`
	Builds []struct {
		Jobs []struct {
			JobId int `json:"id"`
		} `json:"jobs"`
	} `json:"builds"`
}

var (
	setkey      string
	org         string
	timeout     int
	concurrency int
	delay       int
	// GetMembers is a flag with the options: list or scan
	GetMembers string
	targets    []string
)

func init() {
	flag.Usage = lib.Usage
	flag.StringVar(&setkey, "setkey", "", "Set API Key for api.travis-ci.org")
	flag.StringVar(&org, "t", "", "Target organization")
	flag.IntVar(&timeout, "timeout", 30, "Timeout for the tool in seconds")
	flag.IntVar(&concurrency, "c", 5, "Number of concurrent fetchers")
	flag.StringVar(&GetMembers, "members", "", "Get GitHub Org Members")
	flag.IntVar(&delay, "delay", 600, "delay between requests + random delay/2 jitter")
}

// Job struct holds a JobID and the Organization name.
type Job struct {
	ID  int
	Org string
}

func main() {
	flag.Parse()

	lib.Secretz.Timeout = time.Duration(timeout) * time.Second
	if setkey != "" {
		lib.SetAPIKey(setkey)
	}
	if len(org) < 1 || setkey != "" {
		log.Fatalf("Usage: %s -t <organization>\n", os.Args[0])
	}
	if GetMembers == "" {
		targets = append(targets, org)
	} else {
		switch GetMembers {
		case "list":
			GHMem := lib.OrgMembers(org)
			for _, Member := range *GHMem {
				fmt.Println(Member.Login)
			}
			os.Exit(0)
		case "scan":
			GHMem := lib.OrgMembers(org)
			for _, Member := range *GHMem {
				targets = append(targets, Member.Login)
			}
			break
		default:
			log.Fatalf("Invalid option specified in -members flag!\n")
		}
	}
	for _, org := range targets {
		lib.CreateOrg(org)
		log.Printf("Fetching repos for %s\n", org)
		ParseResponse(org)
		log.Printf("Fetching builds for %s's repos\n", org)

		jobChan := make(chan *Job)
		finishedChan := make(chan string)
		var wg, wg2 sync.WaitGroup
		wg2.Add(1)

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				SaveLogs(jobChan, finishedChan)
			}()
		}

		go func() {
			defer wg2.Done()
			for str := range finishedChan {
				log.Printf("%s\n", str)
			}
		}()

		go func() {
			for _, slug := range repos {
				builds := GetBuilds(slug)
				for _, job := range builds {
					jobChan <- &Job{ID: job, Org: org}
				}
			}
			close(jobChan)
		}()

		wg.Wait()
		close(finishedChan)
		wg2.Wait()
	}
}

// ParseResponse gets the JSON response from Travis and parses it for repo slugs
func ParseResponse(org string) {
	for {
		api := fmt.Sprintf("https://api.travis-ci.org/owner/%s/repos?limit=100&offset=%d", org, len(repos))
		res := lib.QueryApi(api)

		ciResp := new(TravisCI)
		err := jsoniter.Unmarshal(res, ciResp)
		if err != nil {
			log.Fatalf("Could not decode json: %s\n", err)
		}

		if len(ciResp.Repositories) == 0 {
			break
		}
		for _, repo := range ciResp.Repositories {
			repos = append(repos, repo.Slug)
		}
	}
}

// GetBuilds gets all JobId's from builds of a repo.
func GetBuilds(slug string) []int {

	builds := new(Builds)
	jobs := []int{}
	offset := 0

	for {
		log.Printf("Fetching builds %s [offset: %d]\n", slug, offset)
		api := fmt.Sprintf("https://api.travis-ci.org/repo/%s/builds?limit=100&offset=%d", url.QueryEscape(slug), offset)
		res := lib.QueryApi(api)

		// delay + jitter
		if delay > 0 {
			time.Sleep((time.Duration(delay + rand.Intn(delay/2))) * time.Millisecond)
		}

		err := jsoniter.Unmarshal(res, builds)
		if err != nil {
			log.Fatalf("Could not decode json: %s\n", err)
		}
		if builds.Type == "error" {
			break
		}

		loop := len(builds.Builds)
		i := 0
		for i < loop {
			for _, z := range builds.Builds[i].Jobs {
				jobs = append(jobs, z.JobId)
			}
			i++
		}

		if builds.Pagination.IsLast {
			break
		}
		offset += 100
	}

	return jobs
}

// SaveLogs saves build logs for given job ids
func SaveLogs(jobChan chan *Job, resultChan chan string) {
	for job := range jobChan {
		api := fmt.Sprintf("https://api.travis-ci.org/job/%d/log.txt", job.ID)
		// delay + jitter
		if delay > 0 {
			time.Sleep((time.Duration(delay + rand.Intn(delay/2))) * time.Millisecond)
		}
		res := lib.QueryApi(api)
		file := fmt.Sprintf("output/%s/%d.txt", job.Org, job.ID)
		err := ioutil.WriteFile(file, res, 0644)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		resultChan <- fmt.Sprintf("Wrote log %d to %s", job.ID, file)
	}
}
