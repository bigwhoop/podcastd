package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Folder     string        `yaml:"folder"`
	FeedFormat string        `yaml:"feed_format"`
	ItemFormat string        `yaml:"item_format"`
	Interval   time.Duration `yaml:"interval"`
	Feeds      []FeedConfig  `yaml:"feeds,flow"`
}

type FeedConfig struct {
	FeedFormat   string            `yaml:"feed_format"`
	ItemFormat   string            `yaml:"item_format"`
	File         string            `yaml:"file"`
	Url          string            `yaml:"url"`
	Placeholders map[string]string `yaml:"placeholders"`
}

func loadConfig() Config {
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			logger.Panicf("Failed getting working dir: %v", err)
		}
		path = filepath.Join(cwd, "podcastd.yml")
	}

	configData, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Fatalf("Failed reading config.yml: %v", err)
	}

	config := Config{}

	if err = yaml.Unmarshal(configData, &config); err != nil {
		logger.Fatalf("Failed parsing config.yml: %v", err)
	}

	return config
}
