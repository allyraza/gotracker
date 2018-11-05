package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Config - collector config
type Config struct {
	Filepath string
	Store    string `json:"store"`
	Logfile  string `json:"logfile"`
	Addr     string `json:"addr"`
	Stats    string `json:"stats"`
	Kafka    struct {
		Brokers []string `json:"brokers"`
	}
	Verbose bool `json:"verbose"`
}

func (c *Config) ParseFile() {

	blob, err := ioutil.ReadFile(c.Filepath)
	if err != nil {
		log.Fatalf("Config file: %v\n", err)
	}

	if err := json.Unmarshal(blob, c); err != nil {
		log.Fatalf("Config parse: %v\n", err)
	}

}
