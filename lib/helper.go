package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

var (
	home = HomeDir()
	path = home + "/.config/secretz/config.json"
	c    Config
)

type Config struct {
	TravisCIOrg string `json:"TravisCIOrgKey"`
}

func Usage() {
	help := `Usage: secretz -t <organization> [options]

  -c int
		Number of concurrent fetchers (default 5)
		
  -delay int
		delay between requests + random delay/2 jitter (default 600)
		
  -members string
		Retrieve members of Github Org parameters: [list | scan]
		
  -setkey string
		Set API Key for api.travis-ci.org
		
  -t string
		Target organization
		
  -timeout int
		Timeout for the tool in seconds (default 30)`
	fmt.Printf(help)

}
func GetAPIKey() string {
	if Exists(path) {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("Error opening config file: %v", err)
		}
		mErr := json.Unmarshal(data, &c)
		if mErr != nil {
			c.TravisCIOrg = ""
		}
	}
	CreateOutputDir(home + "/.config/secretz")
	return c.TravisCIOrg
}
func SetAPIKey(key string) {
	c := Config{key}
	bytes, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("Error marshalling data: %v", err)
	}
	CreateOutputDir(home + "/.config/secretz/")
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		log.Fatalf("Error creating config file: %v\n", err)
	}
	fmt.Printf("Created config file: %s\n", path)
}

func Exists(loc string) bool {
	_, err := os.Stat(loc)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true

}
func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Error: %v", err)
		os.Exit(1)
	}
	return usr.HomeDir
}
func CreateOrg(dir string) {
	CreateOutputDir("output")
	dir = fmt.Sprintf("output/%s", dir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			log.Fatalf("Error: %v", err)
			os.Exit(1)
		}
	}
}
func CreateOutputDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
}
