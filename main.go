package main

import (
	"github.com/bigwhoop/podcastd/feeds"
	"github.com/nareix/curl"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	VERSION           = "0.9.0"
	NUM_FEED_WORKERS  = 3
	NUM_SIM_DOWNLOADS = 4
)

type downloadItemJob struct {
	feed       feeds.Feed
	feedPath   string
	config     Config
	feedConfig FeedConfig
	item       feeds.FeedItem
}

func run() {
	config := loadConfig()

	wg := new(sync.WaitGroup)
	wg.Add(1)

	feedConfigChan := make(chan FeedConfig)
	downloadItemChan := make(chan downloadItemJob)

	for i := 0; i < NUM_FEED_WORKERS; i++ {
		go downloadFeedWorker(config, feedConfigChan, downloadItemChan)
	}

	for i := 0; i < NUM_SIM_DOWNLOADS; i++ {
		go downloadFeedItemWorker(downloadItemChan)
	}

	go func() {
		for {
			for _, feedConfig := range config.Feeds {
				feedConfigChan <- feedConfig
			}

			logger.Printf("Checking feeds again in %s", config.Interval)
			time.Sleep(config.Interval)
		}
	}()

	wg.Wait()
}

func downloadFeedWorker(config Config, feedConfigChan chan FeedConfig, downloadItemChan chan downloadItemJob) {
	logger.Printf("Started feed worker ...")

	for feedConfig := range feedConfigChan {
		logger.Printf("Checking feed '%s'", feedConfig.Url)
		ui.SetFeedStatus(feedConfig.Url, true)

		feed, err := feeds.Download(feedConfig.Url)
		ui.SetFeedStatus(feedConfig.Url, false)

		if err != nil {
			logger.Printf("Failed downloading feed %s: %v", feedConfig.Url, err)
			continue
		}

		feedPath := filepath.Join(config.Folder, getFeedPath(config, feedConfig, feed))

		for _, item := range feed.Items {
			job := downloadItemJob{
				feed,
				feedPath,
				config,
				feedConfig,
				item,
			}
			go func() { downloadItemChan <- job }()
		}
	}
}

func downloadFeedItemWorker(downloadItemChan chan downloadItemJob) {
	logger.Printf("Started download worker ...")

	for job := range downloadItemChan {
		filePath := filepath.Join(job.feedPath, getFilePath(job.config, job.feedConfig, job.feed, job.item))
		folderPath := filepath.Dir(filePath)

		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			os.MkdirAll(folderPath, 0755)
		}

		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			if err = downloadFile(job.item.Title, job.item.FileUrl, filePath); err != nil {
				logger.Printf("Failed downloading file: %v", err)
			}
		} else if err != nil {
			logger.Printf("Failed writing file: %v", err)
		}
	}
}

func getFeedPath(config Config, feedConfig FeedConfig, feed feeds.Feed) string {
	values := map[string]string{
		"feed.title": feed.Title,
	}

	for k, v := range feedConfig.Placeholders {
		values[k] = v
	}

	format := feedConfig.FeedFormat
	if format == "" {
		format = config.FeedFormat
	}

	return replacePathPlaceholders(values, format)
}

func getFilePath(config Config, feedConfig FeedConfig, feed feeds.Feed, item feeds.FeedItem) string {
	values := map[string]string{
		"feed.title":    feed.Title,
		"item.title":    item.Title,
		"item.category": item.Category,
		"item.author":   item.Author,
		"item.date":     item.PubDate.Format("2006-01-02"),
		"item.time":     item.PubDate.Format("150405"),
	}

	for k, v := range feedConfig.Placeholders {
		values[k] = v
	}

	format := feedConfig.ItemFormat
	if format == "" {
		format = config.ItemFormat
	}

	return replacePathPlaceholders(values, format) + filepath.Ext(item.FileUrl)
}

func downloadFile(title, src, dst string) error {
	logger.Printf("Downloading '%s' from '%s' to '%s'", title, src, dst)

	tmpDst := dst + ".part"

	req := curl.Get(src)
	req.SaveToFile(tmpDst)
	req.Progress(func(p curl.ProgressStatus) {
		if p.Stat == curl.Downloading {
			ui.SetActiveDownload(src, p)
		} else if p.Stat == curl.Closed {
			ui.RemoveActiveDownload(src)
		}
	}, time.Millisecond*100)

	if _, err := req.Do(); err != nil {
		return err
	}

	if err := os.Rename(tmpDst, dst); err != nil {
		return err
	}

	logger.Printf("Downloaded '%s' to '%s", title, dst)

	return nil
}

func replacePathPlaceholders(placeholders map[string]string, s string) string {
	for k, v := range placeholders {
		v = sanitizeForFileName(v)
		s = strings.Replace(s, "%"+k+"%", v, -1)
	}
	return s
}

func sanitizeForFileName(s string) string {
	s = strings.Replace(s, ": ", " - ", -1)
	r := regexp.MustCompile("[^0-9a-zA-Z-.,;_() äöüÄÖÜß]")
	s = r.ReplaceAllLiteralString(s, "")
	s = strings.TrimSpace(s)
	s = strings.Replace(s, "  ", " ", -1)

	return s
}
