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

// Config struct is used to unmarshal our config file into.
type Config struct {
	TravisCIOrg string `json:"TravisCIOrgKey"`
}

// GetConfig opens the config file and returns it
func GetConfig() Config {
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
	return c
}

// GetAPIKey opens the config file and returns an API key
func GetAPIKey() string {
	return GetConfig().TravisCIOrg
}

// SetAPIKey takes an API key and saves it to a config file.
func SetAPIKey(key string) {
	c := GetConfig()
	c.TravisCIOrg = key

	_ = writeConfig(c)
}

func writeConfig(c Config) error {
	bytes, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("Error marshalling data: %v", err)
		return err
	}

	alreadyExists := Exists(path)

	CreateOutputDir(home + "/.config/secretz/")
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		if alreadyExists {
			log.Fatalf("Error updating config file: %v\n", err)
		} else {
			log.Fatalf("Error creating config file: %v\n", err)
		}
		return err
	}

	if alreadyExists {
		fmt.Printf("Updated config file: %s\n", path)
	} else {
		fmt.Printf("Created config file: %s\n", path)
	}

	return nil
}

// Exists checks if the location provided exists or not.
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

// HomeDir returns the user's home directory.
func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Error: %v", err)
		os.Exit(1)
	}
	return usr.HomeDir
}

// CreateOrg creates a directory for a given Organization to hold the build logs.
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

// CreateOutputDir creates a directory from a given path
func CreateOutputDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
}
