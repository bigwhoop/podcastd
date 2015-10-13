package main

import (
	"log"
	"os"
	"path/filepath"
)

var (
	logger *log.Logger
)

func init() {
	path := filepath.Join(os.TempDir(), "podcastd.log")

	loggerFile, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}

	logger = log.New(loggerFile, "", log.LstdFlags)
}
